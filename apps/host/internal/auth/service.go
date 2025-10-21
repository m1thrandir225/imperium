package auth

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
)

// AuthService is the authentication service for the host application.
// It is responsible for communication with the external auth-provider service.
type Service interface {
	GetAuthURL() string
	GetAuthenticatedClient() *httpclient.Client
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	GetHostByName(ctx context.Context, hostname string) (*Host, error)
	CreateHost(ctx context.Context, req CreateHostRequest) (*Host, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	RegisterHost(ctx context.Context, hostname, ipAddress string, port int) (*Host, error)
	GetOrCreateHost(ctx context.Context, hostname, ipAddress string, port int) (*Host, error)
}
