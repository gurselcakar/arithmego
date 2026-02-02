package operations

import (
	"fmt"

	"github.com/gurselcakar/arithmego/internal/game"
)

func init() {
	Register(&Subtraction{})
}

// Subtraction implements the subtraction operation.
type Subtraction struct{}

func (s *Subtraction) Name() string            { return "Subtraction" }
func (s *Subtraction) Symbol() string          { return "−" }
func (s *Subtraction) Arity() game.Arity       { return game.Binary }
func (s *Subtraction) Category() game.Category { return game.CategoryBasic }

func (s *Subtraction) Apply(operands []int) int {
	return operands[0] - operands[1]
}

func (s *Subtraction) Format(operands []int) string {
	return fmt.Sprintf("%d − %d", operands[0], operands[1])
}

// ScoreDifficulty computes a difficulty score based on cognitive load factors.
//
// Scoring weights rationale:
//   - Digit count (+0.5/digit): More digits require tracking more information in working memory.
//   - Borrows (+1.5/borrow): Like carries, borrows require holding intermediate state while
//     computing. Weighted equally to carries in addition.
//   - Zeros crossed (+1.0/zero): Borrowing across zeros (e.g., 1000-456) requires chained
//     borrows, adding extra cognitive steps.
//   - Nice numbers (-1.0): Round number subtraction has familiar patterns (100-25=75).
//   - Negative result (+1.0): Tracking sign while computing adds cognitive overhead.
//
// Weights are initial estimates subject to tuning based on playtesting.
func (s *Subtraction) ScoreDifficulty(operands []int, answer int) float64 {
	op1, op2 := operands[0], operands[1]
	score := 1.0

	digits1 := countDigits(op1)
	digits2 := countDigits(op2)
	score += float64(digits1-1) * 0.5
	score += float64(digits2-1) * 0.5

	borrows := countBorrows(op1, op2)
	score += float64(borrows) * 1.5

	// Borrowing across zeros (e.g., 1000-456) requires chained borrows
	if op1 > op2 {
		zeros := countZerosCrossed(op1)
		if borrows > 0 && zeros > 0 {
			score += float64(minInt(zeros, borrows)) * 1.0
		}
	}

	if isNiceNumber(op1) && isNiceNumber(op2) {
		score -= 1.0
	}

	if answer < 0 {
		score += 1.0
	}

	return clampScore(score)
}

func (s *Subtraction) Generate(diff game.Difficulty) game.Question {
	return generateWithFallback(s, diff, s.makeCandidate, s.makeCandidateRelaxed)
}

// makeCandidate generates a candidate with standard operand ranges.
func (s *Subtraction) makeCandidate(diff game.Difficulty) (Candidate, bool) {
	var op1, op2 int
	switch diff {
	case game.Beginner:
		op1 = randomInRange(2, 9)
		op2 = randomInRange(1, op1) // Ensure positive result for beginners
	case game.Easy:
		op1 = randomInRange(20, 99)
		op2 = randomInRange(10, op1)
	case game.Medium:
		op1 = randomInRange(50, 300)
		op2 = randomInRange(20, op1)
	case game.Hard:
		op1 = randomInRange(100, 999)
		op2 = randomInRange(50, op1)
	case game.Expert:
		op1 = randomInRange(500, 9999)
		op2 = randomInRange(200, op1)
	default:
		op1 = randomInRange(2, 9)
		op2 = randomInRange(1, op1)
	}

	operands := []int{op1, op2}
	return Candidate{Operands: operands, Answer: s.Apply(operands)}, true
}

// makeCandidateRelaxed generates a candidate with expanded operand ranges.
func (s *Subtraction) makeCandidateRelaxed(diff game.Difficulty) (Candidate, bool) {
	var min1, max1 int
	switch diff {
	case game.Beginner:
		min1, max1 = 2, 20
	case game.Easy:
		min1, max1 = 10, 150
	case game.Medium:
		min1, max1 = 30, 500
	case game.Hard:
		min1, max1 = 50, 1500
	case game.Expert:
		min1, max1 = 200, 9999
	default:
		min1, max1 = 2, 20
	}

	op1 := randomInRange(min1, max1)
	op2 := randomInRange(1, op1)
	operands := []int{op1, op2}
	return Candidate{Operands: operands, Answer: s.Apply(operands)}, true
}
