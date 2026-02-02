package operations

import (
	"math"
	"math/rand"

	"github.com/gurselcakar/arithmego/internal/game"
)

// countDigits returns the number of digits in n.
func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	if n < 0 {
		n = -n
	}
	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count
}

// countCarries returns the number of carries in a + b.
func countCarries(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	carries := 0
	carry := 0
	for a > 0 || b > 0 {
		sum := a%10 + b%10 + carry
		if sum >= 10 {
			carries++
			carry = 1
		} else {
			carry = 0
		}
		a /= 10
		b /= 10
	}
	return carries
}

// countBorrows returns the number of borrows in a - b (assumes a >= b).
func countBorrows(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	// Ensure a >= b for counting borrows
	if a < b {
		a, b = b, a
	}
	borrows := 0
	borrow := 0
	for a > 0 || b > 0 {
		diff := a%10 - b%10 - borrow
		if diff < 0 {
			borrows++
			borrow = 1
		} else {
			borrow = 0
		}
		a /= 10
		b /= 10
	}
	return borrows
}

// countZerosCrossed counts zeros that need to be borrowed across (e.g., 1000 - 456).
func countZerosCrossed(a int) int {
	if a <= 0 {
		return 0
	}
	zeros := 0
	for a > 0 {
		if a%10 == 0 {
			zeros++
		}
		a /= 10
	}
	return zeros
}

// isNiceNumber returns true if n is a "nice" number (multiple of 5, 10, 25, etc.)
func isNiceNumber(n int) bool {
	if n == 0 {
		return true
	}
	if n < 0 {
		n = -n
	}
	return n%10 == 0 || n%25 == 0 || (n%5 == 0 && n < 100)
}

// isRoundNumber returns true if n is a round number (multiple of 10, 100, etc.)
func isRoundNumber(n int) bool {
	if n == 0 {
		return true
	}
	if n < 0 {
		n = -n
	}
	return n%10 == 0
}

// crossesBoundary returns true if a + b crosses a power of 10 boundary.
func crossesBoundary(a, b, result int) bool {
	if a <= 0 {
		return false
	}
	aBoundary := powerOf10Above(a)
	return result >= aBoundary && a < aBoundary
}

// powerOf10Above returns the next power of 10 above n.
func powerOf10Above(n int) int {
	if n <= 0 {
		return 10
	}
	p := 10
	for p <= n {
		p *= 10
	}
	return p
}

// clampScore ensures the score stays within 1.0 to 10.0.
func clampScore(score float64) float64 {
	if score < 1.0 {
		return 1.0
	}
	if score > 10.0 {
		return 10.0
	}
	return score
}

// maxAcceptableDistance is the maximum allowed deviation from the target
// difficulty range when falling back to the closest match. If the best
// question found exceeds this threshold, a second pass with relaxed
// constraints should be attempted.
const maxAcceptableDistance = 1.5

// distanceFromRange returns how far a score is from a range.
// Returns 0 if within range.
func distanceFromRange(score, min, max float64) float64 {
	if score >= min && score <= max {
		return 0
	}
	if score < min {
		return min - score
	}
	return score - max
}

// isAcceptableFallback returns true if the fallback question's distance
// from the target range is within acceptable bounds.
func isAcceptableFallback(distance float64) bool {
	return distance <= maxAcceptableDistance
}

// randomInRange returns a random integer in [min, max].
func randomInRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}
	return min + rand.Intn(max-min+1)
}

// minInt returns the smaller of a or b.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isTimesTableFact returns true if divisor and quotient are both <= 12.
func isTimesTableFact(divisor, quotient int) bool {
	return divisor <= 12 && quotient <= 12
}

// intPow computes base^exp for non-negative integer exponents.
// Panics if exp < 0 since negative exponents produce non-integer results.
func intPow(base, exp int) int {
	if exp < 0 {
		panic("intPow: negative exponent")
	}
	result := 1
	for exp > 0 {
		result *= base
		exp--
	}
	return result
}

// factorial computes n! for non-negative integers.
// Panics if n < 0 since factorial is undefined for negative numbers.
func factorial(n int) int {
	if n < 0 {
		panic("factorial: negative input")
	}
	if n <= 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// Candidate represents a generated question candidate before validation.
type Candidate struct {
	Operands []int
	Answer   int
}

// CandidateFunc generates a candidate for a given difficulty.
// Returns the candidate and whether it's valid. Some operations (like percentage)
// may generate invalid candidates that should be skipped.
type CandidateFunc func(diff game.Difficulty) (Candidate, bool)

// generateWithFallback generates a question using the provided candidate generator.
// It tries up to 100 candidates, tracking the closest match. If no exact match is
// found and the best candidate exceeds maxAcceptableDistance, it runs a second pass
// with the relaxed generator (if provided) for 50 more attempts.
//
// Parameters:
//   - op: The operation (used for scoring and formatting)
//   - diff: Target difficulty level
//   - primary: Main candidate generator for this difficulty
//   - relaxed: Optional relaxed generator with wider constraints (can be nil)
func generateWithFallback(
	op game.Operation,
	diff game.Difficulty,
	primary CandidateFunc,
	relaxed CandidateFunc,
) game.Question {
	minScore, maxScore := diff.ScoreRange()

	var bestQuestion game.Question
	bestDistance := math.MaxFloat64

	// Primary pass: 100 attempts
	for attempts := 0; attempts < 100; attempts++ {
		candidate, valid := primary(diff)
		if !valid {
			continue
		}

		score := op.ScoreDifficulty(candidate.Operands, candidate.Answer)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  candidate.Operands,
				Operation: op,
				Answer:    candidate.Answer,
				Display:   op.Format(candidate.Operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  candidate.Operands,
				Operation: op,
				Answer:    candidate.Answer,
				Display:   op.Format(candidate.Operands),
			}
		}
	}

	// If fallback is acceptable or no relaxed generator, return best match
	if isAcceptableFallback(bestDistance) || relaxed == nil {
		return bestQuestion
	}

	// Relaxed pass: 50 more attempts with wider constraints
	for attempts := 0; attempts < 50; attempts++ {
		candidate, valid := relaxed(diff)
		if !valid {
			continue
		}

		score := op.ScoreDifficulty(candidate.Operands, candidate.Answer)

		if score >= minScore && score <= maxScore {
			return game.Question{
				Operands:  candidate.Operands,
				Operation: op,
				Answer:    candidate.Answer,
				Display:   op.Format(candidate.Operands),
			}
		}

		dist := distanceFromRange(score, minScore, maxScore)
		if dist < bestDistance {
			bestDistance = dist
			bestQuestion = game.Question{
				Operands:  candidate.Operands,
				Operation: op,
				Answer:    candidate.Answer,
				Display:   op.Format(candidate.Operands),
			}
		}
	}

	return bestQuestion
}

