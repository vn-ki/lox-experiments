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

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens, 0, nil}
}

func (p *Parser) Parse() ast.Expr {
	expr, err := p.experssion()
	if err != nil {
		return nil
	}
	return expr
}

func (p *Parser) experssion() (ast.Expr, error) {
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

/*
Grammar

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
		expr, err := p.experssion()
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
	if p.peek().Type == token {
		p.advance()
	}
	return p.err(p.peek(), message)
}

func (p *Parser) err(tok token.Token, message string) error {
	log.Printf("error at line %d", tok.Line)
	if p.ErrorHandler != nil {
		p.ErrorHandler(tok, message)
	}
	return errors.New(message)
}

func (p *Parser) check(tt token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return tt == p.peek().Type
}

func (p *Parser) match(tts ...token.TokenType) bool {
	for _, tokenType := range tts {
		if p.peek().Type == tokenType {
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
	p.current++
}
