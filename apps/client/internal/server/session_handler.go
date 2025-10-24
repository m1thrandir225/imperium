package server

import "context"

type SessionHandler interface {
	CreateSession(ctx context.Context)
	GetSession(ctx context.Context)
	StartSession(ctx context.Context)
	EndSession(ctx context.Context)
	CancelSession(ctx context.Context)
}
