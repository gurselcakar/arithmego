package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Base styles
	Normal = lipgloss.NewStyle()
	Bold   = lipgloss.NewStyle().Bold(true)
	// Dim is for de-emphasized UI chrome (borders, separators).
	// Subtle is for secondary text content (descriptions, hints).
	// Currently identical, but kept separate for future styling flexibility.
	Dim    = lipgloss.NewStyle().Faint(true)
	Subtle = lipgloss.NewStyle().Faint(true)
	Accent = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // Cyan

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
