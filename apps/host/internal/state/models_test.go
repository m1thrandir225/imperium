package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	assert.NotEmpty(t, state.UserSession.AccessToken)
	assert.NotEmpty(t, state.UserInfo.ID)
	assert.NotEmpty(t, state.HostInfo.ID)
	assert.NotEmpty(t, state.Settings.FFmpegPath)
}

func TestUserSession_Expiration(t *testing.T) {
	now := time.Now()
	session := UserSession{
		AccessToken:          "token",
		AccessTokenExpiresAt: now.Add(1 * time.Hour),
	}

	assert.True(t, session.AccessTokenExpiresAt.After(now))
}

func TestSettings_DefaultValues(t *testing.T) {
	settings := Settings{
		FFmpegPath:    "/usr/bin/ffmpeg",
		Encoder:       "h264",
		Framerate:     30,
		Bitrate:       "1000k",
		ServerAddress: "https://api.example.com",
	}

	assert.Equal(t, "h264", settings.Encoder)
	assert.Equal(t, 30, settings.Framerate)
	assert.Equal(t, "1000k", settings.Bitrate)
}
