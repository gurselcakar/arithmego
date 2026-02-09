package styles

import "github.com/charmbracelet/lipgloss"

// ANSI color constants (standard 16-color palette for terminal compatibility)
const (
	colorRed         = lipgloss.Color("1")
	colorGreen       = lipgloss.Color("2")
	colorYellow      = lipgloss.Color("3")
	colorMagenta     = lipgloss.Color("5")
	colorCyan        = lipgloss.Color("6")
	colorWhite       = lipgloss.Color("7")
	colorBrightBlue  = lipgloss.Color("12")
	colorBrightWhite = lipgloss.Color("15")
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
	Accent = lipgloss.NewStyle().Foreground(colorCyan)

	// Feedback (brief flashes)
	Correct   = lipgloss.NewStyle().Foreground(colorGreen)
	Incorrect = lipgloss.NewStyle().Foreground(colorRed)

	// Selection
	Selected   = lipgloss.NewStyle().Bold(true)
	Unselected = lipgloss.NewStyle().Faint(true)

	// Scoring - Score display
	Score      = lipgloss.NewStyle().Bold(true)
	ScoreLarge = lipgloss.NewStyle().Bold(true).Foreground(colorBrightWhite)

	// Scoring - Multiplier
	Multiplier = lipgloss.NewStyle().Foreground(colorYellow)

	// Scoring - Streak tiers (progressively more intense)
	StreakNone        = lipgloss.NewStyle().Faint(true)
	StreakBuilding    = lipgloss.NewStyle().Foreground(colorWhite)
	StreakActive      = lipgloss.NewStyle().Foreground(colorGreen)
	StreakMax         = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	StreakBlazing     = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
	StreakUnstoppable = lipgloss.NewStyle().Foreground(colorMagenta).Bold(true)
	StreakLegendary   = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)

	// Scoring - Milestone announcements
	Milestone = lipgloss.NewStyle().Bold(true).Foreground(colorYellow)

	// Branding
	Logo    = lipgloss.NewStyle().Foreground(colorBrightBlue)
	Tagline = lipgloss.NewStyle().Foreground(colorWhite)
)
