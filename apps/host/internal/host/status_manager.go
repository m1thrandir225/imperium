package host

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/session"
)

type StatusManager struct {
	hostID            string
	authServerBaseURL string
	httpClient        *httpclient.Client
	statusChan        chan Status
	stopChan          chan struct{}
	sessionService    interface{ GetCurrentSession() *session.Session }
}

func NewStatusManager(
	hostID,
	authServerBaseURL string,
	httpClient *httpclient.Client,
	sessionService interface{ GetCurrentSession() *session.Session },
) *StatusManager {
	return &StatusManager{
		hostID:            hostID,
		authServerBaseURL: authServerBaseURL,
		httpClient:        httpClient,
		statusChan:        make(chan Status, 10),
		stopChan:          make(chan struct{}),
		sessionService:    sessionService,
	}
}

func (sm *StatusManager) Start(ctx context.Context) {
	go sm.statusUpdateLoop(ctx)
}

func (sm *StatusManager) Stop() {
	close(sm.stopChan)
}

func (sm *StatusManager) UpdateStatus(status Status) {
	select {
	case sm.statusChan <- status:
	default:
		log.Printf("Status channel is full, dropping status update")
	}
}

func (sm *StatusManager) statusUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sm.stopChan:
			return
		case status := <-sm.statusChan:
			sm.sendStatusUpdate(ctx, status)
		case <-ticker.C:
			//skip if there is a session
			if sm.sessionService != nil && sm.sessionService.GetCurrentSession() != nil {
				continue
			}
			sm.sendStatusUpdate(ctx, StatusAvailable)
		}
	}
}

func (sm *StatusManager) sendStatusUpdate(ctx context.Context, status Status) {
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
