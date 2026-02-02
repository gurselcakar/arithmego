package operations

import "testing"

func TestCountDigits(t *testing.T) {
	tests := []struct {
		n      int
		expect int
	}{
		{0, 1},
		{1, 1},
		{9, 1},
		{10, 2},
		{99, 2},
		{100, 3},
		{999, 3},
		{1000, 4},
		{-5, 1},
		{-99, 2},
		{-100, 3},
	}

	for _, tt := range tests {
		if got := countDigits(tt.n); got != tt.expect {
			t.Errorf("countDigits(%d) = %d, want %d", tt.n, got, tt.expect)
		}
	}
}

func TestCountCarries(t *testing.T) {
	tests := []struct {
		a, b   int
		expect int
	}{
		{3, 5, 0},      // 3 + 5 = 8, no carry
		{7, 8, 1},      // 7 + 8 = 15, one carry
		{23, 14, 0},    // 23 + 14 = 37, no carry
		{47, 35, 1},    // 47 + 35 = 82, one carry
		{99, 1, 2},     // 99 + 1 = 100, two carries
		{789, 456, 3},  // 789 + 456 = 1245, three carries (9+6=15, 8+5+1=14, 7+4+1=12)
		{999, 999, 3},  // 999 + 999 = 1998, three carries
	}

	for _, tt := range tests {
		if got := countCarries(tt.a, tt.b); got != tt.expect {
			t.Errorf("countCarries(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expect)
		}
	}
}

func TestCountBorrows(t *testing.T) {
	tests := []struct {
		a, b   int
		expect int
	}{
		{8, 3, 0},     // 8 - 3 = 5, no borrow
		{15, 7, 1},    // 15 - 7 = 8, one borrow
		{58, 23, 0},   // 58 - 23 = 35, no borrow
		{52, 37, 1},   // 52 - 37 = 15, one borrow
		{100, 1, 2},   // 100 - 1 = 99, two borrows
		{1000, 456, 3}, // 1000 - 456 = 544, three borrows
	}

	for _, tt := range tests {
		if got := countBorrows(tt.a, tt.b); got != tt.expect {
			t.Errorf("countBorrows(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expect)
		}
	}
}

func TestIsNiceNumber(t *testing.T) {
	tests := []struct {
		n      int
		expect bool
	}{
		{0, true},
		{5, true},
		{10, true},
		{15, true},
		{20, true},
		{25, true},
		{50, true},
		{100, true},
		{7, false},
		{13, false},
		{37, false},
		{105, false}, // 105 % 5 == 0 but 105 >= 100, so not nice unless multiple of 10 or 25
	}

	for _, tt := range tests {
		if got := isNiceNumber(tt.n); got != tt.expect {
			t.Errorf("isNiceNumber(%d) = %v, want %v", tt.n, got, tt.expect)
		}
	}
}

func TestIsRoundNumber(t *testing.T) {
	tests := []struct {
		n      int
		expect bool
	}{
		{0, true},
		{10, true},
		{20, true},
		{100, true},
		{5, false},
		{15, false},
		{99, false},
	}

	for _, tt := range tests {
		if got := isRoundNumber(tt.n); got != tt.expect {
			t.Errorf("isRoundNumber(%d) = %v, want %v", tt.n, got, tt.expect)
		}
	}
}

func TestClampScore(t *testing.T) {
	tests := []struct {
		score  float64
		expect float64
	}{
		{0.5, 1.0},
		{1.0, 1.0},
		{5.0, 5.0},
		{10.0, 10.0},
		{15.0, 10.0},
	}

	for _, tt := range tests {
		if got := clampScore(tt.score); got != tt.expect {
			t.Errorf("clampScore(%v) = %v, want %v", tt.score, got, tt.expect)
		}
	}
}

func TestDistanceFromRange(t *testing.T) {
	tests := []struct {
		score, min, max float64
		expect          float64
	}{
		{3.0, 2.0, 4.0, 0.0},  // Within range
		{2.0, 2.0, 4.0, 0.0},  // At min
		{4.0, 2.0, 4.0, 0.0},  // At max
		{1.0, 2.0, 4.0, 1.0},  // Below range
		{5.0, 2.0, 4.0, 1.0},  // Above range
	}

	for _, tt := range tests {
		if got := distanceFromRange(tt.score, tt.min, tt.max); got != tt.expect {
			t.Errorf("distanceFromRange(%v, %v, %v) = %v, want %v", tt.score, tt.min, tt.max, got, tt.expect)
		}
	}
}

func TestIntPow(t *testing.T) {
	tests := []struct {
		base, exp int
		expect    int
	}{
		{2, 0, 1},
		{2, 1, 2},
		{2, 2, 4},
		{2, 3, 8},
		{2, 10, 1024},
		{3, 3, 27},
		{10, 3, 1000},
	}

	for _, tt := range tests {
		if got := intPow(tt.base, tt.exp); got != tt.expect {
			t.Errorf("intPow(%d, %d) = %d, want %d", tt.base, tt.exp, got, tt.expect)
		}
	}
}

func TestIntPowPanicsOnNegativeExponent(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for negative exponent")
		}
	}()
	intPow(2, -1)
}

func TestFactorialHelper(t *testing.T) {
	tests := []struct {
		n      int
		expect int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{4, 24},
		{5, 120},
		{6, 720},
		{7, 5040},
		{10, 3628800},
	}

	for _, tt := range tests {
		if got := factorial(tt.n); got != tt.expect {
			t.Errorf("factorial(%d) = %d, want %d", tt.n, got, tt.expect)
		}
	}
}

func TestFactorialPanicsOnNegativeInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for negative factorial input")
		}
	}()
	factorial(-1)
}

func TestIsAcceptableFallback(t *testing.T) {
	tests := []struct {
		distance float64
		expect   bool
	}{
		{0.0, true},              // Exact match
		{0.5, true},              // Small deviation
		{1.0, true},              // Moderate deviation
		{1.5, true},              // At threshold
		{1.6, false},             // Just over threshold
		{2.0, false},             // Over threshold
		{5.0, false},             // Way over threshold
	}

	for _, tt := range tests {
		if got := isAcceptableFallback(tt.distance); got != tt.expect {
			t.Errorf("isAcceptableFallback(%v) = %v, want %v", tt.distance, got, tt.expect)
		}
	}
}

func TestMaxAcceptableDistanceConstant(t *testing.T) {
	// Verify the constant is set to a reasonable value
	if maxAcceptableDistance <= 0 {
		t.Error("maxAcceptableDistance should be positive")
	}
	if maxAcceptableDistance > 3.0 {
		t.Error("maxAcceptableDistance should not be too large (would accept very poor matches)")
	}
}

