package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStateManager(t *testing.T) {
	tests := []struct {
		name     string
		appName  string
		wantErr  bool
		setup    func() string
		teardown func(string)
	}{
		{
			name:    "valid app name",
			appName: "test-app",
			wantErr: false,
			setup: func() string {
				return "test-app"
			},
			teardown: func(appName string) {
				configDir, _ := os.UserConfigDir()
				testDir := filepath.Join(configDir, appName)
				os.RemoveAll(testDir)
			},
		},
		{
			name:    "empty app name",
			appName: "",
			wantErr: false,
			setup: func() string {
				return ""
			},
			teardown: func(appName string) {
				configDir, _ := os.UserConfigDir()
				testDir := filepath.Join(configDir, appName)
				os.RemoveAll(testDir)
			},
		},
	}

	for _, tt := range tests {
		appName := tt.setup()

		defer tt.teardown(appName)

		manager, err := NewStateManager(appName)

		if tt.wantErr {
			assert.Error(t, err)
			assert.Nil(t, manager)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, manager)
			assert.NotNil(t, manager.state)
		}
	}
}

func TestStateManager_Load(t *testing.T) {
	manager, err := NewStateManager("test-load")
	require.NoError(t, err)
	defer func() {
		configDir, _ := os.UserConfigDir()
		testDir := filepath.Join(configDir, "test-load")
		os.RemoveAll(testDir)
	}()

	err = manager.Load()
	assert.NoError(t, err)
	assert.NotNil(t, manager.state)
}

func TestStateManager_Save(t *testing.T) {
	manager, err := NewStateManager("test-save")
	require.NoError(t, err)
	defer func() {
		configDir, _ := os.UserConfigDir()
		testDir := filepath.Join(configDir, "test-save")
		os.RemoveAll(testDir)
	}()

	state := AppState{
		UserSession: UserSession{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
		},
		UserInfo: UserInfo{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
		},
		HostInfo: HostInfo{
			ID:   "test-host-id",
			Name: "Test Host",
			IP:   "127.0.0.1",
			Port: 8080,
		},
		Settings: Settings{
			FFmpegPath:    "/usr/bin/ffmpeg",
			Encoder:       "h264",
			Framerate:     30,
			Bitrate:       "1000k",
			ServerAddress: "https://example.com",
		},
	}
	manager.state = &state
	err = manager.Save()

	assert.NoError(t, err)

	err = manager.Load()
	loadedState := manager.state
	assert.NoError(t, err)
	assert.Equal(t, state.UserInfo.ID, loadedState.UserInfo.ID)
	assert.Equal(t, state.UserInfo.Name, loadedState.UserInfo.Name)
	assert.Equal(t, state.UserInfo.Email, loadedState.UserInfo.Email)
	assert.Equal(t, state.UserSession.AccessToken, loadedState.UserSession.AccessToken)
	assert.Equal(t, state.UserSession.RefreshToken, loadedState.UserSession.RefreshToken)
	assert.Equal(t, state.HostInfo.ID, loadedState.HostInfo.ID)
	assert.Equal(t, state.HostInfo.Name, loadedState.HostInfo.Name)
	assert.Equal(t, state.HostInfo.IP, loadedState.HostInfo.IP)
	assert.Equal(t, state.HostInfo.Port, loadedState.HostInfo.Port)
	assert.Equal(t, state.Settings.FFmpegPath, loadedState.Settings.FFmpegPath)
	assert.Equal(t, state.Settings.Encoder, loadedState.Settings.Encoder)
	assert.Equal(t, state.Settings.Framerate, loadedState.Settings.Framerate)
	assert.Equal(t, state.Settings.Bitrate, loadedState.Settings.Bitrate)
	assert.Equal(t, state.Settings.ServerAddress, loadedState.Settings.ServerAddress)
}

func TestStateManager_Get(t *testing.T) {
	manager, err := NewStateManager("test-get")
	require.NoError(t, err)
	defer func() {
		configDir, _ := os.UserConfigDir()
		testDir := filepath.Join(configDir, "test-get")
		os.RemoveAll(testDir)
	}()

	expectedState := AppState{
		UserSession: UserSession{
			AccessToken: "test-token",
		},
		UserInfo: UserInfo{
			ID: "test-id",
		},
	}
	manager.state = &expectedState

	state := manager.Get()

	assert.Equal(t, expectedState, state)
}

func TestStateManager_Update(t *testing.T) {
	manager, err := NewStateManager("test-update")
	require.NoError(t, err)

	defer func() {
		configDir, _ := os.UserConfigDir()
		testDir := filepath.Join(configDir, "test-update")
		os.RemoveAll(testDir)
	}()

	manager.state = &AppState{
		UserInfo: UserInfo{
			ID:   "original-id",
			Name: "Original Name",
		},
	}

	err = manager.Update(func(state *AppState) {
		state.UserInfo.Name = "Updated Name"
		state.UserInfo.Email = "updated@example.com"
	})
	assert.NoError(t, err)

	assert.Equal(t, "Updated Name", manager.state.UserInfo.Name)
	assert.Equal(t, "updated@example.com", manager.state.UserInfo.Email)
	assert.Equal(t, "original-id", manager.state.UserInfo.ID)
}

func TestStateManager_Update_NoState(t *testing.T) {
	manager := &StateManager{state: nil}

	err := manager.Update(func(state *AppState) {
		state.UserInfo.Name = "Test"
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no state initialized")
}
