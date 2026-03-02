package types

// AppState represents the current view state
type AppState int

const (
	StateOnboarding AppState = iota
	StateCommandCenter
	StateNewBranch
	StateSwitchBranch
	StateCommit
	StateError
	StatePushPrompt
)
