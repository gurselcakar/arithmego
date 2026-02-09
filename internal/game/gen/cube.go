package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type CubeGen struct{}

func (g *CubeGen) Label() string { return "Cube" }

func (g *CubeGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(cubePatterns, diff, g.Label(), 100)
}

var cubePatterns = PatternSet{
	game.Beginner: {
		{cubeSingle, 10},
	},
	game.Easy: {
		{cubeSingle, 10},
	},
	game.Medium: {
		{cubeSingle, 7},
		{cubeCompositeAdd, 3},
	},
	game.Hard: {
		{cubeSingle, 4},
		{cubeCompositeAdd, 3},
		{cubeCompositeSub, 3},
	},
	game.Expert: {
		{cubeSingle, 3},
		{cubeCompositeAdd, 4},
		{cubeCompositeSub, 3},
	},
}

func cubeSingle(diff game.Difficulty) (expr.Expr, bool) {
	r := CubeRanges[diff]
	n := RandomInRange(r.Min, r.Max)
	return &expr.UnarySuffix{Op: expr.OpCube, Operand: &expr.Num{Value: n}}, true
}

func cubeCompositeAdd(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Medium:
		n = RandomInRange(2, 5)
		m = RandomInRange(2, 5)
	case game.Hard:
		n = RandomInRange(3, 7)
		m = RandomInRange(2, 5)
	default: // Expert
		n = RandomInRange(4, 8)
		m = RandomInRange(2, 6)
	}
	return &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.UnarySuffix{Op: expr.OpCube, Operand: &expr.Num{Value: n}},
		Right: &expr.UnarySuffix{Op: expr.OpCube, Operand: &expr.Num{Value: m}},
	}, true
}

func cubeCompositeSub(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Hard:
		n = RandomInRange(4, 7)
		m = RandomInRange(2, n-1)
	default: // Expert
		n = RandomInRange(5, 8)
		m = RandomInRange(2, n-1)
	}
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.UnarySuffix{Op: expr.OpCube, Operand: &expr.Num{Value: n}},
		Right: &expr.UnarySuffix{Op: expr.OpCube, Operand: &expr.Num{Value: m}},
	}, true
}
