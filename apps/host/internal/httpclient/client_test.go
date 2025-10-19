package httpclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://example.com"
	tokenGetter := GetTokenGetter()
	tokenRefresher := GetTokenRefresher()
	testCases := []struct {
		name        string
		expectError bool
		buildClient func() (*Client, error)
	}{
		{
			name:        "valid",
			expectError: false,
			buildClient: func() (*Client, error) {
				return NewClient(baseURL, tokenGetter, tokenRefresher)
			},
		},
		{
			name:        "invalid-url",
			expectError: true,
			buildClient: func() (*Client, error) {
				return NewClient("", tokenGetter, tokenRefresher)
			},
		},
		{
			name:        "invalid-token-getter",
			expectError: true,
			buildClient: func() (*Client, error) {
				return NewClient(baseURL, nil, tokenRefresher)
			},
		},
		{
			name:        "invalid-token-refresher",
			expectError: true,
			buildClient: func() (*Client, error) {
				return NewClient(baseURL, tokenGetter, nil)
			},
		},
	}

	for _, tc := range testCases {
		client, err := tc.buildClient()
		if tc.expectError {
			require.Error(t, err)
			require.Empty(t, client)
		} else {
			require.NoError(t, err)
			require.NotEmpty(t, client)
			require.Equal(t, baseURL, client.GetBaseURL())
		}

	}
}
