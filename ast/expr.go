package ast

import "github.com/vn-ki/go-lox/token"

type Expr interface {
	Accept(ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinary(Binary) interface{}
	VisitGrouping(Grouping) interface{}
	VisitLiteral(Literal) interface{}
	VisitUnary(Unary) interface{}
	VisitVariable(Evariable) interface{}
	VisitAssign(Eassign) interface{}
	VisitLogical(Elogical) interface{}
}

type Binary struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value interface{}
}

type Unary struct {
	Op    token.Token
	Right Expr
}

type Evariable struct {
	Name token.Token
}

type Eassign struct {
	Name  token.Token
	Value Expr
}

type Elogical struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

func (b Binary) Accept(e ExprVisitor) interface{}    { return e.VisitBinary(b) }
func (g Grouping) Accept(e ExprVisitor) interface{}  { return e.VisitGrouping(g) }
func (l Literal) Accept(e ExprVisitor) interface{}   { return e.VisitLiteral(l) }
func (u Unary) Accept(e ExprVisitor) interface{}     { return e.VisitUnary(u) }
func (u Evariable) Accept(e ExprVisitor) interface{} { return e.VisitVariable(u) }
func (u Eassign) Accept(e ExprVisitor) interface{}   { return e.VisitAssign(u) }
func (u Elogical) Accept(e ExprVisitor) interface{}  { return e.VisitLogical(u) }
