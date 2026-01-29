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

func (c *CubeRoot) Name() string           { return "Cube Root" }
func (c *CubeRoot) Symbol() string         { return "∛" }
func (c *CubeRoot) Arity() game.Arity      { return game.Unary }
func (c *CubeRoot) Category() game.Category { return game.CategoryPower }

func (c *CubeRoot) Apply(operands []int) int {
	return int(math.Cbrt(float64(operands[0])))
}

func (c *CubeRoot) Format(operands []int) string {
	return fmt.Sprintf("∛%d", operands[0])
}

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
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	for attempts := 0; attempts < 100; attempts++ {
		// Generate perfect cube by picking the result first
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
		operands := []int{operand}
		score := c.ScoreDifficulty(operands, result)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: c,
				Answer:    result,
				Display:   c.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: c,
				Answer:    result,
				Display:   c.Format(operands),
			}
		}
	}

	return bestQuestion
}
