// Package token is based on the code from the book
// Writing an Interpreter in Go 1.7 by Thorsten Ball
// https://thorstenball.com/books/

package token

type TokenKind string // For performance reasons this should be an integer instead

type Token struct {
	Kind    TokenKind
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL    = "ILLEGAL"
	EOF        = "EOF"
	STRING_LIT = "STRING_LIT"
	NUMBER_LIT = "NUMBER_LIT"
	BACKSLASH  = "\\"
	COMMA      = ","
	DOT        = "."
	COLON      = ":"
	LBRACE     = "{"
	RBRACE     = "}"
	LBRACKET   = "["
	RBRACKET   = "]"
	TRUE       = "true"
	FALSE      = "false"
	NULL       = "null"
)

var keywordsMap = map[string]TokenKind{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
}

func LookupIdent(ident string) TokenKind {
	if tok, ok := keywordsMap[ident]; ok {
		return tok
	}
	return ILLEGAL
}
