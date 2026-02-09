package modes

import (
	"testing"
	"time"

	"github.com/gurselcakar/arithmego/internal/game"

	// Import gen to trigger init() which registers all generators
	_ "github.com/gurselcakar/arithmego/internal/game/gen"
)

func init() {
	// Register presets for testing
	RegisterPresets()
}

func TestRegisterPresets(t *testing.T) {
	modes := All()
	if len(modes) != 16 {
		t.Errorf("expected 16 preset modes, got %d", len(modes))
	}
}

func TestGetMode(t *testing.T) {
	tests := []struct {
		id       string
		expected string
	}{
		// Basic operations
		{IDAddition, "Addition"},
		{IDSubtraction, "Subtraction"},
		{IDMultiplication, "Multiplication"},
		{IDDivision, "Division"},
		// Power operations
		{IDSquares, "Squares"},
		{IDCubes, "Cubes"},
		{IDSquareRoots, "Square Roots"},
		{IDCubeRoots, "Cube Roots"},
		// Advanced operations
		{IDExponents, "Exponents"},
		{IDRemainders, "Remainders"},
		{IDPercentages, "Percentages"},
		{IDFactorials, "Factorials"},
		// Mixed modes
		{IDMixedBasics, "Mixed Basics"},
		{IDMixedPowers, "Mixed Powers"},
		{IDMixedAdvanced, "Mixed Advanced"},
		{IDAnythingGoes, "Anything Goes"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			mode, ok := Get(tt.id)
			if !ok {
				t.Fatalf("mode %q not found", tt.id)
			}
			if mode.Name != tt.expected {
				t.Errorf("expected name %q, got %q", tt.expected, mode.Name)
			}
		})
	}
}

func TestGetModeNotFound(t *testing.T) {
	_, ok := Get("nonexistent-mode")
	if ok {
		t.Error("expected mode not to be found")
	}
}

func TestAllModesHaveGeneratorLabel(t *testing.T) {
	for _, mode := range All() {
		if mode.GeneratorLabel == "" {
			t.Errorf("mode %q has no generator label", mode.Name)
		}
	}
}

func TestModeCategoryString(t *testing.T) {
	if CategorySprint.String() != "Sprint" {
		t.Errorf("expected Sprint, got %s", CategorySprint.String())
	}
	if CategoryChallenge.String() != "Challenge" {
		t.Errorf("expected Challenge, got %s", CategoryChallenge.String())
	}
}

func TestAllowedDurations(t *testing.T) {
	if len(AllowedDurations) != 4 {
		t.Errorf("expected 4 allowed durations, got %d", len(AllowedDurations))
	}

	expectedDurations := []time.Duration{
		30 * time.Second,
		60 * time.Second,
		90 * time.Second,
		2 * time.Minute,
	}

	for i, expected := range expectedDurations {
		if AllowedDurations[i].Value != expected {
			t.Errorf("duration %d: expected %v, got %v", i, expected, AllowedDurations[i].Value)
		}
	}
}

func TestFindDurationIndex(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected int
	}{
		{30 * time.Second, 0},
		{60 * time.Second, 1},
		{90 * time.Second, 2},
		{2 * time.Minute, 3},
		{45 * time.Second, 0}, // Not found, defaults to 0
	}

	for _, tt := range tests {
		idx := FindDurationIndex(tt.duration)
		if idx != tt.expected {
			t.Errorf("FindDurationIndex(%v): expected %d, got %d", tt.duration, tt.expected, idx)
		}
	}
}

func TestAllModesHaveDefaultDifficulty(t *testing.T) {
	validDifficulties := map[game.Difficulty]bool{
		game.Beginner: true,
		game.Easy:     true,
		game.Medium:   true,
		game.Hard:     true,
		game.Expert:   true,
	}

	for _, mode := range All() {
		if !validDifficulties[mode.DefaultDifficulty] {
			t.Errorf("mode %q has invalid default difficulty: %v", mode.Name, mode.DefaultDifficulty)
		}
	}
}

func TestAllModesHaveDescription(t *testing.T) {
	for _, mode := range All() {
		if mode.Description == "" {
			t.Errorf("mode %q has no description", mode.Name)
		}
	}
}
