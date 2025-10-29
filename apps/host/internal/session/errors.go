package session

import "errors"

var (
	ErrInvalidProgramService = errors.New("invalid ProgramService")
	ErrInvalidVideoRecorder  = errors.New("invalid VideoRecorder")
	ErrInvalidWebrtcStreamer = errors.New("invalid WebrtcStreamer")
	ErrInvalidAuthService    = errors.New("invalid AuthService")
	ErrInvalidAuthBaseURL    = errors.New("invalid AuthBaseURL")
)
