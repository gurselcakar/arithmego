package modes

import (
	"time"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/operations"
)

// Preset mode IDs
const (
	IDAdditionSprint       = "addition-sprint"
	IDSubtractionSprint    = "subtraction-sprint"
	IDMultiplicationSprint = "multiplication-sprint"
	IDDivisionSprint       = "division-sprint"
	IDMixedOperations      = "mixed-operations"
	IDSpeedRound           = "speed-round"
	IDEndurance            = "endurance"
)

// RegisterPresets registers all built-in modes.
// Must be called after operations are registered.
func RegisterPresets() {
	// Sprint modes - single operation focus
	Register(&Mode{
		ID:                IDAdditionSprint,
		Name:              "Addition Sprint",
		Description:       "Master addition with rapid-fire problems",
		Operations:        getOperations("Addition"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDSubtractionSprint,
		Name:              "Subtraction Sprint",
		Description:       "Sharpen your subtraction skills",
		Operations:        getOperations("Subtraction"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDMultiplicationSprint,
		Name:              "Multiplication Sprint",
		Description:       "Multiply your way to victory",
		Operations:        getOperations("Multiplication"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDDivisionSprint,
		Name:              "Division Sprint",
		Description:       "Divide and conquer",
		Operations:        getOperations("Division"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	// Challenge modes - multiple operations
	Register(&Mode{
		ID:                IDMixedOperations,
		Name:              "Mixed Operations",
		Description:       "All four basic operators in one session",
		Operations:        getOperations("Addition", "Subtraction", "Multiplication", "Division"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDSpeedRound,
		Name:              "Speed Round",
		Description:       "30 seconds of intense arithmetic",
		Operations:        getOperations("Addition", "Subtraction", "Multiplication", "Division"),
		DefaultDifficulty: game.Easy,
		DefaultDuration:   30 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDEndurance,
		Name:              "Endurance",
		Description:       "Two minutes of sustained focus",
		Operations:        getOperations("Addition", "Subtraction", "Multiplication", "Division"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   2 * time.Minute,
		Category:          CategoryChallenge,
	})
}

// getOperations retrieves operations by name from the registry.
// Panics if any operation is not found (indicates a bug in preset definitions).
func getOperations(names ...string) []game.Operation {
	ops := make([]game.Operation, 0, len(names))
	for _, name := range names {
		op, ok := operations.Get(name)
		if !ok {
			panic("getOperations: unknown operation " + name)
		}
		ops = append(ops, op)
	}
	return ops
}
