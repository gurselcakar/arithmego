package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type ModuloGen struct{}

func (g *ModuloGen) Label() string { return "Modulo" }

func (g *ModuloGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(moduloPatterns, diff, g.Label(), 100)
}

var moduloPatterns = PatternSet{
	game.Beginner: {
		{modSingle, 10},
	},
	game.Easy: {
		{modSingle, 10},
	},
	game.Medium: {
		{modSingle, 10},
	},
	game.Hard: {
		{modSingle, 10},
	},
	game.Expert: {
		{modSingle, 7},
		{modCompositeAdd, 3},
	},
}

func modSingle(diff game.Difficulty) (expr.Expr, bool) {
	r := ModuloRanges[diff]
	divisor := RandomInRange(r[0].Min, r[0].Max)
	dividend := RandomInRange(r[1].Min, r[1].Max)
	if dividend <= divisor {
		dividend = divisor + RandomInRange(1, divisor*2)
	}
	return &expr.BinOp{Op: expr.OpMod, Left: &expr.Num{Value: dividend}, Right: &expr.Num{Value: divisor}}, true
}

// modCompositeAdd: a mod b + c (Expert only)
func modCompositeAdd(diff game.Difficulty) (expr.Expr, bool) {
	r := ModuloRanges[diff]
	divisor := RandomInRange(r[0].Min, r[0].Max)
	dividend := RandomInRange(r[1].Min, r[1].Max)
	if dividend <= divisor {
		dividend = divisor + RandomInRange(1, divisor*2)
	}
	c := RandomInRange(1, 20)
	return &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.BinOp{Op: expr.OpMod, Left: &expr.Num{Value: dividend}, Right: &expr.Num{Value: divisor}},
		Right: &expr.Num{Value: c},
	}, true
}
