package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type MultiplicationGen struct{}

func (g *MultiplicationGen) Label() string { return "Multiplication" }

func (g *MultiplicationGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(multiplicationPatterns, diff, g.Label(), 100)
}

var multiplicationPatterns = PatternSet{
	game.Beginner: {
		{mulTwo, 10},
	},
	game.Easy: {
		{mulTwo, 8},
		{mulThree, 2},
	},
	game.Medium: {
		{mulTwo, 6},
		{mulThree, 4},
	},
	game.Hard: {
		{mulTwo, 5},
		{mulThree, 5},
	},
	game.Expert: {
		{mulTwo, 4},
		{mulThree, 4},
		{mulFour, 2},
	},
}

func mulTwo(diff game.Difficulty) (expr.Expr, bool) {
	r := MultiplicationRanges[diff]
	a := RandomInRange(r[0].Min, r[0].Max)
	b := RandomInRange(r[1].Min, r[1].Max)
	return &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}, true
}

func mulThree(diff game.Difficulty) (expr.Expr, bool) {
	r := MultiplicationRanges[diff]
	mr, ok := MultiplicationMultiRanges[diff]
	if !ok {
		return mulTwo(diff)
	}
	a := RandomInRange(r[0].Min, r[0].Max)
	b := RandomInRange(mr.Min, mr.Max)
	c := RandomInRange(mr.Min, mr.Max)
	return &expr.BinOp{
		Op:    expr.OpMul,
		Left:  &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
		Right: &expr.Num{Value: c},
	}, true
}

func mulFour(diff game.Difficulty) (expr.Expr, bool) {
	mr, ok := MultiplicationMultiRanges[diff]
	if !ok {
		return mulTwo(diff)
	}
	a := RandomInRange(mr.Min, mr.Max)
	b := RandomInRange(mr.Min, mr.Max)
	c := RandomInRange(mr.Min, mr.Max)
	d := RandomInRange(mr.Min, mr.Max)
	return &expr.BinOp{
		Op: expr.OpMul,
		Left: &expr.BinOp{
			Op:    expr.OpMul,
			Left:  &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
			Right: &expr.Num{Value: c},
		},
		Right: &expr.Num{Value: d},
	}, true
}
