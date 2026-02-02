package modes

import (
	"testing"
	"time"

	"github.com/gurselcakar/arithmego/internal/game"

	// Register operations for testing
	_ "github.com/gurselcakar/arithmego/internal/game/operations"
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

func TestModesByCategory(t *testing.T) {
	sprintModes := ByCategory(CategorySprint)
	if len(sprintModes) != 12 {
		t.Errorf("expected 12 sprint modes, got %d", len(sprintModes))
	}

	challengeModes := ByCategory(CategoryChallenge)
	if len(challengeModes) != 4 {
		t.Errorf("expected 4 challenge modes, got %d", len(challengeModes))
	}
}

func TestSprintModesHaveSingleOperation(t *testing.T) {
	sprintModes := ByCategory(CategorySprint)
	for _, mode := range sprintModes {
		if !mode.IsSingleOperation() {
			t.Errorf("sprint mode %q should have single operation, has %d", mode.Name, len(mode.Operations))
		}
	}
}

func TestMixedBasicsHasMultipleOperations(t *testing.T) {
	mode, ok := Get(IDMixedBasics)
	if !ok {
		t.Fatal("Mixed Basics mode not found")
	}
	if mode.IsSingleOperation() {
		t.Error("Mixed Basics should have multiple operations")
	}
	if len(mode.Operations) != 4 {
		t.Errorf("Mixed Basics should have 4 operations, got %d", len(mode.Operations))
	}
}

func TestMixedPowersHasMultipleOperations(t *testing.T) {
	mode, ok := Get(IDMixedPowers)
	if !ok {
		t.Fatal("Mixed Powers mode not found")
	}
	if mode.IsSingleOperation() {
		t.Error("Mixed Powers should have multiple operations")
	}
	if len(mode.Operations) != 4 {
		t.Errorf("Mixed Powers should have 4 operations, got %d", len(mode.Operations))
	}
}

func TestMixedAdvancedHasMultipleOperations(t *testing.T) {
	mode, ok := Get(IDMixedAdvanced)
	if !ok {
		t.Fatal("Mixed Advanced mode not found")
	}
	if mode.IsSingleOperation() {
		t.Error("Mixed Advanced should have multiple operations")
	}
	if len(mode.Operations) != 4 {
		t.Errorf("Mixed Advanced should have 4 operations, got %d", len(mode.Operations))
	}
}

func TestAnythingGoesHasAllOperations(t *testing.T) {
	mode, ok := Get(IDAnythingGoes)
	if !ok {
		t.Fatal("Anything Goes mode not found")
	}
	if mode.IsSingleOperation() {
		t.Error("Anything Goes should have multiple operations")
	}
	if len(mode.Operations) != 12 {
		t.Errorf("Anything Goes should have 12 operations, got %d", len(mode.Operations))
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

func TestModeOperationNames(t *testing.T) {
	mode, ok := Get(IDMixedBasics)
	if !ok {
		t.Fatal("Mixed Basics mode not found")
	}
	names := mode.OperationNames()
	if len(names) != 4 {
		t.Errorf("expected 4 operation names, got %d", len(names))
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

func TestAllModesHaveOperations(t *testing.T) {
	for _, mode := range All() {
		if len(mode.Operations) == 0 {
			t.Errorf("mode %q has no operations", mode.Name)
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
