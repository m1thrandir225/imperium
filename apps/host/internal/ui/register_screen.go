package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/app"
)

type RegisterScreen struct {
	manager *uiManager
}

func NewRegisterScreen(manager *uiManager) *RegisterScreen {
	return &RegisterScreen{
		manager: manager,
	}
}

func (s *RegisterScreen) Name() string {
	return REGISTER_SCREEN
}

func (s *RegisterScreen) Render(w fyne.Window) fyne.CanvasObject {
	firstNameEntry := widget.NewEntry()
	firstNameEntry.SetPlaceHolder("First name")

	lastNameEntry := widget.NewEntry()
	lastNameEntry.SetPlaceHolder("Last name")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	backToLoginBtn := widget.NewButton("Already have an account? Login", func() {
		s.manager.ShowScreen(LOGIN_SCREEN)
	})
	form := &widget.Form{
		OnSubmit: func() {
			if firstNameEntry.Text == "" || lastNameEntry.Text == "" || emailEntry.Text == "" || passwordEntry.Text == "" {
				dialog.ShowError(errors.New("validation: First name, last name, email and password are required"), w)
				return
			}
			s.manager.publish(app.EventRegisterRequested, app.RegisterRequestedPayload{
				Email:     emailEntry.Text,
				Password:  passwordEntry.Text,
				FirstName: firstNameEntry.Text,
				LastName:  lastNameEntry.Text,
			})
		},
		SubmitText: "Register",
	}

	form.Append("First name", firstNameEntry)
	form.Append("Last name", lastNameEntry)
	form.Append("Email", emailEntry)
	form.Append("Password", passwordEntry)

	return container.NewVBox(
		container.NewCenter(Logo(true)),
		form,
		backToLoginBtn,
		Copyright(),
	)
}
