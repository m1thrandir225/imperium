package host

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/session"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

type InMemoryStatusManager struct {
	hostID            string
	authServerBaseURL string
	httpClient        *httpclient.Client
	statusChan        chan Status
	stopChan          chan struct{}
	sessionService    interface{ GetCurrentSession() *session.Session }
}

func NewInMemoryStatusManager(
	hostID,
	authServerBaseURL string,
	httpClient *httpclient.Client,
	sessionService interface{ GetCurrentSession() *session.Session },
) (StatusManager, error) {

	if strings.TrimSpace(hostID) == "" {
		return nil, ErrInvalidHostID
	}

	if !util.ValidURL(authServerBaseURL) {
		return nil, ErrInvalidAuthServerBaseURL
	}

	if httpClient == nil {
		return nil, ErrInvalidHttpClient
	}
	if sessionService == nil {
		return nil, ErrInvalidSessionService
	}

	return &InMemoryStatusManager{
		hostID:            hostID,
		authServerBaseURL: authServerBaseURL,
		httpClient:        httpClient,
		statusChan:        make(chan Status, 10),
		stopChan:          make(chan struct{}),
		sessionService:    sessionService,
	}, nil
}

func (sm *InMemoryStatusManager) Start(ctx context.Context) {
	go sm.statusUpdateLoop(ctx)
}

func (sm *InMemoryStatusManager) Stop() {
	close(sm.stopChan)
}

func (sm *InMemoryStatusManager) UpdateStatus(status Status) {
	select {
	case sm.statusChan <- status:
	default:
		log.Printf("Status channel is full, dropping status update")
	}
}

func (sm *InMemoryStatusManager) statusUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sm.stopChan:
			return
		case status := <-sm.statusChan:
			sm.SendStatusUpdate(ctx, status)
		case <-ticker.C:
			//skip if there is a session
			if sm.sessionService != nil && sm.sessionService.GetCurrentSession() != nil {
				continue
			}
			sm.SendStatusUpdate(ctx, StatusAvailable)
		}
	}
}

func (sm *InMemoryStatusManager) SendStatusUpdate(ctx context.Context, status Status) {
	url := fmt.Sprintf("/api/v1/hosts/%s/status", sm.hostID)

	// Create request body with status
	requestBody := status.toAPIEnum()

	resp, err := sm.httpClient.Patch(ctx, url, requestBody, make(map[string]string), true, make(map[string]string))
	if err != nil {
		log.Printf("Failed to send status update: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status update failed with status: %d", resp.StatusCode)
		log.Printf("Response body: %s", string(resp.Body))
		return
	}

	log.Printf("Status updated to: %s", status)
}
