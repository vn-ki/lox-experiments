package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/lexer"
	"github.com/vn-ki/go-lox/parser"
	"github.com/vn-ki/go-lox/token"
)

func report(err error) {

}

func logTokens(tokens []token.Token) {
	for _, tk := range tokens {
		log.Println(tk)
	}
}

func run(src string) error {
	log.Printf("src: '%s'\n", src)
	lexer := lexer.NewLexer(src)
	tokens := lexer.ScanTokens()
	logTokens(tokens)

	parser := parser.NewParser(tokens)
	expr := parser.Parse()
	if expr != nil {
		log.Printf("AST: %s", ast.NewAstPrinter().PrintExpr(expr))
	}
	return nil
}

func runFile(path string) {

}

func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(">> ")
	for scanner.Scan() {
		err := run(scanner.Text())
		if err != nil {
			report(err)
		}
		fmt.Print(">> ")
	}
}

func setupLogger() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	setupLogger()
	argsLen := len(os.Args)
	if argsLen > 2 {
		fmt.Printf("Usage: %s [script]\n", os.Args[0])
		return
	} else if argsLen == 2 {
		//
	} else {
		//
		log.Println("repl")
		runREPL()
	}
}
