package app

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
	EventLoginRequested    = "auth.login.requested"
	EventLoginSucceeded    = "auth.login.succeeded"
	EventLoginFailed       = "auth.login.failed"
	EventLogoutRequested   = "auth.logout.requested"
	EventLogoutCompleted   = "auth.logout.completed"
	EventRegisterRequested = "auth.register.requested"
	EventRegisterCompleted = "auth.register.completed"
	EventRegisterFailed    = "auth.register.failed"

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
