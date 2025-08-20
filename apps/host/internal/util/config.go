package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
	"github.com/m1thrandir225/imperium/apps/host/internal/host"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"github.com/spf13/viper"
)

type Config struct {
	VideoConfig   *video.Config `mapstructure:"video"`
	ServerAddress string        `mapstructure:"server_address"`
	AuthConfig    *auth.Config  `mapstructure:"auth_config"`
	HostConfig    *host.Config  `mapstructure:"host_config"`
}

func (c *Config) SetServerAddress(serverAddress string) {
	c.ServerAddress = serverAddress
}

func (c *Config) GetAuthConfig() *auth.Config {
	return c.AuthConfig
}

func (c *Config) GetVideoConfig() *video.Config {
	return c.VideoConfig
}

func (c *Config) GetHostConfig() *host.Config {
	return c.HostConfig
}

func LoadConfig() (*Config, error) {
	configDir := GetConfigDir()
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
			VideoConfig: &video.Config{
				Encoder:    "libx264",
				FPS:        30,
				FFMPEGPath: "",
			},
			ServerAddress: "",
			AuthConfig: auth.LoadConfig(
				auth.User{},
				"",
				"",
				time.Time{},
				time.Time{},
			),
			HostConfig: &host.Config{
				HostName:  "",
				IPAddress: "",
				Port:      0,
				Status:    string(host.StatusAvailable),
				UniqueID:  "",
			},
		}

		viper.Set("video", config.VideoConfig)
		viper.Set("server_address", config.ServerAddress)
		viper.Set("auth_config", config.AuthConfig)

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
	viper.Set("video", config.VideoConfig)
	viper.Set("server_address", config.ServerAddress)
	viper.Set("auth_config", config.AuthConfig)
	viper.Set("host_config", config.HostConfig)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func SaveConfigSections(config *Config, sections ...string) error {
	// Read current config to preserve other sections
	currentConfig, err := LoadConfig()
	if err != nil {
		// If we can't load current config, save everything
		return SaveConfig(config)
	}

	// Create a merged config
	mergedConfig := &Config{
		VideoConfig:   currentConfig.VideoConfig,
		ServerAddress: currentConfig.ServerAddress,
		AuthConfig:    currentConfig.AuthConfig,
		HostConfig:    currentConfig.HostConfig,
	}

	// Update only the specified sections
	for _, section := range sections {
		switch section {
		case "video":
			mergedConfig.VideoConfig = config.VideoConfig
		case "server_address":
			mergedConfig.ServerAddress = config.ServerAddress
		case "auth_config":
			mergedConfig.AuthConfig = config.AuthConfig
		case "host_config":
			mergedConfig.HostConfig = config.HostConfig
		}
	}

	return SaveConfig(mergedConfig)
}

func (c *Config) SetAuthConfig(config *auth.Config) {
	c.AuthConfig = config
	log.Println("config", c.AuthConfig)
	SaveConfig(c)
}

func (c *Config) SetHostConfig(config *host.Config) {
	c.HostConfig = config
	SaveConfig(c)
}

func (c *Config) SetVideoConfig(config *video.Config) {
	c.VideoConfig = config

	SaveConfig(c)
}
