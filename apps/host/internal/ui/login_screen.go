package ui

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
)

type LoginScreen struct {
	manager     *Manager
	authService *auth.AuthService
}

func NewLoginScreen(manager *Manager, authService *auth.AuthService) *LoginScreen {
	return &LoginScreen{
		manager:     manager,
		authService: authService,
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

	loginBtn := widget.NewButton("Login", func() {
		req := auth.LoginRequest{
			Email:    emailEntry.Text,
			Password: passwordEntry.Text,
		}

		resp, err := s.authService.Login(context.Background(), req)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		log.Printf("Login successful for user: %s", resp.User.Email)
		s.manager.OnLoginSuccess()
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
