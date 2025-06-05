package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (manager *UIManager) CreateWindow() fyne.Window {
	w := manager.app.NewWindow("Imperium")
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)
	return w
}

func (manager *UIManager) RunUI() {
	mainWindow := manager.CreateWindow()

	mainWindow.SetContent(widget.NewLabel("Hello Native UI"))
	mainWindow.Show()

	manager.app.Run()
}
