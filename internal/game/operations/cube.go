package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Cube{})
}

// Cube implements the cube operation (n³).
type Cube struct{}

func (c *Cube) Name() string            { return "Cube" }
func (c *Cube) Symbol() string          { return "³" }
func (c *Cube) Arity() game.Arity       { return game.Unary }
func (c *Cube) Category() game.Category { return game.CategoryPower }

func (c *Cube) Apply(operands []int) int {
	n := operands[0]
	return n * n * n
}

func (c *Cube) Format(operands []int) string {
	return fmt.Sprintf("%d³", operands[0])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Base number range (+1.0 to +7.0): Cubes of small numbers (1-5) are commonly memorized.
//     Larger bases require multi-step multiplication (n × n × n), increasing cognitive load.
//   - Common cubes (-1.0): Values like 1³, 2³, 3³, 10³ are frequently encountered and memorized.
//   - Round numbers (-0.5): Multiples of 10 have predictable patterns (10³ = 1000).
//
// Weights are initial estimates subject to tuning based on playtesting.
func (c *Cube) ScoreDifficulty(operands []int, answer int) float64 {
	n := operands[0]
	score := 1.0

	// Common memorized cubes (1-5) are easier
	if n <= 5 {
		score += 1.0
	} else if n <= 10 {
		// 6³ to 10³ - need computation
		score += 3.0
	} else if n <= 15 {
		score += 5.0
	} else {
		score += 7.0
	}

	// Very common cubes get a bonus
	commonCubes := map[int]bool{1: true, 2: true, 3: true, 10: true}
	if commonCubes[n] {
		score -= 1.0
	}

	// Round numbers are easier
	if n%10 == 0 {
		score -= 0.5
	}

	return clampScore(score)
}

func (c *Cube) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(c, diff, c.makeCandidate, c.makeCandidateRelaxed)
}

// makeCandidate generates a candidate with standard operand ranges.
func (c *Cube) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var n int
	switch diff {
	case game.Beginner:
		n = randomInRange(2, 5)
	case game.Easy:
		n = randomInRange(2, 7)
	case game.Medium:
		n = randomInRange(4, 10)
	case game.Hard:
		n = randomInRange(6, 12)
	case game.Expert:
		n = randomInRange(8, 15)
	default:
		n = randomInRange(2, 5)
	}

	operands := []int{n}
	return Candidate{Operands: operands, Answer: c.Apply(operands)}, true
}

// makeCandidateRelaxed generates a candidate with expanded operand ranges.
func (c *Cube) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var min, max int
	switch diff {
	case game.Beginner:
		min, max = 2, 7
	case game.Easy:
		min, max = 2, 10
	case game.Medium:
		min, max = 3, 12
	case game.Hard:
		min, max = 4, 15
	case game.Expert:
		min, max = 5, 20
	default:
		min, max = 2, 7
	}

	n := randomInRange(min, max)
	operands := []int{n}
	return Candidate{Operands: operands, Answer: c.Apply(operands)}, true
}
