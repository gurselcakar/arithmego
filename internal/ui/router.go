package ui

// Screen represents the current screen in the application.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenPlay       // Unified play screen (Phase 11)
	ScreenGame
	ScreenPause
	ScreenResults
	ScreenPractice     // Phase 7
	ScreenStatistics   // Phase 5
	ScreenSettings     // Phase 8
	ScreenOnboarding   // Phase 9
	ScreenQuitConfirm  // Phase 11
	ScreenFeatureTour  // Post-onboarding feature introduction
)

// StartMode determines how the app should start (used by CLI commands).
type StartMode int

const (
	// StartModeMenu starts at the main menu (default behavior).
	StartModeMenu StartMode = iota
	// StartModeQuickPlay starts a quick play session immediately.
	StartModeQuickPlay
	// StartModeStatistics opens the statistics screen directly.
	StartModeStatistics
	// StartModeOnboarding starts the onboarding flow.
	StartModeOnboarding
)
