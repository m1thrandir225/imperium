package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type EncoderScreen struct {
	manager *Manager
}

func (s *EncoderScreen) Name() string {
	return "Encoders"
}

func (s *EncoderScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabel("Encoder screen"),
	))
}
