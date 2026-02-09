package gen

import (
	"math/rand"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type MixedAdvancedGen struct{}

func (g *MixedAdvancedGen) Label() string { return "Mixed Advanced" }

func (g *MixedAdvancedGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(mixedAdvancedPatterns, diff, g.Label(), 100)
}

var mixedAdvancedPatterns = PatternSet{
	game.Beginner: {
		{maSingleRandom, 10},
	},
	game.Easy: {
		{maSingleRandom, 8},
		{maSimpleComposite, 2},
	},
	game.Medium: {
		{maSingleRandom, 5},
		{maComposite, 5},
	},
	game.Hard: {
		{maComposite, 7},
		{maSingleRandom, 3},
	},
	game.Expert: {
		{maComplex, 8},
		{maSingleRandom, 2},
	},
}

// maSingleRandom: single random advanced operation (modulo, factorial, percentage, power)
func maSingleRandom(diff game.Difficulty) (expr.Expr, bool) {
	generators := []func(game.Difficulty) (expr.Expr, bool){
		modSingle, factSingle, pctSingle, powSingle,
	}
	return generators[rand.Intn(len(generators))](diff)
}

// maSimpleComposite: simple combo like n! + m
func maSimpleComposite(diff game.Difficulty) (expr.Expr, bool) {
	r := FactorialRanges[diff]
	n := RandomInRange(r.Min, r.Max)
	m := RandomInRange(1, 20)
	if Factorial(n) > MaxPowerResult {
		return nil, false
	}
	return &expr.BinOp{
		Op:    randomAddSub(),
		Left:  &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: n}},
		Right: &expr.Num{Value: m},
	}, true
}

// maComposite: n! ÷ m!, 2⁴ + 3!, a mod b + c, etc.
func maComposite(diff game.Difficulty) (expr.Expr, bool) {
	switch rand.Intn(3) {
	case 0:
		return factDivision(diff)
	case 1:
		return maFactPlusPow(diff)
	default:
		return maModPlusConst(diff)
	}
}

// maFactPlusPow: n! + aⁿ or aⁿ + n!
func maFactPlusPow(diff game.Difficulty) (expr.Expr, bool) {
	fr := FactorialRanges[diff]
	n := RandomInRange(fr.Min, fr.Max)
	if Factorial(n) > MaxPowerResult {
		return nil, false
	}
	pr := PowerRanges[diff]
	base := RandomInRange(pr[0].Min, pr[0].Max)
	exp := RandomInRange(pr[1].Min, pr[1].Max)
	if WouldOverflow(base, exp, MaxPowerResult) {
		return nil, false
	}
	factExpr := &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: n}}
	powExpr := &expr.Pow{Base: &expr.Num{Value: base}, Exp: &expr.Num{Value: exp}}
	return &expr.BinOp{Op: randomAddSub(), Left: factExpr, Right: powExpr}, true
}

// maModPlusConst: a mod b + c
func maModPlusConst(diff game.Difficulty) (expr.Expr, bool) {
	mr := ModuloRanges[diff]
	divisor := RandomInRange(mr[0].Min, mr[0].Max)
	dividend := RandomInRange(mr[1].Min, mr[1].Max)
	if dividend <= divisor {
		dividend = divisor + RandomInRange(1, divisor*2)
	}
	c := RandomInRange(1, 20)
	modExpr := &expr.BinOp{Op: expr.OpMod, Left: &expr.Num{Value: dividend}, Right: &expr.Num{Value: divisor}}
	return &expr.BinOp{Op: randomAddSub(), Left: modExpr, Right: &expr.Num{Value: c}}, true
}

// maComplex: complex composites for Expert
func maComplex(diff game.Difficulty) (expr.Expr, bool) {
	switch rand.Intn(3) {
	case 0:
		// n! ÷ m! + a²
		divExpr, ok := factDivision(diff)
		if !ok {
			return maSingleRandom(diff)
		}
		n := RandomInRange(3, 8)
		sqExpr := &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: n}}
		return &expr.BinOp{Op: randomAddSub(), Left: divExpr, Right: sqExpr}, true
	case 1:
		// aⁿ mod b + c!
		pr := PowerRanges[diff]
		base := RandomInRange(pr[0].Min, pr[0].Max)
		exp := RandomInRange(pr[1].Min, pr[1].Max)
		if WouldOverflow(base, exp, MaxPowerResult) {
			return nil, false
		}
		modB := RandomInRange(3, 20)
		fr := FactorialRanges[diff]
		factN := RandomInRange(fr.Min, fr.Max)
		if Factorial(factN) > MaxPowerResult {
			return nil, false
		}
		powExpr := &expr.Pow{Base: &expr.Num{Value: base}, Exp: &expr.Num{Value: exp}}
		modExpr := &expr.BinOp{Op: expr.OpMod, Left: powExpr, Right: &expr.Num{Value: modB}}
		factExpr := &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: factN}}
		return &expr.BinOp{Op: randomAddSub(), Left: modExpr, Right: factExpr}, true
	default:
		return maFactPlusPow(diff)
	}
}
