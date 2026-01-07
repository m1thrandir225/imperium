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
	manager *uiManager
}

func NewLoginScreen(manager *uiManager) *LoginScreen {
	return &LoginScreen{
		manager: manager,
	}
}

func (s *LoginScreen) Name() string {
	return LOGIN_SCREEN
}

func (s *LoginScreen) Render(w fyne.Window) fyne.CanvasObject {
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	passwordEntry := widget.NewPasswordEntry()

	passwordEntry.SetPlaceHolder("Password")

	registerBtn := widget.NewButton("Don't have an account? Register", func() {
		s.manager.showScreen(REGISTER_SCREEN)
	})

	form := &widget.Form{
		OnSubmit: func() {
			if emailEntry.Text == "" || passwordEntry.Text == "" {
				dialog.ShowError(errors.New("validation: Email and password are required"), w)
				return
			}
			s.manager.publish(app.EventLoginRequested, app.LoginRequestedPayload{
				Email:    emailEntry.Text,
				Password: passwordEntry.Text,
			})
		},
		SubmitText: "Login",
	}
	form.Append("Email", emailEntry)
	form.Append("Password", passwordEntry)

	return container.NewVBox(
		container.NewCenter(Logo(true)),
		form,
		registerBtn,
		Copyright(),
	)
}
