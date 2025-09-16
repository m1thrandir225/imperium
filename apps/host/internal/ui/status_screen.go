package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	uapp "github.com/m1thrandir225/imperium/apps/host/internal/app"
)

type StatusScreen struct {
	manager         *Manager
	statusLabel     *widget.Label
	sessionLabel    *widget.Label
	hostInfoLabel   *widget.Label
	lastUpdateLabel *widget.Label
	subscribed      bool
	currentStatus   string
}

func NewStatusScreen(manager *Manager) *StatusScreen {
	return &StatusScreen{
		manager: manager,
	}
}

func (s *StatusScreen) Name() string {
	return STATUS_SCREEN
}

func (s *StatusScreen) Render(w fyne.Window) fyne.CanvasObject {
	if !s.subscribed {
		s.subscribeToEvents()
		s.subscribed = true
	}

	// Create labels
	s.statusLabel = widget.NewLabel("Status: Unknown")
	s.statusLabel.TextStyle.Bold = true

	s.sessionLabel = widget.NewLabel("Session: No active session")
	s.hostInfoLabel = widget.NewLabel("Host: Not initialized")
	s.lastUpdateLabel = widget.NewLabel("Last update: Never")

	s.updateDisplay()

	refreshBtn := widget.NewButton("Refresh", func() {
		s.updateDisplay()
	})

	backBtn := widget.NewButton("Back to Main Menu", func() {
		s.manager.ShowScreen(MAIN_MENU_SCREEN)
	})

	content := container.NewVBox(
		widget.NewLabel("Host Status"),
		widget.NewSeparator(),
		s.statusLabel,
		s.sessionLabel,
		s.hostInfoLabel,
		s.lastUpdateLabel,
		widget.NewSeparator(),
		container.NewHBox(refreshBtn, backBtn),
	)

	return container.NewBorder(
		nil,
		nil,
		nil, nil,
		content,
	)
}

func (s *StatusScreen) subscribeToEvents() {
	statusCh := s.manager.bus.Subscribe(uapp.EventHostStatusChanged)
	go func() {
		for evt := range statusCh {
			fmt.Printf("DEBUG: Received EventHostStatusChanged: %v\n", evt)
			payload, ok := evt.(uapp.HostStatusChangedPayload)
			if !ok {
				fmt.Printf("DEBUG: Failed to cast to HostStatusChangedPayload\n")
				continue
			}
			s.currentStatus = payload.Status
			fmt.Printf("DEBUG: Updated currentStatus to: %s\n", s.currentStatus)
			fyne.Do(func() {
				s.updateDisplay()
			})
		}
	}()

	hostInitCh := s.manager.bus.Subscribe(uapp.EventHostInitialized)
	go func() {
		for range hostInitCh {
			fyne.Do(func() {
				s.updateDisplay()
			})
		}
	}()

	stateCh := s.manager.bus.Subscribe(uapp.EventStateUpdated)
	go func() {
		for range stateCh {
			fyne.Do(func() {
				s.updateDisplay()
			})
		}
	}()
}

func (s *StatusScreen) updateDisplay() {
	state := s.manager.GetState()

	if state.HostInfo.ID != "" {
		s.hostInfoLabel.SetText(fmt.Sprintf("Host: %s (%s:%d)",
			state.HostInfo.Name,
			state.HostInfo.IP,
			state.HostInfo.Port))
	} else {
		s.hostInfoLabel.SetText("Host: Not initialized")
	}

	currentSession := s.manager.GetCurrentSession()

	statusText := "Unknown"
	if currentSession != nil {
		statusText = "In Use"
	} else if state.HostInfo.ID != "" {
		statusText = "Available"
	} else {
		statusText = "Not initialized"
	}

	s.statusLabel.SetText(fmt.Sprintf("Status: %s", statusText))

	if currentSession != nil {
		sessionText := fmt.Sprintf("Session: %s (Client: %s)",
			currentSession.WindowTitle,
			currentSession.ClientName)
		s.sessionLabel.SetText(sessionText)
	} else {
		s.sessionLabel.SetText("Session: No active session")
	}

	s.lastUpdateLabel.SetText(fmt.Sprintf("Last update: %s",
		time.Now().Format("15:04:05")))
}
