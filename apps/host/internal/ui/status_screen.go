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
	return STATUS_SCREEN
}

func (s *StatusScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewBorder(
		nil,
		widget.NewButton("Back to Main Menu", func() { s.manager.SetScreen(MAIN_MENU_SCREEN) }),
		nil, nil,
		container.NewCenter(widget.NewLabel("Status: No one connected")),
	)
}
