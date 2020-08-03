package ast

import "github.com/vn-ki/go-lox/token"

type Stmt interface {
	Accept(StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitExpression(Sexpression) interface{}
	VisitPrint(Sprint) interface{}
	VisitVar(Svar) interface{}
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

func (t Sexpression) Accept(s StmtVisitor) interface{} { return s.VisitExpression(t) }
func (t Sprint) Accept(s StmtVisitor) interface{}      { return s.VisitPrint(t) }
func (t Svar) Accept(s StmtVisitor) interface{}        { return s.VisitVar(t) }
