package interpreter

import (
	"testing"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/lexer"
	"github.com/vn-ki/go-lox/parser"
)

func parse(src string) ([]ast.Stmt, bool) {
	lexer := lexer.NewLexer(src)
	tokens := lexer.ScanTokens()

	parser := parser.NewParser(tokens)
	return parser.Parse()
}

func TestWhile(t *testing.T) {

	src := `
	var i=0;
	while (i < 10) {
		i = i + 1;
		print i;
	}
	`
	stmts, _ := parse(src)
	// got := ast.NewAstPrinter().PrintStatement(stmts[0])
	// expected := "(- (1) (* (2) (3)))"

	interp := NewInterpreter()
	interp.Interpret(stmts)

	// if got != expected {
	//
	// }
}
