package lexer

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/hultan/per/internal/token"
)

func TestNewLexer_EmptyFile(t *testing.T) {
	reader := strings.NewReader("")
	lexer := NewLexer("test.per", reader)
	assert.NotNil(t, lexer)

	tok := lexer.NextToken()
	assert.Equal(t, token.Token{
		Kind:     "EOF",
		Literal:  "EOF",
		FileName: "test.per",
		Line:     1,
		Column:   0,
	}, tok)
}

func TestNewLexer_BoolAssignment(t *testing.T) {
	input := `bool t=true;
bool f=5<7;`

	tests := []struct {
		expectedKind     token.TokenKind
		expectedLiteral  string
		expectedFileName string
		expectedLine     int
		expectedColumn   int
	}{
		{token.BOOL, "bool", "test.per", 1, 0},
		{token.IDENT, "t", "test.per", 1, 5},
		{token.ASSIGN, "=", "test.per", 1, 6},
		{token.TRUE, "true", "test.per", 1, 7},
		{token.SEMICOLON, ";", "test.per", 1, 11},

		{token.BOOL, "bool", "test.per", 2, 0},
		{token.IDENT, "f", "test.per", 2, 5},
		{token.ASSIGN, "=", "test.per", 2, 6},
		{token.INT_LIT, "5", "test.per", 2, 7},
		{token.LT, "<", "test.per", 2, 8},
		{token.INT_LIT, "7", "test.per", 2, 9},
		{token.SEMICOLON, ";", "test.per", 2, 10},

		{token.EOF, "EOF", "test.per", 2, 11},
	}

	tokenTest(t, input, tests)
}

func TestNewLexer_IntegerAssignment(t *testing.T) {
	input := `int c=5;
int x,y=5+7;`

	tests := []struct {
		expectedKind     token.TokenKind
		expectedLiteral  string
		expectedFileName string
		expectedLine     int
		expectedColumn   int
	}{
		{token.INT, "int", "test.per", 1, 0},
		{token.IDENT, "c", "test.per", 1, 4},
		{token.ASSIGN, "=", "test.per", 1, 5},
		{token.INT_LIT, "5", "test.per", 1, 6},
		{token.SEMICOLON, ";", "test.per", 1, 7},

		{token.INT, "int", "test.per", 2, 0},
		{token.IDENT, "x", "test.per", 2, 4},
		{token.COMMA, ",", "test.per", 2, 5},
		{token.IDENT, "y", "test.per", 2, 6},
		{token.ASSIGN, "=", "test.per", 2, 7},
		{token.INT_LIT, "5", "test.per", 2, 8},
		{token.PLUS, "+", "test.per", 2, 9},
		{token.INT_LIT, "7", "test.per", 2, 10},
		{token.SEMICOLON, ";", "test.per", 2, 11},

		{token.EOF, "EOF", "test.per", 2, 12},
	}

	tokenTest(t, input, tests)
}

func TestNewLexer_FloatNumbers(t *testing.T) {
	input := `5.0;5;5.1.1;`

	tests := []struct {
		expectedKind     token.TokenKind
		expectedLiteral  string
		expectedFileName string
		expectedLine     int
		expectedColumn   int
	}{
		{token.FLOAT_LIT, "5.0", "test.per", 1, 0},
		{token.SEMICOLON, ";", "test.per", 1, 3},
		{token.INT_LIT, "5", "test.per", 1, 4},
		{token.SEMICOLON, ";", "test.per", 1, 5},
		{token.FLOAT_LIT, "5.1", "test.per", 1, 6},
		{token.DOT, ".", "test.per", 1, 9},
		{token.INT_LIT, "1", "test.per", 1, 10},
		{token.SEMICOLON, ";", "test.per", 1, 11},

		{token.EOF, "EOF", "test.per", 1, 12},
	}

	tokenTest(t, input, tests)
}

func TestNewLexer_FloatAssignment(t *testing.T) {
	input := `float f=5.0;
float x,y=7.5;`

	tests := []struct {
		expectedKind     token.TokenKind
		expectedLiteral  string
		expectedFileName string
		expectedLine     int
		expectedColumn   int
	}{
		{token.FLOAT, "float", "test.per", 1, 0},
		{token.IDENT, "f", "test.per", 1, 6},
		{token.ASSIGN, "=", "test.per", 1, 7},
		{token.FLOAT_LIT, "5.0", "test.per", 1, 8},
		{token.SEMICOLON, ";", "test.per", 1, 11},

		{token.FLOAT, "float", "test.per", 2, 0},
		{token.IDENT, "x", "test.per", 2, 6},
		{token.COMMA, ",", "test.per", 2, 7},
		{token.IDENT, "y", "test.per", 2, 8},
		{token.ASSIGN, "=", "test.per", 2, 9},
		{token.FLOAT_LIT, "7.5", "test.per", 2, 10},
		{token.SEMICOLON, ";", "test.per", 2, 13},

		{token.EOF, "EOF", "test.per", 2, 14},
	}

	tokenTest(t, input, tests)
}

func TestNewLexer_StringAssignment(t *testing.T) {
	input := `string c="per";
string x,y="per";`

	tests := []struct {
		expectedKind     token.TokenKind
		expectedLiteral  string
		expectedFileName string
		expectedLine     int
		expectedColumn   int
	}{
		{token.STRING, "string", "test.per", 1, 0},
		{token.IDENT, "c", "test.per", 1, 7},
		{token.ASSIGN, "=", "test.per", 1, 8},
		{token.STRING_LIT, "per", "test.per", 1, 9},
		{token.SEMICOLON, ";", "test.per", 1, 14},

		{token.STRING, "string", "test.per", 2, 0},
		{token.IDENT, "x", "test.per", 2, 7},
		{token.COMMA, ",", "test.per", 2, 8},
		{token.IDENT, "y", "test.per", 2, 9},
		{token.ASSIGN, "=", "test.per", 2, 10},
		{token.STRING_LIT, "per", "test.per", 2, 11},
		{token.SEMICOLON, ";", "test.per", 2, 16},

		{token.EOF, "EOF", "test.per", 2, 17},
	}

	tokenTest(t, input, tests)
}

// func TestNewLexer_StringInterpolation(t *testing.T) {
// 	input := `string c=$"per {test}:{{test2}}";
// string d=$"{test}";`
//
// 	tests := []struct {
// 		expectedKind     token.TokenKind
// 		expectedLiteral  string
// 		expectedFileName string
// 		expectedLine     int
// 		expectedColumn   int
// 	}{
// 		{token.STRING, "string", "test.per", 1, 0},
// 		{token.IDENT, "c", "test.per", 1, 7},
// 		{token.ASSIGN, "=", "test.per", 1, 8},
// 		{token.DOLLAR, "$", "test.per", 1, 9},
// 		{token.STRING_LIT, "per ", "test.per", 1, 11},
// 		{token.LBRACE, "{", "test.per", 1, 15},
// 		{token.IDENT, "test", "test.per", 1, 16},
// 		{token.RBRACE, "}", "test.per", 1, 20},
// 		{token.STRING_LIT, ":{test2}", "test.per", 1, 21},
// 		{token.SEMICOLON, ";", "test.per", 1, 32},
//
// 		{token.STRING, "string", "test.per", 2, 0},
// 		{token.IDENT, "d", "test.per", 2, 7},
// 		{token.ASSIGN, "=", "test.per", 2, 8},
// 		{token.DOLLAR, "$", "test.per", 2, 9},
// 		{token.LBRACE, "{", "test.per", 2, 11},
// 		{token.IDENT, "test", "test.per", 2, 12},
// 		{token.RBRACE, "}", "test.per", 2, 16},
// 		{token.SEMICOLON, ";", "test.per", 2, 18},
//
// 		{token.EOF, "EOF", "test.per", 2, 19},
// 	}
//
// 	tokenTest(t, input, tests)
// }

func tokenTest(t *testing.T, input string, tests []struct {
	expectedKind     token.TokenKind
	expectedLiteral  string
	expectedFileName string
	expectedLine     int
	expectedColumn   int
}) {
	reader := strings.NewReader(input)
	l := NewLexer("test.mon", reader)

	for i, tt := range tests {
		tok := l.NextToken()

		fail := false
		if tok.Kind != tt.expectedKind || tok.Literal != tt.expectedLiteral ||
			tok.Line != tt.expectedLine || tok.Column != tt.expectedColumn {
			fail = true
		}

		if tok.Kind != tt.expectedKind {
			log.Printf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedKind, tok.Kind)
			fail = true
		}

		if tok.Literal != tt.expectedLiteral {
			log.Printf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
			fail = true
		}

		if tok.Line != tt.expectedLine {
			log.Printf("tests[%d] - Line number wrong. expected=%q, got=%q",
				i, tt.expectedLine, tok.Line)
			fail = true
		}

		if tok.Column != tt.expectedColumn {
			log.Printf("tests[%d] - Column number wrong. expected=%q, got=%q",
				i, tt.expectedColumn, tok.Column)
			fail = true
		}

		if fail {
			fmt.Printf("[%-10v] %-20v\t(%s ,%v, %v)\n", tok.Kind, tok.Literal, tok.FileName, tok.Line, tok.Column)
			spew.Dump(tok)
			t.FailNow()
		}
	}
}
