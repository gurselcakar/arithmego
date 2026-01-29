package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Base styles
	Normal     = lipgloss.NewStyle()
	Bold       = lipgloss.NewStyle().Bold(true)
	Dim        = lipgloss.NewStyle().Faint(true)

	// Feedback (brief flashes)
	Correct   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	Incorrect = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red

	// Layout
	Centered = lipgloss.NewStyle().Align(lipgloss.Center)

	// Selection
	Selected   = lipgloss.NewStyle().Bold(true)
	Unselected = lipgloss.NewStyle().Faint(true)
)

// Phase 4: Add scoring styles (multipliers, streaks, animations)
// Phase 8: Add onboarding-specific styles if needed
