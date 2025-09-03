// Package auth provides the authentication service for the host application.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
)

type AuthService struct {
	config             *Config
	authServiceBaseURL string
	httpClient         *httpclient.Client
	saveConfig         func(config *Config)
}

func (s *AuthService) GetConfig() *Config {
	return s.config
}

func (s *AuthService) GetAuthURL() string {
	return s.authServiceBaseURL
}

func (s *AuthService) GetAuthenticatedClient() *httpclient.Client {
	return s.httpClient
}

func NewAuthService(authServiceBaseURL string, config *Config, saveConfig func(config *Config)) *AuthService {
	authService := &AuthService{
		authServiceBaseURL: authServiceBaseURL,
		config:             config,
		saveConfig:         saveConfig,
	}

	authService.httpClient = httpclient.NewClient(
		authServiceBaseURL,
		config,
		authService,
	)

	return authService
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	url := "/api/v1/auth/login"

	resp, err := s.httpClient.Post(ctx, url, req, make(map[string]string), false, make(map[string]string))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %d", resp.StatusCode)
	}

	var result LoginResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
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
	url := "/api/v1/auth/register"

	resp, err := s.httpClient.Post(ctx, url, req, make(map[string]string), false, make(map[string]string))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register failed: %d", resp.StatusCode)
	}

	var result RegisterResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AuthService) GetHostByName(ctx context.Context, hostname string) (*Host, error) {
	url := "/api/v1/hosts/name/" + hostname

	resp, err := s.httpClient.Get(ctx, url, make(map[string]string), make(map[string]string), true)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get host by name failed: %d", resp.StatusCode)
	}

	var result Host
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) CreateHost(ctx context.Context, req CreateHostRequest) (*Host, error) {
	url := "/api/v1/hosts"

	resp, err := s.httpClient.Post(ctx, url, req, make(map[string]string), true, make(map[string]string))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create host failed: %d", resp.StatusCode)
	}

	var result Host
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) RefreshToken(ctx context.Context) error {
	if s.config.GetRefreshToken() == "" {
		return fmt.Errorf("refresh token is empty")
	}

	requestBody := RefreshTokenRequest{
		RefreshToken: s.config.GetRefreshToken(),
	}

	url := "/api/v1/auth/refresh"

	// Refresh token endpoint doesn't need Authorization header (uses refresh token in body)
	resp, err := s.httpClient.Post(ctx, url, requestBody, make(map[string]string), false, make(map[string]string))
	if err != nil {
		return err
	}

	var response RefreshTokenResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return err
	}

	s.config.SetAccessToken(response.AccessToken)
	s.config.SetAccessTokenExpiresAt(response.ExpiresAt)
	s.saveConfig(s.config)

	return nil
}

func (s *AuthService) RegisterHost(ctx context.Context, hostname, ipAddress string, port int) (*Host, error) {
	req := CreateHostRequest{
		Name:      hostname,
		IPAddress: ipAddress,
		Port:      port,
	}
	log.Printf("Registering host with request: %+v", req)
	if jsonData, err := json.Marshal(req); err == nil {
		log.Printf("Request JSON: %s", string(jsonData))
	}
	return s.CreateHost(ctx, req)
}

func (s *AuthService) GetOrCreateHost(ctx context.Context, hostname, ipAddress string, port int) (*Host, error) {
	url := "/api/v1/hosts/get-or-create"

	req := CreateHostRequest{
		Name:      hostname,
		IPAddress: ipAddress,
		Port:      port,
	}

	resp, err := s.httpClient.Post(ctx, url, req, make(map[string]string), true, make(map[string]string))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get or create host failed: %d", resp.StatusCode)
	}

	var result Host
	if err := json.Unmarshal(resp.Body, &result); err != nil {
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
