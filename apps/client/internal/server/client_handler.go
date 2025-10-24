package server

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type ClientHandler interface {
	GetClientInfo(ctx context.Context)
}

type clientHandler struct {
	config  *config.Config
	service services.ClientService
}

func NewClientHandler(
	cfg *config.Config,
	service services.ClientService,
) (ClientHandler, error) {
	return &clientHandler{
		config:  cfg,
		service: service,
	}, nil
}

func (h *clientHandler) GetClientInfo(ctx context.Context) {}
