package token

import "fmt"

type TokenType int

const (
	// Single character tokens
	TleftParen TokenType = iota
	TrightParen
	TleftBrace
	TrightBrace
	Tcomma
	Tdot
	Tminus
	Tplus
	Tsemicolon
	Tslash
	Tstar

	// one or two character tokens
	Tbang
	TbangEqual
	Tequal
	TequalEqual
	Tgreater
	TgreaterEqual
	Tless
	TlessEqual

	// Literals
	// TODO: Prefix with L?
	Tidentifier
	Tstring
	Tnumber

	// Keywords
	// TODO: Prefix with K?
	Tand
	Tclass
	Telse
	Tfalse
	Tfun
	Tfor
	Tif
	Tnil
	Tor
	Tprint
	Treturn
	Tsuper
	Tthis
	Ttrue
	Tvar
	Twhile

	Teof
)

type Token struct {
	Type   TokenType
	Lexeme string
	// XXX: Not sure what this is
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("\033[0;36m%s\033[0m\t%s", tokenReverseLookup(t.Type), t.Lexeme)
}

func tokenReverseLookup(ty TokenType) string {
	var tokenNames = []string{
		"LeftParen",
		"RightParen",
		"LeftBrace",
		"RightBrace",
		"Comma",
		"Dot",
		"Minus",
		"Plus",
		"Semicolon",
		"Slash",
		"Star",
		"Bang",
		"BangEq",
		"Equal",
		"EqualEqual",
		"Greater",
		"GreaterEqual",
		"Less",
		"LessEqual",
		"Identifier",
		"String",
		"Number",
		"Keyword And",
		"Keyword Class",
		"Keyword Else",
		"Keyword False",
		"Keyword Fun",
		"Keyword for",
		"Keyword if",
		"Keyword nil",
		"Keyword or",
		"Keyword print",
		"Keyword return",
		"Keyword super",
		"Keyword this",
		"Keyword true",
		"Keyword var",
		"Keyword while",
		"EOF",
	}
	return tokenNames[ty]
}
