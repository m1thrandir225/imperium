package ui

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

	imageURL := "https://github.com/m1thrandir225/imperium/blob/master/assets/imperium_horizontal_fill_logo.png?raw=true"
	imageResource, err := fyne.LoadResourceFromURLString(imageURL)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to load logo: %w", err), w)
	}
	logo := canvas.NewImageFromResource(imageResource)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(300, 120)) // Adjust for your layout

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	passwordEntry := widget.NewPasswordEntry()

	passwordEntry.SetPlaceHolder("Password")

	registerBtn := widget.NewButton("Don't have an account? Register", func() {
		s.manager.ShowScreen(REGISTER_SCREEN)
	})

	form := &widget.Form{
		OnSubmit: func() {
			if emailEntry.Text == "" || passwordEntry.Text == "" {
				dialog.ShowError(errors.New("validation: Email and password are required"), w)
				return
			}
			s.manager.Publish(app.EventLoginRequested, app.LoginRequestedPayload{
				Email:    emailEntry.Text,
				Password: passwordEntry.Text,
			})
		},
		SubmitText: "Login",
	}
	form.Append("Email", emailEntry)
	form.Append("Password", passwordEntry)

	return container.NewVBox(
		container.NewCenter(logo),
		form,
		registerBtn,
	)
}
