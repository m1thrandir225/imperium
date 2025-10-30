// Package app
package app

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
	"github.com/m1thrandir225/imperium/apps/host/internal/events"
	"github.com/m1thrandir225/imperium/apps/host/internal/host"
	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/httpserver"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
	"github.com/m1thrandir225/imperium/apps/host/internal/session"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
	tokenrefresher "github.com/m1thrandir225/imperium/apps/host/internal/tokenrefresher"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
)

type App struct {
	Name           string
	Bus            events.EventBroker
	State          state.StateManager
	AuthBaseURL    string
	ProgramService programs.Service
	AuthService    auth.Service
	StatusManager  host.StatusManager
	HTTPClient     *httpclient.Client
	tokenRefresher tokenrefresher.Refresher
	SessionService session.Service
	HTTPServer     *httpserver.Server
}

func New(name string) (*App, error) {
	sm, err := state.NewPersistedStateManager(name)
	if err != nil {
		return nil, err
	}

	st := sm.Get()

	a := &App{
		Name:        name,
		Bus:         events.NewInMemoryBroker(),
		State:       sm,
		AuthBaseURL: st.Settings.ServerAddress,
	}

	a.buildClients()

	return a, nil
}

func (a *App) Start() {
	a.Bus.Publish(EventStateLoaded, a.State.Get())
	a.Bus.Publish(EventAppStarted, nil)
	a.Bus.Publish(EventHostInitRequested, nil)

	a.checkAndRefreshTokensAtStartup()
	log.Println("App started")
}

func (a *App) NeedsLogin() bool {
	st := a.State.Get()
	if st.UserInfo.ID == "" || st.UserInfo.Email == "" {
		return true
	}

	if st.UserSession.AccessToken == "" {
		return true
	}

	return time.Now().After(st.UserSession.AccessTokenExpiresAt)
}

// Rebuild clients/services from current state
func (a *App) buildClients() {
	tok := &stateTokens{
		sm:        a.State,
		baseURLFn: func() string { return a.AuthBaseURL },
	}

	httpClient, err := httpclient.NewClient(a.AuthBaseURL, tok, tok)
	if err != nil {
		log.Printf("failed to create a new http client")
		panic(err) // should panic
	}

	a.HTTPClient = httpClient

	authService, err := auth.NewService(
		a.AuthBaseURL,
		httpClient,
	)
	if err != nil {
		panic(err)
	}

	a.AuthService = authService

	configDir, err := util.GetConfigDir(a.Name)
	if err != nil {
		log.Printf("failed to get config directory: %v", err)
	}

	dbPath := filepath.Join(configDir, "programs.db")
	programService, err := programs.NewService(
		dbPath,
		a.State.Get().Settings.RawgAPIKey,
	)
	if err != nil {
		panic(err)
	}

	a.ProgramService = programService

	a.stopTokenRefresher()

	tokenRefresher, err := tokenrefresher.NewTokenRefresher(tok, tok)
	if err != nil {
		panic(err) // FIXME: should panic here??
	}

	a.tokenRefresher = tokenRefresher
	a.tokenRefresher.Start(context.Background())

	a.startStatusManagerIfReady()
}

func (a *App) startStatusManagerIfReady() {
	if a.StatusManager != nil {
		a.StatusManager.Stop()
		a.StatusManager = nil
	}

	st := a.State.Get()

	if st.UserInfo.ID == "" || st.UserInfo.Email == "" {
		return
	}

	if st.UserSession.AccessToken == "" {
		return
	}

	if st.HostInfo.ID == "" {
		a.Bus.Publish(EventHostInitRequested, nil)
		return
	}

	statusManager, err := host.NewInMemoryStatusManager(
		st.HostInfo.ID,
		a.AuthBaseURL,
		a.HTTPClient,
		a.SessionService,
	)
	if err != nil {
		panic(err) // should panic, invalid status manager
	}
	a.StatusManager = statusManager
	a.StatusManager.Start(context.Background())
}

func (a *App) stopTokenRefresher() {
	log.Println("Stopping token refresher")
	if a.tokenRefresher != nil {
		a.tokenRefresher.Stop()
		a.tokenRefresher = nil
	}
}

func (a *App) buildSessionService() {
	recorder, err := video.NewRecorder(
		&video.Config{
			Encoder:    a.State.Get().Settings.Encoder,
			FPS:        a.State.Get().Settings.Framerate,
			FFMPEGPath: a.State.Get().Settings.FFmpegPath,
		},
	)

	if err != nil {
		panic(err)
	}

	sessionService, err := session.NewService(
		a.AuthBaseURL,
		a.State.Get().UserSession.AccessToken,
		a.AuthService,
		a.ProgramService,
		recorder,
		nil,
	)
	if err != nil {
		panic(err) // should panic invalid SessionService
	}

	a.SessionService = sessionService
}

func (a *App) checkAndRefreshTokensAtStartup() {
	st := a.State.Get()

	if st.UserInfo.ID == "" || st.UserInfo.Email == "" {
		return
	}

	if st.UserSession.AccessToken == "" || st.UserSession.RefreshToken == "" {
		return
	}

	if time.Now().After(st.UserSession.RefreshTokenExpiresAt) {
		log.Println("Refresh token expired, user needs to login again.")
		return
	}

	accessTokenExpiresAt := st.UserSession.AccessTokenExpiresAt

	if time.Now().After(accessTokenExpiresAt) {
		log.Println("Access token expired, attempting to refresh.")

		tok := &stateTokens{
			sm:        a.State,
			baseURLFn: func() string { return a.AuthBaseURL },
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := tok.RefreshToken(ctx); err != nil {
			log.Printf("failed to refresh token: %v", err)
		} else {
			log.Println("Token refreshed successfully.")
		}
	}
}

func (a *App) Stop() {
	a.stopTokenRefresher()

	if a.StatusManager != nil {
		log.Println("Stopping status manager")
		a.StatusManager.Stop()
		a.StatusManager = nil
	}

	if a.SessionService != nil {
		log.Println("Stopping session service")
		a.SessionService.EndSession()
		a.SessionService = nil
	}

	if a.HTTPServer != nil {
		log.Println("Stopping HTTP server")
		a.HTTPServer.Stop()
		a.HTTPServer = nil
	}
}
