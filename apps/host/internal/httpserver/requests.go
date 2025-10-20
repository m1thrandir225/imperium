package httpserver

import "time"

// Deprecated: StartSessionRequest is the request the host sends to the AuthServer when the
// host initiates a new session (rarely used, as the client usually starts the
// session)
type StartSessionRequest struct {
	SessionID string `json:"sessionId"`
	ProgramID string `json:"programId"`
	ClientID  string `json:"clientId"`
}

// Deprecated: StartSessionResponse the response the host gets from the AuthServer when the
// host initiates a new session
type StartSessionResponse struct {
	Status string `json:"status"`
}

// StatusResponse is the response of the current status of the host
type StatusResponse struct {
	Status string `json:"status"`
}

// AuthServerSessionRequest represents the request the AuthServer sends to the
// host when the client starts a new session
type AuthServerSessionRequest struct {
	SessionID    string    `json:"id"`
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

// AuthServerSessionResponse represents the response from the AuthServer when
// starting/ending a session
type AuthServerSessionResponse struct {
	Success      bool   `json:"success"`
	WebrtcAnswer string `json:"webrtc_answer,omitempty"`
	Error        string `json:"error,omitempty"`
}
