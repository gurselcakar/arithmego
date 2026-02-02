package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Multiplication{})
}

// Multiplication implements the multiplication operation.
type Multiplication struct{}

func (m *Multiplication) Name() string            { return "Multiplication" }
func (m *Multiplication) Symbol() string          { return "×" }
func (m *Multiplication) Arity() game.Arity       { return game.Binary }
func (m *Multiplication) Category() game.Category { return game.CategoryBasic }

func (m *Multiplication) Apply(operands []int) int {
	return operands[0] * operands[1]
}

func (m *Multiplication) Format(operands []int) string {
	return fmt.Sprintf("%d × %d", operands[0], operands[1])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Digit combinations: Single×single uses memorized times tables (no penalty).
//     Each additional digit dramatically increases mental computation steps.
//     Double×double (+4.5) is harder than single×triple (+3.5) because both operands
//     require decomposition.
//   - Easy multipliers (-0.5): ×10 is just appending zero, ×5 is half of ×10,
//     ×2 is simple doubling. These shortcuts reduce cognitive load.
//   - ×11 pattern (-0.5): 11×n has a memorable pattern (e.g., 11×12=132).
//   - Both odd (+0.5): Even numbers allow halving tricks; odd×odd has no shortcuts.
//   - Squares (-0.3): Small squares (up to 12²) are often memorized.
//
// Weights are initial estimates subject to tuning based on playtesting.
func (m *Multiplication) ScoreDifficulty(operands []int, answer int) float64 {
	op1, op2 := operands[0], operands[1]
	score := 1.0

	digits1 := countDigits(op1)
	digits2 := countDigits(op2)

	// Digit combination scoring based on mental computation complexity
	if digits1 == 1 && digits2 == 1 {
		score += 0.0 // Times table: memorized
	} else if (digits1 == 1 && digits2 == 2) || (digits1 == 2 && digits2 == 1) {
		score += 2.0 // Single × double: one decomposition
	} else if digits1 == 2 && digits2 == 2 {
		score += 4.5 // Double × double: both operands need decomposition
	} else if (digits1 == 1 && digits2 == 3) || (digits1 == 3 && digits2 == 1) {
		score += 3.5 // Single × triple: extended but one operand simple
	} else {
		score += float64(digits1+digits2) * 1.0
	}

	// Easy multipliers have mental shortcuts
	if op1 == 10 || op2 == 10 {
		score -= 0.5
	}
	if op1 == 5 || op2 == 5 {
		score -= 0.5
	}
	if op1 == 2 || op2 == 2 {
		score -= 0.5
	}

	if op1 == 11 || op2 == 11 {
		score -= 0.5
	}

	// Odd × odd has no halving shortcuts
	if op1%2 == 1 && op2%2 == 1 {
		score += 0.5
	}

	// Small squares are often memorized
	if op1 == op2 && op1 <= 12 {
		score -= 0.3
	}

	return clampScore(score)
}

func (m *Multiplication) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(m, diff, m.makeCandidate, m.makeCandidateRelaxed)
}

// makeCandidate generates a candidate with standard operand ranges.
func (m *Multiplication) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var op1, op2 int
	switch diff {
	case game.Beginner:
		op1 = randomInRange(2, 9)
		op2 = randomInRange(2, 9)
	case game.Easy:
		op1 = randomInRange(2, 12)
		op2 = randomInRange(10, 20)
	case game.Medium:
		op1 = randomInRange(5, 15)
		op2 = randomInRange(10, 30)
	case game.Hard:
		op1 = randomInRange(10, 30)
		op2 = randomInRange(10, 50)
	case game.Expert:
		op1 = randomInRange(15, 50)
		op2 = randomInRange(20, 99)
	default:
		op1 = randomInRange(2, 9)
		op2 = randomInRange(2, 9)
	}

	operands := []int{op1, op2}
	return Candidate{Operands: operands, Answer: m.Apply(operands)}, true
}

// makeCandidateRelaxed generates a candidate with expanded operand ranges.
func (m *Multiplication) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var min1, max1, min2, max2 int
	switch diff {
	case game.Beginner:
		min1, max1 = 2, 12
		min2, max2 = 2, 12
	case game.Easy:
		min1, max1 = 2, 15
		min2, max2 = 5, 30
	case game.Medium:
		min1, max1 = 3, 25
		min2, max2 = 5, 50
	case game.Hard:
		min1, max1 = 5, 50
		min2, max2 = 5, 75
	case game.Expert:
		min1, max1 = 10, 75
		min2, max2 = 10, 99
	default:
		min1, max1 = 2, 12
		min2, max2 = 2, 12
	}

	op1 := randomInRange(min1, max1)
	op2 := randomInRange(min2, max2)
	operands := []int{op1, op2}
	return Candidate{Operands: operands, Answer: m.Apply(operands)}, true
}
