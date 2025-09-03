package ui

import (
	"context"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
	"github.com/m1thrandir225/imperium/apps/host/internal/config"
	"github.com/m1thrandir225/imperium/apps/host/internal/host"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
	"github.com/m1thrandir225/imperium/apps/host/internal/session"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"github.com/m1thrandir225/imperium/apps/host/internal/webrtc"
)

type Manager struct {
	app            fyne.App
	window         fyne.Window
	screens        map[string]Screen
	config         *config.Config
	authService    *auth.AuthService
	programService *programs.ProgramService
	sessionService *session.SessionService
	hostManager    *host.HostManager
}

func NewUIManager(config *config.Config) *Manager {

	manager := &Manager{
		app:     app.NewWithID("imperium"),
		screens: make(map[string]Screen),
		config:  config,
	}
	manager.window = manager.app.NewWindow("Imperium")

	dbPath := fmt.Sprintf("%s/programs.db", util.GetConfigDir())

	// Initialize services
	manager.authService = auth.NewAuthService(
		config.ServerAddress,
		config.GetAuthConfig(),
		config.SetAuthConfig,
	)

	manager.programService = programs.NewProgramService(
		config.ServerAddress,
		manager.authService.GetConfig().GetAccessToken(),
		manager.authService,
		dbPath,
	)

	videoRecorder := video.NewRecorder(&video.Config{
		FFMPEGPath: config.VideoConfig.FFMPEGPath,
		FPS:        config.VideoConfig.FPS,
		Encoder:    config.VideoConfig.Encoder,
	})

	webrtcStreamer, err := webrtc.NewStreamer()
	if err != nil {
		log.Fatalf("Failed to create webrtc streamer: %v", err)
	}

	manager.sessionService = session.NewSessionService(
		config.ServerAddress,
		manager.authService.GetConfig().GetAccessToken(),
		manager.authService,
		manager.programService,
		videoRecorder,
		webrtcStreamer,
	)

	hostConfig := config.GetHostConfig()
	if hostConfig == nil {
		hostConfig = &host.Config{
			HostName:  "",
			IPAddress: "",
			Port:      8080,
			UniqueID:  "",
			Status:    string(host.StatusAvailable),
		}
		config.SetHostConfig(hostConfig)
	}

	manager.hostManager = host.NewHostManager(
		hostConfig,
		manager.authService,
		manager.programService,
		config.SetHostConfig,
	)

	log.Println("shouldShowSetup", shouldShowSetup(config))
	log.Println("shouldShowLogin", shouldShowLogin(config))

	manager.initializeScreens()

	return manager
}

func (m *Manager) initializeScreens() {
	m.screens[SETTINGS_SCREEN] = NewSettingsScreen(m, m.config)
	m.screens[MAIN_MENU_SCREEN] = NewMainMenuScreen(m, m.authService)
	m.screens[STATUS_SCREEN] = NewStatusScreen(m, m.sessionService)
	m.screens[PROGRAMS_SCREEN] = NewProgramsScreen(m, m.programService)

	if shouldShowLogin(m.config) {
		m.screens[LOGIN_SCREEN] = NewLoginScreen(m, m.authService)
		m.screens[REGISTER_SCREEN] = NewRegisterScreen(m, m.authService)
	}

	if shouldShowSetup(m.config) {
		m.screens[SETUP_SCREEN] = NewSetupScreen(m.config, config.SaveConfigSections, func() {
			m.OnSetupSuccess()
		})
	}
}

func (m *Manager) AddScreen(screen Screen) {
	m.screens[screen.Name()] = screen
}

func (m *Manager) ShowScreen(name string) {
	if screen, ok := m.screens[name]; ok {
		m.window.SetContent(screen.Render(m.window))
	}
}

func (m *Manager) RunUI() {
	log.Println("Starting UI with auth check:", shouldShowLogin(m.config))
	log.Println("Starting UI with setup check:", shouldShowSetup(m.config))

	if shouldShowSetup(m.config) {
		m.ShowScreen(SETUP_SCREEN)
	} else if shouldShowLogin(m.config) {
		m.ShowScreen(LOGIN_SCREEN)
	} else {
		m.ShowScreen(MAIN_MENU_SCREEN)
	}

	m.window.Resize(fyne.NewSize(800, 600))
	m.window.SetFixedSize(true)

	m.window.ShowAndRun()
}

func shouldShowLogin(cfg *config.Config) bool {
	authConfig := cfg.GetAuthConfig()

	if authConfig == nil {
		return true
	}

	user := authConfig.GetCurrentUser()
	if user.ID == "" || user.Email == "" {
		return true
	}

	// If the refresh token is expired, we need to login again
	if authConfig.IsRefreshTokenExpired() {
		return true
	}

	return false
}

func (m *Manager) OnSetupSuccess() {
	delete(m.screens, SETUP_SCREEN)

	if shouldShowLogin(m.config) {
		if _, exists := m.screens[LOGIN_SCREEN]; !exists {
			m.screens[LOGIN_SCREEN] = NewLoginScreen(m, m.authService)
			m.screens[REGISTER_SCREEN] = NewRegisterScreen(m, m.authService)
		}
		m.ShowScreen(LOGIN_SCREEN)
	} else {
		m.ShowScreen(MAIN_MENU_SCREEN)
	}
}

func (m *Manager) OnLoginSuccess() {
	delete(m.screens, LOGIN_SCREEN)
	delete(m.screens, REGISTER_SCREEN)

	ctx := context.Background()
	if err := m.hostManager.Initialize(ctx); err != nil {
		log.Printf("Failed to initialize host manager: %v", err)
	}

	m.ShowScreen(MAIN_MENU_SCREEN)
}

func (m *Manager) OnLogout() {
	m.authService.Logout(context.Background())

	if m.hostManager != nil {
		m.hostManager.Shutdown()
	}

	m.ResetScreens()
	m.ShowScreen(LOGIN_SCREEN)
}

func (m *Manager) ResetScreens() {
	m.screens = make(map[string]Screen)
	m.initializeScreens()
}

func shouldShowSetup(cfg *config.Config) bool {

	if cfg.VideoConfig == nil || cfg.VideoConfig.FFMPEGPath == "" {
		return true
	}

	if _, err := os.Stat(cfg.VideoConfig.FFMPEGPath); err != nil {
		return true
	}

	// Check if server address is configured
	if cfg.ServerAddress == "" {
		return true
	}

	return false
}
