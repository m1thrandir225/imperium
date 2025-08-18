package internal

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type RegisterClientRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	Token string `json:"newToken"`
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
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type RegisterUserResponse struct {
	ID        string      `json:"id"`
	Email     string      `json:"email"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Client    []ClientDTO `json:"client"`
	Hosts     []HostDTO   `json:"hosts"`
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
