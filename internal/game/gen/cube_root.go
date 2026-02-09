package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type CubeRootGen struct{}

func (g *CubeRootGen) Label() string { return "Cube Root" }

func (g *CubeRootGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(cubeRootPatterns, diff, g.Label(), 100)
}

var cubeRootPatterns = PatternSet{
	game.Beginner: {
		{cbrtSingle, 10},
	},
	game.Easy: {
		{cbrtSingle, 10},
	},
	game.Medium: {
		{cbrtSingle, 7},
		{cbrtCompositeAdd, 3},
	},
	game.Hard: {
		{cbrtSingle, 4},
		{cbrtCompositeAdd, 3},
		{cbrtCompositeSub, 3},
	},
	game.Expert: {
		{cbrtSingle, 3},
		{cbrtCompositeAdd, 4},
		{cbrtCompositeSub, 3},
	},
}

// cbrtSingle generates ∛(n³) by picking the root value first.
func cbrtSingle(diff game.Difficulty) (expr.Expr, bool) {
	r := CubeRootRanges[diff]
	result := RandomInRange(r.Min, r.Max)
	operand := result * result * result
	return &expr.UnaryPrefix{Op: expr.OpCbrt, Operand: &expr.Num{Value: operand}}, true
}

// cbrtCompositeAdd: ∛a + ∛b
func cbrtCompositeAdd(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Medium:
		n = RandomInRange(2, 5)
		m = RandomInRange(2, 5)
	case game.Hard:
		n = RandomInRange(3, 7)
		m = RandomInRange(3, 7)
	default: // Expert
		n = RandomInRange(5, 10)
		m = RandomInRange(5, 10)
	}
	return &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.UnaryPrefix{Op: expr.OpCbrt, Operand: &expr.Num{Value: n * n * n}},
		Right: &expr.UnaryPrefix{Op: expr.OpCbrt, Operand: &expr.Num{Value: m * m * m}},
	}, true
}

// cbrtCompositeSub: ∛a − ∛b (a > b guaranteed)
func cbrtCompositeSub(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Hard:
		n = RandomInRange(4, 7)
		m = RandomInRange(2, n-1)
	default: // Expert
		n = RandomInRange(6, 10)
		m = RandomInRange(2, n-1)
	}
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.UnaryPrefix{Op: expr.OpCbrt, Operand: &expr.Num{Value: n * n * n}},
		Right: &expr.UnaryPrefix{Op: expr.OpCbrt, Operand: &expr.Num{Value: m * m * m}},
	}, true
}
