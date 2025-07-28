package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SettingsScreen struct {
	manager *Manager
}

func (s *SettingsScreen) Name() string {
	return "Settings"
}

func (s *SettingsScreen) Render(w fyne.Window) fyne.CanvasObject {
	backButton := widget.NewButton("Back to Encoders", func() {
		s.manager.SetScreen("Encoders")
	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabel("This is the Settings Screen"),
		backButton,
	))
}
