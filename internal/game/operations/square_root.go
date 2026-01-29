package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&SquareRoot{})
}

// SquareRoot implements the square root operation (√n).
type SquareRoot struct{}

func (s *SquareRoot) Name() string           { return "Square Root" }
func (s *SquareRoot) Symbol() string         { return "√" }
func (s *SquareRoot) Arity() game.Arity      { return game.Unary }
func (s *SquareRoot) Category() game.Category { return game.CategoryPower }

func (s *SquareRoot) Apply(operands []int) int {
	return int(math.Sqrt(float64(operands[0])))
}

func (s *SquareRoot) Format(operands []int) string {
	return fmt.Sprintf("√%d", operands[0])
}

func (s *SquareRoot) ScoreDifficulty(operands []int, answer int) float64 {
	score := 1.0

	// Common perfect squares (1-12) are easier
	if answer <= 12 {
		score += 0.5
	} else if answer <= 20 {
		score += 1.5
	} else if answer <= 30 {
		score += 2.5
	} else {
		score += 3.5
	}

	// Very common squares get a bonus (√4, √9, √16, √25, √36, √49, √64, √81, √100)
	commonSquares := map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true}
	if commonSquares[answer] {
		score -= 0.5
	}

	return clampScore(score)
}

func (s *SquareRoot) Generate(diff game.Difficulty) game.Question {
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	for attempts := 0; attempts < 100; attempts++ {
		// Generate perfect square by picking the result first
		var result int
		switch diff {
		case game.Beginner:
			result = randomInRange(2, 10)
		case game.Easy:
			result = randomInRange(5, 15)
		case game.Medium:
			result = randomInRange(10, 25)
		case game.Hard:
			result = randomInRange(15, 35)
		case game.Expert:
			result = randomInRange(25, 50)
		default:
			result = randomInRange(2, 10)
		}

		operand := result * result
		operands := []int{operand}
		score := s.ScoreDifficulty(operands, result)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: s,
				Answer:    result,
				Display:   s.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: s,
				Answer:    result,
				Display:   s.Format(operands),
			}
		}
	}

	return bestQuestion
}
