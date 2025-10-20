package httpserver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
		buildServer func() (*Server, error)
	}{}

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

func TestServer_Serve(t *testing.T) {}

func TestServer_Stop(t *testing.T) {}
