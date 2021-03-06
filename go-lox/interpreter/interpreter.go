package interpreter

import (
	"errors"
	"fmt"
	"log"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/env"
	"github.com/vn-ki/go-lox/token"
)

type Interpreter struct {
	ErrorHandler func(token token.Token, msg string)
	env          *env.Environemnt
	globals      *env.Environemnt
}

type runtimeError struct {
	error
	token token.Token
}
type returnError struct {
	Value interface{}
}

func NewInterpreter() *Interpreter {
	globals := env.NewEnvironment(nil)
	globals.Define("clock", FnClock{})
	globals.DumpEnv()
	return &Interpreter{ErrorHandler: nil, env: globals, globals: globals}
}

func (i *Interpreter) Evaluate(e ast.Expr) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) {
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

	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *Interpreter) execute(s ast.Stmt) {
	s.Accept(i)
}

func (i *Interpreter) VisitExpression(s ast.Sexpression) interface{} {
	return i.Evaluate(s.Expression)
}

func (i *Interpreter) VisitCall(c ast.Ecall) (returnVal interface{}) {
	callee := i.Evaluate(c.Callee)
	args := make([]interface{}, 0)
	for _, arg := range c.Args {
		args = append(args, i.Evaluate(arg))
	}
	if fun, ok := callee.(LoxCallable); ok {
		if len(args) != fun.Arity() {
			i.err("arity doesn't match", c.Paren)
		}
		defer func() {
			if r := recover(); r != nil {
				if w, ok := r.(returnError); ok {
					// log.Printf("return value: %v\n", w.Value)
					// assign the value to the named return value
					returnVal = w.Value
					return
				} else {
					panic(r)
				}
			}
		}()
		fun.Call(i, args)
		return
	}
	i.err("not callable", c.Paren)
	return nil
}

func (i *Interpreter) VisitReturn(r ast.Sreturn) interface{} {
	panic(returnError{i.Evaluate(r.Value)})
}

func (i *Interpreter) VisitFunction(f ast.Sfunction) interface{} {
	log.Printf("getting defined %s\n", f.Name.Lexeme)
	i.env.Define(f.Name.Lexeme, NewLoxFunctionFromAst(f, i.env))
	i.env.DumpEnv()
	return nil
}

func (i *Interpreter) VisitVariable(v ast.Evariable) interface{} {
	val, ok := i.env.Get(v.Name.Lexeme)
	if !ok {
		panic(runtimeError{
			errors.New(fmt.Sprintf("variable '%s' not defined", v.Name.Lexeme)),
			v.Name,
		})
	}
	return val
}

func (i *Interpreter) VisitVar(v ast.Svar) interface{} {
	var val interface{}
	if v.Expression != nil {
		val = i.Evaluate(v.Expression)
	}
	log.Printf("defining '%s' with '%v'", v.Name.Lexeme, val)
	i.env.Define(v.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitWhile(s ast.Swhile) interface{} {
	for i.isTruthy(i.Evaluate(s.Condition)) {
		i.execute(s.Body)
	}
	return nil
}

func (i *Interpreter) VisitIf(s ast.Sif) interface{} {
	cond := i.Evaluate(s.Condition)
	if i.isTruthy(cond) {
		i.execute(s.ThenBranch)
	} else {
		if s.ElseBranch != nil {
			i.execute(s.ElseBranch)
		}
	}
	return nil
}

func (i *Interpreter) VisitPrint(s ast.Sprint) interface{} {
	val := i.Evaluate(s.Expression)
	fmt.Println(val)
	return nil
}

func (i *Interpreter) VisitBlock(s ast.Sblock) interface{} {
	i.ExecuteBlock(s.Stmts, env.NewEnvironment(i.env))
	return nil
}

func (i *Interpreter) ExecuteBlock(stmts []ast.Stmt, env *env.Environemnt) {
	prevEnv := i.env
	defer func() { i.env = prevEnv }()
	i.env = env

	for _, stmt := range stmts {
		stmt.Accept(i)
	}
}

func (i *Interpreter) checkNumberOperand(op token.Token, operand interface{}) {
	if _, ok := operand.(float64); !ok {
		panic(runtimeError{errors.New("the operand should be a number"), op})
	}
}

func (i *Interpreter) checkNumberOperands(op token.Token, left interface{}, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	panic(runtimeError{errors.New("both operands should be number"), op})
}

func (i *Interpreter) VisitAssign(e ast.Eassign) interface{} {
	if i.env.Assign(e.Name.Lexeme, i.Evaluate(e.Value)) {
		// i.env.DumpEnv(0)
		return e.Value
	}
	panic(runtimeError{errors.New("Undefined variable"), e.Name})
}

func (i *Interpreter) VisitLiteral(e ast.Literal) interface{} {
	return e.Value
}

func (i *Interpreter) VisitGrouping(e ast.Grouping) interface{} {
	return i.Evaluate(e.Expression)
}

func (i *Interpreter) VisitUnary(e ast.Unary) interface{} {
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

func (i *Interpreter) VisitLogical(e ast.Elogical) interface{} {
	left := i.Evaluate(e.Left)
	switch e.Op.Type {
	case token.Tand:
		if !i.isTruthy(left) {
			return left
		}
	case token.Tfalse:
		if i.isTruthy(left) {
			return left
		}
	}
	return i.Evaluate(e.Right)
}

func (i *Interpreter) VisitBinary(e ast.Binary) interface{} {
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

func (i *Interpreter) isTruthy(e interface{}) bool {
	switch v := e.(type) {
	case nil:
		return false
	case bool:
		return v
	}
	return true
}

func (i *Interpreter) err(msg string, token token.Token) {
	panic(runtimeError{errors.New(msg), token})
}
