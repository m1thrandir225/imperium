package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/client/internal/models"
)

type SessionService interface {
	CreateSession(ctx context.Context, req models.CreateSessionRequest, token string) (*models.Session, error)
	GetSession(ctx context.Context, sessionID string, token string) (*models.Session, error)
	StartSession(ctx context.Context, sessionID string, req models.StartSessionRequest, token string) (*models.Session, error)
	EndSession(ctx context.Context, sessionID string, req models.EndSessionRequest, token string) (*models.Session, error)
	CancelSession(ctx context.Context, sessionID string, reason string, token string) (*models.Session, error)
}

type sessionService struct {
	baseURL    string
	httpClient *http.Client
}

func NewSessionService(baseURL string, httpClient *http.Client) SessionService {
	return &sessionService{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (s *sessionService) GetBaseURL() string {
	return s.baseURL
}

func (s *sessionService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *sessionService) CreateSession(ctx context.Context, req models.CreateSessionRequest, token string) (*models.Session, error) {
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

	var result models.Session
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *sessionService) GetSession(ctx context.Context, sessionID string, token string) (*models.Session, error) {
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

	var result models.Session
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *sessionService) StartSession(ctx context.Context, sessionID string, req models.StartSessionRequest, token string) (*models.Session, error) {
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

	var result models.Session
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *sessionService) EndSession(ctx context.Context, sessionID string, req models.EndSessionRequest, token string) (*models.Session, error) {
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

	var result models.Session
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *sessionService) CancelSession(ctx context.Context, sessionID string, reason string, token string) (*models.Session, error) {
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

	var result models.Session
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateBaseURL updates the base URL
func (s *sessionService) UpdateBaseURL(baseURL string) {
	s.baseURL = baseURL
}
