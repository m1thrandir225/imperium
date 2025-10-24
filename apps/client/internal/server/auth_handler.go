package server

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type AuthHandler interface {
	Login(ctx context.Context)
	Register(ctx context.Context)
	RefreshToken(ctx context.Context)
}

type authHandler struct {
	config  *config.Config
	service services.AuthService
}

func NewAuthHandler(
	cfg *config.Config,
	service services.AuthService,
) (AuthHandler, error) {
	return &authHandler{
		config:  cfg,
		service: service,
	}, nil
}

func (h *authHandler) Login(ctx context.Context)        {}
func (h *authHandler) Register(ctx context.Context)     {}
func (h *authHandler) RefreshToken(ctx context.Context) {}
