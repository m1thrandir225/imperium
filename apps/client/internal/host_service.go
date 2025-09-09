package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HostService struct {
	hostServerBaseURL string
	httpClient        *http.Client
}

func NewHostService(hostServerBaseURL string) *HostService {
	return &HostService{
		hostServerBaseURL: hostServerBaseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *HostService) GetHostServerBaseURL() string {
	return s.hostServerBaseURL
}

func (s *HostService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *HostService) GetUserHostsURL() string {
	return s.GetHostServerBaseURL()
}

func (s *HostService) GetHostURL(hostID string) string {
	return fmt.Sprintf("%s/%s", s.GetHostServerBaseURL(), hostID)
}

func (s *HostService) GetHostProgramsURL(hostID string) string {
	return fmt.Sprintf("%s/%s/programs", s.GetHostServerBaseURL(), hostID)
}

func (s *HostService) GetUserHosts(ctx context.Context, token string) ([]SimpleHostDTO, error) {
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
		return nil, ErrUnauthorized
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hosts []SimpleHostDTO
	err = json.Unmarshal(body, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (s *HostService) GetHost(ctx context.Context, hostID, token string) (*HostDTO, error) {
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

	var host HostDTO
	err = json.Unmarshal(body, &host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (s *HostService) GetHostPrograms(ctx context.Context, hostID, token string) ([]ProgramDTO, error) {
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

	var programs []ProgramDTO
	err = json.Unmarshal(body, &programs)
	if err != nil {
		return nil, err
	}

	return programs, nil
}
