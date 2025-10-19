package httpclient

import (
	"context"
	"errors"
)

// MockTokenRefresher implements TokenRefresher interface
type MockTokenRefresher struct {
	refreshToken string
}

func NewMockTokenRefresher(refreshToken string) TokenRefresher {
	return &MockTokenRefresher{
		refreshToken: refreshToken,
	}
}

func (m *MockTokenRefresher) RefreshToken(ctx context.Context) error {
	if ctx.Value("have-errror") != nil {
		return errors.New("failed refresh token")
	}
	return nil
}
