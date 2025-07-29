package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type EncoderScreen struct {
}

func NewEncoderScreen() *EncoderScreen {
	return &EncoderScreen{}
}

func (s *EncoderScreen) Name() string {
	return ENCODER_SCREEN
}

func (s *EncoderScreen) Render(w fyne.Window) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabel("Encoder screen"),
	))
}
