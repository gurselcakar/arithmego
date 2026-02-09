package expr

import "fmt"

// superscript maps digits to their Unicode superscript equivalents.
var superscript = map[rune]rune{
	'0': '⁰', '1': '¹', '2': '²', '3': '³', '4': '⁴',
	'5': '⁵', '6': '⁶', '7': '⁷', '8': '⁸', '9': '⁹',
}

// toSuperscript converts a non-negative integer to superscript digits.
func toSuperscript(n int) string {
	s := fmt.Sprintf("%d", n)
	result := make([]rune, len(s))
	for i, c := range s {
		if sup, ok := superscript[c]; ok {
			result[i] = sup
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func (n *Num) Format() string {
	return fmt.Sprintf("%d", n.Value)
}

func (b *BinOp) Format() string {
	left := b.formatChild(b.Left, true)
	right := b.formatChild(b.Right, false)
	return left + " " + b.Op.Symbol() + " " + right
}

// formatChild wraps a child expression in parens if PEMDAS requires it.
func (b *BinOp) formatChild(child Expr, isLeft bool) string {
	switch c := child.(type) {
	case *BinOp:
		needsParens := false
		if c.Op.Precedence() < b.Op.Precedence() {
			// Lower precedence child needs parens: (5 + 3) × 2
			needsParens = true
		} else if !isLeft && c.Op.Precedence() == b.Op.Precedence() && (b.Op == OpSub || b.Op == OpDiv || b.Op == OpMod) {
			// Right child with same precedence needs parens for non-commutative ops:
			// 10 − (3 − 1), 12 ÷ (6 ÷ 2)
			needsParens = true
		}
		if needsParens {
			return "(" + c.Format() + ")"
		}
		return c.Format()
	default:
		return child.Format()
	}
}

func (p *Paren) Format() string {
	return "(" + p.Inner.Format() + ")"
}

func (u *UnaryPrefix) Format() string {
	switch u.Op {
	case OpSqrt:
		return "√" + u.Operand.Format()
	case OpCbrt:
		return "∛" + u.Operand.Format()
	default:
		return u.Operand.Format()
	}
}

func (u *UnarySuffix) Format() string {
	val := u.Operand.Format()
	switch u.Op {
	case OpSquare:
		return val + "²"
	case OpCube:
		return val + "³"
	case OpFactorial:
		return val + "!"
	default:
		return val
	}
}

func (p *Pow) Format() string {
	base := p.Base.Format()
	exp := p.Exp.Eval()
	return base + toSuperscript(exp)
}
