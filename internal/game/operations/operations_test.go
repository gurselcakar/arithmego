package operations

import (
	"testing"

	"github.com/gurselcakar/arithmego/internal/game"
)

// TestAddition tests the Addition operation.
func TestAddition(t *testing.T) {
	op := &Addition{}

	t.Run("metadata", func(t *testing.T) {
		if op.Name() != "Addition" {
			t.Errorf("Name() = %q, want Addition", op.Name())
		}
		if op.Symbol() != "+" {
			t.Errorf("Symbol() = %q, want +", op.Symbol())
		}
		if op.Arity() != game.Binary {
			t.Errorf("Arity() = %v, want Binary", op.Arity())
		}
		if op.Category() != game.CategoryBasic {
			t.Errorf("Category() = %v, want CategoryBasic", op.Category())
		}
	})

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{3, 5}, 8},
			{[]int{0, 0}, 0},
			{[]int{100, 200}, 300},
			{[]int{-5, 10}, 5},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{5, 3}); got != "5 + 3" {
			t.Errorf("Format([5, 3]) = %q, want '5 + 3'", got)
		}
	})

	t.Run("generate", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			q := op.Generate(diff)
			if q.Answer != op.Apply(q.Operands) {
				t.Errorf("Generate(%v): answer %d != Apply(operands) %d", diff, q.Answer, op.Apply(q.Operands))
			}
			if q.Operation != op {
				t.Errorf("Generate(%v): operation mismatch", diff)
			}
		}
	})
}

// TestSubtraction tests the Subtraction operation.
func TestSubtraction(t *testing.T) {
	op := &Subtraction{}

	t.Run("metadata", func(t *testing.T) {
		if op.Name() != "Subtraction" {
			t.Errorf("Name() = %q, want Subtraction", op.Name())
		}
		if op.Symbol() != "−" {
			t.Errorf("Symbol() = %q, want −", op.Symbol())
		}
	})

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{8, 3}, 5},
			{[]int{10, 10}, 0},
			{[]int{5, 10}, -5},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("generate", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			q := op.Generate(diff)
			if q.Answer != op.Apply(q.Operands) {
				t.Errorf("Generate(%v): answer %d != Apply(operands) %d", diff, q.Answer, op.Apply(q.Operands))
			}
		}
	})
}

// TestMultiplication tests the Multiplication operation.
func TestMultiplication(t *testing.T) {
	op := &Multiplication{}

	t.Run("metadata", func(t *testing.T) {
		if op.Name() != "Multiplication" {
			t.Errorf("Name() = %q, want Multiplication", op.Name())
		}
		if op.Symbol() != "×" {
			t.Errorf("Symbol() = %q, want ×", op.Symbol())
		}
	})

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{6, 7}, 42},
			{[]int{0, 100}, 0},
			{[]int{12, 12}, 144},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("generate", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			q := op.Generate(diff)
			if q.Answer != op.Apply(q.Operands) {
				t.Errorf("Generate(%v): answer %d != Apply(operands) %d", diff, q.Answer, op.Apply(q.Operands))
			}
		}
	})
}

// TestDivision tests the Division operation.
func TestDivision(t *testing.T) {
	op := &Division{}

	t.Run("metadata", func(t *testing.T) {
		if op.Name() != "Division" {
			t.Errorf("Name() = %q, want Division", op.Name())
		}
		if op.Symbol() != "÷" {
			t.Errorf("Symbol() = %q, want ÷", op.Symbol())
		}
	})

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{20, 4}, 5},
			{[]int{56, 8}, 7},
			{[]int{100, 10}, 10},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("generate produces clean division", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			for i := 0; i < 10; i++ {
				q := op.Generate(diff)
				dividend, divisor := q.Operands[0], q.Operands[1]
				if dividend%divisor != 0 {
					t.Errorf("Generate(%v): %d %% %d != 0 (not clean division)", diff, dividend, divisor)
				}
				if q.Answer != dividend/divisor {
					t.Errorf("Generate(%v): answer mismatch", diff)
				}
			}
		}
	})

	t.Run("apply panics on division by zero", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for division by zero")
			}
		}()
		op.Apply([]int{10, 0})
	})
}

// TestSquare tests the Square operation.
func TestSquare(t *testing.T) {
	op := &Square{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{5}, 25},
			{[]int{12}, 144},
			{[]int{1}, 1},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{7}); got != "7²" {
			t.Errorf("Format([7]) = %q, want '7²'", got)
		}
	})

	t.Run("generate", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			q := op.Generate(diff)
			if q.Answer != q.Operands[0]*q.Operands[0] {
				t.Errorf("Generate(%v): answer mismatch", diff)
			}
		}
	})
}

// TestCube tests the Cube operation.
func TestCube(t *testing.T) {
	op := &Cube{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{3}, 27},
			{[]int{4}, 64},
			{[]int{5}, 125},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{4}); got != "4³" {
			t.Errorf("Format([4]) = %q, want '4³'", got)
		}
	})
}

// TestSquareRoot tests the SquareRoot operation.
func TestSquareRoot(t *testing.T) {
	op := &SquareRoot{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{49}, 7},
			{[]int{144}, 12},
			{[]int{1}, 1},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{49}); got != "√49" {
			t.Errorf("Format([49]) = %q, want '√49'", got)
		}
	})

	t.Run("generate produces perfect squares", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			for i := 0; i < 10; i++ {
				q := op.Generate(diff)
				operand := q.Operands[0]
				if q.Answer*q.Answer != operand {
					t.Errorf("Generate(%v): %d is not a perfect square (√%d != %d)", diff, operand, operand, q.Answer)
				}
			}
		}
	})
}

// TestCubeRoot tests the CubeRoot operation.
func TestCubeRoot(t *testing.T) {
	op := &CubeRoot{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{64}, 4},
			{[]int{125}, 5},
			{[]int{1000}, 10},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("generate produces perfect cubes", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			for i := 0; i < 10; i++ {
				q := op.Generate(diff)
				operand := q.Operands[0]
				if q.Answer*q.Answer*q.Answer != operand {
					t.Errorf("Generate(%v): %d is not a perfect cube", diff, operand)
				}
			}
		}
	})
}

// TestModulo tests the Modulo operation.
func TestModulo(t *testing.T) {
	op := &Modulo{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{17, 5}, 2},
			{[]int{20, 4}, 0},
			{[]int{10, 3}, 1},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{17, 5}); got != "17 mod 5" {
			t.Errorf("Format([17, 5]) = %q, want '17 mod 5'", got)
		}
	})

	t.Run("generate", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			q := op.Generate(diff)
			if q.Answer != q.Operands[0]%q.Operands[1] {
				t.Errorf("Generate(%v): answer mismatch", diff)
			}
		}
	})

	t.Run("apply panics on modulo by zero", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for modulo by zero")
			}
		}()
		op.Apply([]int{10, 0})
	})
}

// TestPower tests the Power operation.
func TestPower(t *testing.T) {
	op := &Power{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{2, 4}, 16},
			{[]int{3, 3}, 27},
			{[]int{10, 2}, 100},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{2, 4}); got != "2^4" {
			t.Errorf("Format([2, 4]) = %q, want '2^4'", got)
		}
	})

	t.Run("generate", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			q := op.Generate(diff)
			expected := intPow(q.Operands[0], q.Operands[1])
			if q.Answer != expected {
				t.Errorf("Generate(%v): answer %d != expected %d", diff, q.Answer, expected)
			}
		}
	})

	t.Run("generate respects max result limit", func(t *testing.T) {
		// Run many iterations to ensure overflow protection works
		for i := 0; i < 100; i++ {
			for _, diff := range game.AllDifficulties() {
				q := op.Generate(diff)
				if q.Answer > 1000000 {
					t.Errorf("Generate(%v): answer %d exceeds max limit of 1000000", diff, q.Answer)
				}
			}
		}
	})
}

// TestWouldOverflow tests the overflow detection helper.
func TestWouldOverflow(t *testing.T) {
	tests := []struct {
		base, exp, max int
		expect         bool
	}{
		{2, 10, 1000000, false},  // 1024 < 1000000
		{10, 6, 1000000, false},  // 1000000 == 1000000
		{10, 7, 1000000, true},   // 10000000 > 1000000
		{100, 4, 1000000, true},  // 100000000 > 1000000
		{2, 20, 1000000, true},   // 1048576 > 1000000
		{1, 100, 1000000, false}, // 1 < 1000000
		{5, 8, 1000000, false},   // 390625 < 1000000
		{5, 9, 1000000, true},    // 1953125 > 1000000
	}

	for _, tt := range tests {
		if got := wouldOverflow(tt.base, tt.exp, tt.max); got != tt.expect {
			t.Errorf("wouldOverflow(%d, %d, %d) = %v, want %v", tt.base, tt.exp, tt.max, got, tt.expect)
		}
	}
}

// TestPercentage tests the Percentage operation.
func TestPercentage(t *testing.T) {
	op := &Percentage{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{25, 80}, 20},  // 25% of 80 = 20
			{[]int{50, 100}, 50}, // 50% of 100 = 50
			{[]int{10, 200}, 20}, // 10% of 200 = 20
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{25, 80}); got != "25% of 80" {
			t.Errorf("Format([25, 80]) = %q, want '25%% of 80'", got)
		}
	})

	t.Run("generate produces integer results", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			for i := 0; i < 10; i++ {
				q := op.Generate(diff)
				percent, value := q.Operands[0], q.Operands[1]
				// Verify the result is exact (no rounding errors)
				expected := (percent * value) / 100
				if (percent*value)%100 != 0 {
					t.Errorf("Generate(%v): %d%% of %d doesn't produce clean integer", diff, percent, value)
				}
				if q.Answer != expected {
					t.Errorf("Generate(%v): answer mismatch", diff)
				}
			}
		}
	})
}

// TestFactorial tests the Factorial operation.
func TestFactorial(t *testing.T) {
	op := &Factorial{}

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			operands []int
			expect   int
		}{
			{[]int{5}, 120},
			{[]int{0}, 1},
			{[]int{1}, 1},
			{[]int{6}, 720},
		}
		for _, tt := range tests {
			if got := op.Apply(tt.operands); got != tt.expect {
				t.Errorf("Apply(%v) = %d, want %d", tt.operands, got, tt.expect)
			}
		}
	})

	t.Run("format", func(t *testing.T) {
		if got := op.Format([]int{5}); got != "5!" {
			t.Errorf("Format([5]) = %q, want '5!'", got)
		}
	})

	t.Run("generate", func(t *testing.T) {
		for _, diff := range game.AllDifficulties() {
			q := op.Generate(diff)
			expected := factorial(q.Operands[0])
			if q.Answer != expected {
				t.Errorf("Generate(%v): answer %d != expected %d", diff, q.Answer, expected)
			}
			// Factorial should be limited to n <= 10
			if q.Operands[0] > 10 {
				t.Errorf("Generate(%v): n=%d exceeds limit of 10", diff, q.Operands[0])
			}
		}
	})
}

// TestAllOperationsGenerateValidQuestions ensures all operations generate valid questions.
func TestAllOperationsGenerateValidQuestions(t *testing.T) {
	for _, op := range All() {
		t.Run(op.Name(), func(t *testing.T) {
			for _, diff := range game.AllDifficulties() {
				for i := 0; i < 5; i++ {
					q := op.Generate(diff)

					// Question should have correct operation
					if q.Operation != op {
						t.Errorf("Generate(%v): operation mismatch", diff)
					}

					// Question should have non-empty display
					if q.Display == "" {
						t.Errorf("Generate(%v): empty display", diff)
					}

					// Answer should match Apply
					computed := op.Apply(q.Operands)
					if q.Answer != computed {
						t.Errorf("Generate(%v): answer %d != computed %d", diff, q.Answer, computed)
					}

					// Operands should match arity
					if op.Arity() == game.Unary && len(q.Operands) != 1 {
						t.Errorf("Generate(%v): unary operation has %d operands", diff, len(q.Operands))
					}
					if op.Arity() == game.Binary && len(q.Operands) != 2 {
						t.Errorf("Generate(%v): binary operation has %d operands", diff, len(q.Operands))
					}
				}
			}
		})
	}
}
