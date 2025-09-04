package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/app"
)

type LoginScreen struct {
	manager *Manager
}

func NewLoginScreen(manager *Manager) *LoginScreen {
	return &LoginScreen{manager: manager}
}

func (s *LoginScreen) Name() string {
	return LOGIN_SCREEN
}

func (s *LoginScreen) Render(w fyne.Window) fyne.CanvasObject {
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	loginBtn := widget.NewButton("Login", func() {
		if emailEntry.Text == "" || passwordEntry.Text == "" {
			dialog.ShowError(errors.New("validation: Email and password are required"), w)
			return
		}
		s.manager.Publish(app.EventLoginRequested, app.LoginRequestedPayload{
			Email:    emailEntry.Text,
			Password: passwordEntry.Text,
		})
	})

	registerBtn := widget.NewButton("Register", func() {
		s.manager.ShowScreen(REGISTER_SCREEN)
	})

	backBtn := widget.NewButton("Back to Main Menu", func() {
		s.manager.ShowScreen(MAIN_MENU_SCREEN)
	})

	form := container.NewVBox(
		widget.NewLabel("Login"),
		emailEntry,
		passwordEntry,
		loginBtn,
		registerBtn,
		backBtn,
	)
	return container.NewCenter(form)
}
