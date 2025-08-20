package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
)

type MainMenuScreen struct {
	manager     *Manager
	authService *auth.AuthService
}

func NewMainMenuScreen(manager *Manager, authService *auth.AuthService) *MainMenuScreen {
	return &MainMenuScreen{
		manager:     manager,
		authService: authService,
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
			widget.NewButton("Settings", func() {
				s.manager.ShowScreen(SETTINGS_SCREEN)
			}),
			widget.NewButton("Logout", func() {
				s.manager.OnLogout()
			}),
		),
	)
}
