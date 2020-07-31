package parser

import (
	"testing"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/lexer"
)

func TestParser(t *testing.T) {
	src := "1-2*3"
	lexer := lexer.NewLexer(src)
	tokens := lexer.ScanTokens()

	parser := NewParser(tokens)
	expr := parser.Parse()
	got := ast.NewAstPrinter().PrintExpr(expr)
	expected := "(- (1) (* (2) (3)))"
	if got != expected {

	}
}
