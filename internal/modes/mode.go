package modes

import (
	"time"

	"github.com/gurselcakar/arithmego/internal/game"
)

// Mode represents a game mode configuration.
type Mode struct {
	// Identification
	ID          string
	Name        string
	Description string

	// Generator label — maps to a registered generator in game/gen
	GeneratorLabel string

	// Defaults (can be overridden at launch)
	DefaultDifficulty game.Difficulty
	DefaultDuration   time.Duration

	// Category for UI grouping
	Category ModeCategory
}

// ModeCategory groups modes in the UI.
type ModeCategory int

const (
	CategorySprint ModeCategory = iota
	CategoryChallenge
)

// String returns the category display name.
func (c ModeCategory) String() string {
	switch c {
	case CategorySprint:
		return "Sprint"
	case CategoryChallenge:
		return "Challenge"
	default:
		return "Unknown"
	}
}
