package expr

// Expr represents a node in an expression tree.
type Expr interface {
	Eval() int
	Format() string
	Key() string
}

// BinOpKind identifies binary operators.
type BinOpKind int

const (
	OpAdd BinOpKind = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPct // "% of"
)

// Precedence returns the operator precedence (1=low, 2=high).
func (op BinOpKind) Precedence() int {
	switch op {
	case OpAdd, OpSub:
		return 1
	case OpMul, OpDiv, OpMod, OpPct:
		return 2
	default:
		return 0
	}
}

// Symbol returns the display symbol for the operator.
func (op BinOpKind) Symbol() string {
	switch op {
	case OpAdd:
		return "+"
	case OpSub:
		return "−"
	case OpMul:
		return "×"
	case OpDiv:
		return "÷"
	case OpMod:
		return "mod"
	case OpPct:
		return "% of"
	default:
		return "?"
	}
}

// KeySymbol returns the canonical symbol for dedup keys.
func (op BinOpKind) KeySymbol() string {
	switch op {
	case OpAdd:
		return "+"
	case OpSub:
		return "-"
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpMod:
		return "%"
	case OpPct:
		return "pct"
	default:
		return "?"
	}
}

// Num is a leaf node representing a literal integer.
type Num struct {
	Value int
}

// BinOp is a binary operation node.
type BinOp struct {
	Op          BinOpKind
	Left, Right Expr
}

// Paren wraps an expression in explicit parentheses for display.
// Semantically identical to Inner for evaluation and dedup.
type Paren struct {
	Inner Expr
}

// UnaryOp identifies unary operators.
type UnaryOp int

const (
	OpSqrt UnaryOp = iota
	OpCbrt
	OpSquare
	OpCube
	OpFactorial
)

// UnaryPrefix is a prefix unary operator (e.g., √49, ∛27).
type UnaryPrefix struct {
	Op      UnaryOp
	Operand Expr
}

// UnarySuffix is a suffix unary operator (e.g., 7², 5!, 3³).
type UnarySuffix struct {
	Op      UnaryOp
	Operand Expr
}

// Pow represents exponentiation (base^exp) with superscript display.
type Pow struct {
	Base, Exp Expr
}
