package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Percentage{})
}

// Percentage implements the percentage operation (a% of b).
type Percentage struct{}

func (p *Percentage) Name() string            { return "Percentage" }
func (p *Percentage) Symbol() string          { return "% of" }
func (p *Percentage) Arity() game.Arity       { return game.Binary }
func (p *Percentage) Category() game.Category { return game.CategoryAdvanced }

func (p *Percentage) Apply(operands []int) int {
	// operands[0]% of operands[1]
	return (operands[0] * operands[1]) / 100
}

func (p *Percentage) Format(operands []int) string {
	return fmt.Sprintf("%d%% of %d", operands[0], operands[1])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Percentage type (+0.5 to +4.0): Common percentages (50%, 25%, 10%, 20%, 100%) have
//     mental shortcuts (halving, quartering, moving decimal). Multiples of 10% and 5% are
//     progressively harder. Odd percentages require full multiplication.
//   - Value digits (+0.5/digit): Larger values require more mental arithmetic steps.
//   - Round values (-0.3/-0.5): Multiples of 10 or 100 simplify multiplication.
//   - Small values (-0.3): Values ≤100 are easier to work with mentally.
//
// Weights are initial estimates subject to tuning based on playtesting.
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

// Percent pools by difficulty tier
var (
	percentEasy   = []int{10, 20, 25, 50, 100}
	percentMedium = []int{5, 10, 15, 20, 25, 30, 40, 50, 75}
	percentHard   = []int{5, 10, 12, 15, 20, 25, 30, 35, 40, 45, 50, 60, 75, 80}
	// Expert uses "odd" percentages that aren't multiples of 5 for higher difficulty
	percentExpert = []int{2, 3, 4, 6, 7, 8, 9, 11, 12, 13, 14, 16, 17, 18, 19,
		21, 22, 23, 24, 26, 27, 28, 29, 32, 33, 34, 36, 37, 38, 39}
	// Combined pool for relaxed generation
	percentAll = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 45, 50, 60, 75, 80}
)

func (p *Percentage) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(p, diff, p.makeCandidate, p.makeCandidateRelaxed)
}

// makeCandidate generates a candidate with standard percent pools.
// Returns invalid if the result wouldn't be a clean integer.
func (p *Percentage) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var percent, value int
	switch diff {
	case game.Beginner:
		percent = percentEasy[randomInRange(0, len(percentEasy)-1)]
		value = randomInRange(2, 20) * (100 / gcd(percent, 100))
	case game.Easy:
		percent = percentEasy[randomInRange(0, len(percentEasy)-1)]
		value = randomInRange(10, 100)
		value = alignToCleanDivision(value, percent, 100)
	case game.Medium:
		percent = percentMedium[randomInRange(0, len(percentMedium)-1)]
		value = randomInRange(20, 200)
		value = alignToCleanDivision(value, percent, 100)
	case game.Hard:
		percent = percentHard[randomInRange(0, len(percentHard)-1)]
		value = randomInRange(50, 500)
		value = alignToCleanDivision(value, percent, 200)
	case game.Expert:
		percent = percentExpert[randomInRange(0, len(percentExpert)-1)]
		value = randomInRange(100, 1000)
		value = alignToCleanDivision(value, percent, 100)
	default:
		percent = percentEasy[randomInRange(0, len(percentEasy)-1)]
		value = randomInRange(2, 20) * (100 / gcd(percent, 100))
	}

	// Ensure clean integer result
	if (percent*value)%100 != 0 {
		return Candidate{}, false
	}

	operands := []int{percent, value}
	return Candidate{Operands: operands, Answer: p.Apply(operands)}, true
}

// makeCandidateRelaxed generates a candidate with a combined percent pool.
func (p *Percentage) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var minVal, maxVal, fallback int
	switch diff {
	case game.Beginner:
		minVal, maxVal, fallback = 20, 100, 100
	case game.Easy:
		minVal, maxVal, fallback = 50, 200, 100
	case game.Medium:
		minVal, maxVal, fallback = 100, 400, 200
	case game.Hard:
		minVal, maxVal, fallback = 200, 800, 400
	case game.Expert:
		minVal, maxVal, fallback = 200, 1000, 500
	default:
		minVal, maxVal, fallback = 20, 100, 100
	}

	percent := percentAll[randomInRange(0, len(percentAll)-1)]
	value := randomInRange(minVal, maxVal)
	value = alignToCleanDivision(value, percent, fallback)

	if (percent*value)%100 != 0 {
		return Candidate{}, false
	}

	operands := []int{percent, value}
	return Candidate{Operands: operands, Answer: p.Apply(operands)}, true
}

// alignToCleanDivision adjusts value so that (percent * value) % 100 == 0.
// Returns the adjusted value, or fallback if adjustment yields zero.
func alignToCleanDivision(value, percent, fallback int) int {
	divisor := 100 / gcd(percent, 100)
	aligned := (value / divisor) * divisor
	if aligned == 0 {
		return fallback
	}
	return aligned
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
