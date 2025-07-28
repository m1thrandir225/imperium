package ui

import "fyne.io/fyne/v2"

type Screen interface {
	Render(w fyne.Window) fyne.CanvasObject
	Name() string
}
