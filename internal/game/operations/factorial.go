package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Factorial{})
}

// Factorial implements the factorial operation (n!).
type Factorial struct{}

func (f *Factorial) Name() string            { return "Factorial" }
func (f *Factorial) Symbol() string          { return "!" }
func (f *Factorial) Arity() game.Arity       { return game.Unary }
func (f *Factorial) Category() game.Category { return game.CategoryAdvanced }

func (f *Factorial) Apply(operands []int) int {
	return factorial(operands[0])
}

func (f *Factorial) Format(operands []int) string {
	return fmt.Sprintf("%d!", operands[0])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Input value n (+0.5 to +9.5): Factorials grow extremely fast, making larger values
//     impractical for mental math. Common breakpoints:
//     • n ≤ 3: Trivial/memorized (1!=1, 2!=2, 3!=6)
//     • n = 4,5: Commonly known (4!=24, 5!=120)
//     • n = 6,7: Require computation (6!=720, 7!=5040)
//     • n ≥ 8: Challenging multi-step multiplication
//
// Weights are initial estimates subject to tuning based on playtesting.
func (f *Factorial) ScoreDifficulty(operands []int, answer int) float64 {
	n := operands[0]
	score := 1.0

	switch {
	case n <= 3:
		score += 0.5
	case n == 4:
		score += 1.5
	case n == 5:
		score += 2.5
	case n == 6:
		score += 4.0
	case n == 7:
		score += 5.5
	case n == 8:
		score += 7.0
	case n == 9:
		score += 8.5
	case n >= 10:
		score += 9.5
	}

	return clampScore(score)
}

func (f *Factorial) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(f, diff, f.makeCandidate, f.makeCandidateRelaxed)
}

// makeCandidate generates a candidate with standard operand ranges.
// Factorials are limited - max 10! = 3,628,800.
func (f *Factorial) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var n int
	switch diff {
	case game.Beginner:
		n = randomInRange(1, 4)
	case game.Easy:
		n = randomInRange(3, 5)
	case game.Medium:
		n = randomInRange(4, 6)
	case game.Hard:
		n = randomInRange(5, 8)
	case game.Expert:
		n = randomInRange(7, 10)
	default:
		n = randomInRange(1, 4)
	}

	operands := []int{n}
	return Candidate{Operands: operands, Answer: f.Apply(operands)}, true
}

// makeCandidateRelaxed generates a candidate with expanded operand ranges.
func (f *Factorial) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var min, max int
	switch diff {
	case game.Beginner:
		min, max = 1, 5
	case game.Easy:
		min, max = 2, 6
	case game.Medium:
		min, max = 3, 7
	case game.Hard:
		min, max = 4, 9
	case game.Expert:
		min, max = 5, 10
	default:
		min, max = 1, 5
	}

	n := randomInRange(min, max)
	operands := []int{n}
	return Candidate{Operands: operands, Answer: f.Apply(operands)}, true
}
