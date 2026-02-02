package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Square{})
}

// Square implements the square operation (n²).
type Square struct{}

func (s *Square) Name() string            { return "Square" }
func (s *Square) Symbol() string          { return "²" }
func (s *Square) Arity() game.Arity       { return game.Unary }
func (s *Square) Category() game.Category { return game.CategoryPower }

func (s *Square) Apply(operands []int) int {
	return operands[0] * operands[0]
}

func (s *Square) Format(operands []int) string {
	return fmt.Sprintf("%d²", operands[0])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Base number range (+0.5 to +5.5): Squares up to 12² are typically memorized from
//     multiplication tables. Larger bases require mental multiplication strategies.
//   - Common squares (-0.5): Values like 1², 2², 3², 4², 5², 10² are instantly recalled.
//   - Round numbers (-0.5): Multiples of 10 have simple patterns (20² = 400).
//   - Numbers ending in 5 (-0.3): Have a known shortcut (25² = 625: compute 2×3=6, append 25).
//
// Weights are initial estimates subject to tuning based on playtesting.
func (s *Square) ScoreDifficulty(operands []int, answer int) float64 {
	n := operands[0]
	score := 1.0

	// Common memorized squares (1-12) are easier
	if n <= 12 {
		score += 0.5
	} else if n <= 15 {
		// 13², 14², 15² - somewhat common
		score += 1.5
	} else if n <= 20 {
		// 16-20 - need computation
		score += 2.5
	} else if n <= 25 {
		score += 3.5
	} else if n <= 30 {
		score += 4.5
	} else {
		score += 5.5
	}

	// Very common squares get a bonus
	commonSquares := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 10: true}
	if commonSquares[n] {
		score -= 0.5
	}

	// Round numbers are easier (10, 20, 30, etc.)
	if n%10 == 0 {
		score -= 0.5
	}

	// Numbers ending in 5 have a pattern (25² = 625, just 2×3 and append 25)
	if n%10 == 5 && n > 5 {
		score -= 0.3
	}

	return clampScore(score)
}

func (s *Square) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(s, diff, s.makeCandidate, s.makeCandidateRelaxed)
}

// makeCandidate generates a candidate with standard operand ranges.
func (s *Square) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var n int
	switch diff {
	case game.Beginner:
		n = randomInRange(2, 10)
	case game.Easy:
		n = randomInRange(5, 15)
	case game.Medium:
		n = randomInRange(10, 20)
	case game.Hard:
		n = randomInRange(15, 30)
	case game.Expert:
		n = randomInRange(20, 50)
	default:
		n = randomInRange(2, 10)
	}

	operands := []int{n}
	return Candidate{Operands: operands, Answer: s.Apply(operands)}, true
}

// makeCandidateRelaxed generates a candidate with expanded operand ranges.
func (s *Square) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var min, max int
	switch diff {
	case game.Beginner:
		min, max = 2, 15
	case game.Easy:
		min, max = 3, 20
	case game.Medium:
		min, max = 5, 30
	case game.Hard:
		min, max = 10, 45
	case game.Expert:
		min, max = 15, 60
	default:
		min, max = 2, 15
	}

	n := randomInRange(min, max)
	operands := []int{n}
	return Candidate{Operands: operands, Answer: s.Apply(operands)}, true
}
