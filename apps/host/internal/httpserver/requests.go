package httpserver

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
	SessionID    string `json:"session_id"`
	SessionToken string `json:"session_token"`
	ClientID     string `json:"client_id"`
	ProgramID    string `json:"program_id"`
}

type AuthServerSessionResponse struct {
	Status      bool   `json:"status"`
	WebrtcOffer string `json:"webrtc_offer,omitempty"`
	Error       string `json:"error,omitempty"`
}
