package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainMenuScreen struct {
	manager *Manager
}

func (s *MainMenuScreen) Name() string {
	return MAIN_MENU_SCREEN
}

func (s *MainMenuScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Imperium Host", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Current Status", func() { s.manager.SetScreen("Status") }),
			widget.NewButton("Encoder List", func() { s.manager.SetScreen("Encoders") }),
			widget.NewButton("Settings", func() { s.manager.SetScreen("Settings") }),
			widget.NewButton("Login", func() { s.manager.SetScreen("Login") }),
		),
	)
}
