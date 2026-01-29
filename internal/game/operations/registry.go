package operations

import "github.com/gurselcakar/arithmego/internal/game"

var registry = make(map[string]game.Operation)

// Register adds an operation to the registry.
func Register(op game.Operation) {
	registry[op.Name()] = op
}

// Get retrieves an operation by name.
func Get(name string) (game.Operation, bool) {
	op, ok := registry[name]
	return op, ok
}

// All returns all registered operations.
func All() []game.Operation {
	ops := make([]game.Operation, 0, len(registry))
	for _, op := range registry {
		ops = append(ops, op)
	}
	return ops
}

// ByCategory returns operations filtered by category.
func ByCategory(cat game.Category) []game.Operation {
	var ops []game.Operation
	for _, op := range registry {
		if op.Category() == cat {
			ops = append(ops, op)
		}
	}
	return ops
}

// BasicOperations returns the four basic operations (+, -, ×, ÷).
func BasicOperations() []game.Operation {
	return ByCategory(game.CategoryBasic)
}

// PowerOperations returns the power operations (squares, cubes, roots).
func PowerOperations() []game.Operation {
	return ByCategory(game.CategoryPower)
}

// AdvancedOperations returns the advanced operations.
func AdvancedOperations() []game.Operation {
	return ByCategory(game.CategoryAdvanced)
}
