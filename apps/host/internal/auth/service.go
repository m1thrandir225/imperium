// Package auth provides the authentication service for the host application.
package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthService struct {
	config             *Config
	authServiceBaseURL string
	httpClient         *http.Client
	saveConfig         func(config *Config)
}

func (s *AuthService) GetConfig() *Config {
	return s.config
}

func (s *AuthService) GetAuthURL() string {
	return s.authServiceBaseURL
}

func (s *AuthService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func NewAuthService(authServiceBaseURL string, config *Config, saveConfig func(config *Config)) *AuthService {
	return &AuthService{
		authServiceBaseURL: authServiceBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		config:     config,
		saveConfig: saveConfig,
	}
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/auth/login", s.GetAuthURL())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
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
	config := s.config
	config.SetAccessToken(result.AccessToken)
	config.SetAccessTokenExpiresAt(result.AccessTokenExpiresAt)
	config.SetRefreshToken(result.RefreshToken)
	config.SetRefreshTokenExpiresAt(result.RefreshTokenExpiresAt)
	config.SetCurrentUser(result.User)

	log.Println("result", result)
	log.Println("config", config)

	s.saveConfig(config)

	return &result, nil
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/auth/register", s.GetAuthURL())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register failed: %s", resp.Status)
	}

	var result RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AuthService) CreateHost(ctx context.Context, req CreateHostRequest) (*Host, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/hosts", s.GetAuthURL())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.config.GetAccessToken())

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create host failed: %s", resp.Status)
	}

	var result Host
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	s.config.SetAccessToken("")
	s.config.SetRefreshToken("")
	s.config.SetAccessTokenExpiresAt(time.Time{})
	s.config.SetRefreshTokenExpiresAt(time.Time{})
	s.config.SetCurrentUser(User{})
	s.saveConfig(s.config)
	return nil
}
