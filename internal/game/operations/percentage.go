package operations

import (
	"fmt"
	"math"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Percentage{})
}

// Percentage implements the percentage operation (a% of b).
type Percentage struct{}

func (p *Percentage) Name() string           { return "Percentage" }
func (p *Percentage) Symbol() string         { return "% of" }
func (p *Percentage) Arity() game.Arity      { return game.Binary }
func (p *Percentage) Category() game.Category { return game.CategoryAdvanced }

func (p *Percentage) Apply(operands []int) int {
	// operands[0]% of operands[1]
	return (operands[0] * operands[1]) / 100
}

func (p *Percentage) Format(operands []int) string {
	return fmt.Sprintf("%d%% of %d", operands[0], operands[1])
}

func (p *Percentage) ScoreDifficulty(operands []int, answer int) float64 {
	percent, value := operands[0], operands[1]
	score := 1.0

	// Easy percentages (50%, 25%, 10%, 20%)
	easyPercents := map[int]bool{50: true, 25: true, 10: true, 20: true, 100: true}
	if easyPercents[percent] {
		score += 0.5
	} else if percent%10 == 0 {
		// Multiples of 10% are easier
		score += 1.5
	} else if percent%5 == 0 {
		// Multiples of 5%
		score += 2.5
	} else {
		// Odd percentages
		score += 4.0
	}

	// Value complexity
	valueDigits := countDigits(value)
	score += float64(valueDigits-1) * 0.5

	// Round values are easier
	if value%10 == 0 {
		score -= 0.3
	}
	if value%100 == 0 {
		score -= 0.5
	}

	// Small values are easier
	if value <= 100 {
		score -= 0.3
	}

	return clampScore(score)
}

func (p *Percentage) Generate(diff game.Difficulty) game.Question {
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	// Define clean percentages that divide evenly
	easyPercents := []int{10, 20, 25, 50, 100}
	mediumPercents := []int{5, 10, 15, 20, 25, 30, 40, 50, 75}
	hardPercents := []int{5, 10, 12, 15, 20, 25, 30, 35, 40, 45, 50, 60, 75, 80}

	for attempts := 0; attempts < 100; attempts++ {
		var percent, value int
		switch diff {
		case game.Beginner:
			percent = easyPercents[randomInRange(0, len(easyPercents)-1)]
			value = randomInRange(2, 20) * (100 / gcd(percent, 100))
		case game.Easy:
			percent = easyPercents[randomInRange(0, len(easyPercents)-1)]
			value = randomInRange(10, 100)
			// Ensure clean division
			value = (value / (100 / gcd(percent, 100))) * (100 / gcd(percent, 100))
			if value == 0 {
				value = 100
			}
		case game.Medium:
			percent = mediumPercents[randomInRange(0, len(mediumPercents)-1)]
			value = randomInRange(20, 200)
			value = (value / (100 / gcd(percent, 100))) * (100 / gcd(percent, 100))
			if value == 0 {
				value = 100
			}
		case game.Hard:
			percent = hardPercents[randomInRange(0, len(hardPercents)-1)]
			value = randomInRange(50, 500)
			value = (value / (100 / gcd(percent, 100))) * (100 / gcd(percent, 100))
			if value == 0 {
				value = 200
			}
		case game.Expert:
			percent = randomInRange(1, 99)
			value = randomInRange(100, 1000)
			// For expert, we allow non-clean division but ensure integer result
			// Adjust value to make percent * value divisible by 100
			remainder := (percent * value) % 100
			if remainder != 0 {
				value += (100 - remainder) / percent
				if (percent*value)%100 != 0 {
					value = (value / 100) * 100
					if value == 0 {
						value = 100
					}
				}
			}
		default:
			percent = easyPercents[randomInRange(0, len(easyPercents)-1)]
			value = randomInRange(2, 20) * (100 / gcd(percent, 100))
		}

		// Ensure we get a clean integer result
		if (percent*value)%100 != 0 {
			continue
		}

		operands := []int{percent, value}
		answer := p.Apply(operands)
		score := p.ScoreDifficulty(operands, answer)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  operands,
				Operation: p,
				Answer:    answer,
				Display:   p.Format(operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  operands,
				Operation: p,
				Answer:    answer,
				Display:   p.Format(operands),
			}
		}
	}

	return bestQuestion
}

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
