package game

// Category groups operations by type for mode selection.
type Category string

const (
	CategoryBasic    Category = "basic"    // +, -, ×, ÷
	CategoryPower    Category = "power"    // squares, cubes, roots
	CategoryAdvanced Category = "advanced" // modulo, factorial, percentage, power
)
