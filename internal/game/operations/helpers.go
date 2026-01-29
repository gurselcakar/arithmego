package operations

import (
	"math/rand"
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

// intPow computes base^exp for integers.
func intPow(base, exp int) int {
	if exp < 0 {
		return 0
	}
	result := 1
	for exp > 0 {
		result *= base
		exp--
	}
	return result
}

// factorial computes n!.
func factorial(n int) int {
	if n < 0 {
		return 0
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

