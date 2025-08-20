package host

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
)

type HostManager struct {
	config         *Config
	authService    *auth.AuthService
	programService *programs.ProgramService
	statusManager  *StatusManager
	hostID         string
}

func NewHostManager(config *Config, authService *auth.AuthService, programService *programs.ProgramService) *HostManager {

	manager := &HostManager{
		config:         config,
		authService:    authService,
		programService: programService,
		hostID:         config.HostName,
	}

	return manager
}

func (hm *HostManager) Initialize(ctx context.Context) error {
	// Check if we have valid authentication
	if hm.authService.GetConfig().GetAccessToken() == "" {
		return fmt.Errorf("no access token available")
	}

	if hm.authService.GetConfig().IsAccessTokenExpired() {
		return fmt.Errorf("access token is expired")
	}

	// Get hostname and IP
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	ipAddress := getLocalIPAddress()
	port := 8080 // Make this configurable

	// Register host with auth server
	host, err := hm.authService.RegisterHost(ctx, hostname, ipAddress, port)
	if err != nil {
		return fmt.Errorf("failed to register host: %w", err)
	}

	hm.hostID = host.ID
	log.Printf("Host registered with ID: %s", hm.hostID)

	// Start status manager
	hm.statusManager = NewStatusManager(
		hm.hostID,
		hm.authService.GetAuthURL(),
		hm.authService.GetAuthenticatedClient(),
	)
	hm.statusManager.Start(ctx)

	// Discover and save programs
	go func() {
		if err := hm.programService.DiscoverAndSavePrograms(); err != nil {
			log.Printf("Failed to discover programs: %v", err)
		}
	}()

	return nil
}

func (hm *HostManager) GetPrograms() ([]*programs.Program, error) {
	return hm.programService.GetLocalPrograms()
}

func (hm *HostManager) UpdateStatus(status Status) {
	if hm.statusManager != nil {
		hm.statusManager.UpdateStatus(status)
	}
}

func (hm *HostManager) Shutdown() {
	if hm.statusManager != nil {
		hm.statusManager.Stop()
	}
}

// Helper functions
func getLocalIPAddress() string {
	// Implement IP address detection
	// This is a simplified version - you might want to use a more robust approach
	return "127.0.0.1"
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return fmt.Sprintf("%s/Documents/imperium", home)
}
