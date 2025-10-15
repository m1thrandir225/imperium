package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetHostname(t *testing.T) {
	expectedHostname, err := os.Hostname()
	require.NoError(t, err)
	require.NotEmpty(t, expectedHostname)

	hostname, err := GetHostname()
	require.NoError(t, err)

	require.Equal(t, expectedHostname, hostname)
}
