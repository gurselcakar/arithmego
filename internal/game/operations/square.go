package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Square{})
}

// Square implements the square operation (n²).
type Square struct{}

func (s *Square) Name() string           { return "Square" }
func (s *Square) Symbol() string         { return "²" }
func (s *Square) Arity() game.Arity      { return game.Unary }
func (s *Square) Category() game.Category { return game.CategoryPower }

func (s *Square) Apply(operands []int) int {
	return operands[0] * operands[0]
}

func (s *Square) Format(operands []int) string {
	return fmt.Sprintf("%d²", operands[0])
}

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
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	for attempts := 0; attempts < 100; attempts++ {
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
		answer := s.Apply(operands)
		score := s.ScoreDifficulty(operands, answer)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: s,
				Answer:    answer,
				Display:   s.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: s,
				Answer:    answer,
				Display:   s.Format(operands),
			}
		}
	}

	return bestQuestion
}
