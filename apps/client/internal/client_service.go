package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClientService struct {
	baseURL    string
	httpClient *http.Client
}

func NewClientService(authServerBaseURL string) *ClientService {
	return &ClientService{
		baseURL: authServerBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ClientService) GetBaseURL() string {
	return s.baseURL
}

func (s *ClientService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *ClientService) RegisterOrUpdateClient(ctx context.Context, req RegisterClientRequest, token string) (*RegisterClientResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/upsert", s.GetBaseURL())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register or update client failed: %s", resp.Status)
	}

	var result RegisterClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
