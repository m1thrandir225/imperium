package state

// StateManager manages the application's state and provides methods to load and save the state
type StateManager interface {
	Load() error
	Save() error
	Get() AppState
	Update(fn func(*AppState)) error
}
