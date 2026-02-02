package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&CubeRoot{})
}

// CubeRoot implements the cube root operation (∛n).
type CubeRoot struct{}

func (c *CubeRoot) Name() string            { return "Cube Root" }
func (c *CubeRoot) Symbol() string          { return "∛" }
func (c *CubeRoot) Arity() game.Arity       { return game.Unary }
func (c *CubeRoot) Category() game.Category { return game.CategoryPower }

func (c *CubeRoot) Apply(operands []int) int {
	return int(math.Cbrt(float64(operands[0])))
}

func (c *CubeRoot) Format(operands []int) string {
	return fmt.Sprintf("∛%d", operands[0])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Answer magnitude (+1.0 to +5.5): Smaller cube roots (∛8=2, ∛27=3) are commonly known.
//     Larger roots require recognizing less familiar perfect cubes or estimation.
//   - Common cubes (-0.5): Roots like ∛8, ∛27, ∛64, ∛125, ∛1000 are frequently memorized.
//
// Weights are initial estimates subject to tuning based on playtesting.
func (c *CubeRoot) ScoreDifficulty(operands []int, answer int) float64 {
	score := 1.0

	// Common perfect cubes (1-5) are easier
	if answer <= 5 {
		score += 1.0
	} else if answer <= 10 {
		score += 2.5
	} else if answer <= 15 {
		score += 4.0
	} else {
		score += 5.5
	}

	// Very common cubes get a bonus (∛8, ∛27, ∛64, ∛125)
	commonCubes := map[int]bool{2: true, 3: true, 4: true, 5: true, 10: true}
	if commonCubes[answer] {
		score -= 0.5
	}

	return clampScore(score)
}

func (c *CubeRoot) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(c, diff, c.makeCandidate, c.makeCandidateRelaxed)
}

// makeCandidate generates a perfect cube by picking the result first.
func (c *CubeRoot) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var result int
	switch diff {
	case game.Beginner:
		result = randomInRange(2, 5)
	case game.Easy:
		result = randomInRange(3, 7)
	case game.Medium:
		result = randomInRange(5, 10)
	case game.Hard:
		result = randomInRange(7, 15)
	case game.Expert:
		result = randomInRange(10, 20)
	default:
		result = randomInRange(2, 5)
	}

	operand := result * result * result
	return Candidate{Operands: []int{operand}, Answer: result}, true
}

// makeCandidateRelaxed generates a candidate with expanded result ranges.
func (c *CubeRoot) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var min, max int
	switch diff {
	case game.Beginner:
		min, max = 2, 7
	case game.Easy:
		min, max = 2, 10
	case game.Medium:
		min, max = 3, 15
	case game.Hard:
		min, max = 5, 20
	case game.Expert:
		min, max = 7, 25
	default:
		min, max = 2, 7
	}

	result := randomInRange(min, max)
	operand := result * result * result
	return Candidate{Operands: []int{operand}, Answer: result}, true
}
