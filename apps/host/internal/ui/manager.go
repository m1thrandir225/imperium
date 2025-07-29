package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"os"
)

type Manager struct {
	app     fyne.App
	window  fyne.Window
	screens map[string]Screen
	config  *util.Config
}

func NewUIManager(config *util.Config) *Manager {
	manager := &Manager{
		app:     app.New(),
		screens: make(map[string]Screen),
		config:  config,
	}
	manager.window = manager.app.NewWindow("Imperium")

	if shouldShowSetup(config) {
		manager.screens[SETUP_SCREEN] = NewSetupScreen(config, util.SaveConfig, func() {
			manager.ShowScreen(SETTINGS_SCREEN)
		})
	}
	manager.screens[ENCODER_SCREEN] = NewEncoderScreen()
	manager.screens[SETTINGS_SCREEN] = NewSettingsScreen()
	manager.screens[LOGIN_SCREEN] = NewLoginScreen()
	manager.screens[MAIN_MENU_SCREEN] = NewMainMenuScreen()
	manager.screens[STATUS_SCREEN] = NewStatusScreen()

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
