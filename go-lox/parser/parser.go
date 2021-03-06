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
		stmt, err := p.declaration()
		if err != nil {
			if w, ok := err.(parserError); ok {
				hadError = true
				if p.ErrorHandler != nil {
					p.ErrorHandler(w.tok, w.Error())
				}
				p.synchorize()
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

program     → declaration* EOF ;

declaration → varDecl
			| funcDecl
			| statement ;

funcDecl -> "fun" function;
function -> IDENTIFIER "(" parameters? ")" block ;
parameters -> IDENTIFIER ( "," IDENTIFIER )* ;

statement   → exprStmt
			| ifStmt
            | printStmt
			| whileStmt
			| returnStmt
			| forStmt
			| block ;

returnStmt -> RETURN expression? ";" ;

forStmt -> "for" "(" (varDecl | exprStmt | ";")
					expression? ;
					expression? ")" statement ;

ifStmt -> "if" "(" expression ")" statement ( "else" statement )? ;

whileStmt -> "while" "(" expression ")" statement;

block  -> "{" declaration* "}";

varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;

exprStmt  → expression ";" ;
printStmt → "print" expression ";" ;

expression     → assignment ;
assignment -> IDENTIFIER "=" assignment
			| logic_or;
logic_or -> logic_and ("or" logic_and)* ;
logic_and -> equality ("and" equality)* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
			   | call ;

call -> primary ( "(" arguments? ")" )* ;
arguments -> expression ("," expression)* ;
primary        → NUMBER | STRING | "false" | "true" | "nil"
			   | "(" expression ")"
			   | IDENTIFIER;
*/

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(token.Tvar) {
		return p.varDecl()
	}
	if p.match(token.Tfun) {
		return p.funcDecl()
	}
	return p.statement()
}

func (p *Parser) funcDecl() (ast.Stmt, error) {
	if name := p.peek(); p.match(token.Tidentifier) {
		err := p.consume(token.TleftParen, "expected ( after function name")
		if err != nil {
			return nil, err
		}
		params := make([]token.Token, 0)
		if !p.match(token.TrightParen) {
			for {
				if p.check(token.Tidentifier) {
					params = append(params, p.peek())
					p.advance()
				} else {
					return nil, p.err(p.peek(), "expected identifier")
				}
				if !p.match(token.Tcomma) {
					break
				}
			}
			err = p.consume(token.TrightParen, "expected ) after parameters")
			if err != nil {
				return nil, err
			}
		}
		p.consume(token.TleftBrace, "Expected { before body")
		body, err := p.block()
		if err != nil {
			return nil, err
		}
		return ast.Sfunction{Name: name, Params: params, Body: body.(ast.Sblock).Stmts}, nil
	}
	return nil, p.err(p.peek(), "expected funciton indentifier")
}

func (p *Parser) varDecl() (ast.Stmt, error) {
	iden := p.peek()

	// consume current token, and confirm it is an identifier
	err := p.consume(token.Tidentifier, "Expected identifier")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.Tequal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	return ast.Svar{Name: iden, Expression: initializer}, p.consume(token.Tsemicolon, "Expected a semicolon")

}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.Tprint) {
		return p.printStatement()
	}
	if p.match(token.Tif) {
		return p.ifStmt()
	}
	if p.match(token.Twhile) {
		return p.whileStmt()
	}
	if p.match(token.Tfor) {
		return p.forStmt()
	}
	if p.match(token.TleftBrace) {
		return p.block()
	}
	if p.match(token.Treturn) {
		return p.returnStmt()
	}
	return p.exprStatement()
}

func (p *Parser) returnStmt() (ast.Stmt, error) {
	keyword := p.previous()
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	return ast.Sreturn{Value: value, Keyword: keyword}, p.consume(token.Tsemicolon, "Expected semicolon after return")
}

func (p *Parser) forStmt() (ast.Stmt, error) {
	var initializer ast.Stmt
	var err error

	p.consume(token.TleftParen, "Expected ( before expression")
	if p.match(token.Tsemicolon) {
	} else if p.match(token.Tvar) {
		initializer, err = p.varDecl()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.exprStatement()
		if err != nil {
			return nil, err
		}
	}
	// initializer done

	// condition check
	var cond ast.Expr
	if !p.match(token.Tsemicolon) {
		cond, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.Tsemicolon, "expected semicolon after loop condition")

	// increment
	var increment ast.Expr
	if !p.match(token.TrightParen) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.TrightParen, "Expected ) after loop")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	whileStmt := ast.Swhile{
		Condition: cond,
		Body:      ast.Sblock{Stmts: []ast.Stmt{body, ast.Sexpression{Expression: increment}}},
	}

	if initializer != nil {
		return ast.Sblock{Stmts: []ast.Stmt{initializer, whileStmt}}, nil
	}
	return whileStmt, nil
}

func (p *Parser) whileStmt() (ast.Stmt, error) {
	p.consume(token.TleftParen, "Expected ( before expression")
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.TrightParen, "Expected ) after expression")

	body, err := p.statement()
	return ast.Swhile{Body: body, Condition: cond}, err
}

func (p *Parser) ifStmt() (ast.Stmt, error) {
	p.consume(token.TleftParen, "Expected ( before expression")

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.TrightParen, "Expected ) after expression")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	if p.match(token.Telse) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return ast.Sif{ThenBranch: thenBranch, ElseBranch: elseBranch, Condition: cond}, nil
}

func (p *Parser) block() (ast.Stmt, error) {
	stmts := make([]ast.Stmt, 0)
	for !p.check(token.TrightBrace) {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return ast.Sblock{Stmts: stmts}, p.consume(token.TrightBrace, "Expected closing brace")
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
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.logic_or()
	if err != nil {
		return nil, err
	}

	if p.match(token.Tequal) {
		if w, ok := expr.(ast.Evariable); ok {
			rval, err := p.assignment()
			if err != nil {
				return nil, err
			}
			return ast.Eassign{Name: w.Name, Value: rval}, nil
		}
		return nil, p.err(p.previous(), "lvalue of assignment is wrong")
	}
	return expr, nil
}

func (p *Parser) logic_or() (ast.Expr, error) {
	expr, err := p.logic_and()
	if err != nil {
		return nil, err
	}

	for p.match(token.Tor) {
		op := p.previous()
		right, err := p.logic_and()
		if err != nil {
			return nil, err
		}
		expr = ast.Elogical{Left: expr, Op: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) logic_and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.Tand) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = ast.Elogical{Left: expr, Op: op, Right: right}
	}
	return expr, nil
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
	return p.call()
}

func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(token.TleftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(expr ast.Expr) (ast.Expr, error) {
	args := make([]ast.Expr, 0)
	if !p.check(token.TrightParen) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
			if !p.match(token.Tcomma) {
				break
			}
		}
	}
	return ast.Ecall{Callee: expr, Paren: p.peek(), Args: args},
		p.consume(token.TrightParen, "Expected ) after call")
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
	if p.match(token.Tidentifier) {
		return ast.Evariable{p.previous()}, nil
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

func (p *Parser) synchorize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.Tsemicolon {
			return
		}

		switch p.peek().Type {
		case token.Tclass:
		case token.Tfun:
		case token.Tvar:
		case token.Tfor:
		case token.Tif:
		case token.Twhile:
		case token.Tprint:
		case token.Treturn:
			return
		}
		p.advance()
	}
}
