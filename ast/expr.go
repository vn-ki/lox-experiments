package ast

type Expr interface {
	Accept(ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinary(Binary) interface{}
	VisitGrouping(Grouping) interface{}
	VisitLiteral(Literal) interface{}
	VisitUnary(Unary) interface{}
}
