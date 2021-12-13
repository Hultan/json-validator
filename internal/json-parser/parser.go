// Package parser is based on the code from the book
// Writing an Interpreter in Go 1.7 by Thorsten Ball
// https://thorstenball.com/books/

package parser

import (
	"os"

	"github.com/hultan/json-validator/internal/json-lexer"
	"github.com/hultan/json-validator/internal/token"
)

type ParserError struct {
	Message string
	Token   token.Token
}

type Parser struct {
	l      *json_lexer.Lexer
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
	p.l = json_lexer.NewLexer(file)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	// fmt.Printf("%20s%20s (%v,%v)\n", p.curToken.Kind, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() bool {
	return p.parseObject()
}

func (p *Parser) parseObject() bool {
	first := true

	if p.curToken.Kind != token.LBRACE {
		p.addError("Expected left brace - {")
		return false
	}
	p.nextToken()

	for {
		// If the next token is the ending brace, then we are done parsing the object
		if p.curToken.Kind == token.RBRACE {
			break
		}
		// Accept a comma if this is not the first item
		if !first && p.curToken.Kind == token.COMMA {
			p.nextToken()
		}
		if p.curToken.Kind != token.STRING_LIT {
			p.addError("Expected string literal")
			return false
		}
		p.nextToken()
		if p.curToken.Kind != token.COLON {
			p.addError("Expected a colon - :")
			return false

		}
		p.nextToken()

		ok := p.parseValue()
		if !ok {
			return false
		}

		first = false
	}

	if p.curToken.Kind != token.RBRACE {
		p.addError("Expected right brace - }")
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
		p.addError("Unknown value")
		return false
	}
	return true
}

func (p *Parser) parseArray() bool {
	first := true

	if p.curToken.Kind != token.LBRACKET {
		p.addError("Expected left bracket - [")
		return false
	}
	p.nextToken()

	for {
		// If the next token is the ending bracket, then we are done parsing the array
		if p.curToken.Kind == token.RBRACKET {
			break
		}
		// Accept a comma if this is not the first item
		if !first && p.curToken.Kind == token.COMMA {
			p.nextToken()
		}

		ok := p.parseValue()
		if !ok {
			return false
		}

		first = false
	}

	if p.curToken.Kind != token.RBRACKET {
		p.addError("Expected right bracket - ]")
		return false
	}
	p.nextToken()

	return true
}
func (p *Parser) addError(message string) {
	err := ParserError{
		Message: message,
		Token:   p.curToken,
	}
	p.Errors = append(p.Errors, err)
}
