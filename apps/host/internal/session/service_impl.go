// Package session provides the session service for the host application.
package session

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/input"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"github.com/m1thrandir225/imperium/apps/host/internal/webrtc"
)

type sessionService struct {
	authServerBaseURL string
	httpClient        *httpclient.Client
	token             string
	programService    programs.Service
	videoRecorder     *video.Recorder
	webrtcStreamer    *webrtc.Streamer
	currentSession    *Session
	mu                sync.Mutex
}

func NewService(
	authServerBaseURL string,
	token string,
	authService interface{ GetAuthenticatedClient() *httpclient.Client },
	programService programs.Service,
	videoRecorder *video.Recorder,
	webrtcStreamer *webrtc.Streamer,
) (Service, error) {
	if !util.ValidURL(authServerBaseURL) {
		return nil, InvalidAuthBaseURL
	}

	if programService == nil {
		return nil, InvalidProgramService
	}

	if videoRecorder == nil {
		return nil, InvalidVideoRecorder
	}

	// if webrtcStreamer == nil {
	// 	return nil, InvalidWebrtcStreamer
	// }

	if authService == nil {
		return nil, InvalidAuthService
	}

	return &sessionService{
		authServerBaseURL: authServerBaseURL,
		programService:    programService,
		videoRecorder:     videoRecorder,
		webrtcStreamer:    webrtcStreamer,
		token:             token,
		httpClient:        authService.GetAuthenticatedClient(),
	}, nil
}

// Deprecated: WebRTCStreamer
func (s *sessionService) WebRTCStreamer() *webrtc.Streamer {
	return s.webrtcStreamer
}

func (s *sessionService) StartSession(ctx context.Context, cmd StartSessionCommand) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	program, err := s.programService.GetLocalProgramByID(cmd.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	programCmd, err := s.programService.LaunchProgram(program.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to launch program: %w", err)
	}

	// Create WebRTC streamer
	streamer, err := webrtc.NewStreamer()
	if err != nil {
		programCmd.Process.Kill()
		return nil, fmt.Errorf("failed to create WebRTC streamer: %w", err)
	}

	// Start video recording
	videoStream, err := s.videoRecorder.RecordScreen(nil)
	if err != nil {
		programCmd.Process.Kill()
		streamer.Close()
		return nil, fmt.Errorf("failed to start video recording: %w", err)
	}

	configFPS := s.videoRecorder.GetFPS()
	log.Printf("Starting video stream at %d FPS", configFPS)
	streamer.StartStream(videoStream, configFPS)

	session := &Session{
		ID:           cmd.SessionID,
		ProgramID:    cmd.ProgramID,
		HostID:       cmd.HostID,
		HostName:     cmd.HostName,
		ClientID:     cmd.ClientID,
		ClientName:   cmd.ClientName,
		Status:       cmd.Status,
		Process:      programCmd,
		WindowTitle:  program.Name,
		SessionToken: cmd.SessionToken,
		CreatedAt:    cmd.CreatedAt,
		StartedAt:    cmd.StartedAt,
	}

	s.currentSession = session
	s.webrtcStreamer = streamer

	return session, nil
}

func (s *sessionService) EndSession() error {
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

	s.currentSession = nil
	return nil
}

func (s *sessionService) GetCurrentSession() *Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.currentSession
}

func (s *sessionService) ProcessInputCommand(cmd input.InputCommand) {
	input.HandleCommand(cmd)
}

func (s *sessionService) GetPrograms() ([]*programs.Program, error) {
	if s.programService == nil {
		return nil, fmt.Errorf("program service not initialized")
	}

	return s.programService.GetLocalPrograms()
}

func (s *sessionService) GenerateWebRTCAnswer(offer string) (string, error) {
	if s.webrtcStreamer == nil {
		return "", fmt.Errorf("webrtc streamer not initialized")
	}

	answer, err := s.webrtcStreamer.HandleOffer(offer)
	if err != nil {
		return "", fmt.Errorf("failed to generate webRTC answer: %w", err)
	}

	return answer, nil
}

// Deprecated: UpdateVideoConfig
func (s *sessionService) UpdateVideoConfig(cfg *video.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.videoRecorder = video.NewRecorder(cfg)
}
