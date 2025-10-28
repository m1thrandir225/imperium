package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidURL(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "valid-url",
			value:    "https://google.com",
			expected: true,
		},
		{
			name:     "valid-url-with-subpath",
			value:    "https://github.com/m1thrandir225",
			expected: true,
		},
		{
			name:     "invalid-url",
			value:    "httt",
			expected: false,
		},
		{
			name:     "invalid-url-with-spaces",
			value:    "htt lorem impsum",
			expected: false,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, ValidURL(tc.value))
	}
}
