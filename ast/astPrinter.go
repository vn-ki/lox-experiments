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
	return e.accept(a).(string)
}

func (a AstPrinter) visitBinary(e Binary) interface{} {
	return a.parenthesize(e.Op.Lexeme, e.Left, e.Right)
}

func (a AstPrinter) visitUnary(e Unary) interface{} {
	return a.parenthesize(e.Op.Lexeme, e.Right)
}

func (a AstPrinter) visitLiteral(e Literal) interface{} {
	if e.Value == nil {
		return "nil"
	}
	return a.parenthesize(fmt.Sprintf("%v", e.Value))
}

func (a AstPrinter) visitGrouping(e Grouping) interface{} {
	return a.parenthesize("group", e.Expression)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	ret := []string{"(", name}
	for _, expr := range exprs {
		ret = append(ret, " ")
		ret = append(ret, expr.accept(a).(string))
	}
	ret = append(ret, ")")
	return strings.Join(ret, "")
}
