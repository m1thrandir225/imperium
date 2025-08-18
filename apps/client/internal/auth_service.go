package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthService struct {
	authServerBaseURL string
	httpClient        *http.Client
}

func NewAuthService(authServerBaseURL string) *AuthService {
	return &AuthService{
		authServerBaseURL: authServerBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *AuthService) GetAuthServerBaseURL() string {
	return s.authServerBaseURL
}

func (s *AuthService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *AuthService) GetLoginURL() string {
	return fmt.Sprintf("%s/login", s.authServerBaseURL)
}

func (s *AuthService) GetRegisterURL() string {
	return fmt.Sprintf("%s/register", s.authServerBaseURL)
}

func (s *AuthService) GetRefreshTokenURL() string {
	return fmt.Sprintf("%s/refresh-token", s.authServerBaseURL)
}

func (s *AuthService) Login(ctx context.Context, cmd LoginRequest) (*LoginResponse, error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	url := s.GetLoginURL()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", resp.Status)
	}

	var result LoginResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, cmd RegisterUserRequest) (*RegisterUserResponse, error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	url := s.GetRegisterURL()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register user failed: %s", resp.Status)
	}

	var result RegisterUserResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, cmd RefreshTokenRequest) (*RefreshTokenResponse, error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	url := s.GetRefreshTokenURL()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh token failed: %s", resp.Status)
	}

	var result RefreshTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
