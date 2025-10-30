package video

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFFMPEGWrapper(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
		build       func() (*FFMPEGWrapper, error)
	}{
		{
			name:        "valid",
			expectError: false,
			build: func() (*FFMPEGWrapper, error) {
				return NewFFMPEGWrapper("~/example/path")
			},
		},
		{
			name:        "invalid-path",
			expectError: true,
			build: func() (*FFMPEGWrapper, error) {
				return NewFFMPEGWrapper("")
			},
		},
	}

	for _, tc := range testCases {
		wrapper, err := tc.build()
		if tc.expectError {
			require.Error(t, err)
			require.Empty(t, wrapper)
		} else {
			require.NoError(t, err)
			require.NotNil(t, wrapper)
			require.NotEmpty(t, wrapper)
		}

	}
}

// TODO: implement
func TestWrapper_Execute(t *testing.T) {
	context.TODO()
}

// TODO: implement
func TestWrapper_ExecuteWithoutOutput(t *testing.T) {
	context.TODO()
}

// TODO: implement
func TestWrapper_Version(t *testing.T) {
	context.TODO()
}

// TODO: implement
func TestWrapper_Stop(t *testing.T) {
	context.TODO()
}

// TODO: implement
func TestWrapper_ExecuteWithStdout(t *testing.T) {
	context.TODO()
}

// TODO: implement
func TestFFMPEGStream_Read(t *testing.T) {
	context.TODO()
}

// TODO: implement
func TestFFMPEGStream_Close(t *testing.T) {
	context.TODO()
}
