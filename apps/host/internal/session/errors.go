package session

import "errors"

var (
	InvalidProgramService = errors.New("invalid ProgramService")
	InvalidVideoRecorder  = errors.New("invalid VideoRecorder")
	InvalidWebrtcStreamer = errors.New("invalid WebrtcStreamer")
	InvalidAuthService    = errors.New("invalid AuthService")
	InvalidAuthBaseURL    = errors.New("invalid AuthBaseURL")
)
