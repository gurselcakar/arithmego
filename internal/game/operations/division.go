package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Division{})
}

// Division implements the division operation.
type Division struct{}

func (d *Division) Name() string            { return "Division" }
func (d *Division) Symbol() string          { return "÷" }
func (d *Division) Arity() game.Arity       { return game.Binary }
func (d *Division) Category() game.Category { return game.CategoryBasic }

// Apply computes operands[0] / operands[1].
//
// Preconditions:
//   - operands must have exactly 2 elements (enforced by Arity)
//   - operands[1] must not be zero
//
// Panics if operands[1] is zero. The Generate method guarantees valid operands
// for game use. Direct callers must validate inputs.
func (d *Division) Apply(operands []int) int {
	if operands[1] == 0 {
		panic("division by zero")
	}
	return operands[0] / operands[1]
}

func (d *Division) Format(operands []int) string {
	return fmt.Sprintf("%d ÷ %d", operands[0], operands[1])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Times table inverse (no penalty): 56÷8 is recalled as "what times 8 is 56?"
//     If both divisor and quotient are ≤12, it's a memorized fact.
//   - Digit combinations: Similar to multiplication, more digits require more
//     mental long division steps. 2-digit÷2-digit (+3.5) requires trial division.
//   - Easy divisors (-0.5): Dividing by 2 is halving, by 5 is doubling then /10,
//     by 10 is removing a zero. These have mental shortcuts.
//   - Quotient >12 (+0.5): Beyond times tables requires actual computation.
//   - Quotient >20 (+0.5 more): Large quotients require estimation and iteration.
//
// Weights are initial estimates subject to tuning based on playtesting.
func (d *Division) ScoreDifficulty(operands []int, answer int) float64 {
	dividend, divisor := operands[0], operands[1]
	score := 1.0

	dividendDigits := countDigits(dividend)
	divisorDigits := countDigits(divisor)

	// Times table inverses are memorized facts
	if isTimesTableFact(divisor, answer) {
		score += 0.0
	} else {
		// Digit-based scoring for non-memorized division
		if divisorDigits == 1 && dividendDigits <= 2 {
			score += 1.5
		} else if divisorDigits == 1 && dividendDigits == 3 {
			score += 2.5
		} else if divisorDigits == 2 && dividendDigits == 2 {
			score += 3.5
		} else if divisorDigits == 2 && dividendDigits >= 3 {
			score += 4.5
		}
	}

	// Easy divisors have mental shortcuts
	if divisor == 2 || divisor == 5 || divisor == 10 {
		score -= 0.5
	}

	// Large quotients require computation beyond recall
	if answer > 12 {
		score += 0.5
	}
	if answer > 20 {
		score += 0.5
	}

	return clampScore(score)
}

func (d *Division) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(d, diff, d.makeCandidate, d.makeCandidateRelaxed)
}

// makeCandidate generates a candidate by working backwards: quotient × divisor = dividend.
func (d *Division) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var divisor, quotient int
	switch diff {
	case game.Beginner:
		divisor = randomInRange(2, 9)
		quotient = randomInRange(2, 9)
	case game.Easy:
		divisor = randomInRange(2, 12)
		quotient = randomInRange(2, 12)
	case game.Medium:
		divisor = randomInRange(3, 15)
		quotient = randomInRange(5, 20)
	case game.Hard:
		divisor = randomInRange(5, 20)
		quotient = randomInRange(10, 30)
	case game.Expert:
		divisor = randomInRange(10, 30)
		quotient = randomInRange(15, 50)
	default:
		divisor = randomInRange(2, 9)
		quotient = randomInRange(2, 9)
	}

	dividend := divisor * quotient
	return Candidate{Operands: []int{dividend, divisor}, Answer: quotient}, true
}

// makeCandidateRelaxed generates a candidate with expanded ranges.
func (d *Division) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var minDiv, maxDiv, minQuot, maxQuot int
	switch diff {
	case game.Beginner:
		minDiv, maxDiv = 2, 12
		minQuot, maxQuot = 2, 12
	case game.Easy:
		minDiv, maxDiv = 2, 15
		minQuot, maxQuot = 2, 15
	case game.Medium:
		minDiv, maxDiv = 2, 20
		minQuot, maxQuot = 3, 30
	case game.Hard:
		minDiv, maxDiv = 3, 30
		minQuot, maxQuot = 5, 50
	case game.Expert:
		minDiv, maxDiv = 5, 40
		minQuot, maxQuot = 10, 75
	default:
		minDiv, maxDiv = 2, 12
		minQuot, maxQuot = 2, 12
	}

	divisor := randomInRange(minDiv, maxDiv)
	quotient := randomInRange(minQuot, maxQuot)
	dividend := divisor * quotient
	return Candidate{Operands: []int{dividend, divisor}, Answer: quotient}, true
}
