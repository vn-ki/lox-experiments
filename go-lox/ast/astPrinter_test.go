package ast

import (
	"testing"

	"github.com/vn-ki/go-lox/token"
)

func TestAstPrinter(t *testing.T) {
	expr := Binary{
		Unary{
			token.Token{token.Tminus, "-", nil, 1},
			Literal{123},
		},
		token.Token{token.Tstar, "*", nil, 1},
		Grouping{
			Literal{45.67},
		},
	}
	astPrinter := NewAstPrinter()
	got := astPrinter.PrintExpr(expr)
	expected := "(* (- (123)) (group (45.67)))"
	if got != expected {
		t.Errorf("Expected: %s, got: %s", expected, got)
	}
}
