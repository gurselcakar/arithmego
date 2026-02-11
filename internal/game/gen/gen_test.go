package gen

import (
	"sort"
	"testing"

	"github.com/gurselcakar/arithmego/internal/game"
	"github.com/gurselcakar/arithmego/internal/game/expr"
)

// ---------------------------------------------------------------------------
// Helper function tests
// ---------------------------------------------------------------------------

func TestRandomInRange(t *testing.T) {
	tests := []struct {
		name     string
		min, max int
	}{
		{"normal range", 1, 10},
		{"same value", 5, 5},
		{"swapped range", 10, 1},
		{"negative range", -5, 5},
		{"both negative", -10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lo, hi := tt.min, tt.max
			if lo > hi {
				lo, hi = hi, lo
			}
			for i := 0; i < 100; i++ {
				v := RandomInRange(tt.min, tt.max)
				if v < lo || v > hi {
					t.Fatalf("RandomInRange(%d, %d) = %d, out of [%d, %d]", tt.min, tt.max, v, lo, hi)
				}
			}
		})
	}

	// Same value returns exactly that value.
	for i := 0; i < 10; i++ {
		if v := RandomInRange(7, 7); v != 7 {
			t.Fatalf("RandomInRange(7, 7) = %d, want 7", v)
		}
	}
}

func TestIntPow(t *testing.T) {
	tests := []struct {
		base, exp, want int
	}{
		{2, 0, 1},
		{2, 1, 2},
		{2, 10, 1024},
		{3, 3, 27},
		{5, 4, 625},
		{1, 100, 1},
		{0, 5, 0},
		{10, 3, 1000},
		{7, -1, 0}, // negative exponent returns 0
	}
	for _, tt := range tests {
		got := IntPow(tt.base, tt.exp)
		if got != tt.want {
			t.Errorf("IntPow(%d, %d) = %d, want %d", tt.base, tt.exp, got, tt.want)
		}
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		n, want int
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
		got := Factorial(tt.n)
		if got != tt.want {
			t.Errorf("Factorial(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{12, 8, 4},
		{7, 13, 1},
		{100, 75, 25},
		{0, 5, 5},
		{5, 0, 5},
		{0, 0, 0},
		{-12, 8, 4},
		{12, -8, 4},
		{-12, -8, 4},
		{1, 1, 1},
		{17, 17, 17},
	}
	for _, tt := range tests {
		got := GCD(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestWouldOverflow(t *testing.T) {
	tests := []struct {
		base, exp, max int
		want           bool
	}{
		{2, 10, 1024, false},
		{2, 10, 1023, true},
		{2, 20, 1_000_000, true},
		{10, 6, 1_000_000, false},
		{10, 7, 1_000_000, true},
		{1, 100, 10, false},  // base <= 1 always false
		{0, 100, 10, false},  // base <= 1 always false
		{-1, 5, 10, false},   // base <= 1 always false
		{3, 0, 1, false},     // 3^0 = 1, not exceeding 1
		{100, 3, 999999, true},
	}
	for _, tt := range tests {
		got := WouldOverflow(tt.base, tt.exp, tt.max)
		if got != tt.want {
			t.Errorf("WouldOverflow(%d, %d, %d) = %v, want %v", tt.base, tt.exp, tt.max, got, tt.want)
		}
	}
}

func TestAlignToCleanDivision(t *testing.T) {
	tests := []struct {
		value, percent, fallback, want int
	}{
		{100, 50, 10, 100}, // GCD(50,100)=50, divisor=2, aligned=100
		{100, 25, 10, 100}, // GCD(25,100)=25, divisor=4, aligned=100
		{47, 10, 20, 40},   // GCD(10,100)=10, divisor=10, aligned=40
		{3, 50, 20, 2},     // GCD(50,100)=50, divisor=2, aligned=2
		{1, 50, 20, 20},    // divisor=2, aligned=0 -> fallback=20
		{0, 50, 20, 20},    // aligned=0 -> fallback
		{200, 20, 10, 200}, // GCD(20,100)=20, divisor=5, aligned=200
	}
	for _, tt := range tests {
		got := AlignToCleanDivision(tt.value, tt.percent, tt.fallback)
		if got != tt.want {
			t.Errorf("AlignToCleanDivision(%d, %d, %d) = %d, want %d",
				tt.value, tt.percent, tt.fallback, got, tt.want)
		}
	}
}

func TestPickFrom(t *testing.T) {
	choices := []int{10, 20, 30}
	seen := make(map[int]bool)
	for i := 0; i < 200; i++ {
		v := PickFrom(choices)
		seen[v] = true
	}
	for _, c := range choices {
		if !seen[c] {
			t.Errorf("PickFrom never returned %d after 200 calls", c)
		}
	}
}

// ---------------------------------------------------------------------------
// BuildQuestion test
// ---------------------------------------------------------------------------

func TestBuildQuestion(t *testing.T) {
	// 3 + 5 = 8
	e := &expr.BinOp{
		Op:    expr.OpAdd,
		Left:  &expr.Num{Value: 3},
		Right: &expr.Num{Value: 5},
	}
	q := BuildQuestion(e, "Addition")

	if q == nil {
		t.Fatal("BuildQuestion returned nil")
	}
	if q.Answer != 8 {
		t.Errorf("Answer = %d, want 8", q.Answer)
	}
	if q.Display != e.Format() {
		t.Errorf("Display = %q, want %q", q.Display, e.Format())
	}
	if q.Key != e.Key() {
		t.Errorf("Key = %q, want %q", q.Key, e.Key())
	}
	if q.OpLabel != "Addition" {
		t.Errorf("OpLabel = %q, want %q", q.OpLabel, "Addition")
	}
	if q.Expression != e {
		t.Error("Expression field does not point to original expr")
	}
}

// ---------------------------------------------------------------------------
// Registry tests
// ---------------------------------------------------------------------------

func TestRegistryAllGenerators(t *testing.T) {
	expectedLabels := []string{
		"Addition",
		"Subtraction",
		"Multiplication",
		"Division",
		"Square",
		"Cube",
		"Square Root",
		"Cube Root",
		"Power",
		"Modulo",
		"Percentage",
		"Factorial",
		"Mixed Basics",
		"Mixed Powers",
		"Mixed Advanced",
		"Anything Goes",
	}

	all := All()
	if len(all) != len(expectedLabels) {
		t.Fatalf("All() returned %d generators, want %d", len(all), len(expectedLabels))
	}

	labels := make([]string, len(all))
	for i, g := range all {
		labels[i] = g.Label()
	}
	sort.Strings(labels)
	sort.Strings(expectedLabels)

	for i := range expectedLabels {
		if labels[i] != expectedLabels[i] {
			t.Errorf("label[%d] = %q, want %q", i, labels[i], expectedLabels[i])
		}
	}
}

func TestRegistryGet(t *testing.T) {
	// Known label returns the correct generator.
	g, ok := Get("Addition")
	if !ok {
		t.Fatal("Get(\"Addition\") returned false")
	}
	if g.Label() != "Addition" {
		t.Errorf("Get(\"Addition\").Label() = %q, want \"Addition\"", g.Label())
	}

	// Unknown label returns nil, false.
	g, ok = Get("Nonexistent Mode")
	if ok {
		t.Error("Get(\"Nonexistent Mode\") returned true, want false")
	}
	if g != nil {
		t.Error("Get(\"Nonexistent Mode\") returned non-nil generator")
	}
}

// ---------------------------------------------------------------------------
// PickPattern test
// ---------------------------------------------------------------------------

func TestPickPattern(t *testing.T) {
	dummyExpr := &expr.Num{Value: 1}
	called := [3]int{}

	patterns := []WeightedPattern{
		{Pattern: func(d game.Difficulty) (expr.Expr, bool) { called[0]++; return dummyExpr, true }, Weight: 1},
		{Pattern: func(d game.Difficulty) (expr.Expr, bool) { called[1]++; return dummyExpr, true }, Weight: 3},
		{Pattern: func(d game.Difficulty) (expr.Expr, bool) { called[2]++; return dummyExpr, true }, Weight: 6},
	}

	// Call PickPattern many times, then call the returned pattern.
	n := 10000
	for i := 0; i < n; i++ {
		p := PickPattern(patterns)
		p(game.Medium) // just to trigger the counter
	}

	// With weights 1:3:6, expect roughly 10%, 30%, 60%.
	for i, c := range called {
		pct := float64(c) / float64(n) * 100
		if pct < 1 { // Sanity: each should be picked at least sometimes.
			t.Errorf("Pattern %d was picked %.1f%% of the time, suspiciously low", i, pct)
		}
	}

	// Pattern index 2 (weight 6) should be picked more than pattern 0 (weight 1).
	if called[2] < called[0] {
		t.Errorf("Pattern 2 (weight 6, count %d) should be picked more than pattern 0 (weight 1, count %d)",
			called[2], called[0])
	}
}

func TestPickPatternZeroTotalWeight(t *testing.T) {
	dummyExpr := &expr.Num{Value: 42}
	first := func(d game.Difficulty) (expr.Expr, bool) { return dummyExpr, true }
	second := func(d game.Difficulty) (expr.Expr, bool) { return &expr.Num{Value: 99}, true }

	patterns := []WeightedPattern{
		{Pattern: first, Weight: 0},
		{Pattern: second, Weight: 0},
	}

	// Should return the first pattern when total weight is 0.
	p := PickPattern(patterns)
	e, ok := p(game.Beginner)
	if !ok {
		t.Fatal("returned pattern indicated invalid")
	}
	if e.Eval() != 42 {
		t.Errorf("expected first pattern (eval 42), got eval %d", e.Eval())
	}
}

// ---------------------------------------------------------------------------
// Generator smoke tests (all 16 generators x 5 difficulties)
// ---------------------------------------------------------------------------

func TestAllGeneratorsSmoke(t *testing.T) {
	difficulties := game.AllDifficulties()
	generators := All()

	if len(generators) == 0 {
		t.Fatal("no generators registered")
	}

	for _, g := range generators {
		g := g // capture
		t.Run(g.Label(), func(t *testing.T) {
			for _, diff := range difficulties {
				diff := diff // capture
				t.Run(diff.String(), func(t *testing.T) {
					for i := 0; i < 10; i++ {
						q := g.Generate(diff)
						if q == nil {
							t.Fatalf("Generate(%s) returned nil on attempt %d", diff, i)
						}
						if q.Display == "" {
							t.Error("Question.Display is empty")
						}
						if q.Key == "" {
							t.Error("Question.Key is empty")
						}
						if q.OpLabel == "" {
							t.Error("Question.OpLabel is empty")
						}
						if q.OpLabel != g.Label() {
							t.Errorf("OpLabel = %q, want %q", q.OpLabel, g.Label())
						}
						if q.Expression == nil {
							t.Fatal("Question.Expression is nil")
						}
						if q.Answer != q.Expression.Eval() {
							t.Errorf("Answer = %d, but Expression.Eval() = %d", q.Answer, q.Expression.Eval())
						}
					}
				})
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TryGenerate test
// ---------------------------------------------------------------------------

func TestTryGenerate(t *testing.T) {
	// Valid pattern set should produce a question.
	patterns := PatternSet{
		game.Medium: {
			{Pattern: func(d game.Difficulty) (expr.Expr, bool) {
				return &expr.BinOp{
					Op:    expr.OpAdd,
					Left:  &expr.Num{Value: 2},
					Right: &expr.Num{Value: 3},
				}, true
			}, Weight: 1},
		},
	}

	q := TryGenerate(patterns, game.Medium, "Test", 10)
	if q == nil {
		t.Fatal("TryGenerate returned nil for valid pattern")
	}
	if q.Answer != 5 {
		t.Errorf("Answer = %d, want 5", q.Answer)
	}

	// Missing difficulty returns nil.
	q = TryGenerate(patterns, game.Expert, "Test", 10)
	if q != nil {
		t.Error("TryGenerate should return nil for missing difficulty")
	}

	// Pattern that always fails returns nil.
	failPatterns := PatternSet{
		game.Easy: {
			{Pattern: func(d game.Difficulty) (expr.Expr, bool) {
				return nil, false
			}, Weight: 1},
		},
	}
	q = TryGenerate(failPatterns, game.Easy, "Fail", 10)
	if q != nil {
		t.Error("TryGenerate should return nil when all attempts fail")
	}
}
