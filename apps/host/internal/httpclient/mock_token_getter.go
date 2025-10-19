package httpclient

import (
	"time"
)

// MockTokenGetter implements TokenGetter interface
type MockTokenGetter struct {
	token     string
	expiresAt time.Time
}

func NewMockTokenGetter(token string, expiresAt time.Time) TokenGetter {
	return &MockTokenGetter{
		token:     token,
		expiresAt: expiresAt,
	}
}

func (m *MockTokenGetter) GetAccessToken() string {
	return m.token
}

func (m *MockTokenGetter) IsAccessTokenExpired() bool {
	return time.Now().After(m.expiresAt)
}

func (m *MockTokenGetter) GetAccessTokenExpiresAt() time.Time {
	return m.expiresAt
}
