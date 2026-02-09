package gen

import (
	"math/rand"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type MixedPowersGen struct{}

func (g *MixedPowersGen) Label() string { return "Mixed Powers" }

func (g *MixedPowersGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(mixedPowersPatterns, diff, g.Label(), 100)
}

var mixedPowersPatterns = PatternSet{
	game.Beginner: {
		{mpSingleRandom, 10},
	},
	game.Easy: {
		{mpSingleRandom, 7},
		{mpSimpleComposite, 3},
	},
	game.Medium: {
		{mpSingleRandom, 4},
		{mpSumDiff, 6},
	},
	game.Hard: {
		{mpSumDiffMul, 5},
		{mpSumDiff, 5},
	},
	game.Expert: {
		{mpComplexComposite, 8},
		{mpSingleRandom, 2},
	},
}

// randomPowerSuffix picks a random power/root unary operation.
func randomPowerSuffix(n int) expr.Expr {
	switch rand.Intn(4) {
	case 0:
		return &expr.UnarySuffix{Op: expr.OpSquare, Operand: &expr.Num{Value: n}}
	case 1:
		return &expr.UnarySuffix{Op: expr.OpCube, Operand: &expr.Num{Value: n}}
	case 2:
		return &expr.UnaryPrefix{Op: expr.OpSqrt, Operand: &expr.Num{Value: n * n}}
	default:
		return &expr.UnaryPrefix{Op: expr.OpCbrt, Operand: &expr.Num{Value: n * n * n}}
	}
}

// mpSingleRandom: a single random power/root operation
func mpSingleRandom(diff game.Difficulty) (expr.Expr, bool) {
	generators := []func(game.Difficulty) (expr.Expr, bool){
		squareSingle, cubeSingle, sqrtSingle, cbrtSingle,
	}
	return generators[rand.Intn(len(generators))](diff)
}

// mpSimpleComposite: n² + m² or √a + √b
func mpSimpleComposite(diff game.Difficulty) (expr.Expr, bool) {
	n := RandomInRange(2, 8)
	m := RandomInRange(2, 8)
	a := randomPowerSuffix(n)
	b := randomPowerSuffix(m)
	return &expr.BinOp{Op: expr.OpAdd, Left: a, Right: b}, true
}

// mpSumDiff: power/root ± power/root
func mpSumDiff(diff game.Difficulty) (expr.Expr, bool) {
	var n, m int
	switch diff {
	case game.Medium:
		n = RandomInRange(3, 10)
		m = RandomInRange(3, 10)
	case game.Hard:
		n = RandomInRange(4, 12)
		m = RandomInRange(3, 10)
	default: // Expert
		n = RandomInRange(5, 15)
		m = RandomInRange(3, 12)
	}
	a := randomPowerSuffix(n)
	b := randomPowerSuffix(m)
	return &expr.BinOp{Op: randomAddSub(), Left: a, Right: b}, true
}

// mpSumDiffMul: power/root × power/root or similar
func mpSumDiffMul(diff game.Difficulty) (expr.Expr, bool) {
	n := RandomInRange(2, 6)
	m := RandomInRange(2, 6)
	a := randomPowerSuffix(n)
	b := randomPowerSuffix(m)

	ops := []expr.BinOpKind{expr.OpAdd, expr.OpSub, expr.OpMul}
	op := ops[rand.Intn(len(ops))]
	return &expr.BinOp{Op: op, Left: a, Right: b}, true
}

// mpComplexComposite: three power/root terms
func mpComplexComposite(diff game.Difficulty) (expr.Expr, bool) {
	n := RandomInRange(3, 10)
	m := RandomInRange(2, 8)
	p := RandomInRange(2, 6)
	a := randomPowerSuffix(n)
	b := randomPowerSuffix(m)
	c := randomPowerSuffix(p)
	return &expr.BinOp{
		Op:    randomAddSub(),
		Left:  &expr.BinOp{Op: randomAddSub(), Left: a, Right: b},
		Right: c,
	}, true
}
