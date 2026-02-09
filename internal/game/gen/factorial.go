package gen

import (
	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

type FactorialGen struct{}

func (g *FactorialGen) Label() string { return "Factorial" }

func (g *FactorialGen) Generate(diff game.Difficulty) *game.Question {
	return TryGenerate(factorialPatterns, diff, g.Label(), 100)
}

var factorialPatterns = PatternSet{
	game.Beginner: {
		{factSingle, 10},
	},
	game.Easy: {
		{factSingle, 10},
	},
	game.Medium: {
		{factSingle, 6},
		{factDivision, 4},
	},
	game.Hard: {
		{factSingle, 4},
		{factDivision, 4},
		{factAddition, 2},
	},
	game.Expert: {
		{factDivision, 5},
		{factAddition, 3},
		{factSingle, 2},
	},
}

func factSingle(diff game.Difficulty) (expr.Expr, bool) {
	r := FactorialRanges[diff]
	n := RandomInRange(r.Min, r.Max)
	return &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: n}}, true
}

// factDivision: n! ÷ m! (simplifies to product of n*(n-1)*...*(m+1))
func factDivision(diff game.Difficulty) (expr.Expr, bool) {
	r := FactorialRanges[diff]
	n := RandomInRange(r.Min, r.Max)
	// m must be < n and gap should be reasonable (≤ 3 for Medium, larger for Expert)
	maxGap := 3
	if diff >= game.Hard {
		maxGap = 4
	}
	minM := n - maxGap
	if minM < 1 {
		minM = 1
	}
	m := RandomInRange(minM, n-1)
	if m < 1 || m >= n {
		return factSingle(diff)
	}

	// Verify result is reasonable (n!/m! = n*(n-1)*...*(m+1))
	result := 1
	for i := m + 1; i <= n; i++ {
		result *= i
		if result > MaxPowerResult {
			return nil, false
		}
	}

	return &expr.BinOp{
		Op:    expr.OpDiv,
		Left:  &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: n}},
		Right: &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: m}},
	}, true
}

// factAddition: n! + m!
func factAddition(diff game.Difficulty) (expr.Expr, bool) {
	r := FactorialRanges[diff]
	n := RandomInRange(r.Min, r.Max)
	m := RandomInRange(r.Min, r.Max)
	// Check both factorials are reasonable
	if Factorial(n) > MaxPowerResult || Factorial(m) > MaxPowerResult {
		return nil, false
	}
	return &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: n}},
		Right: &expr.UnarySuffix{Op: expr.OpFactorial, Operand: &expr.Num{Value: m}},
	}, true
}
