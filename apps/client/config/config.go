package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AuthServerBaseURL string `mapstructure:"auth_server_base_url"`
}

var (
	cfg *Config
)

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: Could not get user home directory: %v", err)
		return ""
	}
	return filepath.Join(homeDir, "Documents", "imperium")
}

func ConfigExists() bool {
	configPath := GetConfigPath()
	if configPath == "" {
		return false
	}
	configFile := filepath.Join(configPath, "client.json")
	_, err := os.Stat(configFile)
	return err == nil
}

func IsConfigured() bool {
	if !ConfigExists() {
		return false
	}
	config := Load()
	return config.AuthServerBaseURL != ""
}

func InitConfig() error {
	viper.SetConfigName("client")
	viper.SetConfigType("json")

	configPath := GetConfigPath()
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("IMPERIUM")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found - configuration required")
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	return nil
}

func Load() *Config {
	if cfg != nil {
		return cfg
	}

	if err := InitConfig(); err != nil {
		log.Printf("Error initializing config: %v", err)
		return &Config{
			AuthServerBaseURL: "",
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Printf("Error unmarshaling config: %v", err)
		return &Config{
			AuthServerBaseURL: "",
		}
	}

	return cfg
}

func SaveConfig(authServerURL string) error {
	configPath := GetConfigPath()
	if configPath == "" {
		return fmt.Errorf("could not determine config path")
	}

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	viper.Set("auth_server_base_url", authServerURL)

	configFile := filepath.Join(configPath, "client.json")
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}
	cfg = &Config{
		AuthServerBaseURL: authServerURL,
	}

	log.Printf("Config saved to: %s", configFile)
	return nil
}

func Get() *Config {
	if cfg == nil {
		return Load()
	}

	return cfg
}

func GetAuthServerURL() string {
	config := Get()

	return config.AuthServerBaseURL
}

func Reset() {
	cfg = nil
	viper.Reset()
}
