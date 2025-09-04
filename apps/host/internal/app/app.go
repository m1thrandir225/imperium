package app

import (
	"log"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/state"
)

type App struct {
	Name        string
	Bus         *EventBus
	State       *state.StateManager
	AuthBaseURL string
}

func New(name string) (*App, error) {
	sm, err := state.NewStateManager(name)
	if err != nil {
		return nil, err
	}

	st := sm.Get()

	return &App{
		Name:        name,
		Bus:         NewEventBus(),
		State:       sm,
		AuthBaseURL: st.Settings.ServerAddress,
	}, nil
}

func (a *App) Start() {
	a.Bus.Publish(EventStateLoaded, a.State.Get())
	a.Bus.Publish(EventAppStarted, nil)
	log.Println("App started")
}

func (a *App) NeedsLogin() bool {
	st := a.State.Get()
	if st.UserInfo.ID == "" || st.UserInfo.Email == "" {
		return true
	}

	if st.UserSession.AccessToken == "" {
		return true
	}

	return time.Now().After(st.UserSession.Expiry)

}
