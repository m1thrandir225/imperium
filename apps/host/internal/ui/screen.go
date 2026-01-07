// Package ui provides the UI for the host application.
package ui

import "fyne.io/fyne/v2"

type screen interface {
	Render(w fyne.Window) fyne.CanvasObject
	Name() string
}
