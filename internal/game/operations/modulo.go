package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Modulo{})
}

// Modulo implements the modulo operation (a mod b).
type Modulo struct{}

func (m *Modulo) Name() string           { return "Modulo" }
func (m *Modulo) Symbol() string         { return "mod" }
func (m *Modulo) Arity() game.Arity      { return game.Binary }
func (m *Modulo) Category() game.Category { return game.CategoryAdvanced }

func (m *Modulo) Apply(operands []int) int {
	if operands[1] == 0 {
		panic("modulo by zero")
	}
	return operands[0] % operands[1]
}

func (m *Modulo) Format(operands []int) string {
	return fmt.Sprintf("%d mod %d", operands[0], operands[1])
}

func (m *Modulo) ScoreDifficulty(operands []int, answer int) float64 {
	dividend, divisor := operands[0], operands[1]
	score := 1.0

	dividendDigits := countDigits(dividend)
	divisorDigits := countDigits(divisor)

	// Digit-based scoring
	if dividendDigits == 1 && divisorDigits == 1 {
		score += 0.5
	} else if dividendDigits == 2 && divisorDigits == 1 {
		score += 1.5
	} else if dividendDigits == 2 && divisorDigits == 2 {
		score += 2.5
	} else if dividendDigits == 3 && divisorDigits == 1 {
		score += 2.0
	} else if dividendDigits == 3 && divisorDigits == 2 {
		score += 3.5
	} else {
		score += float64(dividendDigits+divisorDigits) * 0.8
	}

	// Easy divisors (2, 5, 10) are easier
	if divisor == 2 || divisor == 5 || divisor == 10 {
		score -= 0.5
	}

	// Small quotient means easier mental division
	quotient := dividend / divisor
	if quotient <= 5 {
		score -= 0.3
	} else if quotient > 10 {
		score += 0.5
	}

	// Zero remainder is easier (just confirm division is clean)
	if answer == 0 {
		score -= 0.3
	}

	return clampScore(score)
}

func (m *Modulo) Generate(diff game.Difficulty) game.Question {
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	for attempts := 0; attempts < 100; attempts++ {
		var dividend, divisor int
		switch diff {
		case game.Beginner:
			divisor = randomInRange(2, 9)
			dividend = randomInRange(divisor+1, divisor*5)
		case game.Easy:
			divisor = randomInRange(2, 12)
			dividend = randomInRange(divisor+1, 50)
		case game.Medium:
			divisor = randomInRange(3, 15)
			dividend = randomInRange(20, 100)
		case game.Hard:
			divisor = randomInRange(5, 25)
			dividend = randomInRange(50, 200)
		case game.Expert:
			divisor = randomInRange(10, 50)
			dividend = randomInRange(100, 500)
		default:
			divisor = randomInRange(2, 9)
			dividend = randomInRange(divisor+1, divisor*5)
		}

		operands := []int{dividend, divisor}
		answer := m.Apply(operands)
		score := m.ScoreDifficulty(operands, answer)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: m,
				Answer:    answer,
				Display:   m.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: m,
				Answer:    answer,
				Display:   m.Format(operands),
			}
		}
	}

	return bestQuestion
}
