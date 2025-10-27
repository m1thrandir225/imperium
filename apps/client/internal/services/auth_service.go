package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/client/internal/models"
)

type AuthService interface {
	Service
	Login(ctx context.Context, cmd models.LoginRequest) (*models.LoginResponse, error)
	RegisterUser(ctx context.Context, cmd models.RegisterUserRequest) (*models.RegisterUserResponse, error)
	RefreshToken(ctx context.Context, cmd models.RefreshTokenRequest) (*models.RefreshTokenResponse, error)
}

type authService struct {
	authServerBaseURL string
	httpClient        *http.Client
}

func NewAuthService(authServerBaseURL string, httpClient *http.Client) AuthService {
	return &authService{
		authServerBaseURL: authServerBaseURL,
		httpClient:        httpClient,
	}
}

func (s *authService) GetAuthServerBaseURL() string {
	return s.authServerBaseURL
}

func (s *authService) GetHTTPClient() *http.Client {
	return s.httpClient
}

func (s *authService) GetLoginURL() string {
	return fmt.Sprintf("%s/login", s.authServerBaseURL)
}

func (s *authService) GetRegisterURL() string {
	return fmt.Sprintf("%s/register", s.authServerBaseURL)
}

func (s *authService) GetRefreshTokenURL() string {
	return fmt.Sprintf("%s/refresh", s.authServerBaseURL)
}

func (s *authService) Login(ctx context.Context, cmd models.LoginRequest) (*models.LoginResponse, error) {
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

	var result models.LoginResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *authService) RegisterUser(ctx context.Context, cmd models.RegisterUserRequest) (*models.RegisterUserResponse, error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	url := s.GetRegisterURL()
	log.Println("url", url)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Println("error creating request", err)
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

	var result models.RegisterUserResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *authService) RefreshToken(ctx context.Context, cmd models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
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

	var result models.RefreshTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateBaseURL updates the auth server base URL
func (s *authService) UpdateBaseURL(url string) {
	s.authServerBaseURL = url
}
