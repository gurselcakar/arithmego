package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Factorial{})
}

// Factorial implements the factorial operation (n!).
type Factorial struct{}

func (f *Factorial) Name() string           { return "Factorial" }
func (f *Factorial) Symbol() string         { return "!" }
func (f *Factorial) Arity() game.Arity      { return game.Unary }
func (f *Factorial) Category() game.Category { return game.CategoryAdvanced }

func (f *Factorial) Apply(operands []int) int {
	return factorial(operands[0])
}

func (f *Factorial) Format(operands []int) string {
	return fmt.Sprintf("%d!", operands[0])
}

func (f *Factorial) ScoreDifficulty(operands []int, answer int) float64 {
	n := operands[0]
	score := 1.0

	// Factorial difficulty by n
	// 1! = 1, 2! = 2, 3! = 6 - trivial/memorized
	// 4! = 24, 5! = 120 - commonly known
	// 6! = 720, 7! = 5040 - need computation
	// 8! and above - challenging

	switch {
	case n <= 3:
		score += 0.5
	case n == 4:
		score += 1.5
	case n == 5:
		score += 2.5
	case n == 6:
		score += 4.0
	case n == 7:
		score += 5.5
	case n == 8:
		score += 7.0
	case n == 9:
		score += 8.5
	case n >= 10:
		score += 9.5
	}

	return clampScore(score)
}

func (f *Factorial) Generate(diff game.Difficulty) game.Question {
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	// Factorials are limited - max 10! = 3,628,800
	for attempts := 0; attempts < 100; attempts++ {
		var n int
		switch diff {
		case game.Beginner:
			n = randomInRange(1, 4)
		case game.Easy:
			n = randomInRange(3, 5)
		case game.Medium:
			n = randomInRange(4, 6)
		case game.Hard:
			n = randomInRange(5, 8)
		case game.Expert:
			n = randomInRange(7, 10)
		default:
			n = randomInRange(1, 4)
		}

		operands := []int{n}
		answer := f.Apply(operands)
		score := f.ScoreDifficulty(operands, answer)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: f,
				Answer:    answer,
				Display:   f.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: f,
				Answer:    answer,
				Display:   f.Format(operands),
			}
		}
	}

	return bestQuestion
}
