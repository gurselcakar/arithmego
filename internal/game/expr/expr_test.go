package expr

import (
	"testing"
)

// ---------------------------------------------------------------------------
// BinOpKind method tests
// ---------------------------------------------------------------------------

func TestBinOpKind_Precedence(t *testing.T) {
	tests := []struct {
		op   BinOpKind
		want int
	}{
		{OpAdd, 1},
		{OpSub, 1},
		{OpMul, 2},
		{OpDiv, 2},
		{OpMod, 2},
		{OpPct, 2},
		{BinOpKind(99), 0}, // unknown operator
	}
	for _, tt := range tests {
		if got := tt.op.Precedence(); got != tt.want {
			t.Errorf("BinOpKind(%d).Precedence() = %d, want %d", tt.op, got, tt.want)
		}
	}
}

func TestBinOpKind_Symbol(t *testing.T) {
	tests := []struct {
		op   BinOpKind
		want string
	}{
		{OpAdd, "+"},
		{OpSub, "−"},
		{OpMul, "×"},
		{OpDiv, "÷"},
		{OpMod, "mod"},
		{OpPct, "% of"},
		{BinOpKind(99), "?"},
	}
	for _, tt := range tests {
		if got := tt.op.Symbol(); got != tt.want {
			t.Errorf("BinOpKind(%d).Symbol() = %q, want %q", tt.op, got, tt.want)
		}
	}
}

func TestBinOpKind_KeySymbol(t *testing.T) {
	tests := []struct {
		op   BinOpKind
		want string
	}{
		{OpAdd, "+"},
		{OpSub, "-"},
		{OpMul, "*"},
		{OpDiv, "/"},
		{OpMod, "%"},
		{OpPct, "pct"},
		{BinOpKind(99), "?"},
	}
	for _, tt := range tests {
		if got := tt.op.KeySymbol(); got != tt.want {
			t.Errorf("BinOpKind(%d).KeySymbol() = %q, want %q", tt.op, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Eval tests
// ---------------------------------------------------------------------------

func TestNum_Eval(t *testing.T) {
	tests := []struct {
		val  int
		want int
	}{
		{0, 0},
		{42, 42},
		{-7, -7},
	}
	for _, tt := range tests {
		n := &Num{Value: tt.val}
		if got := n.Eval(); got != tt.want {
			t.Errorf("Num{%d}.Eval() = %d, want %d", tt.val, got, tt.want)
		}
	}
}

func TestBinOp_Eval(t *testing.T) {
	tests := []struct {
		name        string
		op          BinOpKind
		left, right int
		want        int
	}{
		{"add", OpAdd, 5, 3, 8},
		{"add negative", OpAdd, -2, 7, 5},
		{"sub", OpSub, 10, 4, 6},
		{"sub negative result", OpSub, 3, 8, -5},
		{"mul", OpMul, 6, 7, 42},
		{"mul by zero", OpMul, 9, 0, 0},
		{"mul negatives", OpMul, -3, -4, 12},
		{"div", OpDiv, 15, 3, 5},
		{"div truncation", OpDiv, 7, 2, 3},
		{"div by zero", OpDiv, 10, 0, 0},
		{"mod", OpMod, 10, 3, 1},
		{"mod exact", OpMod, 9, 3, 0},
		{"mod by zero", OpMod, 10, 0, 0},
		{"pct", OpPct, 50, 200, 100},     // 50% of 200 = 100
		{"pct small", OpPct, 25, 80, 20}, // 25% of 80 = 20
		{"pct 10 of 50", OpPct, 10, 50, 5},
		{"unknown op", BinOpKind(99), 5, 3, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinOp{Op: tt.op, Left: &Num{tt.left}, Right: &Num{tt.right}}
			if got := b.Eval(); got != tt.want {
				t.Errorf("BinOp{%d %s %d}.Eval() = %d, want %d",
					tt.left, tt.op.Symbol(), tt.right, got, tt.want)
			}
		})
	}
}

func TestParen_Eval(t *testing.T) {
	inner := &BinOp{Op: OpAdd, Left: &Num{3}, Right: &Num{4}}
	p := &Paren{Inner: inner}
	if got := p.Eval(); got != 7 {
		t.Errorf("Paren{3+4}.Eval() = %d, want 7", got)
	}
}

func TestUnaryPrefix_Eval(t *testing.T) {
	tests := []struct {
		name string
		op   UnaryOp
		val  int
		want int
	}{
		{"sqrt 49", OpSqrt, 49, 7},
		{"sqrt 0", OpSqrt, 0, 0},
		{"sqrt 1", OpSqrt, 1, 1},
		{"sqrt 100", OpSqrt, 100, 10},
		{"sqrt non-perfect", OpSqrt, 10, 3}, // int(sqrt(10)) = 3
		{"cbrt 27", OpCbrt, 27, 3},
		{"cbrt 0", OpCbrt, 0, 0},
		{"cbrt 1", OpCbrt, 1, 1},
		{"cbrt 64", OpCbrt, 64, 4},
		{"cbrt 125", OpCbrt, 125, 5},
		{"unknown prefix op", UnaryOp(99), 42, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnaryPrefix{Op: tt.op, Operand: &Num{tt.val}}
			if got := u.Eval(); got != tt.want {
				t.Errorf("UnaryPrefix{%d, %d}.Eval() = %d, want %d", tt.op, tt.val, got, tt.want)
			}
		})
	}
}

func TestUnarySuffix_Eval(t *testing.T) {
	tests := []struct {
		name string
		op   UnaryOp
		val  int
		want int
	}{
		{"square 5", OpSquare, 5, 25},
		{"square 0", OpSquare, 0, 0},
		{"square -3", OpSquare, -3, 9},
		{"square 1", OpSquare, 1, 1},
		{"cube 3", OpCube, 3, 27},
		{"cube 0", OpCube, 0, 0},
		{"cube -2", OpCube, -2, -8},
		{"cube 1", OpCube, 1, 1},
		{"factorial 0", OpFactorial, 0, 1},
		{"factorial 1", OpFactorial, 1, 1},
		{"factorial 5", OpFactorial, 5, 120},
		{"factorial 6", OpFactorial, 6, 720},
		{"factorial 10", OpFactorial, 10, 3628800},
		{"unknown suffix op", UnaryOp(99), 7, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnarySuffix{Op: tt.op, Operand: &Num{tt.val}}
			if got := u.Eval(); got != tt.want {
				t.Errorf("UnarySuffix{%d, %d}.Eval() = %d, want %d", tt.op, tt.val, got, tt.want)
			}
		})
	}
}

func TestPow_Eval(t *testing.T) {
	tests := []struct {
		name      string
		base, exp int
		want      int
	}{
		{"2^3", 2, 3, 8},
		{"2^0", 2, 0, 1},
		{"2^1", 2, 1, 2},
		{"5^2", 5, 2, 25},
		{"3^4", 3, 4, 81},
		{"10^3", 10, 3, 1000},
		{"0^5", 0, 5, 0},
		{"1^100", 1, 100, 1},
		{"negative exp", 2, -1, 0},
		{"negative base", -2, 3, -8},
		{"negative base even", -2, 4, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pow{Base: &Num{tt.base}, Exp: &Num{tt.exp}}
			if got := p.Eval(); got != tt.want {
				t.Errorf("Pow{%d, %d}.Eval() = %d, want %d", tt.base, tt.exp, got, tt.want)
			}
		})
	}
}

func TestEval_NestedExpressions(t *testing.T) {
	// (5 + 3) * 2 = 16
	expr1 := &BinOp{
		Op:    OpMul,
		Left:  &Paren{Inner: &BinOp{Op: OpAdd, Left: &Num{5}, Right: &Num{3}}},
		Right: &Num{2},
	}
	if got := expr1.Eval(); got != 16 {
		t.Errorf("(5+3)*2 = %d, want 16", got)
	}

	// 5 + 3 * 2 = 11 (tree structure determines evaluation, not PEMDAS)
	expr2 := &BinOp{
		Op:   OpAdd,
		Left: &Num{5},
		Right: &BinOp{Op: OpMul, Left: &Num{3}, Right: &Num{2}},
	}
	if got := expr2.Eval(); got != 11 {
		t.Errorf("5+3*2 = %d, want 11", got)
	}

	// √(4²) = √16 = 4
	expr3 := &UnaryPrefix{
		Op:      OpSqrt,
		Operand: &UnarySuffix{Op: OpSquare, Operand: &Num{4}},
	}
	if got := expr3.Eval(); got != 4 {
		t.Errorf("√(4²) = %d, want 4", got)
	}
}

// ---------------------------------------------------------------------------
// Format tests
// ---------------------------------------------------------------------------

func TestNum_Format(t *testing.T) {
	tests := []struct {
		val  int
		want string
	}{
		{0, "0"},
		{42, "42"},
		{-7, "-7"},
		{1000, "1000"},
	}
	for _, tt := range tests {
		n := &Num{Value: tt.val}
		if got := n.Format(); got != tt.want {
			t.Errorf("Num{%d}.Format() = %q, want %q", tt.val, got, tt.want)
		}
	}
}

func TestBinOp_Format_Simple(t *testing.T) {
	tests := []struct {
		name        string
		op          BinOpKind
		left, right int
		want        string
	}{
		{"add", OpAdd, 5, 3, "5 + 3"},
		{"sub", OpSub, 10, 4, "10 − 4"},
		{"mul", OpMul, 6, 7, "6 × 7"},
		{"div", OpDiv, 15, 3, "15 ÷ 3"},
		{"mod", OpMod, 10, 3, "10 mod 3"},
		{"pct", OpPct, 25, 80, "25 % of 80"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinOp{Op: tt.op, Left: &Num{tt.left}, Right: &Num{tt.right}}
			if got := b.Format(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBinOp_Format_PEMDAS(t *testing.T) {
	tests := []struct {
		name string
		expr Expr
		want string
	}{
		{
			// (5 + 3) × 2 — lower precedence on left gets parens
			"lower prec left",
			&BinOp{
				Op:    OpMul,
				Left:  &BinOp{Op: OpAdd, Left: &Num{5}, Right: &Num{3}},
				Right: &Num{2},
			},
			"(5 + 3) × 2",
		},
		{
			// 2 × (5 + 3) — lower precedence on right gets parens
			"lower prec right",
			&BinOp{
				Op:   OpMul,
				Left: &Num{2},
				Right: &BinOp{Op: OpAdd, Left: &Num{5}, Right: &Num{3}},
			},
			"2 × (5 + 3)",
		},
		{
			// 5 + 3 × 2 — higher precedence child, no parens needed
			"higher prec no parens",
			&BinOp{
				Op:   OpAdd,
				Left: &Num{5},
				Right: &BinOp{Op: OpMul, Left: &Num{3}, Right: &Num{2}},
			},
			"5 + 3 × 2",
		},
		{
			// 3 × 2 + 5 — higher precedence child on left, no parens
			"higher prec left no parens",
			&BinOp{
				Op:    OpAdd,
				Left:  &BinOp{Op: OpMul, Left: &Num{3}, Right: &Num{2}},
				Right: &Num{5},
			},
			"3 × 2 + 5",
		},
		{
			// 10 − (3 − 1) — right-associativity for sub
			"sub right assoc",
			&BinOp{
				Op:   OpSub,
				Left: &Num{10},
				Right: &BinOp{Op: OpSub, Left: &Num{3}, Right: &Num{1}},
			},
			"10 − (3 − 1)",
		},
		{
			// 10 − 3 + 1 — left child same prec no parens for sub
			"sub left no parens",
			&BinOp{
				Op:    OpSub,
				Left:  &BinOp{Op: OpAdd, Left: &Num{10}, Right: &Num{3}},
				Right: &Num{1},
			},
			"10 + 3 − 1",
		},
		{
			// 12 ÷ (6 ÷ 2) — right-associativity for div
			"div right assoc",
			&BinOp{
				Op:   OpDiv,
				Left: &Num{12},
				Right: &BinOp{Op: OpDiv, Left: &Num{6}, Right: &Num{2}},
			},
			"12 ÷ (6 ÷ 2)",
		},
		{
			// 10 mod (5 mod 3) — right-associativity for mod
			"mod right assoc",
			&BinOp{
				Op:   OpMod,
				Left: &Num{10},
				Right: &BinOp{Op: OpMod, Left: &Num{5}, Right: &Num{3}},
			},
			"10 mod (5 mod 3)",
		},
		{
			// 5 + 3 + 2 — same prec, commutative, no parens on right
			"add same prec right no parens",
			&BinOp{
				Op:   OpAdd,
				Left: &Num{5},
				Right: &BinOp{Op: OpAdd, Left: &Num{3}, Right: &Num{2}},
			},
			"5 + 3 + 2",
		},
		{
			// 2 × 3 × 4 — same prec, commutative, no parens on right
			"mul same prec right no parens",
			&BinOp{
				Op:   OpMul,
				Left: &Num{2},
				Right: &BinOp{Op: OpMul, Left: &Num{3}, Right: &Num{4}},
			},
			"2 × 3 × 4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.expr.Format(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParen_Format(t *testing.T) {
	p := &Paren{Inner: &BinOp{Op: OpAdd, Left: &Num{3}, Right: &Num{4}}}
	want := "(3 + 4)"
	if got := p.Format(); got != want {
		t.Errorf("Paren.Format() = %q, want %q", got, want)
	}
}

func TestUnaryPrefix_Format(t *testing.T) {
	tests := []struct {
		name string
		op   UnaryOp
		val  int
		want string
	}{
		{"sqrt", OpSqrt, 49, "√49"},
		{"cbrt", OpCbrt, 27, "∛27"},
		{"unknown", UnaryOp(99), 5, "5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnaryPrefix{Op: tt.op, Operand: &Num{tt.val}}
			if got := u.Format(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnarySuffix_Format(t *testing.T) {
	tests := []struct {
		name string
		op   UnaryOp
		val  int
		want string
	}{
		{"square", OpSquare, 7, "7²"},
		{"cube", OpCube, 3, "3³"},
		{"factorial", OpFactorial, 5, "5!"},
		{"unknown", UnaryOp(99), 9, "9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnarySuffix{Op: tt.op, Operand: &Num{tt.val}}
			if got := u.Format(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPow_Format(t *testing.T) {
	tests := []struct {
		name      string
		base, exp int
		want      string
	}{
		{"2^3", 2, 3, "2³"},
		{"5^2", 5, 2, "5²"},
		{"10^0", 10, 0, "10⁰"},
		{"2^10", 2, 10, "2¹⁰"},
		{"3^15", 3, 15, "3¹⁵"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pow{Base: &Num{tt.base}, Exp: &Num{tt.exp}}
			if got := p.Format(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToSuperscript(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "⁰"},
		{1, "¹"},
		{2, "²"},
		{3, "³"},
		{4, "⁴"},
		{5, "⁵"},
		{6, "⁶"},
		{7, "⁷"},
		{8, "⁸"},
		{9, "⁹"},
		{10, "¹⁰"},
		{123, "¹²³"},
	}
	for _, tt := range tests {
		if got := toSuperscript(tt.n); got != tt.want {
			t.Errorf("toSuperscript(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Key tests
// ---------------------------------------------------------------------------

func TestNum_Key(t *testing.T) {
	tests := []struct {
		val  int
		want string
	}{
		{0, "0"},
		{5, "5"},
		{-3, "-3"},
	}
	for _, tt := range tests {
		n := &Num{Value: tt.val}
		if got := n.Key(); got != tt.want {
			t.Errorf("Num{%d}.Key() = %q, want %q", tt.val, got, tt.want)
		}
	}
}

func TestBinOp_Key(t *testing.T) {
	tests := []struct {
		name        string
		op          BinOpKind
		left, right int
		want        string
	}{
		{"add", OpAdd, 5, 3, "(+ 5 3)"},
		{"sub", OpSub, 10, 4, "(- 10 4)"},
		{"mul", OpMul, 6, 7, "(* 6 7)"},
		{"div", OpDiv, 15, 3, "(/ 15 3)"},
		{"mod", OpMod, 10, 3, "(% 10 3)"},
		{"pct", OpPct, 25, 80, "(pct 25 80)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinOp{Op: tt.op, Left: &Num{tt.left}, Right: &Num{tt.right}}
			if got := b.Key(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBinOp_Key_Nested(t *testing.T) {
	// 5 + 3 * 2 → "(+ 5 (* 3 2))"
	expr := &BinOp{
		Op:   OpAdd,
		Left: &Num{5},
		Right: &BinOp{Op: OpMul, Left: &Num{3}, Right: &Num{2}},
	}
	want := "(+ 5 (* 3 2))"
	if got := expr.Key(); got != want {
		t.Errorf("nested key = %q, want %q", got, want)
	}
}

func TestParen_Key_Transparent(t *testing.T) {
	// Paren should be transparent — same key as inner
	inner := &BinOp{Op: OpAdd, Left: &Num{5}, Right: &Num{3}}
	paren := &Paren{Inner: inner}
	if paren.Key() != inner.Key() {
		t.Errorf("Paren.Key() = %q, want same as inner %q", paren.Key(), inner.Key())
	}

	// Nested parens should still be transparent
	doubleParen := &Paren{Inner: &Paren{Inner: &Num{42}}}
	if got := doubleParen.Key(); got != "42" {
		t.Errorf("double Paren.Key() = %q, want %q", got, "42")
	}
}

func TestUnaryPrefix_Key(t *testing.T) {
	tests := []struct {
		name string
		op   UnaryOp
		val  int
		want string
	}{
		{"sqrt", OpSqrt, 49, "(sqrt 49)"},
		{"cbrt", OpCbrt, 27, "(cbrt 27)"},
		{"unknown", UnaryOp(99), 5, "(? 5)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnaryPrefix{Op: tt.op, Operand: &Num{tt.val}}
			if got := u.Key(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnarySuffix_Key(t *testing.T) {
	tests := []struct {
		name string
		op   UnaryOp
		val  int
		want string
	}{
		{"square", OpSquare, 7, "(sq 7)"},
		{"cube", OpCube, 3, "(cb 3)"},
		{"factorial", OpFactorial, 5, "(! 5)"},
		{"unknown", UnaryOp(99), 9, "(? 9)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UnarySuffix{Op: tt.op, Operand: &Num{tt.val}}
			if got := u.Key(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPow_Key(t *testing.T) {
	p := &Pow{Base: &Num{2}, Exp: &Num{3}}
	want := "(^ 2 3)"
	if got := p.Key(); got != want {
		t.Errorf("Pow.Key() = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// Interface compliance — ensure all node types implement Expr
// ---------------------------------------------------------------------------

func TestExprInterface(t *testing.T) {
	var nodes []Expr
	nodes = append(nodes,
		&Num{Value: 1},
		&BinOp{Op: OpAdd, Left: &Num{1}, Right: &Num{2}},
		&Paren{Inner: &Num{3}},
		&UnaryPrefix{Op: OpSqrt, Operand: &Num{4}},
		&UnarySuffix{Op: OpSquare, Operand: &Num{5}},
		&Pow{Base: &Num{2}, Exp: &Num{3}},
	)
	for _, n := range nodes {
		// Just verify the methods don't panic
		_ = n.Eval()
		_ = n.Format()
		_ = n.Key()
	}
}
