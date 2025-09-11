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
	Name      string `json:"client_name" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
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

type LoginResponseWithClient struct {
	LoginResponse
	Client ClientDTO `json:"client"`
}

type RegisterUserResponse struct {
	User UserDTO `json:"user"`
}

type RegisterClientResponse struct {
	ID        string `json:"id"`
	Name      string `json:"client_name"`
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

type ProgramDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	HostID      string `json:"hostId"`
}

type CreateSessionRequest struct {
	HostID    string `json:"host_id" binding:"required,uuid"`
	ClientID  string `json:"client_id" binding:"required,uuid"`
	ProgramID string `json:"program_id" binding:"required"`
}

type StartSessionRequest struct {
	WebrtcOffer string `json:"webrtc_offer" binding:"required"`
}

type EndSessionRequest struct {
	Reason       string `json:"reason,omitempty"`
	WebrtcAnswer string `json:"webrtc_answer,omitempty"`
}

type SessionDTO struct {
	ID           string     `json:"id"`
	HostID       string     `json:"host_id"`
	HostName     string     `json:"host_name"`
	ClientID     string     `json:"client_id"`
	ClientName   string     `json:"client_name"`
	Status       string     `json:"status"`
	SessionToken string     `json:"session_token"`
	WebrtcOffer  *string    `json:"webrtc_offer"`
	WebrtcAnswer *string    `json:"webrtc_answer"`
	ExpiresAt    time.Time  `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at"`
	EndedAt      *time.Time `json:"ended_at"`
	EndReason    *string    `json:"end_reason"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func messageResponse(message string) gin.H {
	return gin.H{"message": message}
}
