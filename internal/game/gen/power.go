package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type PowerGen struct{}

func (g *PowerGen) Label() string { return "Power" }

func (g *PowerGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(powerPatterns, diff, g.Label(), 100)
}

var powerPatterns = PatternSet{
	game.Beginner: {
		{powSingle, 10},
	},
	game.Easy: {
		{powSingle, 10},
	},
	game.Medium: {
		{powSingle, 6},
		{powCompositeAdd, 4},
	},
	game.Hard: {
		{powCompositeAdd, 5},
		{powCompositeSub, 3},
		{powSingle, 2},
	},
	game.Expert: {
		{powCompositeAdd, 4},
		{powCompositeSub, 3},
		{powSingle, 3},
	},
}

func powSingle(diff game.Difficulty) (expr.Expr, bool) {
	r := PowerRanges[diff]
	base := RandomInRange(r[0].Min, r[0].Max)
	exp := RandomInRange(r[1].Min, r[1].Max)
	if WouldOverflow(base, exp, MaxPowerResult) {
		return nil, false
	}
	return &expr.Pow{Base: &expr.Num{Value: base}, Exp: &expr.Num{Value: exp}}, true
}

// powCompositeAdd: aⁿ + bᵐ
func powCompositeAdd(diff game.Difficulty) (expr.Expr, bool) {
	r := PowerRanges[diff]
	b1 := RandomInRange(r[0].Min, r[0].Max)
	e1 := RandomInRange(r[1].Min, r[1].Max)
	b2 := RandomInRange(r[0].Min, r[0].Max)
	e2 := RandomInRange(r[1].Min, r[1].Max)
	if WouldOverflow(b1, e1, MaxPowerResult) || WouldOverflow(b2, e2, MaxPowerResult) {
		return nil, false
	}
	return &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.Pow{Base: &expr.Num{Value: b1}, Exp: &expr.Num{Value: e1}},
		Right: &expr.Pow{Base: &expr.Num{Value: b2}, Exp: &expr.Num{Value: e2}},
	}, true
}

// powCompositeSub: aⁿ − bᵐ
func powCompositeSub(diff game.Difficulty) (expr.Expr, bool) {
	r := PowerRanges[diff]
	b1 := RandomInRange(r[0].Min, r[0].Max)
	e1 := RandomInRange(r[1].Min, r[1].Max)
	b2 := RandomInRange(r[0].Min, r[0].Max)
	e2 := RandomInRange(r[1].Min, r[1].Max)
	if WouldOverflow(b1, e1, MaxPowerResult) || WouldOverflow(b2, e2, MaxPowerResult) {
		return nil, false
	}
	return &expr.BinOp{
		Op:    expr.OpSub,
		Left:  &expr.Pow{Base: &expr.Num{Value: b1}, Exp: &expr.Num{Value: e1}},
		Right: &expr.Pow{Base: &expr.Num{Value: b2}, Exp: &expr.Num{Value: e2}},
	}, true
}
