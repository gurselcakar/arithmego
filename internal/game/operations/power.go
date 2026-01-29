package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Power{})
}

// Power implements the power operation (a^b).
type Power struct{}

func (p *Power) Name() string           { return "Power" }
func (p *Power) Symbol() string         { return "^" }
func (p *Power) Arity() game.Arity      { return game.Binary }
func (p *Power) Category() game.Category { return game.CategoryAdvanced }

func (p *Power) Apply(operands []int) int {
	return intPow(operands[0], operands[1])
}

func (p *Power) Format(operands []int) string {
	return fmt.Sprintf("%d^%d", operands[0], operands[1])
}

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
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	for attempts := 0; attempts < 100; attempts++ {
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

		// Prevent overflow - limit result to reasonable size
		result := intPow(base, exp)
		if result > 1000000 {
			continue
		}

		operands := []int{base, exp}
		score := p.ScoreDifficulty(operands, result)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: p,
				Answer:    result,
				Display:   p.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: p,
				Answer:    result,
				Display:   p.Format(operands),
			}
		}
	}

	return bestQuestion
}
