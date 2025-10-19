package httpclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
		err         error
	}{
		//TODO: implement test cases
	}

	for _, tc := range testCases {
		if tc.expectError {
			require.Error(t, tc.err)
		}

		require.NoError(t, tc.err)
	}
}
