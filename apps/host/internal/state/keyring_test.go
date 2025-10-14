package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveTokens(t *testing.T) {
	// Note: These tests may fail in CI/CD environments without proper keyring setup
	// Consider using build tags to skip on certain platforms

	accessToken := "test-access-token-123"
	refreshToken := "test-refresh-token-456"

	_ = DeleteTokens()

	err := SaveTokens(accessToken, refreshToken)
	if err != nil {
		t.Skipf("Skipping keyring test due to keyring error: %v", err)
	}

	loadedAccess, loadedRefresh, err := LoadTokens()
	require.NoError(t, err)
	assert.Equal(t, accessToken, loadedAccess)
	assert.Equal(t, refreshToken, loadedRefresh)

	err = DeleteTokens()
	assert.NoError(t, err)
}

func TestLoadTokens(t *testing.T) {
	accessToken := "test-access-token-load"
	refreshToken := "test-refresh-token-load"

	err := SaveTokens(accessToken, refreshToken)
	if err != nil {
		t.Skipf("Skipping keyring test due to keyring error: %v", err)
	}
	defer DeleteTokens()

	loadedAccess, loadedRefresh, err := LoadTokens()
	require.NoError(t, err)
	assert.Equal(t, accessToken, loadedAccess)
	assert.Equal(t, refreshToken, loadedRefresh)
}

func TestLoadTokens_NotFound(t *testing.T) {
	_ = DeleteTokens()

	_, _, err := LoadTokens()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load access token")
}

func TestDeleteTokens(t *testing.T) {
	accessToken := "test-access-token-delete"
	refreshToken := "test-refresh-token-delete"

	err := SaveTokens(accessToken, refreshToken)
	if err != nil {
		t.Skipf("Skipping keyring test due to keyring error: %v", err)
	}

	_, _, err = LoadTokens()
	require.NoError(t, err)

	err = DeleteTokens()
	assert.NoError(t, err)

	_, _, err = LoadTokens()
	assert.Error(t, err)
}
