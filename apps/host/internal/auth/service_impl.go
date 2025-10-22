// Package auth provides the authentication service for the host application.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
)

type authService struct {
	authServiceBaseURL string
	httpClient         *httpclient.Client
}

func (s *authService) GetAuthURL() string {
	return s.authServiceBaseURL
}

func (s *authService) GetAuthenticatedClient() *httpclient.Client {
	return s.httpClient
}

func NewService(authServiceBaseURL string, httpClient *httpclient.Client) (Service, error) {
	// if !util.ValidURL(authServiceBaseURL) {
	// 	return nil, ErrInvalidAuthServiceBaseURL
	// }

	if httpClient == nil {
		return nil, ErrInvalidHttpClient
	}

	return &authService{
		authServiceBaseURL: authServiceBaseURL,
		httpClient:         httpClient,
	}, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
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

	log.Println("result", result)

	return &result, nil
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
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

func (s *authService) GetHostByName(ctx context.Context, hostname string) (*Host, error) {
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

func (s *authService) CreateHost(ctx context.Context, req CreateHostRequest) (*Host, error) {
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

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {

	requestBody := RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	url := "/api/v1/auth/refresh"

	// Refresh token endpoint doesn't need Authorization header (uses refresh token in body)
	resp, err := s.httpClient.Post(ctx, url, requestBody, make(map[string]string), false, make(map[string]string))
	if err != nil {
		return nil, err
	}

	var response RefreshTokenResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *authService) RegisterHost(ctx context.Context, hostname, ipAddress string, port int) (*Host, error) {
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

func (s *authService) GetOrCreateHost(ctx context.Context, hostname, ipAddress string, port int) (*Host, error) {
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
