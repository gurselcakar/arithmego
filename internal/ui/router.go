package ui

// Screen represents the current screen in the application.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenModes      // Phase 3
	ScreenLaunch     // Phase 3
	ScreenGame
	ScreenPause
	ScreenResults
	ScreenPractice   // Phase 7
	ScreenStatistics // Phase 5
	ScreenSettings   // Phase 8
	ScreenOnboarding // Phase 8
)

// Phase 3: Add mode selection state
// Phase 6: Add Quick Play logic
// Phase 8: Add first-run detection for onboarding
