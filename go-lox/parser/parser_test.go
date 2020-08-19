package parser

import (
	"testing"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/lexer"
)

func parse(src string) ([]ast.Stmt, bool) {
	lexer := lexer.NewLexer(src)
	tokens := lexer.ScanTokens()

	parser := NewParser(tokens)
	return parser.Parse()
}

func TestParser(t *testing.T) {
	src := "1-2*3;"
	stmts, _ := parse(src)
	got := ast.NewAstPrinter().PrintStatement(stmts[0])
	expected := "(- (1) (* (2) (3)))"
	if got != expected {

	}
}

func TestParserDoubleSemicolon(t *testing.T) {
	src := ";;"
	stmts, _ := parse(src)
	got := ast.NewAstPrinter().PrintStatement(stmts[0])
	expected := "(- (1) (* (2) (3)))"
	if got != expected {

	}
}

func TestParserErr(t *testing.T) {
	src := "1-"
	stmts, hadError := parse(src)
	if hadError {

	}
	_ = ast.NewAstPrinter().PrintStatement(stmts[0])
}
