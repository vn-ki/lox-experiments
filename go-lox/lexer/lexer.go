package lexer

import (
	"log"
	"strconv"
	"unicode"

	"github.com/vn-ki/go-lox/token"
)

var KEYWORDS = map[string]token.TokenType{
	"and":    token.Tand,
	"class":  token.Tclass,
	"else":   token.Telse,
	"false":  token.Tfalse,
	"for":    token.Tfor,
	"fun":    token.Tfun,
	"if":     token.Tif,
	"nil":    token.Tnil,
	"or":     token.Tor,
	"print":  token.Tprint,
	"return": token.Treturn,
	"super":  token.Tsuper,
	"this":   token.Tsuper,
	"true":   token.Ttrue,
	"var":    token.Tvar,
	"while":  token.Twhile,
}

type Lexer struct {
	start   int
	current int
	line    int
	src     []rune
	// This could be a channel in the most golang-y way
	// But following crafting interpreters closely here
	tokens       []token.Token
	ErrorHandler func(int, string)
}

// TODO: use reader instead of string here
func NewLexer(src string) *Lexer {
	return &Lexer{
		0, 0, 1,
		[]rune(src),
		make([]token.Token, 0),
		nil,
	}
}

func (l *Lexer) ScanTokens() []token.Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	l.tokens = append(l.tokens, token.Token{token.Teof, "", nil, l.line})
	return l.tokens
}

func (l *Lexer) scanToken() {
	c := l.advance()
	// log.Printf("Reading character %c\n", c)
	switch c {
	case '(':
		l.addToken(token.TleftParen)
	case ')':
		l.addToken(token.TrightParen)
	case '{':
		l.addToken(token.TleftBrace)
	case '}':
		l.addToken(token.TrightBrace)
	case ',':
		l.addToken(token.Tcomma)
	case '.':
		l.addToken(token.Tdot)
	case '-':
		l.addToken(token.Tminus)
	case '+':
		l.addToken(token.Tplus)
	case ';':
		l.addToken(token.Tsemicolon)
	case '*':
		l.addToken(token.Tstar)

	case '!':
		if l.match('=') {
			l.addToken(token.TbangEqual)
		} else {
			l.addToken(token.Tbang)
		}
	case '=':
		if l.match('=') {
			l.addToken(token.TequalEqual)
		} else {
			l.addToken(token.Tequal)
		}
	case '<':
		if l.match('=') {
			l.addToken(token.TlessEqual)
		} else {
			l.addToken(token.Tless)
		}
	case '>':
		if l.match('=') {
			l.addToken(token.TgreaterEqual)
		} else {
			l.addToken(token.Tgreater)
		}

	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			l.addToken(token.Tslash)
		}

	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		l.line++

	case '"':
		l.parseString()

	default:
		if unicode.IsDigit(c) {
			l.parseNum()
		} else if unicode.IsLetter(c) {
			l.parseIden()
		} else {
			l.err("Unexpected character")
		}
	}
}

func (l *Lexer) parseIden() {
	for unicode.IsLetter(l.peek()) {
		l.advance()
	}
	val := string(l.src[l.start:l.current])
	tokenType, ok := KEYWORDS[val]
	if !ok {
		tokenType = token.Tidentifier
	}
	l.addToken(tokenType)
}

func (l *Lexer) parseString() {
	log.Println("Parsing string")
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}

	if l.isAtEnd() {
		l.err("Unterminated string")
		return
	}

	// last '"'
	l.advance()

	// trim surrounding quotes
	value := string(l.src[l.start+1 : l.current-1])
	l.addTokenWithLiteral(token.Tstring, value)
}

func (l *Lexer) parseNum() {
	log.Println("Parsing a number")
	for unicode.IsDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && unicode.IsDigit(l.peekNext()) {
		l.advance()
		for unicode.IsDigit(l.peek()) {
			l.advance()
		}
	}
	val, err := strconv.ParseFloat(string(l.src[l.start:l.current]), 64)
	if err != nil {
		log.Fatalln("Couldnot parse number: " + string(l.src[l.start:l.current]))
	}
	l.addTokenWithLiteral(token.Tnumber, val)
}

func (l *Lexer) addToken(ty token.TokenType) {
	// log.Printf()
	l.addTokenWithLiteral(ty, nil)
}

func (l *Lexer) addTokenWithLiteral(ty token.TokenType, literal interface{}) {
	text := string(l.src[l.start:l.current])
	l.tokens = append(l.tokens, token.Token{ty, text, literal, l.line})
}

func (l *Lexer) advance() rune {
	l.current += 1
	return l.src[l.current-1]
}

func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return rune(0)
	}
	return l.src[l.current]
}

func (l *Lexer) peekNext() rune {
	if l.current+1 >= len(l.src) {
		return rune(0)
	}
	return l.src[l.current+1]
}

func (l *Lexer) match(expected rune) bool {
	if l.isAtEnd() {
		return false
	}

	if l.src[l.current] != expected {
		return false
	}

	l.current += 1
	return true
}

func (l *Lexer) err(msg string) {
	log.Print("!!! Error: " + msg)
	if l.ErrorHandler != nil {
		l.ErrorHandler(l.line, msg)
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.src)
}
