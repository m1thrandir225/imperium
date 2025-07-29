package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusScreen struct {
	manager *Manager
}

func (s *StatusScreen) Name() string {
	return "Status"
}

func (s *StatusScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewBorder(
		nil,
		widget.NewButton("Back to Main Menu", func() { s.manager.SetScreen("Main Menu") }),
		nil, nil,
		container.NewCenter(widget.NewLabel("Status: No one connected")),
	)
}
