package auth

import "time"

type LoginResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	User                  User      `json:"user"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_in"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_in"`
}

type RegisterResponse struct {
	User User `json:"user"`
}

type Host struct {
	ID        string `json:"id"`
	IPAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
}

type RefreshTokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}
