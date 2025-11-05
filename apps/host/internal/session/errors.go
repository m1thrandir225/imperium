package session

import "errors"

var (
	ErrInvalidProgramService        = errors.New("invalid ProgramService")
	ErrInvalidVideoRecorder         = errors.New("invalid VideoRecorder")
	ErrInvalidWebrtcStreamer        = errors.New("invalid WebrtcStreamer")
	ErrInvalidAuthService           = errors.New("invalid AuthService")
	ErrInvalidAuthBaseURL           = errors.New("invalid AuthBaseURL")
	ErrFailedToGetProgram           = errors.New("failed to get program")
	ErrFailedToLaunchProgram        = errors.New("failed to launch program")
	ErrFailedToCreateWebrtcStreamer = errors.New("failed to create WebRTC streamer")
	ErrFailedToStartRecording       = errors.New("failed to start video recording")
	ErrNotInitializedProgramService = errors.New("program service is not initialized")
	ErrNotInitializedWebRTCStreamer = errors.New("webrtc streamer is not initialized")
	ErrFailedWebRTCOfferGeneration  = errors.New("failed to generate WebRTC offer answer")
)
