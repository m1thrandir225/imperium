// Package state provides a manager for the application's state as well as event definitions
package state

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"github.com/spf13/viper"
)

// StateManager manages the application's state and provides methods to load and save the state
type StateManager struct {
	state *AppState
	mu    sync.Mutex
}

// NewStateManager creates a new StateManager
func NewStateManager(appName string) (*StateManager, error) {
	configDir, err := util.GetConfigDir(appName)
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app directory: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(configDir)
	viper.SafeWriteConfig()

	m := &StateManager{}

	if err := m.Load(); err != nil {
		return nil, err
	}

	return m, nil
}

// Load loads the state from the config file
func (m *StateManager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	var s AppState
	if err := viper.Unmarshal(&s); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}
	m.state = &s
	return nil
}

// Save saves the state to the config file
func (m *StateManager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state == nil {
		return fmt.Errorf("no state to save")
	}

	viper.Set("user_session", m.state.UserSession)
	viper.Set("user_info", m.state.UserInfo)
	viper.Set("host_info", m.state.HostInfo)
	viper.Set("settings", m.state.Settings)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// Get returns a copy of the in-memory state
func (m *StateManager) Get() AppState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return *m.state
}

// Update applies a function to modify the state and persist the changes
func (m *StateManager) Update(fn func(*AppState)) error {

	if m.state == nil {
		return errors.New("no state initialized")
	}

	fn(m.state)

	if err := m.Save(); err != nil {
		return err
	}

	return nil
}
