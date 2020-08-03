package parser

import (
	"errors"
	"log"

	"github.com/vn-ki/go-lox/ast"
	"github.com/vn-ki/go-lox/token"
)

type Parser struct {
	tokens       []token.Token
	current      int
	ErrorHandler func(token.Token, string)
}

type parserError struct {
	error
	tok token.Token
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens, 0, nil}
}

func (p *Parser) Parse() ([]ast.Stmt, bool) {
	stmts := make([]ast.Stmt, 0)
	hadError := false
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			if w, ok := err.(parserError); ok {
				hadError = true
				if p.ErrorHandler != nil {
					p.ErrorHandler(w.tok, w.Error())
				}
			} else {
				panic(err)
			}
			continue
		}
		stmts = append(stmts, stmt)
	}
	return stmts, hadError
}

/*
Grammar

program   → statement* EOF ;

statement → exprStmt
		  | printStmt ;

exprStmt  → expression ";" ;
printStmt → "print" expression ";" ;

expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
			   | primary ;
primary        → NUMBER | STRING | "false" | "true" | "nil"
			   | "(" expression ")" ;
*/

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.Tprint) {
		return p.printStatement()
	}
	return p.exprStatement()
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(token.Tsemicolon, "Expected semicolon")
	if err != nil {
		return nil, err
	}
	return ast.Sprint{Expression: expr}, nil
}

func (p *Parser) exprStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(token.Tsemicolon, "Expected semicolon")
	if err != nil {
		return nil, err
	}
	return ast.Sexpression{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.TbangEqual, token.TequalEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.addition()
	if err != nil {
		return nil, err
	}

	for p.match(token.Tless, token.TlessEqual, token.Tgreater, token.TgreaterEqual) {
		op := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Op: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) addition() (ast.Expr, error) {
	expr, err := p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(token.Tplus, token.Tminus) {
		op := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Op: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) multiplication() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.Tstar, token.Tslash) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = ast.Binary{Left: expr, Op: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.Tbang, token.Tminus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.Unary{Op: op, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.Tfalse) {
		return ast.Literal{Value: false}, nil
	}
	if p.match(token.Ttrue) {
		return ast.Literal{Value: true}, nil
	}
	if p.match(token.Tnil) {
		return ast.Literal{Value: nil}, nil
	}

	if p.match(token.Tnumber, token.Tstring) {
		return ast.Literal{Value: p.previous().Literal}, nil
	}

	if p.match(token.TleftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		err = p.consume(token.TrightParen, "Expected ')' after experssion.")
		if err != nil {
			return nil, err
		}
		return ast.Grouping{Expression: expr}, nil
	}

	return nil, p.err(p.peek(), "Expected expression")
}

func (p *Parser) consume(token token.TokenType, message string) error {
	if p.check(token) {
		p.advance()
		return nil
	}
	return p.err(p.peek(), message)
}

func (p *Parser) err(tok token.Token, message string) error {
	log.Printf("error at line %d %v: %s", tok.Line, tok, message)
	// XXX: Should this be here? Putting here to avoid the double semicolon infinte loop
	p.advance()
	return parserError{errors.New(message), tok}
}

func (p *Parser) check(tt token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return tt == p.peek().Type
}

func (p *Parser) match(tts ...token.TokenType) bool {
	for _, tokenType := range tts {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) previous() token.Token {
	if p.current < 1 {
		panic("this shouldnt happen")
	}
	return p.tokens[p.current-1]
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	if p.peek().Type == token.Teof {
		return true
	}
	return false
}

func (p *Parser) advance() {
	// XXX: this is here so that when advance is called from err
	// it doesnt go out of bounds. Check with the book
	if p.isAtEnd() {
		return
	}
	p.current++
}
