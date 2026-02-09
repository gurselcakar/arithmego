package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type DivisionGen struct{}

func (g *DivisionGen) Label() string { return "Division" }

func (g *DivisionGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(divisionPatterns, diff, g.Label(), 100)
}

var divisionPatterns = PatternSet{
	game.Beginner: {
		{divTwo, 10},
	},
	game.Easy: {
		{divTwo, 10},
	},
	game.Medium: {
		{divTwo, 7},
		{divChainTwo, 3},
	},
	game.Hard: {
		{divTwo, 5},
		{divChainTwo, 5},
	},
	game.Expert: {
		{divTwo, 4},
		{divChainTwo, 6},
	},
}

// divTwo generates a ÷ b using backward generation (divisor × quotient = dividend).
func divTwo(diff game.Difficulty) (expr.Expr, bool) {
	r := DivisionRanges[diff]
	divisor := RandomInRange(r[0].Min, r[0].Max)
	quotient := RandomInRange(r[1].Min, r[1].Max)
	dividend := divisor * quotient
	return &expr.BinOp{Op: expr.OpDiv, Left: &expr.Num{Value: dividend}, Right: &expr.Num{Value: divisor}}, true
}

// divChainTwo generates a ÷ b ÷ c using backward generation.
// Pick final quotient q, divisors d1, d2 → (q × d2 × d1) ÷ d1 ÷ d2 = q
func divChainTwo(diff game.Difficulty) (expr.Expr, bool) {
	r := DivisionRanges[diff]
	// Use smaller ranges for chain division to keep numbers reasonable
	minDiv := r[0].Min
	maxDiv := r[0].Max
	if maxDiv > 15 {
		maxDiv = 15
	}
	d1 := RandomInRange(minDiv, maxDiv)
	d2 := RandomInRange(minDiv, maxDiv)
	q := RandomInRange(2, 10)
	dividend := q * d1 * d2
	return &expr.BinOp{
		Op:    expr.OpDiv,
		Left:  &expr.BinOp{Op: expr.OpDiv, Left: &expr.Num{Value: dividend}, Right: &expr.Num{Value: d1}},
		Right: &expr.Num{Value: d2},
	}, true
}
