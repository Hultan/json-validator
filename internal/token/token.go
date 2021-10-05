package token

type TokenKind string // For performance reasons this should be an integer instead

type Token struct {
	Kind     TokenKind
	Literal  string
	FileName string
	Line     int
	Column   int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT      = "IDENT"      // add, foobar, x, y, ...
	BOOL       = "BOOL"       // true, false
	INT_LIT    = "INT_LIT"    // 1343456
	FLOAT_LIT  = "FLOAT_LIT"  // 12.3455
	STRING_LIT = "STRING_LIT" // "foobar"
	INT        = "int"        // int keyword
	FLOAT      = "float"      // float keyword
	STRING     = "string"     // string keyword

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	INCREASE = "++"
	MINUS    = "-"
	DECREASE = "--"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	DOLLAR   = "$" // String interpolation

	LT    = "<"
	LT_EQ = "<="
	GT    = ">"
	GT_EQ = ">="

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	DOT       = "."
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FOR      = "FOR"
	DO       = "DO"
	WHILE    = "WHILE"
	UNTIL    = "UNTIL"
	CONTINUE = "CONTINUE"
	BREAK    = "BREAK"

	IF     = "IF"
	ELSEIF = "ELSEIF"
	ELSE   = "ELSE"

	SWITCH  = "SWITCH"
	CASE    = "CASE"
	DEFAULT = "DEFAULT"

	ENUM  = "ENUM"
	CONST = "CONST"

	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"

	TRUE  = "TRUE"
	FALSE = "FALSE"
)

var keywordsMap = map[string]TokenKind{
	"int":    INT,
	"bool":   BOOL,
	"float":  FLOAT,
	"string": STRING,

	"for":      FOR,
	"do":       DO,
	"while":    WHILE,
	"until":    UNTIL,
	"continue": CONTINUE,
	"break":    BREAK,

	"if":     IF,
	"elseif": ELSEIF,
	"else":   ELSE,

	"switch":  SWITCH,
	"case":    CASE,
	"default": DEFAULT,

	"enum":  ENUM,
	"const": CONST,

	"function": FUNCTION,
	"return":   RETURN,

	"true":  TRUE,
	"false": FALSE,
}

func LookupIdent(ident string) TokenKind {
	if tok, ok := keywordsMap[ident]; ok {
		return tok
	}
	return IDENT
}