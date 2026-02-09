package gen

import (
	"math/rand"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

// Pattern generates an expression for a given difficulty.
// Returns the expression and whether it's valid.
type Pattern func(diff game.Difficulty) (expr.Expr, bool)

// WeightedPattern pairs a pattern with a selection weight.
type WeightedPattern struct {
	Pattern Pattern
	Weight  int
}

// PatternSet maps difficulties to weighted pattern lists.
type PatternSet map[game.Difficulty][]WeightedPattern

// PickPattern selects a random pattern based on weights.
func PickPattern(patterns []WeightedPattern) Pattern {
	total := 0
	for _, p := range patterns {
		total += p.Weight
	}
	if total == 0 {
		return patterns[0].Pattern
	}

	r := rand.Intn(total)
	for _, p := range patterns {
		r -= p.Weight
		if r < 0 {
			return p.Pattern
		}
	}
	return patterns[len(patterns)-1].Pattern
}
