package host

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/m1thrandir225/imperium/apps/host/internal/auth"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

type HostManager struct {
	authService    *auth.AuthService
	programService *programs.ProgramService
	statusManager  *StatusManager
	hostID         string
}

func NewHostManager(
	authService *auth.AuthService,
	programService *programs.ProgramService,
) *HostManager {

	hostID, err := util.GetHostname()
	if err != nil {
		hostID = ""
	}

	manager := &HostManager{
		authService:    authService,
		programService: programService,
		hostID:         hostID,
	}

	return manager
}

func (hm *HostManager) Initialize(ctx context.Context) error {
	// Check if we have valid authentication
	if hm.authService.GetAccessToken() == "" {
		return fmt.Errorf("no access token available")
	}

	if hm.authService.GetConfig().IsAccessTokenExpired() {
		return fmt.Errorf("access token is expired")
	}

	// Get hostname and IP
	hostname, err := util.GetHostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	ipAddress, err := util.GetIPAddress()
	if err != nil {
		return fmt.Errorf("failed to get IP address: %w", err)
	}

	port := 8080 // Make this configurable

	// Register host with auth server
	host, err := hm.authService.GetOrCreateHost(ctx, hostname, ipAddress, port)
	if err != nil {
		return fmt.Errorf("failed to register host: %w", err)
	}

	hm.config.HostName = host.Name
	hm.config.IPAddress = host.IPAddress
	hm.config.Port = host.Port
	hm.config.UniqueID = host.ID
	hm.config.Status = host.Status

	hm.saveGlobalConfig(hm.config)
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
