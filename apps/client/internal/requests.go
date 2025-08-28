package internal

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterClientRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
}

type RefreshTokenRequest struct {
	Token string `json:"token"`
}

type RefreshTokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type ConnectRequest struct {
	HostID string `json:"hostId"`
}

type DisconnectRequest struct {
	HostID string `json:"hostId"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User                  UserDTO   `json:"user"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresIn  time.Time `json:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Time `json:"refresh_token_expires_in"`
}

type RegisterUserResponse struct {
	User UserDTO `json:"user"`
}

type RegisterClientResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type ConnectResponse struct {
	HostID string `json:"hostId"`
}

type DisconnectResponse struct {
	HostID string `json:"hostId"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
