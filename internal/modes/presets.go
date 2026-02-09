package modes

import (
	"time"

	"github.com/gurselcakar/arithmego/internal/game"

	// Import gen to trigger init() which registers all generators
	_ "github.com/gurselcakar/arithmego/internal/game/gen"
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
func RegisterPresets() {
	// Basic operations
	Register(&Mode{
		ID:                IDAddition,
		Name:              "Addition",
		Description:       "Calculate a + b",
		GeneratorLabel:    "Addition",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDSubtraction,
		Name:              "Subtraction",
		Description:       "Calculate a − b",
		GeneratorLabel:    "Subtraction",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDMultiplication,
		Name:              "Multiplication",
		Description:       "Calculate a × b",
		GeneratorLabel:    "Multiplication",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDDivision,
		Name:              "Division",
		Description:       "Calculate a ÷ b",
		GeneratorLabel:    "Division",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	// Power operations
	Register(&Mode{
		ID:                IDSquares,
		Name:              "Squares",
		Description:       "Calculate n²",
		GeneratorLabel:    "Square",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDCubes,
		Name:              "Cubes",
		Description:       "Calculate n³",
		GeneratorLabel:    "Cube",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDSquareRoots,
		Name:              "Square Roots",
		Description:       "Calculate √n",
		GeneratorLabel:    "Square Root",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDCubeRoots,
		Name:              "Cube Roots",
		Description:       "Calculate ³√n",
		GeneratorLabel:    "Cube Root",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	// Advanced operations
	Register(&Mode{
		ID:                IDExponents,
		Name:              "Exponents",
		Description:       "Calculate aⁿ",
		GeneratorLabel:    "Power",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDRemainders,
		Name:              "Remainders",
		Description:       "Calculate a mod b",
		GeneratorLabel:    "Modulo",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDPercentages,
		Name:              "Percentages",
		Description:       "Calculate percentages",
		GeneratorLabel:    "Percentage",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	Register(&Mode{
		ID:                IDFactorials,
		Name:              "Factorials",
		Description:       "Calculate n!",
		GeneratorLabel:    "Factorial",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategorySprint,
	})

	// Mixed modes
	Register(&Mode{
		ID:                IDMixedBasics,
		Name:              "Mixed Basics",
		Description:       "Random mix of + − × ÷",
		GeneratorLabel:    "Mixed Basics",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDMixedPowers,
		Name:              "Mixed Powers",
		Description:       "Random mix of n² n³ √n ³√n",
		GeneratorLabel:    "Mixed Powers",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDMixedAdvanced,
		Name:              "Mixed Advanced",
		Description:       "Random mix of mod % n! aⁿ",
		GeneratorLabel:    "Mixed Advanced",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})

	Register(&Mode{
		ID:                IDAnythingGoes,
		Name:              "Anything Goes",
		Description:       "Random mix of all operations",
		GeneratorLabel:    "Anything Goes",
		DefaultDifficulty: game.Medium,
		DefaultDuration:   60 * time.Second,
		Category:          CategoryChallenge,
	})
}
