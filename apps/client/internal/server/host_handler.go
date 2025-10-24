package server

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type HostHandler interface {
	GetUserHosts(ctx context.Context)
	GetHost(ctx context.Context)
	GetHostPrograms(ctx context.Context)
}

type hostHandler struct {
	config  *config.Config
	service services.HostService
}

func NewHostHandler(
	cfg *config.Config,
	service services.HostService,
) (HostHandler, error) {
	return &hostHandler{
		config:  cfg,
		service: service,
	}, nil
}

func (h *hostHandler) GetUserHosts(ctx context.Context)    {}
func (h *hostHandler) GetHost(ctx context.Context)         {}
func (h *hostHandler) GetHostPrograms(ctx context.Context) {}
