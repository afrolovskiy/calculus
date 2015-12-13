package calculus

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

const (
	TokenValue = iota
	TokenAdd
	TokenSub
	TokenMul
	TokenDiv
	TokenLeftParen
	TokenRightParen
)

var opers = map[token.Token]TokenType{
	token.ADD: TokenAdd,
	token.SUB: TokenSub,
	token.MUL: TokenMul,
	token.QUO: TokenDiv,
}

var (
	ErrUnknown      = errors.New("Unknown literal")
	ErrInvalidExpr  = errors.New("Invalid expression")
	ErrUnknownToken = errors.New("Unknown token")
	ErrInvalidRPN   = errors.New("Invalid RPN")
)

type tokenStack struct{ Stack }

func (s *tokenStack) Peek() *Token  { return s.Stack.Peek().(*Token) }
func (s *tokenStack) Push(v *Token) { s.Stack.Push(v) }
func (s *tokenStack) Pop() *Token   { return s.Stack.Pop().(*Token) }

type floatStack struct{ Stack }

func (s *floatStack) Peek() float64  { return s.Stack.Peek().(float64) }
func (s *floatStack) Push(v float64) { s.Stack.Push(v) }
func (s *floatStack) Pop() float64   { return s.Stack.Pop().(float64) }

type TokenType int

type Token struct {
	Type  TokenType
	Value float64
}

type RPN struct {
	tokens []*Token
}

// NewRPN build RPN instance from expression.
func NewRPN(expr string) (*RPN, error) {
	rpn := RPN{}

	tree, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, ErrInvalidExpr
	}

	tokens, err := parseExpr(tree)
	if err != nil {
		return nil, err
	}

	rpntokens, err := convertInfixToPostfix(tokens)
	if err != nil {
		return nil, err
	}

	rpn.tokens = rpntokens
	return &rpn, nil
}

func parseExpr(tree ast.Expr) ([]*Token, error) {
	var tokens []*Token

	switch t := tree.(type) {
	case *ast.BinaryExpr:
		ot, known := opers[t.Op]
		if !known {
			return nil, ErrUnknown
		}

		x, err := parseExpr(t.X)
		if err != nil {
			return nil, err
		}

		y, err := parseExpr(t.Y)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, x...)
		tokens = append(tokens, &Token{Type: ot})
		tokens = append(tokens, y...)

	case *ast.BasicLit:
		switch {
		case t.Kind == token.INT:
			v, err := strconv.Atoi(t.Value)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, &Token{Type: TokenValue, Value: float64(v)})

		case t.Kind == token.FLOAT:
			v, err := strconv.ParseFloat(t.Value, 64)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, &Token{Type: TokenValue, Value: float64(v)})

		default:
			return nil, ErrUnknown
		}

	case *ast.ParenExpr:
		ops, err := parseExpr(t.X)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, &Token{Type: TokenLeftParen})
		tokens = append(tokens, ops...)
		tokens = append(tokens, &Token{Type: TokenRightParen})

	default:
		return nil, ErrUnknown
	}

	return tokens, nil
}

func convertInfixToPostfix(tokens []*Token) ([]*Token, error) {
	var rpn []*Token
	var stack tokenStack

	for _, t := range tokens {
		switch {
		case t.Type == TokenValue:
			rpn = append(rpn, t)

		case t.Type == TokenAdd || t.Type == TokenSub:
			for !stack.Empty() && (stack.Peek().Type == TokenAdd || stack.Peek().Type == TokenSub || stack.Peek().Type == TokenMul || stack.Peek().Type == TokenDiv) {
				rpn = append(rpn, stack.Pop())
			}
			stack.Push(t)

		case t.Type == TokenMul || t.Type == TokenDiv:
			for !stack.Empty() && (stack.Peek().Type == TokenMul || stack.Peek().Type == TokenDiv) {
				rpn = append(rpn, stack.Pop())
			}
			stack.Push(t)

		case t.Type == TokenLeftParen:
			stack.Push(t)

		case t.Type == TokenRightParen:
			for !stack.Empty() && stack.Peek().Type != TokenLeftParen {
				rpn = append(rpn, stack.Pop())
			}
			if !stack.Empty() {
				stack.Pop()
			}

		default:
			return nil, ErrUnknownToken
		}
	}

	for !stack.Empty() {
		rpn = append(rpn, stack.Pop())
	}

	return rpn, nil
}

// Calculate evaluates rpn expression.
func (r *RPN) Calculate() (float64, error) {
	var stack floatStack

	for _, t := range r.tokens {
		switch {
		case t.Type == TokenAdd:
			if stack.Len() < 2 {
				return 0, ErrInvalidRPN
			}
			x1 := stack.Pop()
			x2 := stack.Pop()
			stack.Push(x1 + x2)

		case t.Type == TokenSub:
			if stack.Len() < 2 {
				return 0, ErrInvalidRPN
			}
			x1 := stack.Pop()
			x2 := stack.Pop()
			stack.Push(x2 - x1)

		case t.Type == TokenMul:
			if stack.Len() < 2 {
				return 0, ErrInvalidRPN
			}
			x1 := stack.Pop()
			x2 := stack.Pop()
			stack.Push(x1 * x2)

		case t.Type == TokenDiv:
			if stack.Len() < 2 {
				return 0, ErrInvalidRPN
			}

			x1 := stack.Pop()
			if x1 == 0 {
				return 0, ErrInvalidRPN
			}

			x2 := stack.Pop()
			stack.Push(x2 / x1)

		case t.Type == TokenValue:
			stack.Push(t.Value)

		default:
			return 0, ErrInvalidRPN
		}
	}

	if stack.Len() != 1 {
		return 0, ErrInvalidRPN
	}

	return stack.Pop(), nil
}
