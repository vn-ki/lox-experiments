package ast

type Expr interface {
	accept(ExprVisitor) interface{}
}

type ExprVisitor interface {
	visitBinary(Binary) interface{}
	visitGrouping(Grouping) interface{}
	visitLiteral(Literal) interface{}
	visitUnary(Unary) interface{}
}
