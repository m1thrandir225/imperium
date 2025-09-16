package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
)

type stateTokens struct {
	sm        *state.StateManager
	baseURLFn func() string
}

// GetAccessToken returns the access token from the state
func (t *stateTokens) GetAccessToken() string {
	return t.sm.Get().UserSession.AccessToken
}

// IsAccessTokenExpired returns true if the access token is expired
func (t *stateTokens) IsAccessTokenExpired() bool {
	return time.Now().After(t.sm.Get().UserSession.AccessTokenExpiresAt)
}

// GetAccessTokenExpiresAt returns the access token expiry from the state
func (t *stateTokens) GetAccessTokenExpiresAt() time.Time {
	return t.sm.Get().UserSession.AccessTokenExpiresAt
}

// GetRefreshTokenExpiresAt returns the refresh token expiry from the state
func (t *stateTokens) GetRefreshTokenExpiresAt() time.Time {
	return t.sm.Get().UserSession.RefreshTokenExpiresAt
}

// IsRefreshTokenExpired returns true if the refresh token is expired
func (t *stateTokens) IsRefreshTokenExpired() bool {
	return time.Now().After(t.sm.Get().UserSession.RefreshTokenExpiresAt)
}

func (t *stateTokens) NeedsToRefreshToken() bool {
	return t.IsAccessTokenExpired() && !t.IsRefreshTokenExpired()
}

// GetRefreshToken returns the refresh token from the state
func (t *stateTokens) GetRefreshToken() string {
	return t.sm.Get().UserSession.RefreshToken
}

func (t *stateTokens) RefreshToken(ctx context.Context) error {
	baseURL := t.baseURLFn()
	if baseURL == "" {
		return fmt.Errorf("base URL is empty")
	}

	body := map[string]string{
		"token": t.GetRefreshToken(),
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/api/v1/auth/refresh",
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh token: %d", resp.StatusCode)
	}

	var response auth.RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	err = t.sm.Update(func(s *state.AppState) {
		s.UserSession.AccessToken = response.AccessToken
		s.UserSession.AccessTokenExpiresAt = response.ExpiresAt
	})

	if err != nil {
		return fmt.Errorf("failed to update state: %w", err)
	}

	refreshToken := t.GetRefreshToken()
	err = state.SaveTokens(response.AccessToken, refreshToken)
	if err != nil {
		log.Printf("failed to save tokens to keyring: %v", err)
	}

	log.Println("Refreshed tokens successfully.")

	return nil
}

func (a *App) httpClient() *httpclient.Client {
	return a.HTTPClient
}

func (a *App) WireAuthHandlers() {
	loginCh := a.Bus.Subscribe(EventLoginRequested)
	go func() {
		for evt := range loginCh {
			payload, ok := evt.(LoginRequestedPayload)
			if !ok {
				continue
			}

			if a.AuthService == nil || a.HTTPClient == nil {
				a.buildClients()
			}

			if a.AuthBaseURL == "" {
				a.Bus.Publish(EventLoginFailed, fmt.Errorf("auth base URL is empty"))
				return
			}

			req := auth.LoginRequest{
				Email:    payload.Email,
				Password: payload.Password,
			}

			resp, err := a.AuthService.Login(context.Background(), req)
			if err != nil {
				log.Printf("failed to login: %v", err)
				a.Bus.Publish(EventLoginFailed, err)
				continue
			}

			err = a.State.Update(func(s *state.AppState) {
				s.UserInfo = state.UserInfo{
					ID:    resp.User.ID,
					Email: resp.User.Email,
					Name:  resp.User.Name,
				}
				s.UserSession = state.UserSession{
					AccessToken:           resp.AccessToken,
					RefreshToken:          resp.RefreshToken,
					AccessTokenExpiresAt:  resp.AccessTokenExpiresAt,
					RefreshTokenExpiresAt: resp.RefreshTokenExpiresAt,
				}
			})
			if err != nil {
				log.Printf("failed to update state: %v", err)
			}

			err = state.SaveTokens(resp.AccessToken, resp.RefreshToken)
			if err != nil {
				log.Printf("failed to save tokens: %v", err)
			}

			a.Bus.Publish(EventLoginSucceeded, LoginSucceededPayload{
				User: state.UserInfo{
					ID:    resp.User.ID,
					Email: resp.User.Email,
					Name:  resp.User.Name,
				},
				Session: state.UserSession{
					AccessToken:           resp.AccessToken,
					RefreshToken:          resp.RefreshToken,
					AccessTokenExpiresAt:  resp.AccessTokenExpiresAt,
					RefreshTokenExpiresAt: resp.RefreshTokenExpiresAt,
				},
			})

			a.Bus.Publish(EventStateSaved, a.State.Get())
			a.buildClients()
			a.Bus.Publish(EventHostInitRequested, nil)
		}
	}()

	registerCh := a.Bus.Subscribe(EventRegisterRequested)
	go func() {
		for evt := range registerCh {
			payload, ok := evt.(RegisterRequestedPayload)
			if !ok {
				continue
			}
			if a.AuthService == nil || a.HTTPClient == nil {
				a.buildClients()
			}

			if a.AuthBaseURL == "" {
				a.Bus.Publish(EventRegisterFailed, fmt.Errorf("auth base URL is empty"))
				return
			}

			req := auth.RegisterRequest{
				Email:     payload.Email,
				Password:  payload.Password,
				FirstName: payload.FirstName,
				LastName:  payload.LastName,
			}

			_, err := a.AuthService.Register(context.Background(), req)
			if err != nil {
				log.Printf("failed to register: %v", err)
				a.Bus.Publish(EventRegisterFailed, err)
				continue
			}

			a.Bus.Publish(EventRegisterCompleted, nil)
		}
	}()

	logoutCh := a.Bus.Subscribe(EventLogoutRequested)
	go func() {
		for range logoutCh {
			a.stopTokenRefresher()
			err := a.State.Update(func(s *state.AppState) {
				s.UserInfo = state.UserInfo{}
				s.UserSession = state.UserSession{}
			})
			if err != nil {
				log.Printf("failed to update state: %v", err)
			}

			err = state.DeleteTokens()
			if err != nil {
				log.Printf("failed to delete tokens: %v", err)
			}

			a.Bus.Publish(EventLogoutCompleted, nil)
			a.Bus.Publish(EventStateSaved, a.State.Get())
			a.buildClients()
		}
	}()
}
