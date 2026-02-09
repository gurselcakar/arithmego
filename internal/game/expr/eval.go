package expr

import "math"

func (n *Num) Eval() int { return n.Value }

func (b *BinOp) Eval() int {
	left := b.Left.Eval()
	right := b.Right.Eval()
	switch b.Op {
	case OpAdd:
		return left + right
	case OpSub:
		return left - right
	case OpMul:
		return left * right
	case OpDiv:
		if right == 0 {
			return 0
		}
		return left / right
	case OpMod:
		if right == 0 {
			return 0
		}
		return left % right
	case OpPct:
		return (left * right) / 100
	default:
		return 0
	}
}

func (p *Paren) Eval() int { return p.Inner.Eval() }

func (u *UnaryPrefix) Eval() int {
	val := u.Operand.Eval()
	switch u.Op {
	case OpSqrt:
		return int(math.Sqrt(float64(val)))
	case OpCbrt:
		return int(math.Cbrt(float64(val)))
	default:
		return val
	}
}

func (u *UnarySuffix) Eval() int {
	val := u.Operand.Eval()
	switch u.Op {
	case OpSquare:
		return val * val
	case OpCube:
		return val * val * val
	case OpFactorial:
		return factorial(val)
	default:
		return val
	}
}

func (p *Pow) Eval() int {
	base := p.Base.Eval()
	exp := p.Exp.Eval()
	return intPow(base, exp)
}

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

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}
