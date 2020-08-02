package interpreter

import (
	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/token"
)

type Interpreter struct {
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i Interpreter) Evaluate(e ast.Expr) interface{} {
	return e.Accept(i)
}

func (i Interpreter) VisitLiteral(e ast.Literal) interface{} {
	return e.Value
}

func (i Interpreter) VisitGrouping(e ast.Grouping) interface{} {
	return i.Evaluate(e.Expression)
}

func (i Interpreter) VisitUnary(e ast.Unary) interface{} {
	right := i.Evaluate(e.Right)

	switch e.Op.Type {
	case token.Tminus:
		return -right.(float64)
	case token.Tbang:
		return !i.isTruthy(right)
	}

	panic("Unreachable")
}

func (i Interpreter) VisitBinary(e ast.Binary) interface{} {
	right := i.Evaluate(e.Right)
	left := i.Evaluate(e.Left)

	switch e.Op.Type {
	case token.Tminus:
		return left.(float64) - right.(float64)
	case token.Tplus:
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		} else {
			return left.(float64) + right.(float64)
		}
		// error handling?
		panic("both should be string")
	case token.Tstar:
		return left.(float64) * right.(float64)
	case token.Tslash:
		return left.(float64) / right.(float64)
	case token.Tgreater:
		return left.(float64) > right.(float64)
	case token.TgreaterEqual:
		return left.(float64) >= right.(float64)
	case token.Tless:
		return left.(float64) < right.(float64)
	case token.TlessEqual:
		return left.(float64) <= right.(float64)
	case token.TequalEqual:
		// XXX: these arent same as book
		return left == right
	case token.TbangEqual:
		return left != right
	}
	panic("All operators must be one of the above")
}

func (i Interpreter) isTruthy(e interface{}) bool {
	switch v := e.(type) {
	case nil:
		return false
	case bool:
		return v
	}
	return true
}
