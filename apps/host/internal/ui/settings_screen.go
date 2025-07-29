package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SettingsScreen struct {
}

func NewSettingsScreen() *SettingsScreen {
	return &SettingsScreen{}
}

func (s *SettingsScreen) Name() string {
	return SETTINGS_SCREEN
}

func (s *SettingsScreen) Render(w fyne.Window) fyne.CanvasObject {
	backButton := widget.NewButton("Back to Main Menu", func() {

	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabel("This is the Settings Screen"),
		backButton,
	))
}
