package styles

import "github.com/charmbracelet/lipgloss"

// ANSI color constants (standard 16-color palette for terminal compatibility)
const (
	ColorRed         = lipgloss.Color("1")
	ColorGreen       = lipgloss.Color("2")
	ColorYellow      = lipgloss.Color("3")
	ColorBlue        = lipgloss.Color("4")
	ColorMagenta     = lipgloss.Color("5")
	ColorCyan        = lipgloss.Color("6")
	ColorWhite       = lipgloss.Color("7")
	ColorBrightBlue  = lipgloss.Color("12")
	ColorBrightWhite = lipgloss.Color("15")
)

var (
	// Base styles
	Normal = lipgloss.NewStyle()
	Bold   = lipgloss.NewStyle().Bold(true)
	// Dim is for de-emphasized UI chrome (borders, separators).
	// Subtle is for secondary text content (descriptions, hints).
	// Currently identical, but kept separate for future styling flexibility.
	Dim    = lipgloss.NewStyle().Faint(true)
	Subtle = lipgloss.NewStyle().Faint(true)
	Accent = lipgloss.NewStyle().Foreground(ColorCyan)

	// Feedback (brief flashes)
	Correct   = lipgloss.NewStyle().Foreground(ColorGreen)
	Incorrect = lipgloss.NewStyle().Foreground(ColorRed)

	// Layout
	Centered = lipgloss.NewStyle().Align(lipgloss.Center)

	// Selection
	Selected   = lipgloss.NewStyle().Bold(true)
	Unselected = lipgloss.NewStyle().Faint(true)

	// Scoring - Score display
	Score      = lipgloss.NewStyle().Bold(true)
	ScoreLarge = lipgloss.NewStyle().Bold(true).Foreground(ColorBrightWhite)

	// Scoring - Multiplier
	Multiplier = lipgloss.NewStyle().Foreground(ColorYellow)

	// Scoring - Streak tiers (progressively more intense)
	StreakNone        = lipgloss.NewStyle().Faint(true)
	StreakBuilding    = lipgloss.NewStyle().Foreground(ColorWhite)
	StreakActive      = lipgloss.NewStyle().Foreground(ColorGreen)
	StreakMax         = lipgloss.NewStyle().Foreground(ColorGreen).Bold(true)
	StreakBlazing     = lipgloss.NewStyle().Foreground(ColorYellow).Bold(true)
	StreakUnstoppable = lipgloss.NewStyle().Foreground(ColorMagenta).Bold(true)
	StreakLegendary   = lipgloss.NewStyle().Foreground(ColorCyan).Bold(true)

	// Scoring - Milestone announcements
	Milestone = lipgloss.NewStyle().Bold(true).Foreground(ColorYellow)

	// Branding
	Logo    = lipgloss.NewStyle().Foreground(ColorBrightBlue)
	Tagline = lipgloss.NewStyle().Foreground(ColorWhite)
)
