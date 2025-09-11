package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SessionService struct {
	baseURL    string
	httpClient *http.Client
}

func NewSessionService(baseURL string) *SessionService {
	return &SessionService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *SessionService) GetBaseURL() string {
	return s.baseURL
}

func (s *SessionService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *SessionService) CreateSession(ctx context.Context, req CreateSessionRequest, token string) (*SessionDTO, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := s.GetBaseURL()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.GetHTTPClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create session failed: %s", resp.Status)
	}

	var result SessionDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) GetSession(ctx context.Context, sessionID string, token string) (*SessionDTO, error) {
	url := fmt.Sprintf("%s/%s", s.GetBaseURL(), sessionID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.GetHTTPClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get session failed: %s", resp.Status)
	}

	var result SessionDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) StartSession(ctx context.Context, sessionID string, req StartSessionRequest, token string) (*SessionDTO, error) {
	url := fmt.Sprintf("%s/%s/start", s.GetBaseURL(), sessionID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.GetHTTPClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("start session failed: %s", resp.Status)
	}

	var result SessionDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) EndSession(ctx context.Context, sessionID string, req EndSessionRequest, token string) (*SessionDTO, error) {
	url := fmt.Sprintf("%s/%s/end", s.GetBaseURL(), sessionID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.GetHTTPClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("end session failed: %s", resp.Status)
	}

	var result SessionDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SessionService) CancelSession(ctx context.Context, sessionID string, reason string, token string) (*SessionDTO, error) {
	url := fmt.Sprintf("%s/%s/cancel?reason=%s", s.GetBaseURL(), sessionID, reason)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.GetHTTPClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cancel session failed: %s", resp.Status)
	}

	var result SessionDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
