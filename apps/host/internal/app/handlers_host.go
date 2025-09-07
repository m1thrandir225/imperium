package app

import (
	"context"

	"github.com/m1thrandir225/imperium/apps/host/internal/host"
	"github.com/m1thrandir225/imperium/apps/host/internal/state"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

func (a *App) WireHostHandlers() {
	initCh := a.Bus.Subscribe(EventHostInitRequested)
	go func() {
		for range initCh {
			hostName, err := util.GetHostname()
			if err != nil {
				continue
			}

			ip, err := util.GetIPAddress()
			if err != nil {
				continue
			}

			host, err := a.AuthService.GetOrCreateHost(context.Background(), hostName, ip, 8080)
			if err != nil {
				continue
			}

			_ = a.State.Update(func(s *state.AppState) {
				s.HostInfo = state.HostInfo{
					ID:   host.ID,
					Name: host.Name,
					IP:   host.IPAddress,
					Port: host.Port,
				}
			})

			if a.StatusManager != nil {
				a.StatusManager.Stop()
				a.StatusManager = nil
			}

			a.StatusManager.Start(context.Background())

			a.Bus.Publish(EventHostInitialized, HostInitializedPayload{
				Host: a.State.Get().HostInfo,
			})
		}
	}()

	statusCh := a.Bus.Subscribe(EventHostStatusChanged)
	go func() {
		for evt := range statusCh {
			payload, ok := evt.(HostStatusChangedPayload)
			if !ok {
				continue
			}

			if a.StatusManager != nil {
				switch payload.Status {
				case string(host.StatusAvailable):
					a.StatusManager.UpdateStatus(host.StatusAvailable)
				case string(host.StatusOffline):
					a.StatusManager.UpdateStatus(host.StatusOffline)
				case string(host.StatusInuse):
					a.StatusManager.UpdateStatus(host.StatusInuse)
				case string(host.StatusDisabled):
					a.StatusManager.UpdateStatus(host.StatusDisabled)
				case string(host.StatusUnknown):
					a.StatusManager.UpdateStatus(host.StatusUnknown)
				default:
					a.StatusManager.UpdateStatus(host.StatusUnknown)
				}
			}
		}
	}()

	logoutDone := a.Bus.Subscribe(EventLogoutCompleted)
	go func() {
		for range logoutDone {
			if a.StatusManager != nil {
				a.StatusManager.Stop()
				a.StatusManager = nil
			}
		}
	}()
}
