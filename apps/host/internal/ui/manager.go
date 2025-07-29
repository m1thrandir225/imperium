package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/m1thrandir225/imperium/apps/host/config"
)

type Manager struct {
	app     fyne.App
	window  fyne.Window
	screens map[string]Screen
	config  *config.Config
}

func NewUIManager(config *config.Config) *Manager {
	fyneApp := app.New()

	return &Manager{
		app:     fyneApp,
		screens: make(map[string]Screen),
		config:  config,
	}
}

func (m *Manager) AddScreen(screen Screen) {
	m.screens[screen.Name()] = screen
}

func (m *Manager) SetScreen(name string) {
	if screen, ok := m.screens[name]; ok {
		m.window.SetContent(screen.Render(m.window))
	}
}

func (m *Manager) SetupScreens() {
	encoderScreen := &EncoderScreen{manager: m}
	settingsScreen := &SettingsScreen{manager: m}
	loginScreen := &LoginScreen{manager: m}
	mainMenuScreen := &MainMenuScreen{manager: m}
	statusScreen := &StatusScreen{manager: m}

	m.AddScreen(loginScreen)
	m.AddScreen(mainMenuScreen)
	m.AddScreen(statusScreen)
	m.AddScreen(encoderScreen)
	m.AddScreen(settingsScreen)
}

func (m *Manager) RunUI() {
	m.SetupScreens()

	m.window = m.app.NewWindow("Imperium")
	m.window.Resize(fyne.NewSize(800, 600))
	m.window.SetFixedSize(true)

	m.SetScreen(MAIN_MENU_SCREEN)

	m.window.Show()
	m.app.Run()
}
