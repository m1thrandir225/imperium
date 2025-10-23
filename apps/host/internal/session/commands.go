package session

import "time"

// StartSessionCommand is the struct to start a session via the service
type StartSessionCommand struct {
	ProgramID    string
	ClientID     string
	ClientName   string
	SessionToken string
	Status       string
	WebrtcOffer  string
	SessionID    string
	HostID       string
	HostName     string
	StartedAt    time.Time
	CreatedAt    time.Time
}
