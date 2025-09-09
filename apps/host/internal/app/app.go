package app

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
	"github.com/m1thrandir225/imperium/apps/host/internal/host"
	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

type App struct {
	Name           string
	Bus            *EventBus
	State          *state.StateManager
	AuthBaseURL    string
	ProgramService *programs.ProgramService
	AuthService    *auth.AuthService
	StatusManager  *host.StatusManager
	HTTPClient     *httpclient.Client
	tokenRefresher *AuthTokenRefresher
}

func New(name string) (*App, error) {
	sm, err := state.NewStateManager(name)
	if err != nil {
		return nil, err
	}

	st := sm.Get()

	a := &App{
		Name:        name,
		Bus:         NewEventBus(),
		State:       sm,
		AuthBaseURL: st.Settings.ServerAddress,
	}

	a.buildClients()

	return a, nil
}

func (a *App) Start() {
	a.Bus.Publish(EventStateLoaded, a.State.Get())
	a.Bus.Publish(EventAppStarted, nil)
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

	httpClient := httpclient.NewClient(a.AuthBaseURL, tok, tok)

	a.HTTPClient = httpClient

	a.AuthService = auth.NewService(
		a.AuthBaseURL,
		httpClient,
	)

	configDir, err := util.GetConfigDir(a.Name)
	if err != nil {
		log.Printf("failed to get config directory: %v", err)
	}

	dbPath := filepath.Join(configDir, "programs.db")
	a.ProgramService = programs.NewService(
		dbPath,
		a.State.Get().Settings.RawgAPIKey,
	)

	a.stopTokenRefresher()
	a.tokenRefresher = NewAuthTokenRefresher(tok, tok)
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

	a.StatusManager = host.NewStatusManager(
		st.HostInfo.ID,
		a.AuthBaseURL,
		a.HTTPClient,
	)
	a.StatusManager.Start(context.Background())
}

func (a *App) stopTokenRefresher() {
	if a.tokenRefresher != nil {
		a.tokenRefresher.Stop()
		a.tokenRefresher = nil
	}
}
