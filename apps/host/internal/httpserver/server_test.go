package httpserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
		buildServer func() (*Server, error)
	}{
		{
			name:        "valid setup",
			expectError: false,
			buildServer: func() (*Server, error) {
				return nil, nil
			},
		},
		{
			name:        "invalid-session-service",
			expectError: true,
			buildServer: func() (*Server, error) {
				return nil, nil
			},
		},
		{
			name:        "invalid-event-bus",
			expectError: true,
			buildServer: func() (*Server, error) {
				return nil, nil
			},
		},
	}

	for _, tc := range testCases {
		server, err := tc.buildServer()
		if tc.expectError {
			require.Error(t, err)
			require.Empty(t, server)
		} else {
			require.NoError(t, err)
			require.NotEmpty(t, server)
		}
	}
}

// TODO
func TestServer_Serve(t *testing.T) {
	context.TODO()
}

// TODO
func TestServer_Stop(t *testing.T) {
	context.TODO()
}
