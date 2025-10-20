package session

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/host/internal/input"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
)

// TODO: refactor Service implementation to use interface for better mocking
type ServiceInterface interface {
	StartSession(ctx context.Context, cmd StartSessionCommand) (*Session, error)
	EndSession() error
	GetCurrentSession() *Session
	ProcessInputCommand(cmd input.InputCommand)
	GetPrograms() ([]*programs.Program, error)
	GenerateWebRTCAnswer(offer string) (string, error)
}
