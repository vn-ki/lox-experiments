package ast

type Stmt interface {
	Accept(StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitExpression(Sexpression) interface{}
	VisitPrint(Sprint) interface{}
}

type Sexpression struct {
	Expression Expr
}

type Sprint struct {
	Expression Expr
}

func (t Sexpression) Accept(s StmtVisitor) interface{} { return s.VisitExpression(t) }
func (t Sprint) Accept(s StmtVisitor) interface{}      { return s.VisitPrint(t) }
