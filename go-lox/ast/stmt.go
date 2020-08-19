package ast

import "github.com/vn-ki/go-lox/token"

type Stmt interface {
	Accept(StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitExpression(Sexpression) interface{}
	VisitPrint(Sprint) interface{}
	VisitVar(Svar) interface{}
	VisitBlock(Sblock) interface{}
	VisitIf(Sif) interface{}
	VisitWhile(Swhile) interface{}
	VisitFunction(Sfunction) interface{}
	VisitReturn(Sreturn) interface{}
}

type Sexpression struct {
	Expression Expr
}

type Sprint struct {
	Expression Expr
}

type Svar struct {
	Name       token.Token
	Expression Expr
}

type Sblock struct {
	Stmts []Stmt
}

type Sif struct {
	ThenBranch Stmt
	ElseBranch Stmt
	Condition  Expr
}

type Swhile struct {
	Condition Expr
	Body      Stmt
}

type Sfunction struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

type Sreturn struct {
	Value   Expr
	Keyword token.Token
}

func (t Sexpression) Accept(s StmtVisitor) interface{} { return s.VisitExpression(t) }
func (t Sprint) Accept(s StmtVisitor) interface{}      { return s.VisitPrint(t) }
func (t Svar) Accept(s StmtVisitor) interface{}        { return s.VisitVar(t) }
func (t Sblock) Accept(s StmtVisitor) interface{}      { return s.VisitBlock(t) }
func (t Sif) Accept(s StmtVisitor) interface{}         { return s.VisitIf(t) }
func (t Swhile) Accept(s StmtVisitor) interface{}      { return s.VisitWhile(t) }
func (t Sfunction) Accept(s StmtVisitor) interface{}   { return s.VisitFunction(t) }
func (t Sreturn) Accept(s StmtVisitor) interface{}     { return s.VisitReturn(t) }
