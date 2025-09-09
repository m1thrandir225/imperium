package session

import (
	"os/exec"
	"time"
)

type Session struct {
	ID          string     `json:"id"`
	ProgramID   string     `json:"programId"`
	HostID      string     `json:"hostId"`
	ClientID    string     `json:"clientId"`
	Status      string     `json:"status"`
	StartedAt   time.Time  `json:"startedAt"`
	EndedAt     *time.Time `json:"endedAt,omitempty"`
	Process     *exec.Cmd  `json:"-"`
	WindowTitle string     `json:"windowTitle"`
}
