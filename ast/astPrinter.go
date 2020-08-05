package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
	depth int
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{0}
}

// Statements
func (a *AstPrinter) PrintStatement(s Stmt) string {
	return s.Accept(a).(string)
}

func (a *AstPrinter) VisitIf(s Sif) interface{} {
	if s.ElseBranch != nil {
		return fmt.Sprintf(
			"(if %s then %s else %s)",
			a.PrintExpr(s.Condition), a.PrintStatement(s.ThenBranch), a.PrintStatement(s.ElseBranch),
		)
	}
	return fmt.Sprintf("(if %s then %s)", a.PrintExpr(s.Condition), a.PrintStatement(s.ThenBranch))
}

func (a *AstPrinter) VisitPrint(s Sprint) interface{} {
	return a.parenthesize("print", s.Expression)
}
func (a *AstPrinter) VisitExpression(s Sexpression) interface{} {
	return a.PrintExpr(s.Expression)
}

func (a *AstPrinter) VisitVar(s Svar) interface{} {
	if s.Expression == nil {
		return a.parenthesize("var " + s.Name.Lexeme)
	}
	return a.parenthesize("var "+s.Name.Lexeme, s.Expression)
}

func (a *AstPrinter) VisitBlock(s Sblock) interface{} {
	a.depth++
	ret := make([]string, 0)
	ret = append(ret, "(block")
	for _, stmt := range s.Stmts {
		ret = append(ret, stmt.Accept(a).(string))
	}
	ret = append(ret, ")")
	a.depth--
	return strings.Join(ret, "\n"+strings.Repeat(" ", a.depth+1))
}

// Expression

func (a *AstPrinter) PrintExpr(e Expr) string {
	return e.Accept(a).(string)
}

func (a *AstPrinter) VisitVariable(e Evariable) interface{} {
	return a.parenthesize("variable " + e.Name.Lexeme)
}

func (a *AstPrinter) VisitBinary(e Binary) interface{} {
	return a.parenthesize(e.Op.Lexeme, e.Left, e.Right)
}

func (a *AstPrinter) VisitUnary(e Unary) interface{} {
	return a.parenthesize(e.Op.Lexeme, e.Right)
}

func (a *AstPrinter) VisitAssign(e Eassign) interface{} {
	return a.parenthesize("assign "+e.Name.Lexeme, e.Value)
}

func (a *AstPrinter) VisitLiteral(e Literal) interface{} {
	if e.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}

func (a *AstPrinter) VisitGrouping(e Grouping) interface{} {
	return a.parenthesize("group", e.Expression)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	ret := []string{"(", name}
	for _, expr := range exprs {
		ret = append(ret, " ")
		ret = append(ret, expr.Accept(a).(string))
	}
	ret = append(ret, ")")
	return strings.Join(ret, "")
}
