package rawg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name        string
		build       func() (Client, error)
		errExpected bool
	}{
		{
			name: "valid-api-key",
			build: func() (Client, error) {
				return New("test")
			},
			errExpected: false,
		},
		{
			name: "invalid-api-key",
			build: func() (Client, error) {
				return New("")
			},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.build()
			if tc.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRawgClient_SearchGame(t *testing.T) {
	testCases := []struct {
		name        string
		validApiKey bool
		build       func(client Client) ([]RAWGGame, error)
		searchQuery string
		errExpected bool
	}{
		{
			name:        "valid-game-search",
			validApiKey: true,
			build: func(client Client) ([]RAWGGame, error) {
				return client.SearchGame("fortnite")
			},
			errExpected: false,
		},
		{
			name:        "invalid-api-search-with-invalid-api-key",
			validApiKey: false,
			build: func(client Client) ([]RAWGGame, error) {
				return client.SearchGame("fortnite")
			},
			errExpected: true,
		},
		{
			name:        "invalid-api-search",
			validApiKey: true,
			build: func(client Client) ([]RAWGGame, error) {
				return client.SearchGame("****qqqqqq")
			},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var client Client
			if tc.validApiKey {
				rawgClient, err := New(GetRAWGApiKey())
				require.NoError(t, err)
				require.NotNil(t, rawgClient)
				client = rawgClient
			} else {
				rawgClient, err := New("test")
				require.NoError(t, err)
				require.NotNil(t, rawgClient)
				client = rawgClient
			}

			results, err := tc.build(client)

			if tc.errExpected {
				require.Error(t, err)
				require.Empty(t, results)
				require.Equal(t, len(results), 0)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, results)
			}
		})
	}
}
