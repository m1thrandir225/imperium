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

type InputCommand struct {
	Type   string `json:"type"`   // "keyboard", "mouse"
	Action string `json:"action"` // "press", "release", "move", "click"
	Key    string `json:"key,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Button string `json:"button,omitempty"`
}
