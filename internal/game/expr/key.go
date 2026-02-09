package expr

import "fmt"

// Key returns a canonical prefix-notation string for dedup.
// Examples:
//   - Num{5} → "5"
//   - BinOp{+, 5, 3×2} → "(+ 5 (* 3 2))"
//   - Paren wrapping is ignored (display-only, not semantic)

func (n *Num) Key() string {
	return fmt.Sprintf("%d", n.Value)
}

func (b *BinOp) Key() string {
	return fmt.Sprintf("(%s %s %s)", b.Op.KeySymbol(), b.Left.Key(), b.Right.Key())
}

func (p *Paren) Key() string {
	return p.Inner.Key()
}

func (u *UnaryPrefix) Key() string {
	var op string
	switch u.Op {
	case OpSqrt:
		op = "sqrt"
	case OpCbrt:
		op = "cbrt"
	default:
		op = "?"
	}
	return fmt.Sprintf("(%s %s)", op, u.Operand.Key())
}

func (u *UnarySuffix) Key() string {
	var op string
	switch u.Op {
	case OpSquare:
		op = "sq"
	case OpCube:
		op = "cb"
	case OpFactorial:
		op = "!"
	default:
		op = "?"
	}
	return fmt.Sprintf("(%s %s)", op, u.Operand.Key())
}

func (p *Pow) Key() string {
	return fmt.Sprintf("(^ %s %s)", p.Base.Key(), p.Exp.Key())
}
