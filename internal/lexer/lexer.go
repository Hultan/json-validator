package lexer

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/hultan/per/internal/token"
)

type Lexer struct {
	runeList     []rune
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	currentRune  rune // current char under examination
	fileName     string
	line         int
	column       int
}

func NewLexer(fileName string, reader io.Reader) *Lexer {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	text := string(input)
	l := &Lexer{
		fileName: fileName,
		runeList: []rune(text),
		line: 1,
	}
	if text != "" {
		l.currentRune = l.runeList[l.position]
	}
	return l
}

func (l *Lexer) NextToken() token.Token {
	var currentToken token.Token
	advance := 1
	l.skipWhitespace()
	l.skipComment()

	switch l.currentRune {
	case '=':
		if l.peekRune() == '=' {
			l.readRune()
			currentToken = l.newTokenString(token.EQ, "==")
		} else {
			currentToken = l.newTokenString(token.ASSIGN, string(l.currentRune))
		}
	case '+':
		currentToken = l.newTokenString(token.PLUS, string(l.currentRune))
	case '-':
		currentToken = l.newTokenString(token.MINUS, string(l.currentRune))
	case '!':
		if l.peekRune() == '=' {
			l.readRune()
			currentToken = l.newTokenString(token.NOT_EQ, "!=")
		} else {
			currentToken = l.newTokenString(token.BANG, string(l.currentRune))
		}
	case '/':
		currentToken = l.newTokenString(token.SLASH, string(l.currentRune))
	case '*':
		currentToken = l.newTokenString(token.ASTERISK, string(l.currentRune))
	case '<':
		currentToken = l.newTokenString(token.LT, string(l.currentRune))
	case '>':
		currentToken = l.newTokenString(token.GT, string(l.currentRune))
	case ';':
		currentToken = l.newTokenString(token.SEMICOLON, string(l.currentRune))
	case ':':
		currentToken = l.newTokenString(token.COLON, string(l.currentRune))
	case ',':
		currentToken = l.newTokenString(token.COMMA, string(l.currentRune))
	case '.':
		currentToken = l.newTokenString(token.DOT, string(l.currentRune))
	case '{':
		currentToken = l.newTokenString(token.LBRACE, string(l.currentRune))
	case '}':
		currentToken = l.newTokenString(token.RBRACE, string(l.currentRune))
	case '(':
		currentToken = l.newTokenString(token.LPAREN, string(l.currentRune))
	case ')':
		currentToken = l.newTokenString(token.RPAREN, string(l.currentRune))
	case '"':
		lit := l.readString()
		currentToken = l.newTokenString(token.STRING, lit)
		advance = len(lit) + 2 // + 2 because of the two quotation marks
	case '[':
		currentToken = l.newTokenString(token.LBRACKET, string(l.currentRune))
	case ']':
		currentToken = l.newTokenString(token.RBRACKET, string(l.currentRune))
	case 0:
		currentToken = l.newTokenString(token.EOF, "EOF")
	default:
		if l.isLetter(l.currentRune) {
			lit := l.readIdentifier()
			currentToken = l.newTokenString(token.LookupIdent(lit), lit)
			l.column += len(lit)
			return currentToken
		} else if l.isDigit(l.currentRune) {
			lit := l.readNumber()
			if strings.Contains(lit, ".") {
				currentToken = l.newTokenString(token.FLOAT, lit)
			} else {
				currentToken = l.newTokenString(token.INT, lit)
			}
			l.column += len(lit)
			return currentToken
		} else {
			currentToken = l.newTokenString(token.ILLEGAL, string(l.currentRune))
		}
	}

	l.readRune()
	l.column += advance
	return currentToken
}

func (l *Lexer) skipComment() {
	for l.isComment() {
		l.skipCommentLine()
		l.skipWhitespace()
	}
}

func (l *Lexer) skipCommentLine() {
	for keepGoing := true; keepGoing; {
		if l.isNewLine() || l.currentRune == 0 {
			keepGoing = false
			break
		}
		l.readRune()
	}
}

func (l *Lexer) skipWhitespace() {
	for keepGoing := true; keepGoing; {
		switch {
		case l.isNewLine():
			l.line += 1
			l.column = 0
			l.readRune()
		case l.isWhitespace():
			l.column += 1
			l.readRune()
		default:
			keepGoing = false
		}
	}
}

func (l *Lexer) readRune() {
	if l.readPosition >= len(l.runeList) {
		l.currentRune = 0
	} else {
		l.currentRune = l.runeList[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekRune() rune {
	if l.readPosition >= len(l.runeList) {
		return 0
	} else {
		return l.runeList[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	first := true
	for (first && l.isLetter(l.currentRune)) ||
		(!first && l.isLetterOrDigit(l.currentRune)) {
		l.readRune()
		first = false
	}
	return string(l.runeList[position:l.position])
}

func (l *Lexer) readNumber() string {
	position := l.position
	for l.isDigit(l.currentRune) {
		l.readRune()
	}
	if l.currentRune == '.' {
		l.readRune()
	}
	for l.isDigit(l.currentRune) {
		l.readRune()
	}
	return string(l.runeList[position:l.position])
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readRune()
		if l.currentRune == '"' || l.currentRune == 0 {
			break
		}
	}
	return string(l.runeList[position:l.position])
}

func (l *Lexer) isComment() bool {
	return l.currentRune == '/' && l.peekRune() == '/'
}

func (l *Lexer) isLetterOrDigit(r rune) bool {
	return l.isDigit(r) || l.isLetter(r)
}

func (l *Lexer) isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' ||
		'A' <= r && r <= 'Z' ||
		r == '_'
}

func (l *Lexer) isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func (l *Lexer) isWhitespace() bool {
	return l.currentRune == ' ' ||
		l.currentRune == '\t'
}

func (l *Lexer) isNewLine() bool {
	return l.currentRune == '\r' ||
		l.currentRune == '\n'
}

func (l *Lexer) newTokenString(tokenKind token.TokenKind, literal string) token.Token {
	return token.Token{
		Kind:     tokenKind,
		Literal:  literal,
		FileName: l.fileName,
		Line:     l.line,
		Column:   l.column,
	}
}
