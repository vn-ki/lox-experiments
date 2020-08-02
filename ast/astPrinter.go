package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}

func (a AstPrinter) PrintExpr(e Expr) string {
	return e.Accept(a).(string)
}

func (a AstPrinter) VisitBinary(e Binary) interface{} {
	return a.parenthesize(e.Op.Lexeme, e.Left, e.Right)
}

func (a AstPrinter) VisitUnary(e Unary) interface{} {
	return a.parenthesize(e.Op.Lexeme, e.Right)
}

func (a AstPrinter) VisitLiteral(e Literal) interface{} {
	if e.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}

func (a AstPrinter) VisitGrouping(e Grouping) interface{} {
	return a.parenthesize("group", e.Expression)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	ret := []string{"(", name}
	for _, expr := range exprs {
		ret = append(ret, " ")
		ret = append(ret, expr.Accept(a).(string))
	}
	ret = append(ret, ")")
	return strings.Join(ret, "")
}
