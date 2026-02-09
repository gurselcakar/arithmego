package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type SquareGen struct{}

func (g *SquareGen) Label() string { return "Square" }

func (g *SquareGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(squarePatterns, diff, g.Label(), 100)
}

var squarePatterns = PatternSet{
	game.Beginner: {
		{squareSingle, 10},
	},
	game.Easy: {
		{squareSingle, 10},
	},
	game.Medium: {
		{squareSingle, 7},
		{squareCompositeAdd, 3},
	},
	game.Hard: {
		{squareSingle, 4},
		{squareCompositeAdd, 3},
		{squareCompositeSub, 3},
	},
	game.Expert: {
		{squareSingle, 2},
		{squareCompositeAdd, 3},
		{squareCompositeSub, 2},
		{squareTriple, 3},
	},
}

func squareSingle(diff game.Difficulty) (expr.Expr, bool) {
	r := SquareRanges[diff]
	n := RandomInRange(r.Min, r.Max)
	return &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: n}}, true
}

// squareCompositeAdd: n² + m²
func squareCompositeAdd(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Medium:
		n = RandomInRange(3, 10)
		m = RandomInRange(3, 10)
	case game.Hard:
		n = RandomInRange(5, 15)
		m = RandomInRange(5, 15)
	default: // Expert
		n = RandomInRange(5, 20)
		m = RandomInRange(5, 20)
	}
	return &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: n}},
		Right: &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: m}},
	}, true
}

// squareCompositeSub: n² − m² (n > m guaranteed)
func squareCompositeSub(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Hard:
		n = RandomInRange(6, 15)
		m = RandomInRange(3, n-1)
	default: // Expert
		n = RandomInRange(8, 20)
		m = RandomInRange(3, n-1)
	}
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: n}},
		Right: &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: m}},
	}, true
}

// squareTriple: n² + m² + p²
func squareTriple(diff game.Difficulty) (expr.Expr, bool) {
	n := RandomInRange(5, 15)
	m := RandomInRange(3, 10)
	p := RandomInRange(3, 10)
	return &expr.BinOp{
		Op: expr.OpAdd,
		Left: &expr.BinOp{
			Op:    expr.OpAdd,
			Left:  &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: n}},
			Right: &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: m}},
		},
		Right: &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: p}},
	}, true
}
