package state

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	serviceName    = "imperium-host"
	accessKeyName  = "imperium-host-access-key"
	refreshKeyName = "imperium-host-refresh-key"
)

// SaveTokens saves the access and refresh tokens to the keyring
func SaveTokens(accessToken, refreshToken string) error {
	if err := keyring.Set(serviceName, accessKeyName, accessToken); err != nil {
		return fmt.Errorf("failed to save access token: %w", err)
	}

	if err := keyring.Set(serviceName, refreshKeyName, refreshToken); err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// LoadTokens loads the access and refresh tokens from the keyring
func LoadTokens() (string, string, error) {
	accessToken, err := keyring.Get(serviceName, accessKeyName)
	if err != nil {
		return "", "", fmt.Errorf("failed to load access token: %w", err)
	}

	refreshToken, err := keyring.Get(serviceName, refreshKeyName)
	if err != nil {
		return "", "", fmt.Errorf("failed to load refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// DeleteTokens deletes the access and refresh tokens from the keyring
func DeleteTokens() error {
	_ = keyring.Delete(serviceName, accessKeyName)
	_ = keyring.Delete(serviceName, refreshKeyName)

	return nil
}
