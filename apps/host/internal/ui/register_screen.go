package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
)

type RegisterScreen struct {
	manager *Manager
}

func (s *RegisterScreen) Name() string {
	return "Register"
}

func (s *RegisterScreen) Render(w fyne.Window) fyne.CanvasObject {
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
			dialog.NewInformation("Login", "Login logic would run here.", s.manager.window).Show()
		},
	}

	return container.NewBorder(
		nil,
		widget.NewButton("Back to Main Menu", func() { s.manager.ShowScreen("Main Menu") }),
		nil, nil,
		container.NewCenter(container.NewVBox(form)),
	)
}
