package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetIPAddress(t *testing.T) {
	ipv4, err := GetIPAddress()

	require.NoError(t, err)

	require.NotEmpty(t, ipv4)
}
