package httpserver

import "time"

type StartSessionRequest struct {
	SessionID string `json:"sessionId"`
	ProgramID string `json:"programId"`
	ClientID  string `json:"clientId"`
}

type StartSessionResponse struct {
	Status string `json:"status"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type AuthServerSessionRequest struct {
	SessionID    string    `json:"session_id"`
	SessionToken string    `json:"session_token"`
	ClientID     string    `json:"client_id"`
	ClientName   string    `json:"client_name"`
	ProgramID    string    `json:"program_id"`
	WebrtcOffer  string    `json:"webrtc_offer"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	HostID       string    `json:"host_id"`
	HostName     string    `json:"host_name"`
	StartedAt    time.Time `json:"started_at"`
}

type AuthServerSessionResponse struct {
	Success      bool   `json:"success"`
	WebrtcAnswer string `json:"webrtc_answer,omitempty"`
	Error        string `json:"error,omitempty"`
}
