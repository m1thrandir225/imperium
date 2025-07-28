package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Manager struct {
	app     fyne.App
	window  fyne.Window
	screens map[string]Screen
}

func NewUIManager() *Manager {
	fyneApp := app.New()

	return &Manager{
		app:     fyneApp,
		screens: make(map[string]Screen),
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

func (m *Manager) SetupUI() {
	encoderScreen := &EncoderScreen{manager: m}
	settingsScreen := &SettingsScreen{manager: m}

	m.AddScreen(encoderScreen)
	m.AddScreen(settingsScreen)
}

func (m *Manager) RunUI() {
	m.SetupUI()

	m.window = m.app.NewWindow("Imperium")
	m.window.Resize(fyne.NewSize(800, 600))
	m.window.SetFixedSize(true)

	// Set the initial screen
	m.SetScreen("Settings")

	m.window.Show()
	m.app.Run()
}
