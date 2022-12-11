package lexer

import "ylang/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.advance()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch ch := l.ch; ch {
	case '=':
		if l.peekChar() == '=' {
			l.advance()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.advance()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, ch)
		}

	case '-':
		tok = newToken(token.MINUS, ch)
	case '+':
		tok = newToken(token.PLUS, ch)
	case '/':
		tok = newToken(token.SLASH, ch)
	case '*':
		tok = newToken(token.ASTERISK, ch)
	case '<':
		tok = newToken(token.LT, ch)
	case '>':
		tok = newToken(token.GT, ch)
	case ';':
		tok = newToken(token.SEMICOLON, ch)
	case '(':
		tok = newToken(token.LPAREN, ch)
	case ')':
		tok = newToken(token.RPAREN, ch)
	case ',':
		tok = newToken(token.COMMA, ch)
	case '{':
		tok = newToken(token.LBRACE, ch)
	case '}':
		tok = newToken(token.RBRACE, ch)

	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}

	default:
		switch {
		case isLetter(ch):
			ident := l.readIdentifier()
			return token.Token{
				Type:    token.LookupIdent(ident),
				Literal: ident,
			}
		case isDigit(ch):
			return token.Token{
				Type:    token.INT,
				Literal: l.readNumber(),
			}
		default:
			tok = newToken(token.ILLEGAL, ch)
		}
	}

	l.advance()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) advance() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.advance()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.advance()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.ch {
		case ' ', '\t', '\n', '\r':
			l.advance()

		case '/':
			if l.peekChar() == '/' {
				// treating comments as whitespace, sue me
				for l.ch != '\n' && l.ch != 0 {
					l.advance()
				}
			} else {
				return
			}

		default:
			return
		}
	}
}

// utils

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
