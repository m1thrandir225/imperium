package ui

import (
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	uapp "github.com/m1thrandir225/imperium/apps/host/internal/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
)

type Manager struct {
	app     fyne.App
	window  fyne.Window
	screens map[string]Screen
	bus     *uapp.EventBus
	state   *state.StateManager
}

func NewUIManager(stateManager *state.StateManager, bus *uapp.EventBus) *Manager {
	manager := &Manager{
		app:     app.NewWithID("imperium"),
		screens: make(map[string]Screen),
		state:   stateManager,
		bus:     bus,
	}
	manager.window = manager.app.NewWindow("Imperium")

	manager.initializeScreens()
	manager.subscribeNavigation()
	return manager
}

func (m *Manager) subscribeNavigation() {
	loginCh := m.bus.Subscribe(uapp.EventLoginSucceeded)
	go func() {
		for range loginCh {
			m.ShowScreen(MAIN_MENU_SCREEN)
		}
	}()

	logoutCh := m.bus.Subscribe(uapp.EventLogoutCompleted)
	go func() {
		for range logoutCh {
			m.ShowScreen(LOGIN_SCREEN)
		}
	}()

	setupCh := m.bus.Subscribe(uapp.EventSetupCompleted)
	go func() {
		for range setupCh {
			if m.shouldShowLogin() {
				m.ShowScreen(LOGIN_SCREEN)
			} else {
				m.ShowScreen(MAIN_MENU_SCREEN)
			}
		}
	}()
}

func (m *Manager) Publish(topic string, data any) {
	if m.bus != nil {
		m.bus.Publish(topic, data)
	}
}

func (m *Manager) GetState() state.AppState {
	return m.state.Get()
}

func (m *Manager) initializeScreens() {
	m.screens[SETTINGS_SCREEN] = NewSettingsScreen(m)
	m.screens[MAIN_MENU_SCREEN] = NewMainMenuScreen(m)
	m.screens[STATUS_SCREEN] = NewStatusScreen(m)
	m.screens[PROGRAMS_SCREEN] = NewProgramsScreen(m)

	if m.shouldShowLogin() {
		m.screens[LOGIN_SCREEN] = NewLoginScreen(m)
		m.screens[REGISTER_SCREEN] = NewRegisterScreen(m)
	}

	if m.shouldShowSetup() {
		m.screens[SETUP_SCREEN] = NewSetupScreen(m)
		m.OnSetupSuccess()
	}
}

func (m *Manager) AddScreen(screen Screen) {
	m.screens[screen.Name()] = screen
}

func (m *Manager) ShowScreen(name string) {
	if screen, ok := m.screens[name]; ok {
		m.window.SetContent(screen.Render(m.window))
		m.Publish(uapp.EventUIShowScreen, uapp.UIShowScreenPayload{Name: name})
	}
}

func (m *Manager) RunUI() {
	log.Println("Starting UI with auth check:", m.shouldShowLogin())
	log.Println("Starting UI with setup check:", m.shouldShowSetup())

	if m.shouldShowSetup() {
		m.ShowScreen(SETUP_SCREEN)
	} else if m.shouldShowLogin() {
		m.ShowScreen(LOGIN_SCREEN)
	} else {
		m.ShowScreen(MAIN_MENU_SCREEN)
	}

	m.window.Resize(fyne.NewSize(800, 600))
	m.window.SetFixedSize(true)

	m.window.ShowAndRun()
}

func (m *Manager) OnSetupSuccess() {
	delete(m.screens, SETUP_SCREEN)

	if m.shouldShowLogin() {
		if _, exists := m.screens[LOGIN_SCREEN]; !exists {
			m.screens[LOGIN_SCREEN] = NewLoginScreen(m)
			m.screens[REGISTER_SCREEN] = NewRegisterScreen(m)
		}
		m.ShowScreen(LOGIN_SCREEN)
	} else {
		m.ShowScreen(MAIN_MENU_SCREEN)
	}
}

// func (m *Manager) OnLoginSuccess() {
// 	delete(m.screens, LOGIN_SCREEN)
// 	delete(m.screens, REGISTER_SCREEN)

// 	ctx := context.Background()
// 	if err := m.hostManager.Initialize(ctx); err != nil {
// 		log.Printf("Failed to initialize host manager: %v", err)
// 	}

// 	m.ShowScreen(MAIN_MENU_SCREEN)
// }

// func (m *Manager) OnLogout() {
// 	m.authService.Logout(context.Background())

// 	if m.hostManager != nil {
// 		m.hostManager.Shutdown()
// 	}

// 	m.ResetScreens()
// 	m.ShowScreen(LOGIN_SCREEN)
// }

func (m *Manager) ResetScreens() {
	m.screens = make(map[string]Screen)
	m.initializeScreens()
}

func (m *Manager) shouldShowLogin() bool {
	cfg := m.GetState()

	if cfg.UserInfo.ID == "" || cfg.UserInfo.Email == "" {
		return true
	}

	if cfg.UserSession.AccessToken == "" {
		return true
	}

	return time.Now().After(cfg.UserSession.Expiry)

}

func (m *Manager) shouldShowSetup() bool {
	cfg := m.GetState().Settings

	if cfg.FFmpegPath == "" {
		return true
	}

	if _, err := os.Stat(cfg.FFmpegPath); err != nil {
		return true
	}

	// Check if server address is configured
	if cfg.ServerAddress == "" {
		return true
	}

	return false
}
