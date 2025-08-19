package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainMenuScreen struct {
	manager *Manager
}

func NewMainMenuScreen(manager *Manager) *MainMenuScreen {
	return &MainMenuScreen{
		manager: manager,
	}
}

func (s *MainMenuScreen) Name() string {
	return MAIN_MENU_SCREEN
}

func (s *MainMenuScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Imperium Host", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Current Status", func() {
				s.manager.ShowScreen(STATUS_SCREEN)
			}),
			widget.NewButton("Programs", func() {
				s.manager.ShowScreen(PROGRAMS_SCREEN)
			}),
			widget.NewButton("Encoder List", func() {
				s.manager.ShowScreen(ENCODER_SCREEN)
			}),
			widget.NewButton("Settings", func() {
				s.manager.ShowScreen(SETTINGS_SCREEN)
			}),
			widget.NewButton("Login", func() {
				s.manager.ShowScreen(LOGIN_SCREEN)
			}),
		),
	)
}
