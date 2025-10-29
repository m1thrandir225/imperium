// Package session provides the session service for the host application.
package session

import (
	"os/exec"
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	ProgramID    string    `json:"program_id"`
	HostID       string    `json:"host_id"`
	HostName     string    `json:"host_name"`
	ClientID     string    `json:"client_id"`
	ClientName   string    `json:"client_name"`
	Status       string    `json:"status"`
	SessionToken string    `json:"session_token"`
	StartedAt    time.Time `json:"started_at"`
	CreatedAt    time.Time `json:"created_at"`
	Process      *exec.Cmd `json:"-"`
	WindowTitle  string    `json:"window_title"`
}
