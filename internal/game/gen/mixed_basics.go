package gen

import (
	"math/rand"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type MixedBasicsGen struct{}

func (g *MixedBasicsGen) Label() string { return "Mixed Basics" }

func (g *MixedBasicsGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(mixedBasicsPatterns, diff, g.Label(), 100)
}

var mixedBasicsPatterns = PatternSet{
	game.Beginner: {
		{mbSingleOp, 6},
		{mbSamePrecedenceChain, 4},
	},
	game.Easy: {
		{mbSamePrecedenceChain, 4},
		{mbParenthesizedMixed, 6},
	},
	game.Medium: {
		{mbTwoOpPEMDAS, 6},
		{mbThreeOpMixed, 3},
		{mbParenthesizedMixed, 1},
	},
	game.Hard: {
		{mbThreeOpPEMDAS, 5},
		{mbFourOpPEMDAS, 3},
		{mbParallelMulDiv, 2},
	},
	game.Expert: {
		{mbFourOpPEMDAS, 4},
		{mbFiveOpPEMDAS, 2},
		{mbParallelMulDiv, 4},
	},
}

// operandForDiff returns a random operand appropriate for the difficulty.
func operandForDiff(diff game.Difficulty) int {
	switch diff {
	case game.Beginner:
		return RandomInRange(1, 9)
	case game.Easy:
		return RandomInRange(2, 20)
	case game.Medium:
		return RandomInRange(3, 30)
	case game.Hard:
		return RandomInRange(5, 50)
	case game.Expert:
		return RandomInRange(5, 99)
	default:
		return RandomInRange(1, 9)
	}
}

// smallMulOperand returns a small multiplicand appropriate for mixed expressions.
func smallMulOperand(diff game.Difficulty) int {
	switch diff {
	case game.Beginner:
		return RandomInRange(2, 5)
	case game.Easy:
		return RandomInRange(2, 9)
	case game.Medium:
		return RandomInRange(2, 12)
	case game.Hard:
		return RandomInRange(3, 15)
	case game.Expert:
		return RandomInRange(3, 20)
	default:
		return RandomInRange(2, 5)
	}
}

// randomAddSub returns either OpAdd or OpSub randomly.
func randomAddSub() expr.BinOpKind {
	if rand.Intn(2) == 0 {
		return expr.OpAdd
	}
	return expr.OpSub
}

// makeSafeDiv generates a ÷ b as BinOp with backward generation.
func makeSafeDiv(diff game.Difficulty) expr.Expr {
	divisor := RandomInRange(2, smallMulOperand(diff))
	quotient := operandForDiff(diff)
	dividend := divisor * quotient
	return &expr.BinOp{Op: expr.OpDiv, Left: &expr.Num{Value: dividend}, Right: &expr.Num{Value: divisor}}
}

// mbSingleOp: simple a ○ b (Beginner)
func mbSingleOp(diff game.Difficulty) (expr.Expr, bool) {
	ops := []expr.BinOpKind{expr.OpAdd, expr.OpSub, expr.OpMul, expr.OpDiv}
	op := ops[rand.Intn(len(ops))]
	if op == expr.OpDiv {
		return makeSafeDiv(diff), true
	}
	if op == expr.OpMul {
		a := smallMulOperand(diff)
		b := smallMulOperand(diff)
		return &expr.BinOp{Op: op, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}, true
	}
	a := operandForDiff(diff)
	b := operandForDiff(diff)
	return &expr.BinOp{Op: op, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}, true
}

// mbSamePrecedenceChain: a + b − c or a × b × c (same precedence, no PEMDAS needed)
func mbSamePrecedenceChain(diff game.Difficulty) (expr.Expr, bool) {
	if rand.Intn(2) == 0 {
		// Addition/subtraction chain
		a := operandForDiff(diff)
		b := operandForDiff(diff)
		c := operandForDiff(diff)
		op1 := randomAddSub()
		op2 := randomAddSub()
		return &expr.BinOp{
			Op:    op2,
			Left:  &expr.BinOp{Op: op1, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
			Right: &expr.Num{Value: c},
		}, true
	}
	// Multiplication chain
	a := smallMulOperand(diff)
	b := smallMulOperand(diff)
	c := smallMulOperand(diff)
	return &expr.BinOp{
		Op:    expr.OpMul,
		Left:  &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
		Right: &expr.Num{Value: c},
	}, true
}

// mbParenthesizedMixed: (a + b) × c or a × (b + c) — guided by explicit parens
func mbParenthesizedMixed(diff game.Difficulty) (expr.Expr, bool) {
	a := operandForDiff(diff)
	b := operandForDiff(diff)
	c := smallMulOperand(diff)
	addSubOp := randomAddSub()

	if rand.Intn(2) == 0 {
		// (a + b) × c
		inner := &expr.Paren{Inner: &expr.BinOp{Op: addSubOp, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}}
		return &expr.BinOp{Op: expr.OpMul, Left: inner, Right: &expr.Num{Value: c}}, true
	}
	// a ÷ (b + c) — backward generation for safe division
	sum := a + b // ensure this is the divisor part
	if sum == 0 {
		sum = 1
	}
	quotient := RandomInRange(2, smallMulOperand(diff))
	dividend := quotient * sum
	inner := &expr.Paren{Inner: &expr.BinOp{Op: expr.OpAdd, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}}
	return &expr.BinOp{Op: expr.OpDiv, Left: &expr.Num{Value: dividend}, Right: inner}, true
}

// mbTwoOpPEMDAS: a + b × c or a − b × c (requires PEMDAS)
func mbTwoOpPEMDAS(diff game.Difficulty) (expr.Expr, bool) {
	a := operandForDiff(diff)
	b := smallMulOperand(diff)
	c := smallMulOperand(diff)
	addSubOp := randomAddSub()

	if rand.Intn(2) == 0 {
		// a + b × c
		return &expr.BinOp{
			Op:    addSubOp,
			Left:  &expr.Num{Value: a},
			Right: &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: b}, Right: &expr.Num{Value: c}},
		}, true
	}
	// a × b + c
	return &expr.BinOp{
		Op:    addSubOp,
		Left:  &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}},
		Right: &expr.Num{Value: c},
	}, true
}

// mbThreeOpMixed: a + b × c − d
func mbThreeOpMixed(diff game.Difficulty) (expr.Expr, bool) {
	a := operandForDiff(diff)
	b := smallMulOperand(diff)
	c := smallMulOperand(diff)
	d := operandForDiff(diff)
	return &expr.BinOp{
		Op: randomAddSub(),
		Left: &expr.BinOp{
			Op:    randomAddSub(),
			Left:  &expr.Num{Value: a},
			Right: &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: b}, Right: &expr.Num{Value: c}},
		},
		Right: &expr.Num{Value: d},
	}, true
}

// mbThreeOpPEMDAS: a + b × c − d with full PEMDAS
func mbThreeOpPEMDAS(diff game.Difficulty) (expr.Expr, bool) {
	a := operandForDiff(diff)
	b := smallMulOperand(diff)
	c := smallMulOperand(diff)
	d := operandForDiff(diff)
	return &expr.BinOp{
		Op: randomAddSub(),
		Left: &expr.BinOp{
			Op:    randomAddSub(),
			Left:  &expr.Num{Value: a},
			Right: &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: b}, Right: &expr.Num{Value: c}},
		},
		Right: &expr.Num{Value: d},
	}, true
}

// mbFourOpPEMDAS: a + b × c − d ÷ e
func mbFourOpPEMDAS(diff game.Difficulty) (expr.Expr, bool) {
	a := operandForDiff(diff)
	b := smallMulOperand(diff)
	c := smallMulOperand(diff)

	// Safe division for the div part (backward-generated for clean results)
	divisor := RandomInRange(2, smallMulOperand(diff))
	quotient := RandomInRange(2, smallMulOperand(diff))
	dividend := divisor * quotient
	divExpr := &expr.BinOp{Op: expr.OpDiv, Left: &expr.Num{Value: dividend}, Right: &expr.Num{Value: divisor}}

	mulExpr := &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: b}, Right: &expr.Num{Value: c}}
	return &expr.BinOp{
		Op: randomAddSub(),
		Left: &expr.BinOp{
			Op:    randomAddSub(),
			Left:  &expr.Num{Value: a},
			Right: mulExpr,
		},
		Right: divExpr,
	}, true
}

// mbFiveOpPEMDAS: a × b + c − d × e + f
func mbFiveOpPEMDAS(diff game.Difficulty) (expr.Expr, bool) {
	a := smallMulOperand(diff)
	b := smallMulOperand(diff)
	c := operandForDiff(diff)
	d := smallMulOperand(diff)
	e := smallMulOperand(diff)
	f := operandForDiff(diff)

	mul1 := &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}
	mul2 := &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: d}, Right: &expr.Num{Value: e}}

	return &expr.BinOp{
		Op: randomAddSub(),
		Left: &expr.BinOp{
			Op: randomAddSub(),
			Left: &expr.BinOp{
				Op:    randomAddSub(),
				Left:  mul1,
				Right: &expr.Num{Value: c},
			},
			Right: mul2,
		},
		Right: &expr.Num{Value: f},
	}, true
}

// mbParallelMulDiv: a × b + c ÷ d (parallel high-precedence ops)
func mbParallelMulDiv(diff game.Difficulty) (expr.Expr, bool) {
	a := smallMulOperand(diff)
	b := smallMulOperand(diff)
	mulExpr := &expr.BinOp{Op: expr.OpMul, Left: &expr.Num{Value: a}, Right: &expr.Num{Value: b}}

	divExpr := makeSafeDiv(diff)

	return &expr.BinOp{
		Op:    randomAddSub(),
		Left:  mulExpr,
		Right: divExpr,
	}, true
}
