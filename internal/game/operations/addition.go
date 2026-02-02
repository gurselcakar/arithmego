package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Addition{})
}

// Addition implements the addition operation.
type Addition struct{}

func (a *Addition) Name() string            { return "Addition" }
func (a *Addition) Symbol() string          { return "+" }
func (a *Addition) Arity() game.Arity       { return game.Binary }
func (a *Addition) Category() game.Category { return game.CategoryBasic }

func (a *Addition) Apply(operands []int) int {
	return operands[0] + operands[1]
}

func (a *Addition) Format(operands []int) string {
	return fmt.Sprintf("%d + %d", operands[0], operands[1])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Digit count (+0.5/digit): More digits require tracking more information in working memory.
//   - Carries (+1.5/carry): Carries require holding intermediate results while computing the next
//     column, significantly increasing cognitive load. Weighted higher than digits.
//   - Nice numbers (-1.0): Multiples of 5/10/25 have familiar patterns that reduce mental effort.
//   - Boundary crossing (+0.5): Crossing 100/1000 requires adjusting mental magnitude estimates.
//   - Round result (-0.5): Round answers are easier to verify and often indicate simpler problems.
//
// Weights are initial estimates subject to tuning based on playtesting.
func (a *Addition) ScoreDifficulty(operands []int, answer int) float64 {
	op1, op2 := operands[0], operands[1]
	score := 1.0

	digits1 := countDigits(op1)
	digits2 := countDigits(op2)
	score += float64(digits1-1) * 0.5
	score += float64(digits2-1) * 0.5

	carries := countCarries(op1, op2)
	score += float64(carries) * 1.5

	if isNiceNumber(op1) && isNiceNumber(op2) {
		score -= 1.0
	}

	if crossesBoundary(op1, op2, answer) {
		score += 0.5
	}

	if isRoundNumber(answer) {
		score -= 0.5
	}

	return clampScore(score)
}

func (a *Addition) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(a, diff, a.makeCandidate, a.makeCandidateRelaxed)
}

// makeCandidate generates a candidate with standard operand ranges.
func (a *Addition) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var op1, op2 int
	switch diff {
	case game.Beginner:
		op1 = randomInRange(1, 9)
		op2 = randomInRange(1, 9)
	case game.Easy:
		op1 = randomInRange(10, 50)
		op2 = randomInRange(10, 50)
	case game.Medium:
		op1 = randomInRange(20, 200)
		op2 = randomInRange(20, 200)
	case game.Hard:
		op1 = randomInRange(100, 500)
		op2 = randomInRange(100, 500)
	case game.Expert:
		op1 = randomInRange(200, 999)
		op2 = randomInRange(200, 999)
	default:
		op1 = randomInRange(1, 9)
		op2 = randomInRange(1, 9)
	}

	operands := []int{op1, op2}
	return Candidate{Operands: operands, Answer: a.Apply(operands)}, true
}

// makeCandidateRelaxed generates a candidate with expanded operand ranges.
func (a *Addition) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var min1, max1, min2, max2 int
	switch diff {
	case game.Beginner:
		min1, max1 = 1, 20
		min2, max2 = 1, 20
	case game.Easy:
		min1, max1 = 5, 100
		min2, max2 = 5, 100
	case game.Medium:
		min1, max1 = 10, 300
		min2, max2 = 10, 300
	case game.Hard:
		min1, max1 = 50, 700
		min2, max2 = 50, 700
	case game.Expert:
		min1, max1 = 100, 999
		min2, max2 = 100, 999
	default:
		min1, max1 = 1, 20
		min2, max2 = 1, 20
	}

	op1 := randomInRange(min1, max1)
	op2 := randomInRange(min2, max2)
	operands := []int{op1, op2}
	return Candidate{Operands: operands, Answer: a.Apply(operands)}, true
}
