// Package rawg provides a client for the RAWG API.
package rawg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://api.rawg.io/api"

type Client interface {
	SearchGame(query string) ([]RAWGGame, error)
}

type rawgClient struct {
	apiKey string
}

// New returns a new instance of a RAWG client.
func New(apiKey string) (Client, error) {
	return newRawgClient(apiKey)
}

func newRawgClient(apiKey string) (*rawgClient, error) {
	if len(apiKey) == 0 || strings.TrimSpace(apiKey) == "" {
		return nil, ErrInvalidRawgAPIKey
	}
	return &rawgClient{
		apiKey: apiKey,
	}, nil
}

// SearchGame searches for a game by a query parameter which is usually the game title.
func (c *rawgClient) SearchGame(query string) ([]RAWGGame, error) {
	endpoint := fmt.Sprintf("%s/games?key=%s&search=%s",
		baseURL,
		url.QueryEscape(c.apiKey),
		url.QueryEscape(query),
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RAWG API error: %s", resp.Status)
	}

	var searchResp RAWGSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return searchResp.Results, nil
}
