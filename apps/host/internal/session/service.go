package session

import (
	"context"
	"log"
	"sync"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
	"github.com/m1thrandir225/imperium/apps/host/internal/input"
	"github.com/m1thrandir225/imperium/apps/host/internal/programs"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"github.com/m1thrandir225/imperium/apps/host/internal/webrtc"
)

type Service interface {
	StartSession(ctx context.Context, cmd StartSessionCommand) (*Session, error)
	EndSession() error
	GetCurrentSession() *Session
	ProcessInputCommand(cmd input.InputCommand)
	GetPrograms() ([]*programs.Program, error)
	GenerateWebRTCAnswer(offer string) (string, error)
	WebRTCStreamer() webrtc.Streamer
	UpdateVideoConfig(cfg *video.Config)
}

type sessionService struct {
	authServerBaseURL string
	httpClient        *httpclient.Client
	token             string
	programService    programs.Service
	videoRecorder     *video.Recorder
	webrtcStreamer    webrtc.Streamer
	currentSession    *Session
	mu                sync.Mutex
}

// NewService returns a new instance of the session service
func NewService(
	authServerBaseURL string,
	token string,
	authService interface{ GetAuthenticatedClient() *httpclient.Client },
	programService programs.Service,
	videoRecorder *video.Recorder,
	webrtcStreamer webrtc.Streamer,
) (Service, error) {
	return newSessionService(
		authServerBaseURL,
		token,
		authService,
		programService,
		videoRecorder,
		webrtcStreamer,
	)
}

func newSessionService(
	authServerBaseURL string,
	token string,
	authService interface{ GetAuthenticatedClient() *httpclient.Client },
	programService programs.Service,
	videoRecorder *video.Recorder,
	webrtcStreamer webrtc.Streamer,
) (*sessionService, error) {
	if !util.ValidURL(authServerBaseURL) {
		return nil, ErrInvalidAuthBaseURL
	}

	if programService == nil {
		return nil, ErrInvalidProgramService
	}

	if videoRecorder == nil {
		return nil, ErrInvalidVideoRecorder
	}

	// if webrtcStreamer == nil {
	// 	return nil, InvalidWebrtcStreamer
	// }

	if authService == nil {
		return nil, ErrInvalidAuthService
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
func (s *sessionService) WebRTCStreamer() webrtc.Streamer {
	return s.webrtcStreamer
}

// StartSession launches the desired program and starts a new webrtc session
func (s *sessionService) StartSession(ctx context.Context, cmd StartSessionCommand) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	program, err := s.programService.GetLocalProgramByID(cmd.ProgramID)
	if err != nil {
		return nil, ErrFailedToLaunchProgram
	}

	programCmd, err := s.programService.LaunchProgram(program.Path)
	if err != nil {
		return nil, ErrFailedToLaunchProgram
	}

	// Create WebRTC streamer
	streamer, err := webrtc.NewStreamer()
	if err != nil {
		programCmd.Process.Kill()
		return nil, ErrFailedToCreateWebrtcStreamer
	}

	// Start video recording
	videoStream, err := s.videoRecorder.RecordScreen(nil)
	if err != nil {
		programCmd.Process.Kill()
		streamer.Close()
		return nil, ErrFailedToStartRecording
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
		return nil, ErrNotInitializedProgramService
	}

	return s.programService.GetLocalPrograms()
}

func (s *sessionService) GenerateWebRTCAnswer(offer string) (string, error) {
	if s.webrtcStreamer == nil {
		return "", ErrNotInitializedWebRTCStreamer
	}

	answer, err := s.webrtcStreamer.HandleOffer(offer)
	if err != nil {
		return "", ErrFailedWebRTCOfferGeneration
	}

	return answer, nil
}

// Deprecated: UpdateVideoConfig
func (s *sessionService) UpdateVideoConfig(cfg *video.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()

	recorder, err := video.NewRecorder(cfg)
	if err != nil {
		//TODO: fix error handling
		log.Printf("invalid recorder might be a nil reference: %v", err.Error())
	}
	s.videoRecorder = recorder
}
