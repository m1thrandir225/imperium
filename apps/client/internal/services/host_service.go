package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/client/internal/models"
)

type HostService interface {
	GetUserHosts(ctx context.Context, token string) ([]models.SimpleHost, error)
	GetHost(ctx context.Context, hostID, token string) (*models.Host, error)
	GetHostPrograms(ctx context.Context, hostID, token string) ([]models.Program, error)
}

type hostService struct {
	hostServerBaseURL string
	httpClient        *http.Client
}

func NewHostService(hostServerBaseURL string, httpClient *http.Client) HostService {
	return &hostService{
		hostServerBaseURL: hostServerBaseURL,
		httpClient:        httpClient,
	}
}

func (s *hostService) GetHostServerBaseURL() string {
	return s.hostServerBaseURL
}

func (s *hostService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *hostService) GetUserHostsURL() string {
	return s.GetHostServerBaseURL()
}

func (s *hostService) GetHostURL(hostID string) string {
	return fmt.Sprintf("%s/%s", s.GetHostServerBaseURL(), hostID)
}

func (s *hostService) GetHostProgramsURL(hostID string) string {
	return fmt.Sprintf("%s/%s/programs", s.GetHostServerBaseURL(), hostID)
}

func (s *hostService) UpdateHostServerBaseURL(hostServerBaseURL string) {
	s.hostServerBaseURL = hostServerBaseURL
}

func (s *hostService) GetUserHosts(ctx context.Context, token string) ([]models.SimpleHost, error) {
	url := s.GetUserHostsURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("resp.StatusCode", resp.StatusCode)
	log.Println("resp.Status", resp.Status)
	log.Println("resp.Body", resp.Body)
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hosts []models.SimpleHost
	err = json.Unmarshal(body, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (s *hostService) GetHost(ctx context.Context, hostID, token string) (*models.Host, error) {
	url := s.GetHostURL(hostID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var host models.Host
	err = json.Unmarshal(body, &host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (s *hostService) GetHostPrograms(ctx context.Context, hostID, token string) ([]models.Program, error) {
	url := s.GetHostProgramsURL(hostID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	log.Println("url", url)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.GetHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("get host programs failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get host programs failed: %s", string(body))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var programs []models.Program
	err = json.Unmarshal(body, &programs)
	if err != nil {
		return nil, err
	}

	return programs, nil
}
