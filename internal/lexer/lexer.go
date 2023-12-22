package lexer

import "github.com/zsh-eng/go-interpreter/internal/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

const (
	NULL_ASCII = 0
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = NULL_ASCII
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = l.newTwoCharToken(token.EQ, l.ch) // ==
		} else {
			tok = newToken(token.ASSIGN, l.ch) // =
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			tok = l.newTwoCharToken(token.NOT_EQ, l.ch) // !=
		} else {
			tok = newToken(token.BANG, l.ch) // !
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case NULL_ASCII:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
		return tok // return here to avoid reading the next character
	}
	l.readChar()

	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[startPos:l.position]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

// Skip whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	startPos := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[startPos:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Peek at the next character without advancing the lexer
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return NULL_ASCII
	} else {
		return l.input[l.readPosition]
	}
}

// New token for two character tokens
func (l *Lexer) newTwoCharToken(tokenType token.TokenType, char byte) token.Token {
	l.readChar()
	return token.Token{Type: tokenType, Literal: string(char) + string(l.ch)}
}
