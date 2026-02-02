package modes

import (
	"time"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/operations"
)

// Preset mode IDs
const (
	// Basic operations
	IDAddition       = "addition"
	IDSubtraction    = "subtraction"
	IDMultiplication = "multiplication"
	IDDivision       = "division"

	// Power operations
	IDSquares     = "squares"
	IDCubes       = "cubes"
	IDSquareRoots = "square-roots"
	IDCubeRoots   = "cube-roots"

	// Advanced operations
	IDExponents   = "exponents"
	IDRemainders  = "remainders"
	IDPercentages = "percentages"
	IDFactorials  = "factorials"

	// Mixed modes
	IDMixedBasics   = "mixed-basics"
	IDMixedPowers   = "mixed-powers"
	IDMixedAdvanced = "mixed-advanced"
	IDAnythingGoes  = "anything-goes"
)

// RegisterPresets registers all built-in modes.
// Must be called after operations are registered.
func RegisterPresets() {
	// Basic operations
	Register(&Mode{
		ID:                IDAddition,
		Name:              "Addition",
		Description:       "Practice addition problems",
		Operations:        getOperations("Addition"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDSubtraction,
		Name:              "Subtraction",
		Description:       "Practice subtraction problems",
		Operations:        getOperations("Subtraction"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDMultiplication,
		Name:              "Multiplication",
		Description:       "Practice multiplication problems",
		Operations:        getOperations("Multiplication"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDDivision,
		Name:              "Division",
		Description:       "Practice division problems",
		Operations:        getOperations("Division"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	// Power operations
	Register(&Mode{
		ID:                IDSquares,
		Name:              "Squares",
		Description:       "Calculate n²",
		Operations:        getOperations("Square"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDCubes,
		Name:              "Cubes",
		Description:       "Calculate n³",
		Operations:        getOperations("Cube"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDSquareRoots,
		Name:              "Square Roots",
		Description:       "Calculate √n",
		Operations:        getOperations("Square Root"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDCubeRoots,
		Name:              "Cube Roots",
		Description:       "Calculate ³√n",
		Operations:        getOperations("Cube Root"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	// Advanced operations
	Register(&Mode{
		ID:                IDExponents,
		Name:              "Exponents",
		Description:       "Calculate aⁿ",
		Operations:        getOperations("Power"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDRemainders,
		Name:              "Remainders",
		Description:       "Calculate a mod b",
		Operations:        getOperations("Modulo"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDPercentages,
		Name:              "Percentages",
		Description:       "Calculate percentages",
		Operations:        getOperations("Percentage"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDFactorials,
		Name:              "Factorials",
		Description:       "Calculate n!",
		Operations:        getOperations("Factorial"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	// Mixed modes
	Register(&Mode{
		ID:                IDMixedBasics,
		Name:              "Mixed Basics",
		Description:       "Random mix of + − × ÷",
		Operations:        getOperations("Addition", "Subtraction", "Multiplication", "Division"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDMixedPowers,
		Name:              "Mixed Powers",
		Description:       "Random mix of n² n³ √n ³√n",
		Operations:        getOperations("Square", "Cube", "Square Root", "Cube Root"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDMixedAdvanced,
		Name:              "Mixed Advanced",
		Description:       "Random mix of mod % n! aⁿ",
		Operations:        getOperations("Modulo", "Percentage", "Factorial", "Power"),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDAnythingGoes,
		Name:              "Anything Goes",
		Description:       "Random mix of all operations",
		Operations:        operations.All(),
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
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
