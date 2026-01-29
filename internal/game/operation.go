package game

// Arity represents the number of operands an operation requires.
type Arity int

const (
	Unary  Arity = 1
	Binary Arity = 2
)

// Category groups operations by type for mode selection.
type Category string

const (
	CategoryBasic    Category = "basic"    // +, -, ×, ÷
	CategoryPower    Category = "power"    // squares, cubes, roots
	CategoryAdvanced Category = "advanced" // modulo, factorial, percentage, power
)

// Operation defines the interface for all arithmetic operations.
type Operation interface {
	// Metadata
	Name() string       // Human-readable: "Addition", "Square Root"
	Symbol() string     // Display symbol: "+", "√"
	Arity() Arity       // Unary (1) or Binary (2)
	Category() Category // For grouping in modes

	// Computation
	Apply(operands []int) int

	// Difficulty scoring - returns a score from 1.0 to 10.0
	ScoreDifficulty(operands []int, answer int) float64

	// Question generation - operation knows how to generate valid questions
	// Uses ScoreDifficulty internally to match requested difficulty tier
	Generate(diff Difficulty) Question

	// Display formatting - operation knows how to format itself
	Format(operands []int) string // "5 + 3", "√49", "7²"
}
