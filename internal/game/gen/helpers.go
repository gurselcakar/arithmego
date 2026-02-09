package gen

import "math/rand"

// RandomInRange returns a random integer in [min, max].
func RandomInRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}
	return min + rand.Intn(max-min+1)
}

// IntPow computes base^exp for non-negative integer exponents.
func IntPow(base, exp int) int {
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

// Factorial computes n! for non-negative integers.
func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// GCD returns the greatest common divisor of a and b.
func GCD(a, b int) int {
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

// WouldOverflow returns true if base^exp would exceed maxResult.
func WouldOverflow(base, exp, maxResult int) bool {
	if base <= 1 {
		return false
	}
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
		if result > maxResult {
			return true
		}
	}
	return false
}

// MaxPowerResult is the maximum allowed result for power operations.
const MaxPowerResult = 1_000_000

// AlignToCleanDivision adjusts value so that (percent * value) % 100 == 0.
func AlignToCleanDivision(value, percent, fallback int) int {
	divisor := 100 / GCD(percent, 100)
	aligned := (value / divisor) * divisor
	if aligned == 0 {
		return fallback
	}
	return aligned
}

// PickFrom selects a random element from a slice.
func PickFrom(choices []int) int {
	return choices[rand.Intn(len(choices))]
}
