package calculus

import (
	"reflect"
	"testing"
)

var (
	one = &Token{Type: TokenValue, Value: 1}
	two = &Token{Type: TokenValue, Value: 2}

	add  = &Token{Type: TokenAdd}
	sub  = &Token{Type: TokenSub}
	mul  = &Token{Type: TokenMul}
	div  = &Token{Type: TokenDiv}
	lpar = &Token{Type: TokenLeftParen}
	rpar = &Token{Type: TokenRightParen}
)

func TestNewRpn(t *testing.T) {
	rpn := func(a ...*Token) *RPN { return &RPN{a} }

	cases := []struct {
		expr string
		want *RPN
	}{
		// 1 + 2 -> 1 2 +
		{"1+2", rpn(one, two, add)},
		{"1 +  2", rpn(one, two, add)},

		// 1 - 2 -> 1 2 -
		{"1 - 2", rpn(one, two, sub)},

		// 1 * 2 -> 1 2 *
		{"1 * 2", rpn(one, two, mul)},

		// 1 / 2 -> 1 2 /
		{"1 / 2", rpn(one, two, div)},

		// (1 + 2) -> 1 2 +
		{"(1 + 2)", rpn(one, two, add)},

		// (1 + 2) * 1 + (1 - 2) -> 1 2 + 1 * 1 2 - +
		{"(1 + 2) * 1 + (1 - 2)", rpn(one, two, add, one, mul, one, two, sub, add)},
	}

	for _, c := range cases {
		rpn, _ := NewRPN(c.expr)
		if reflect.DeepEqual(rpn, c.want) != true {
			t.Errorf("NewRPN(%v) = %#v; want %#v", c.expr, rpn, c.want)
		}
	}
}

func TestRPNCalculate(t *testing.T) {
	rpn := func(a ...*Token) *RPN { return &RPN{a} }

	cases := []struct {
		rpn  *RPN
		want float32
	}{
		// 1 + 2 -> 3
		{rpn(one, two, add), 3},

		// 1 * 2 -> 2
		{rpn(one, two, mul), 2},

		// 1 / 2 -> 0.5
		{rpn(one, two, div), 0.5},

		// 1 - 2 -> -1
		{rpn(one, two, sub), -1},

		// (1 + 2) * 1 + (1 - 2) -> 1
		{rpn(one, two, add, one, mul, one, two, sub, add), 2},
	}

	for _, c := range cases {
		res, _ := c.rpn.Calculate()
		if res != float64(c.want) {
			t.Errorf("RPN.Calculate(%v) = %#v; want %#v", c.rpn, res, c.want)
		}
	}
}
