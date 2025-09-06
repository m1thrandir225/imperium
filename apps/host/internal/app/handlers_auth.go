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
	return time.Now().After(t.sm.Get().UserSession.Expiry)
}

// GetAccessTokenExpiresAt returns the access token expiry from the state
func (t *stateTokens) GetAccessTokenExpiresAt() time.Time {
	return t.sm.Get().UserSession.Expiry
}

// GetRefreshToken returns the refresh token from the state
func (t *stateTokens) GetRefreshToken() string {
	return t.sm.Get().UserSession.RefreshToken
}

// func (t *stateTokens) GetRefreshTokenExpiresAt() time.Time {
// 	return t.sm.Get().UserSession.RefreshTokenExpiry
// }

func (t *stateTokens) RefreshToken(ctx context.Context) error {
	baseURL := t.baseURLFn()
	if baseURL == "" {
		return fmt.Errorf("base URL is empty")
	}

	body := map[string]string{
		"refresh_token": t.GetRefreshToken(),
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

	var response auth.RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return t.sm.Update(func(s *state.AppState) {
		s.UserSession.AccessToken = response.AccessToken
		s.UserSession.Expiry = response.ExpiresAt
	})
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

			req := auth.LoginRequest{
				Email:    payload.Email,
				Password: payload.Password,
			}

			resp, err := a.AuthService.Login(context.Background(), req)
			if err != nil {
				log.Printf("failed to login: %v", err)
			}

			err = a.State.Update(func(s *state.AppState) {
				s.UserInfo = state.UserInfo{
					ID:    resp.User.ID,
					Email: resp.User.Email,
					Name:  resp.User.Name,
				}
				s.UserSession = state.UserSession{
					AccessToken:  resp.AccessToken,
					RefreshToken: resp.RefreshToken,
					Expiry:       resp.AccessTokenExpiresAt,
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
					AccessToken:  resp.AccessToken,
					RefreshToken: resp.RefreshToken,
					Expiry:       resp.AccessTokenExpiresAt,
				},
			})

			a.Bus.Publish(EventStateSaved, a.State.Get())
			a.buildClients()
		}
	}()

	logoutCh := a.Bus.Subscribe(EventLogoutRequested)
	go func() {
		for range logoutCh {
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
