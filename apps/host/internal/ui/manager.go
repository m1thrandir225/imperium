package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Manager struct {
	app fyne.App
}

func NewUIManager() *Manager {
	fyneApp := app.New()

	return &Manager{
		app: fyneApp,
	}
}
