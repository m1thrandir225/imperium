package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type UIManager struct {
	app fyne.App
}

func NewUIManager() *UIManager {
	app := app.New()

	return &UIManager{
		app: app,
	}
}
