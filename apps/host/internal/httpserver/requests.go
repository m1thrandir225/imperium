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
