package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Power{})
}

// Power implements the power operation (a^b).
type Power struct{}

func (p *Power) Name() string            { return "Power" }
func (p *Power) Symbol() string          { return "^" }
func (p *Power) Arity() game.Arity       { return game.Binary }
func (p *Power) Category() game.Category { return game.CategoryAdvanced }

func (p *Power) Apply(operands []int) int {
	return intPow(operands[0], operands[1])
}

func (p *Power) Format(operands []int) string {
	return fmt.Sprintf("%d^%d", operands[0], operands[1])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Exponent value: Determines the number of multiplications needed.
//     • exp=2: Squares, often memorized up to 12². Larger bases need computation.
//     • exp=3: Cubes, commonly known up to 5³. Requires chaining two multiplications.
//     • exp=4+: Requires multiple chained multiplications with growing intermediate values.
//   - Easy bases (-0.5/-1.0): Base 2 has familiar powers (2,4,8,16,32...). Base 10 is trivial
//     (just append zeros).
//   - Base of 1: Trivial case, always equals 1 regardless of exponent.
//
// Weights are initial estimates subject to tuning based on playtesting.
func (p *Power) ScoreDifficulty(operands []int, answer int) float64 {
	base, exp := operands[0], operands[1]
	score := 1.0

	// Exponent-based scoring
	if exp == 2 {
		// Squares - often memorized
		if base <= 12 {
			score += 0.5
		} else if base <= 20 {
			score += 2.0
		} else {
			score += 3.5
		}
	} else if exp == 3 {
		// Cubes - harder
		if base <= 5 {
			score += 1.5
		} else if base <= 10 {
			score += 3.5
		} else {
			score += 5.5
		}
	} else if exp == 4 {
		if base <= 5 {
			score += 3.0
		} else {
			score += 5.5
		}
	} else {
		// Higher exponents
		score += float64(exp) * 1.5
	}

	// Easy bases (2, 10) are easier
	if base == 2 {
		score -= 0.5
	}
	if base == 10 {
		score -= 1.0
	}

	// Base of 1 is trivial
	if base == 1 {
		score = 1.0
	}

	return clampScore(score)
}

func (p *Power) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(p, diff, p.makeCandidate, p.makeCandidateRelaxed)
}

// maxPowerResult is the maximum allowed result for power operations.
// This prevents integer overflow and keeps answers reasonable for mental math.
const maxPowerResult = 1000000

// makeCandidate generates a candidate with standard operand ranges.
// Returns invalid if the result would exceed maxPowerResult.
func (p *Power) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var base, exp int
	switch diff {
	case game.Beginner:
		base = randomInRange(2, 10)
		exp = 2
	case game.Easy:
		base = randomInRange(2, 12)
		exp = randomInRange(2, 3)
	case game.Medium:
		base = randomInRange(2, 10)
		exp = randomInRange(2, 4)
	case game.Hard:
		base = randomInRange(2, 8)
		exp = randomInRange(3, 5)
	case game.Expert:
		base = randomInRange(2, 6)
		exp = randomInRange(4, 6)
	default:
		base = randomInRange(2, 10)
		exp = 2
	}

	// Check for overflow before computing
	if wouldOverflow(base, exp, maxPowerResult) {
		return Candidate{}, false
	}

	result := intPow(base, exp)
	return Candidate{Operands: []int{base, exp}, Answer: result}, true
}

// makeCandidateRelaxed generates a candidate with expanded operand ranges.
func (p *Power) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var minBase, maxBase, minExp, maxExp int
	switch diff {
	case game.Beginner:
		minBase, maxBase = 2, 12
		minExp, maxExp = 2, 2
	case game.Easy:
		minBase, maxBase = 2, 15
		minExp, maxExp = 2, 3
	case game.Medium:
		minBase, maxBase = 2, 12
		minExp, maxExp = 2, 4
	case game.Hard:
		minBase, maxBase = 2, 10
		minExp, maxExp = 2, 5
	case game.Expert:
		minBase, maxBase = 2, 8
		minExp, maxExp = 3, 7
	default:
		minBase, maxBase = 2, 12
		minExp, maxExp = 2, 2
	}

	base := randomInRange(minBase, maxBase)
	exp := randomInRange(minExp, maxExp)

	// Check for overflow before computing
	if wouldOverflow(base, exp, maxPowerResult) {
		return Candidate{}, false
	}

	result := intPow(base, exp)
	return Candidate{Operands: []int{base, exp}, Answer: result}, true
}

// wouldOverflow returns true if base^exp would exceed maxResult.
// Uses iterative multiplication with early exit to avoid actual overflow.
func wouldOverflow(base, exp, maxResult int) bool {
	if base <= 1 {
		return false
	}
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
		if result > maxResult {
			return true
		}
	}
	return false
}
