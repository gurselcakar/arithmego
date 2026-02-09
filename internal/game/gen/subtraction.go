package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type SubtractionGen struct{}

func (g *SubtractionGen) Label() string { return "Subtraction" }

func (g *SubtractionGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(subtractionPatterns, diff, g.Label(), 100)
}

var subtractionPatterns = PatternSet{
	game.Beginner: {
		{subTwoPositive, 10},
	},
	game.Easy: {
		{subTwoPositive, 8},
		{subThreePositive, 2},
	},
	game.Medium: {
		{subTwo, 6},
		{subThree, 3},
		{subAddMixed, 1},
	},
	game.Hard: {
		{subTwo, 4},
		{subThree, 4},
		{subAddMixed4, 2},
	},
	game.Expert: {
		{subTwo, 3},
		{subThree, 3},
		{subAddMixed4, 3},
		{subFive, 1},
	},
}

func subTwoPositive(diff game.Difficulty) (expr.Expr, bool) {
	r := SubtractionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(1, a)
	return &expr.BinOp{Op: expr.OpSub, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}, true
}

func subTwo(diff game.Difficulty) (expr.Expr, bool) {
	r := SubtractionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(r.Min, r.Max)
	return &expr.BinOp{Op: expr.OpSub, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}, true
}

func subThreePositive(diff game.Difficulty) (expr.Expr, bool) {
	r := SubtractionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(1, a/2)
	c := RandomInRange(1, a-b)
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.BinOp{Op: expr.OpSub, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
		Right: &expr.Num{Value: c},
	}, true
}

func subThree(diff game.Difficulty) (expr.Expr, bool) {
	r := SubtractionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(r.Min/2, r.Max/2)
	c := RandomInRange(r.Min/2, r.Max/2)
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.BinOp{Op: expr.OpSub, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
		Right: &expr.Num{Value: c},
	}, true
}

// subAddMixed: a + b − c (mixed addition and subtraction)
func subAddMixed(diff game.Difficulty) (expr.Expr, bool) {
	r := SubtractionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(r.Min/2, r.Max/2)
	c := RandomInRange(r.Min/2, r.Max/2)
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.BinOp{Op: expr.OpAdd, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
		Right: &expr.Num{Value: c},
	}, true
}

// subAddMixed4: a − b + c − d
func subAddMixed4(diff game.Difficulty) (expr.Expr, bool) {
	r := SubtractionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(r.Min/3, r.Max/3)
	c := RandomInRange(r.Min/3, r.Max/3)
	d := RandomInRange(r.Min/3, r.Max/3)
	return &expr.BinOp{
		Op: expr.OpSub,
		Left: &expr.BinOp{
			Op:    expr.OpAdd,
			Left:  &expr.BinOp{Op: expr.OpSub, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
			Right: &expr.Num{Value: c},
		},
		Right: &expr.Num{Value: d},
	}, true
}

func subFive(diff game.Difficulty) (expr.Expr, bool) {
	r := SubtractionRanges[diff]
	a := RandomInRange(r.Min, r.Max)
	b := RandomInRange(r.Min/4, r.Max/4)
	c := RandomInRange(r.Min/4, r.Max/4)
	d := RandomInRange(r.Min/4, r.Max/4)
	e := RandomInRange(r.Min/4, r.Max/4)
	return &expr.BinOp{
		Op: expr.OpSub,
		Left: &expr.BinOp{
			Op: expr.OpAdd,
			Left: &expr.BinOp{
				Op:    expr.OpSub,
				Left:  &expr.BinOp{Op: expr.OpSub, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
				Right: &expr.Num{Value: c},
			},
			Right: &expr.Num{Value: d},
		},
		Right: &expr.Num{Value: e},
	}, true
}
