package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAppState_Serialization(t *testing.T) {
	now := time.Now()

	state := AppState{
		UserSession: UserSession{
			AccessToken:           "access-token",
			RefreshToken:          "refresh-token",
			AccessTokenExpiresAt:  now,
			RefreshTokenExpiresAt: now.Add(24 * time.Hour),
		},
		UserInfo: UserInfo{
			ID:    "user-123",
			Name:  "Test User",
			Email: "test@example.com",
		},
		HostInfo: HostInfo{
			ID:   "host-456",
			Name: "Test Host",
			IP:   "192.168.1.100",
			Port: 8080,
		},
		Settings: Settings{
			FFmpegPath:         "/usr/bin/ffmpeg",
			Encoder:            "h264",
			Framerate:          60,
			Bitrate:            "2000k",
			ServerAddress:      "https://api.example.com",
			CustomProgramPaths: []string{"/usr/bin/custom1", "/usr/bin/custom2"},
			RawgAPIKey:         "rawg-api-key-123",
		},
	}

	require.NotEmpty(t, state.UserSession.AccessToken)
	require.NotEmpty(t, state.UserInfo.ID)
	require.NotEmpty(t, state.HostInfo.ID)
	require.NotEmpty(t, state.Settings.FFmpegPath)
}

func TestUserSession_Expiration(t *testing.T) {
	now := time.Now()
	session := UserSession{
		AccessToken:          "token",
		AccessTokenExpiresAt: now.Add(1 * time.Hour),
	}

	require.Equal(t, "token", session.AccessToken)
	require.True(t, session.AccessTokenExpiresAt.After(now))
}

func TestSettings_DefaultValues(t *testing.T) {
	settings := Settings{
		FFmpegPath:    "/usr/bin/ffmpeg",
		Encoder:       "h264",
		Framerate:     30,
		Bitrate:       "1000k",
		ServerAddress: "https://api.example.com",
	}

	require.Equal(t, "h264", settings.Encoder)
	require.Equal(t, 30, settings.Framerate)
	require.Equal(t, "1000k", settings.Bitrate)
	require.Equal(t, "https://api.example.com", settings.ServerAddress)
	require.Equal(t, "/usr/bin/ffmpeg", settings.FFmpegPath)
}
