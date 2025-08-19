package ui

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
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
	config         *util.Config
	authService    *auth.AuthService
	programService *programs.ProgramService
	sessionService *session.SessionService
}

func NewUIManager(config *util.Config) *Manager {
	manager := &Manager{
		app:     app.New(),
		screens: make(map[string]Screen),
		config:  config,
	}
	manager.window = manager.app.NewWindow("Imperium")

	// Initialize services
	manager.authService = auth.NewAuthService(config.ServerAddress)
	manager.programService = programs.NewProgramService(config.ServerAddress, manager.authService.GetToken())

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
		manager.authService.GetToken(),
		manager.programService,
		videoRecorder,
		webrtcStreamer,
	)

	if shouldShowSetup(config) {
		manager.screens[SETUP_SCREEN] = NewSetupScreen(config, util.SaveConfig, func() {
			manager.ShowScreen(SETTINGS_SCREEN)
		})
	}

	manager.screens[ENCODER_SCREEN] = NewEncoderScreen()
	manager.screens[SETTINGS_SCREEN] = NewSettingsScreen()
	manager.screens[LOGIN_SCREEN] = NewLoginScreen(manager, manager.authService)
	manager.screens[REGISTER_SCREEN] = NewRegisterScreen(manager, manager.authService)
	manager.screens[MAIN_MENU_SCREEN] = NewMainMenuScreen(manager)
	manager.screens[STATUS_SCREEN] = NewStatusScreen(manager, manager.sessionService)
	manager.screens[PROGRAMS_SCREEN] = NewProgramsScreen(manager, manager.programService)

	return manager
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
	if _, exists := m.screens[SETUP_SCREEN]; exists {
		m.ShowScreen(SETUP_SCREEN)
	} else {
		m.ShowScreen(MAIN_MENU_SCREEN)
	}

	m.window.Resize(fyne.NewSize(800, 600))
	m.window.SetFixedSize(true)

	m.window.ShowAndRun()
}

func shouldShowSetup(cfg *util.Config) bool {
	if cfg.VideoConfig.FFMPEGPath == "" {
		if _, err := os.Stat(cfg.VideoConfig.FFMPEGPath); err == nil {
			return false
		}
	}
	return true
}
