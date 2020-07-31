package ast

import (
	"github.com/vn-ki/go-lox/token"
)

type Binary struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

func (b Binary) accept(e ExprVisitor) interface{} {
	return e.visitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) accept(e ExprVisitor) interface{} {
	return e.visitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) accept(e ExprVisitor) interface{} {
	return e.visitLiteral(l)
}

type Unary struct {
	Op    token.Token
	Right Expr
}

func (u Unary) accept(e ExprVisitor) interface{} {
	return e.visitUnary(u)
}
