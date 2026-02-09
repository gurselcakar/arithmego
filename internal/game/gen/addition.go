package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type AdditionGen struct{}

func (g *AdditionGen) Label() string { return "Addition" }

func (g *AdditionGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(additionPatterns, diff, g.Label(), 100)
}

var additionPatterns = PatternSet{
	game.Beginner: {
		{addTwoOperands, 10},
	},
	game.Easy: {
		{addTwoOperands, 8},
		{addThreeOperands, 2},
	},
	game.Medium: {
		{addTwoOperands, 6},
		{addThreeOperands, 3},
		{addFourOperands, 1},
	},
	game.Hard: {
		{addTwoOperands, 4},
		{addThreeOperands, 4},
		{addFourOperands, 2},
	},
	game.Expert: {
		{addTwoOperands, 2},
		{addThreeOperands, 4},
		{addFourOperands, 3},
		{addFiveOperands, 1},
	},
}

func addTwoOperands(diff game.Difficulty) (expr.Expr, bool) {
	r := AdditionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(r.Min, r.Max)
	return &expr.BinOp{Op: expr.OpAdd, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}, true
}

func addThreeOperands(diff game.Difficulty) (expr.Expr, bool) {
	mr, ok := AdditionMultiRanges[diff]
	if !ok {
		return addTwoOperands(diff)
	}
	a := RandomInRange(mr.Primary.Min, mr.Primary.Max)
	b := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	c := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	return &expr.BinOp{
		Op:   expr.OpAdd,
		Left: &expr.BinOp{Op: expr.OpAdd, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
		Right: &expr.Num{Value: c},
	}, true
}

func addFourOperands(diff game.Difficulty) (expr.Expr, bool) {
	mr, ok := AdditionMultiRanges[diff]
	if !ok {
		return addTwoOperands(diff)
	}
	a := RandomInRange(mr.Primary.Min, mr.Primary.Max)
	b := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	c := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	d := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	return &expr.BinOp{
		Op: expr.OpAdd,
		Left: &expr.BinOp{
			Op:   expr.OpAdd,
			Left: &expr.BinOp{Op: expr.OpAdd, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
			Right: &expr.Num{Value: c},
		},
		Right: &expr.Num{Value: d},
	}, true
}

func addFiveOperands(diff game.Difficulty) (expr.Expr, bool) {
	mr, ok := AdditionMultiRanges[diff]
	if !ok {
		return addTwoOperands(diff)
	}
	a := RandomInRange(mr.Primary.Min, mr.Primary.Max)
	b := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	c := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	d := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	e := RandomInRange(mr.Secondary.Min, mr.Secondary.Max)
	return &expr.BinOp{
		Op: expr.OpAdd,
		Left: &expr.BinOp{
			Op: expr.OpAdd,
			Left: &expr.BinOp{
				Op:   expr.OpAdd,
				Left: &expr.BinOp{Op: expr.OpAdd, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
				Right: &expr.Num{Value: c},
			},
			Right: &expr.Num{Value: d},
		},
		Right: &expr.Num{Value: e},
	}, true
}
