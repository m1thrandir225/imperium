package app

import "github.com/m1thrandir225/imperium/apps/host/internal/state"

const (
	//App Lifecycle Events
	EventAppStarted   = "app.started"
	EventAppStopped   = "app.stopped"
	EventAppRestarted = "app.restarted"
	EventAppError     = "app.error"

	//State Management
	EventStateLoaded  = "state.loaded"
	EventStateSaved   = "state.saved"
	EventStateUpdated = "state.updated"

	//Setup
	EventSetupRequested = "setup.requested"
	EventSetupCompleted = "setup.completed"

	//Auth
	EventLoginRequested  = "auth.login.requested"
	EventLoginSucceeded  = "auth.login.succeeded"
	EventLoginFailed     = "auth.login.failed"
	EventLogoutRequested = "auth.logout.requested"
	EventLogoutCompleted = "auth.logout.completed"

	//Setup
	EventStatusCompleted = "setup.status.completed"

	//UI Navigation
	EventUIShowScreen = "ui.show_screen"

	// Settings
	EventSettingsSaved = "settings.saved"

	//Programs
	EventProgramsDiscoverRequested = "programs.discover.requested"
	EventProgramsDisocvered        = "programs.discovered"
	EventProgramRegisterRequested  = "programs.register.requested"
	EventProgramRegistered         = "programs.registered"

	//Host
	EventHostInitRequested = "host.init.requested"
	EventHostInitialized   = "host.initialized"
	EventHostStatusChanged = "host.status.changed"

	//Session
	EventSessionStarted = "session.started"
	EventSessionEnded   = "session.ended"
)

// Payloads

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
