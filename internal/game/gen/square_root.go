package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type SquareRootGen struct{}

func (g *SquareRootGen) Label() string { return "Square Root" }

func (g *SquareRootGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(squareRootPatterns, diff, g.Label(), 100)
}

var squareRootPatterns = PatternSet{
	game.Beginner: {
		{sqrtSingle, 10},
	},
	game.Easy: {
		{sqrtSingle, 10},
	},
	game.Medium: {
		{sqrtSingle, 7},
		{sqrtCompositeAdd, 3},
	},
	game.Hard: {
		{sqrtSingle, 4},
		{sqrtCompositeAdd, 3},
		{sqrtCompositeSub, 3},
	},
	game.Expert: {
		{sqrtSingle, 3},
		{sqrtCompositeAdd, 4},
		{sqrtCompositeSub, 3},
	},
}

// sqrtSingle generates √(n²) by picking the root value first.
func sqrtSingle(diff game.Difficulty) (expr.Expr, bool) {
	r := SquareRootRanges[diff]
	result := RandomInRange(r.Min, r.Max)
	operand := result * result
	return &expr.UnaryPrefix{Op: expr.OpSqrt, Operand: &expr.Num{Value: operand}}, true
}

// sqrtCompositeAdd: √a + √b
func sqrtCompositeAdd(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Medium:
		n = RandomInRange(3, 10)
		m = RandomInRange(3, 10)
	case game.Hard:
		n = RandomInRange(5, 15)
		m = RandomInRange(5, 15)
	default: // Expert
		n = RandomInRange(8, 20)
		m = RandomInRange(8, 20)
	}
	return &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.UnaryPrefix{Op: expr.OpSqrt, Operand: &expr.Num{Value: n * n}},
		Right: &expr.UnaryPrefix{Op: expr.OpSqrt, Operand: &expr.Num{Value: m * m}},
	}, true
}

// sqrtCompositeSub: √a − √b (a > b guaranteed)
func sqrtCompositeSub(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Hard:
		n = RandomInRange(6, 15)
		m = RandomInRange(2, n-1)
	default: // Expert
		n = RandomInRange(10, 20)
		m = RandomInRange(2, n-1)
	}
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.UnaryPrefix{Op: expr.OpSqrt, Operand: &expr.Num{Value: n * n}},
		Right: &expr.UnaryPrefix{Op: expr.OpSqrt, Operand: &expr.Num{Value: m * m}},
	}, true
}
