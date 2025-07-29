package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

type LoginScreen struct {
}

func NewLoginScreen() *LoginScreen {
	return &LoginScreen{}
}

func (s *LoginScreen) Name() string {
	return LOGIN_SCREEN
}

func (s *LoginScreen) Render(w fyne.Window) fyne.CanvasObject {
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Email", Widget: emailEntry},
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			log.Printf("Login attempt with Email: %s, Password: %s", emailEntry.Text, passwordEntry.Text)
		},
	}

	return container.NewBorder(
		nil,
		widget.NewButton("Back to Main Menu", func() {
		}),
		nil, nil,
		container.NewCenter(container.NewVBox(form)),
	)
}
