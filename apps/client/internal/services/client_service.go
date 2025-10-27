package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/client/internal/models"
	"github.com/m1thrandir225/imperium/apps/client/internal/util"
)

type ClientService interface {
	Service
	GetClientInfo() (hostname, ip string, err error)
	RegisterOrUpdateClient(ctx context.Context, req models.RegisterClientRequest, token string) (*models.RegisterClientResponse, error)
}

type clientService struct {
	baseURL    string
	httpClient *http.Client
}

func NewClientService(authServerBaseURL string, httpClient *http.Client) ClientService {
	return &clientService{
		baseURL:    authServerBaseURL,
		httpClient: httpClient,
	}
}

func (s *clientService) GetBaseURL() string {
	return s.baseURL
}

func (s *clientService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *clientService) GetClientInfo() (string, string, error) {
	hostname, err := util.GetHostname()
	if err != nil {
		return "", "", err
	}

	ip, err := util.GetIPV4Address()
	if err != nil {
		return "", "", err
	}

	return hostname, ip, nil
}

func (s *clientService) RegisterOrUpdateClient(ctx context.Context, req models.RegisterClientRequest, token string) (*models.RegisterClientResponse, error) {
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

	var result models.RegisterClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateBaseURL updates the base URL
func (s *clientService) UpdateBaseURL(baseURL string) {
	s.baseURL = baseURL
}
