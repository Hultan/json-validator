package lexer

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/hultan/json-validator/internal/token"
)

func TestNewLexer_EmptyString(t *testing.T) {
	reader := strings.NewReader("")
	lexer := NewLexer(reader)
	assert.NotNil(t, lexer)

	tok := lexer.NextToken()
	assert.Equal(t, token.Token{
		Kind:    "EOF",
		Literal: "EOF",
		Line:    1,
		Column:  0,
	}, tok)
}

func TestNewLexer_EmptyJSON(t *testing.T) {
	reader := strings.NewReader("	{  } ")
	lexer := NewLexer(reader)
	assert.NotNil(t, lexer)

	tok := lexer.NextToken()
	assert.Equal(t, token.Token{
		Kind:    "{",
		Literal: "{",
		Line:    1,
		Column:  1,
	}, tok)

	tok = lexer.NextToken()
	assert.Equal(t, token.Token{
		Kind:    "}",
		Literal: "}",
		Line:    1,
		Column:  4,
	}, tok)

	tok = lexer.NextToken()
	assert.Equal(t, token.Token{
		Kind:    "EOF",
		Literal: "EOF",
		Line:    1,
		Column:  6,
	}, tok)
}

func TestNewLexer_NumberLiterals(t *testing.T) {
	input := `{
"n1":0,
"n2":1234567890,
"n3":-12,
"n4":-12.345,
"n5":-12.345e12,
"n6":-12.345E-12,
"n7":-12.345E+12,
}`

	tests := []struct {
		expectedKind    token.TokenKind
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LBRACE, "{", 1, 0},

		{token.STRING_LIT, "n1", 2, 0},
		{token.COLON, ":", 2, 4},
		{token.NUMBER_LIT, "0", 2, 5},
		{token.COMMA, ",", 2, 6},

		{token.STRING_LIT, "n2", 3, 0},
		{token.COLON, ":", 3, 4},
		{token.NUMBER_LIT, "1234567890", 3, 5},
		{token.COMMA, ",", 3, 15},

		{token.STRING_LIT, "n3", 4, 0},
		{token.COLON, ":", 4, 4},
		{token.NUMBER_LIT, "-12", 4, 5},
		{token.COMMA, ",", 4, 8},

		{token.STRING_LIT, "n4", 5, 0},
		{token.COLON, ":", 5, 4},
		{token.NUMBER_LIT, "-12.345", 5, 5},
		{token.COMMA, ",", 5, 12},

		{token.STRING_LIT, "n5", 6, 0},
		{token.COLON, ":", 6, 4},
		{token.NUMBER_LIT, "-12.345e12", 6, 5},
		{token.COMMA, ",", 6, 15},

		{token.STRING_LIT, "n6", 7, 0},
		{token.COLON, ":", 7, 4},
		{token.NUMBER_LIT, "-12.345E-12", 7, 5},
		{token.COMMA, ",", 7, 16},

		{token.STRING_LIT, "n7", 8, 0},
		{token.COLON, ":", 8, 4},
		{token.NUMBER_LIT, "-12.345E+12", 8, 5},
		{token.COMMA, ",", 8, 16},

		{token.RBRACE, "}", 9, 0},

		{token.EOF, "EOF", 9, 1},
	}

	tokenTest(t, input, tests)
}

func TestNewLexer_Arrays(t *testing.T) {
	input := `{
"a1":[0,1,2]
}`

	tests := []struct {
		expectedKind    token.TokenKind
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LBRACE, "{", 1, 0},

		{token.STRING_LIT, "a1", 2, 0},
		{token.COLON, ":", 2, 4},
		{token.LBRACKET, "[", 2, 5},
		{token.NUMBER_LIT, "0", 2, 6},
		{token.COMMA, ",", 2, 7},
		{token.NUMBER_LIT, "1", 2, 8},
		{token.COMMA, ",", 2, 9},
		{token.NUMBER_LIT, "2", 2, 10},
		{token.RBRACKET, "]", 2, 11},

		{token.RBRACE, "}", 3, 0},

		{token.EOF, "EOF", 3, 1},
	}

	tokenTest(t, input, tests)
}

func TestNewLexer_AdvancedJSON(t *testing.T) {
	input := `{
"t":true,
"f":false,
"s":"s",
"n":-12.345
}`

	tests := []struct {
		expectedKind    token.TokenKind
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LBRACE, "{", 1, 0},

		{token.STRING_LIT, "t", 2, 0},
		{token.COLON, ":", 2, 3},
		{token.TRUE, "true", 2, 4},
		{token.COMMA, ",", 2, 8},

		{token.STRING_LIT, "f", 3, 0},
		{token.COLON, ":", 3, 3},
		{token.FALSE, "false", 3, 4},
		{token.COMMA, ",", 3, 9},

		{token.STRING_LIT, "s", 4, 0},
		{token.COLON, ":", 4, 3},
		{token.STRING_LIT, "s", 4, 4},
		{token.COMMA, ",", 4, 7},

		{token.STRING_LIT, "n", 5, 0},
		{token.COLON, ":", 5, 3},
		{token.NUMBER_LIT, "-12.345", 5, 4},

		{token.RBRACE, "}", 6, 0},

		{token.EOF, "EOF", 6, 1},
	}

	tokenTest(t, input, tests)
}

func tokenTest(t *testing.T, input string, tests []struct {
	expectedKind    token.TokenKind
	expectedLiteral string
	expectedLine    int
	expectedColumn  int
}) {
	reader := strings.NewReader(input)
	l := NewLexer(reader)

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
			fmt.Printf("[%-10v] %-20v\t(%v, %v)\n", tok.Kind, tok.Literal, tok.Line, tok.Column)
			spew.Dump(tok)
			t.FailNow()
		}
	}
}
