package programs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

const rawgBaseURL = "https://api.rawg.io/api"

type RAWGIntegration interface {
	SearchGame(query string) ([]RAWGGame, error)
}

type rawgClient struct {
	APIKey string
}

type RAWGGame struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Released string `json:"released"`
}

type RAWGSearchResponse struct {
	Results []RAWGGame `json:"results"`
}

func NewRAWGClient(apiKey string) (RAWGIntegration, error) {
	return newRawgClient(apiKey)
}

func newRawgClient(apiKey string) (*rawgClient, error) {
	if util.IsEmptyString(apiKey) {
		return nil, ErrInvalidRawgAPIKey
	}

	return &rawgClient{
		APIKey: apiKey,
	}, nil
}

func (c *rawgClient) SearchGame(query string) ([]RAWGGame, error) {
	endpoint := fmt.Sprintf("%s/games?key=%s&search=%s",
		rawgBaseURL,
		url.QueryEscape(c.APIKey),
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
