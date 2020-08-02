package ast

import (
	"github.com/vn-ki/go-lox/token"
)

type Binary struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

func (b Binary) Accept(e ExprVisitor) interface{} {
	return e.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(e ExprVisitor) interface{} {
	return e.VisitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(e ExprVisitor) interface{} {
	return e.VisitLiteral(l)
}

type Unary struct {
	Op    token.Token
	Right Expr
}

func (u Unary) Accept(e ExprVisitor) interface{} {
	return e.VisitUnary(u)
}
