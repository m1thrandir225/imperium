package models

import "time"

type Session struct {
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

type ConnectRequest struct {
	HostID string `json:"hostId"`
}

type ConnectResponse struct {
	HostID string `json:"hostId"`
}

type DisconnectRequest struct {
	HostID string `json:"hostId"`
}

type DisconnectResponse struct {
	HostID string `json:"hostId"`
}
