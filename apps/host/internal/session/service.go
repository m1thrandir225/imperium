// Package session provides the session service for the host application.
package session

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/input"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"github.com/m1thrandir225/imperium/apps/host/internal/webrtc"
)

type Service struct {
	authServerBaseURL string
	httpClient        *httpclient.Client
	token             string
	programService    *programs.ProgramService
	videoRecorder     *video.Recorder
	webrtcStreamer    *webrtc.Streamer
	wsConn            *websocket.Conn
	currentSession    *Session
	mu                sync.Mutex
}

func NewService(
	authServerBaseURL string,
	token string,
	authService interface{ GetAuthenticatedClient() *httpclient.Client },
	programService *programs.ProgramService,
	videoRecorder *video.Recorder,
	webrtcStreamer *webrtc.Streamer,
) *Service {
	return &Service{
		authServerBaseURL: authServerBaseURL,
		programService:    programService,
		videoRecorder:     videoRecorder,
		webrtcStreamer:    webrtcStreamer,
		token:             token,
		httpClient:        authService.GetAuthenticatedClient(),
	}
}

func (s *Service) WebRTCStreamer() *webrtc.Streamer {
	return s.webrtcStreamer
}

func (s *Service) StartSession(ctx context.Context, programID, clientID string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Launch the program
	//TODO: fix
	program, err := s.programService.GetLocalProgramByPath(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	cmd, err := s.programService.LaunchProgram(program.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to launch program: %w", err)
	}

	// Wait a bit for the program to start and get window title
	time.Sleep(2 * time.Second)
	windowTitle, err := s.programService.GetWindowTitle(program.Name)
	if err != nil {
		windowTitle = program.Name // fallback
	}

	// Create WebRTC streamer
	streamer, err := webrtc.NewStreamer()
	if err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("failed to create WebRTC streamer: %w", err)
	}

	// Start video recording
	videoStream, err := s.videoRecorder.RecordWindow(windowTitle, nil)
	if err != nil {
		cmd.Process.Kill()
		streamer.Close()
		return nil, fmt.Errorf("failed to start video recording: %w", err)
	}

	// Start streaming
	streamer.StartStream(videoStream, 30) // 30 FPS

	session := &Session{
		ID:          generateSessionID(),
		ProgramID:   programID,
		HostID:      "current-host-id", // This should come from auth service
		ClientID:    clientID,
		Status:      "active",
		StartedAt:   time.Now(),
		Process:     cmd,
		WindowTitle: windowTitle,
	}

	s.currentSession = session
	s.webrtcStreamer = streamer

	// Start input handling
	go s.handleInputCommands()

	return session, nil
}

func (s *Service) handleInputCommands() {
	for {
		if s.wsConn == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		var cmd input.InputCommand
		err := s.wsConn.ReadJSON(&cmd)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket error: %v\n", err)
			}
			break
		}

		input.HandleCommand(cmd)
	}
}

func (s *Service) EndSession() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentSession == nil {
		return nil
	}

	// Kill the program process
	if s.currentSession.Process != nil {
		s.currentSession.Process.Process.Kill()
	}

	// Stop video recording
	if s.videoRecorder != nil {
		s.videoRecorder.StopRecording()
	}

	// Close WebRTC connection
	if s.webrtcStreamer != nil {
		s.webrtcStreamer.Close()
	}

	// Close WebSocket connection
	if s.wsConn != nil {
		s.wsConn.Close()
	}

	// Update session status
	now := time.Now()
	s.currentSession.EndedAt = &now
	s.currentSession.Status = "ended"

	s.currentSession = nil
	return nil
}

func (s *Service) GetCurrentSession() *Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.currentSession
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

func (s *Service) SetWebSocketConnection(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wsConn = conn
}

func (s *Service) ProcessInputCommand(cmd input.InputCommand) {
	input.HandleCommand(cmd)
}

func (s *Service) GetPrograms() ([]*programs.Program, error) {
	if s.programService == nil {
		return nil, fmt.Errorf("program service not initialized")
	}

	return s.programService.GetLocalPrograms()
}
