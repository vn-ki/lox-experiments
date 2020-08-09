package interpreter

import (
	"fmt"
	"time"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/env"
	"github.com/vn-ki/go-lox/token"
)

type LoxCallable interface {
	Call(*Interpreter, []interface{}) interface{}
	Arity() int
}

/// Native Function: clock
type FnClock struct{}

func (f FnClock) Arity() int { return 0 }

func (f FnClock) Call(_ *Interpreter, _ []interface{}) interface{} {
	return float64(time.Now().UnixNano())
}

func (f FnClock) String() string { return "<clock native fn>" }

/// Lox Function

type LoxFunction struct {
	Name   token.Token
	Params []token.Token
	Body   []ast.Stmt
	Env    *env.Environemnt
}

func NewLoxFunctionFromAst(f ast.Sfunction, env *env.Environemnt) LoxFunction {
	return LoxFunction{Name: f.Name, Params: f.Params, Body: f.Body, Env: env}
}

func (f LoxFunction) Arity() int { return len(f.Params) }

func (f LoxFunction) Call(i *Interpreter, args []interface{}) interface{} {
	env := env.NewEnvironment(f.Env)

	for idx, arg := range args {
		env.Define(f.Params[idx].Lexeme, arg)
	}
	i.ExecuteBlock(f.Body, env)

	return nil
}

func (f LoxFunction) String() string { return fmt.Sprintf("<fn %s>", f.Name.Lexeme) }
