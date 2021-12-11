// Package parser is based on the code from the book
// Writing an Interpreter in Go 1.7 by Thorsten Ball
// https://thorstenball.com/books/

package parser

import (
	"fmt"
	"os"

	"github.com/hultan/json-validator/internal/lexer"
	"github.com/hultan/json-validator/internal/token"
)

type ParserError struct {
	Message string
	Token   token.Token
}

type Parser struct {
	l      *lexer.Lexer
	Errors []ParserError

	curToken  token.Token
	peekToken token.Token
}

func NewParser(fileName string) *Parser {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	// Create parser and lexer
	p := &Parser{}
	p.l = lexer.NewLexer(file)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	fmt.Printf("%20s%20s (%v,%v)\n", p.curToken.Kind, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenKind) bool {
	return p.curToken.Kind == t
}

func (p *Parser) peekTokenIs(t token.TokenKind) bool {
	return p.peekToken.Kind == t
}

func (p *Parser) expectPeek(t token.TokenKind) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t ...token.TokenKind) {
	ts := p.getTokenString(t)
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		ts, p.peekToken.Kind)
	err := ParserError{
		Message: msg,
		Token:   p.curToken,
	}
	p.Errors = append(p.Errors, err)
}

func (p *Parser) getTokenString(t []token.TokenKind) string {
	var ts string
	for i, kind := range t {
		if i == len(t)-1 {
			ts += " or "
		} else if i > 0 {
			ts += ","
		}
		ts += fmt.Sprintf("%s", kind)
	}
	return ts
}

func (p *Parser) Parse() bool {
	return p.parseObject()
}

func (p *Parser) parseObject() bool {
	count := 0

	if p.curToken.Kind != token.LBRACE {
		err := ParserError{
			Message: "expected {",
			Token:   p.curToken,
		}
		p.Errors = append(p.Errors, err)
		return false
	}
	p.nextToken()

	for {
		// If the next token is the ending brace, then we are done parsing the object
		if p.curToken.Kind == token.RBRACE {
			break
		}
		// Accept a comma if this is not the first item
		if count > 0 && p.curToken.Kind == token.COMMA {
			p.nextToken()
		}
		if p.curToken.Kind != token.STRING_LIT {
			err := ParserError{
				Message: "expected string literal",
				Token:   p.curToken,
			}
			p.Errors = append(p.Errors, err)
			return false
		}
		p.nextToken()
		if p.curToken.Kind != token.COLON {
			err := ParserError{
				Message: "expected :",
				Token:   p.curToken,
			}
			p.Errors = append(p.Errors, err)
			return false

		}
		p.nextToken()

		ok := p.parseValue()
		if !ok {
			return false
		}

		count++
	}

	if p.curToken.Kind != token.RBRACE {
		err := ParserError{
			Message: "expected }",
			Token:   p.curToken,
		}
		p.Errors = append(p.Errors, err)
		return false
	}
	p.nextToken()

	return true
}

func (p *Parser) parseValue() bool {
	switch p.curToken.Kind {
	case token.STRING_LIT:
		p.nextToken()
	case token.NUMBER_LIT:
		p.nextToken()
	case token.LBRACKET:
		return p.parseArray()
	case token.LBRACE:
		return p.parseObject()
	case token.TRUE:
		p.nextToken()
	case token.FALSE:
		p.nextToken()
	case token.NULL:
		p.nextToken()
	default:
		err := ParserError{
			Message: "expected value",
			Token:   p.curToken,
		}
		p.Errors = append(p.Errors, err)
		return false
	}
	return true
}

func (p *Parser) parseArray() bool {
	count := 0

	if p.curToken.Kind != token.LBRACKET {
		err := ParserError{
			Message: "expected [",
			Token:   p.curToken,
		}
		p.Errors = append(p.Errors, err)
		return false
	}
	p.nextToken()

	for {
		// If the next token is the ending bracket, then we are done parsing the array
		if p.curToken.Kind == token.RBRACKET {
			break
		}
		// Accept a comma if this is not the first item
		if count > 0 && p.curToken.Kind == token.COMMA {
			p.nextToken()
		}

		ok := p.parseValue()
		if !ok {
			return false
		}

		count++
	}

	if p.curToken.Kind != token.RBRACKET {
		err := ParserError{
			Message: "expected ]",
			Token:   p.curToken,
		}
		p.Errors = append(p.Errors, err)
		return false
	}
	p.nextToken()

	return true
}
