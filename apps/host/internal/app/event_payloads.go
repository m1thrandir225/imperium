package app

import "github.com/m1thrandir225/imperium/apps/host/internal/state"

type UIShowScreenPayload struct {
	Name string
}

type LoginRequestedPayload struct {
	Email    string
	Password string
}

type LoginSucceededPayload struct {
	User    state.UserInfo
	Session state.UserSession
}

type SettingsSavedPayload struct {
	Settings state.Settings
}

type SetupRequestedPayload struct {
	FFmpegPath    string
	ServerAddress string
}

type SetupCompletedPayload struct {
	FFmpegPath    string
	ServerAddress string
}

type ProgramItem struct {
	ID          string
	Name        string
	Path        string
	Description string
}

type ProgramsDiscoveredPayload struct {
	Programs []ProgramItem
}

type ProgramRegisterRequestedPayload struct {
	Program ProgramItem
}

type ProgramRegisteredPayload struct {
	Program ProgramItem
}

type HostInitializedPayload struct {
	Host state.HostInfo
}

type HostStatusChangedPayload struct {
	Status string
}

type SessionStartedPayload struct {
	SessionID   string
	ProgramName string
	ClientName  string
}

type SessionEndedPayload struct {
	SessionID string
}

type RegisterRequestedPayload struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}
