package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainMenuScreen struct {
}

func NewMainMenuScreen() *MainMenuScreen {
	return &MainMenuScreen{}
}

func (s *MainMenuScreen) Name() string {
	return MAIN_MENU_SCREEN
}

func (s *MainMenuScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Imperium Host", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Current Status", func() {}),
			widget.NewButton("Encoder List", func() {}),
			widget.NewButton("Settings", func() {}),
			widget.NewButton("Login", func() {}),
		),
	)
}
