package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/session"
)

type StatusScreen struct {
	manager        *Manager
	sessionService *session.SessionService
}

func NewStatusScreen(manager *Manager, sessionService *session.SessionService) *StatusScreen {
	return &StatusScreen{
		manager:        manager,
		sessionService: sessionService,
	}
}

func (s *StatusScreen) Name() string {
	return STATUS_SCREEN
}

func (s *StatusScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewBorder(
		nil,
		widget.NewButton("Back to Main Menu", func() {}),
		nil, nil,
		container.NewCenter(widget.NewLabel("Status: No one connected")),
	)
}
