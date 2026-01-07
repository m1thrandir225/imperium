package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/app"
)

type MainMenuScreen struct {
	manager *uiManager
}

func NewMainMenuScreen(manager *uiManager) *MainMenuScreen {
	return &MainMenuScreen{
		manager: manager,
	}
}

func (s *MainMenuScreen) Name() string {
	return MAIN_MENU_SCREEN
}

func (s *MainMenuScreen) Render(w fyne.Window) fyne.CanvasObject {
	imageResource, _ := fyne.LoadResourceFromURLString("https://github.com/m1thrandir225/imperium/blob/master/assets/imperium_horizontal_logo.png?raw=true")

	logo := canvas.NewImageFromResource(imageResource)
	logo.FillMode = canvas.ImageFillContain

	return container.NewCenter(
		container.NewVBox(
			logo,
			widget.NewLabelWithStyle("Imperium Host", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Current Status", func() {
				s.manager.showScreen(STATUS_SCREEN)
			}),
			widget.NewButton("Programs", func() {
				s.manager.showScreen(PROGRAMS_SCREEN)
			}),
			widget.NewButton("Settings", func() {
				s.manager.showScreen(SETTINGS_SCREEN)
			}),
			widget.NewButton("Logout", func() {
				s.manager.publish(app.EventLogoutRequested, nil)
			}),
		),
	)
}
