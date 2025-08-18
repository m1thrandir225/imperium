package internal

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/client/config"
)

type MatchmakeService struct {
}

func NewMatchmakeService(config *config.Config) *MatchmakeService {
	return &MatchmakeService{}
}

func (s *MatchmakeService) Connect(ctx context.Context, cmd ConnectRequest) (*ConnectResponse, error) {
	return &ConnectResponse{}, nil
}

func (s *MatchmakeService) Disconnect(ctx context.Context, cmd DisconnectRequest) (*DisconnectResponse, error) {
	return &DisconnectResponse{}, nil
}
