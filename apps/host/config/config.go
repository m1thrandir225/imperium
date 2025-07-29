package config

import (
	"fmt"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	videoConfig video.Config `mapstructure:"video"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return nil, fmt.Errorf("failed to get user's home directory: %w", err)
	}

	configDir := filepath.Join(home, "Documents", "imperium")
	configName := "config"
	configPath := filepath.Join(configDir, "config.json")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	// if a config file doesn't exist, create it with default values
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := &Config{
			videoConfig: video.Config{
				Encoder:    "libx264",
				FPS:        30,
				FFMPEGPath: "",
			},
		}

		viper.Set("video", config.videoConfig)

		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			return nil, fmt.Errorf("failed to write default config: %w", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	viper.Set("video", config.videoConfig)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
