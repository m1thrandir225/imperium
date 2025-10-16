package video

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDefaultConfig(t *testing.T) {
	config := Config{
		Encoder:    "libx264",
		FPS:        30,
		FFMPEGPath: "",
	}

	defaultConfig := NewDefaultConfig()
	require.Equal(t, config.Encoder, defaultConfig.Encoder)
	require.Equal(t, config.FPS, defaultConfig.FPS)
	require.Equal(t, config.FFMPEGPath, defaultConfig.FFMPEGPath)
}

func TestLoadConfig(t *testing.T) {
	config := Config{
		Encoder:    "libx264",
		FPS:        30,
		FFMPEGPath: "",
	}

	loadedConfig := LoadConfig(
		config.Encoder,
		config.FPS,
		config.FFMPEGPath,
	)

	require.Equal(t, config.Encoder, loadedConfig.Encoder)
	require.Equal(t, config.FPS, loadedConfig.FPS)
	require.Equal(t, config.FFMPEGPath, loadedConfig.FFMPEGPath)
}

func TestSet_Encoder(t *testing.T) {
	config := NewDefaultConfig()

	require.NotEmpty(t, config)

	config.SetEncoder("libx265")

	require.NotEqual(t, config.Encoder, "libx264")
	require.Equal(t, config.Encoder, "libx265")
}

func TestSet_FPS(t *testing.T) {
	config := NewDefaultConfig()

	require.NotEmpty(t, config)

	config.SetFPS(60)

	require.NotEqual(t, config.FPS, 30)
	require.Equal(t, config.FPS, 60)
}

func TestSet_FFMPEGPath(t *testing.T) {
	config := NewDefaultConfig()

	require.NotEmpty(t, config)

	config.SetFFMPEGPath("/usr/bin/ffmpeg")

	require.NotEqual(t, config.FFMPEGPath, "")
	require.Equal(t, config.FFMPEGPath, "/usr/bin/ffmpeg")
}
