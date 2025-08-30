package internal

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HostUriID struct {
	HostID string `uri:"hostId" binding:"required,uuid"`
}

type RegisterUserRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type RegisterClientRequest struct {
	Name      string `json:"name" binding:"required"`
	IPAddress string `json:"ipAddress" binding:"required"`
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

type SetupConfigRequest struct {
	AuthServerBaseURL string `json:"auth_server_base_url",binding:"required"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func messageResponse(message string) gin.H {
	return gin.H{"message": message}
}
