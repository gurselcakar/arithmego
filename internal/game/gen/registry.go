package gen

import "github.com/gurselcakar/arithmego/internal/game"

var registry = make(map[string]game.Generator)

func init() {
	// Single-operation generators
	Register(&AdditionGen{})
	Register(&SubtractionGen{})
	Register(&MultiplicationGen{})
	Register(&DivisionGen{})
	Register(&SquareGen{})
	Register(&CubeGen{})
	Register(&SquareRootGen{})
	Register(&CubeRootGen{})
	Register(&PowerGen{})
	Register(&ModuloGen{})
	Register(&PercentageGen{})
	Register(&FactorialGen{})

	// Mixed mode generators
	Register(&MixedBasicsGen{})
	Register(&MixedPowersGen{})
	Register(&MixedAdvancedGen{})
	Register(&AnythingGoesGen{})
}

// Register adds a generator to the registry.
func Register(g game.Generator) {
	registry[g.Label()] = g
}

// Get retrieves a generator by label.
func Get(label string) (game.Generator, bool) {
	g, ok := registry[label]
	return g, ok
}

// All returns all registered generators.
func All() []game.Generator {
	gens := make([]game.Generator, 0, len(registry))
	for _, g := range registry {
		gens = append(gens, g)
	}
	return gens
}
