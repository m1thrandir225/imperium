package ui

import "fyne.io/fyne/v2"

const (
	LOGIN_SCREEN     = "Login"
	REGISTER_SCREEN  = "Register"
	MAIN_MENU_SCREEN = "MainMenu"
	SETTINGS_SCREEN  = "Settings"
	STATUS_SCREEN    = "Status"
	ENCODER_SCREEN   = "Encoder"
)

type Screen interface {
	Render(w fyne.Window) fyne.CanvasObject
	Name() string
}
