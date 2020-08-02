package interpreter

import (
	"errors"
	"log"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/token"
)

type Interpreter struct {
	ErrorHandler func(token token.Token, msg string)
}

type runtimeError struct {
	error
	token token.Token
}

func NewInterpreter() Interpreter {
	return Interpreter{nil}
}

func (i Interpreter) Evaluate(e ast.Expr) interface{} {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(runtimeError); ok {
				_ = re
				log.Printf("RuntimeError: at %d: %s\n", re.token.Line, re.Error())
				if i.ErrorHandler != nil {
					i.ErrorHandler(re.token, re.Error())
				}
			} else {
				panic(r)
			}
		}
	}()
	return e.Accept(i)
}

func (i Interpreter) checkNumberOperand(op token.Token, operand interface{}) {
	if _, ok := operand.(float64); !ok {
		panic(runtimeError{errors.New("the operand should be a number"), op})
	}
}

func (i Interpreter) checkNumberOperands(op token.Token, left interface{}, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	panic(runtimeError{errors.New("both operands should be number"), op})
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
		i.checkNumberOperand(e.Op, right)
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
		i.checkNumberOperands(e.Op, left, right)
		return left.(float64) - right.(float64)
	case token.Tplus:
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		} else if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		panic(runtimeError{errors.New("Both operands must be either string or number"), e.Op})
	case token.Tstar:
		i.checkNumberOperands(e.Op, left, right)
		return left.(float64) * right.(float64)
	case token.Tslash:
		i.checkNumberOperands(e.Op, left, right)
		return left.(float64) / right.(float64)
	case token.Tgreater:
		i.checkNumberOperands(e.Op, left, right)
		return left.(float64) > right.(float64)
	case token.TgreaterEqual:
		i.checkNumberOperands(e.Op, left, right)
		return left.(float64) >= right.(float64)
	case token.Tless:
		i.checkNumberOperands(e.Op, left, right)
		return left.(float64) < right.(float64)
	case token.TlessEqual:
		i.checkNumberOperands(e.Op, left, right)
		return left.(float64) <= right.(float64)
	case token.TequalEqual:
		i.checkNumberOperands(e.Op, left, right)
		// XXX: these arent same as book
		return left == right
	case token.TbangEqual:
		i.checkNumberOperands(e.Op, left, right)
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
