package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Cube{})
}

// Cube implements the cube operation (n³).
type Cube struct{}

func (c *Cube) Name() string           { return "Cube" }
func (c *Cube) Symbol() string         { return "³" }
func (c *Cube) Arity() game.Arity      { return game.Unary }
func (c *Cube) Category() game.Category { return game.CategoryPower }

func (c *Cube) Apply(operands []int) int {
	n := operands[0]
	return n * n * n
}

func (c *Cube) Format(operands []int) string {
	return fmt.Sprintf("%d³", operands[0])
}

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
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	for attempts := 0; attempts < 100; attempts++ {
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
		answer := c.Apply(operands)
		score := c.ScoreDifficulty(operands, answer)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: c,
				Answer:    answer,
				Display:   c.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: c,
				Answer:    answer,
				Display:   c.Format(operands),
			}
		}
	}

	return bestQuestion
}
