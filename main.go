package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/interpreter"
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

func run(src string, interp *interpreter.Interpreter) error {
	log.Printf("src: '%s'\n", src)
	lexer := lexer.NewLexer(src)
	tokens := lexer.ScanTokens()
	logTokens(tokens)

	parser := parser.NewParser(tokens)
	// parserErrOccured := false
	parser.ErrorHandler = func(tok token.Token, msg string) {
		fmt.Println("parser error")
	}
	expr, hadError := parser.Parse()

	if !hadError {
		for _, stmt := range expr {
			log.Printf("AST: %s", ast.NewAstPrinter().PrintStatement(stmt))
		}
		interp.Interpret(expr)
		// log.Printf("Evaluated value: %v\n", val)
	}

	return nil
}

func runFile(path string) {
	interp := interpreter.NewInterpreter()
	src, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	run(string(src), interp)
}

func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	interp := interpreter.NewInterpreter()

	fmt.Print(">> ")
	for scanner.Scan() {
		err := run(scanner.Text(), interp)
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
		runFile(os.Args[1])
	} else {
		//
		log.Println("repl")
		runREPL()
	}
}
