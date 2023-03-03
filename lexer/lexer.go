package lexer

import "yy/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	line         int  // current line
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.advance()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.ch {
	case '<':
		tok = l.newTokenByType(token.LT)
	case '>':
		tok = l.newTokenByType(token.GT)
	case ';':
		tok = l.newTokenByType(token.SEMICOLON)
	case '(':
		tok = l.newTokenByType(token.LPAREN)
	case ')':
		tok = l.newTokenByType(token.RPAREN)
	case ',':
		tok = l.newTokenByType(token.COMMA)
	case '{':
		tok = l.newTokenByType(token.LBRACE)
	case '}':
		tok = l.newTokenByType(token.RBRACE)
	case '[':
		tok = l.newTokenByType(token.LBRACKET)
	case ']':
		tok = l.newTokenByType(token.RBRACKET)
	case '\\':
		tok = l.newTokenByType(token.BACKSLASH)
	case '%':
		tok = l.newTokenByType(token.PERCENT)

	case '+':
		tok = l.switch2(token.PLUS, token.ADD_ASSIGN)
	case '-':
		tok = l.switch2(token.MINUS, token.SUB_ASSIGN)
	case '*':
		tok = l.switch2(token.ASTERISK, token.MUL_ASSIGN)
	case '/':
		tok = l.switch2(token.SLASH, token.DIV_ASSIGN)
	case '=':
		tok = l.switch2(token.ASSIGN, token.EQ)
	case '!':
		tok = l.switch2(token.BANG, token.NOT_EQ)
	case ':':
		tok = l.switch2(token.COLON, token.WALRUS)

	case '.':
		if l.peek() == '.' {
			l.advance()
			tok = l.newToken(token.RANGE, token.RANGE.String())
		}

	case '@':
		if l.peek() == '\\' {
			l.advance()
			tok = l.newToken(token.MACRO, token.MACRO.String())
		}

	case '"':
		tok = l.newToken(token.STRING, l.readString())

	case 0:
		tok = l.newToken(token.EOF, "EOF")

	default:
		switch {
		case isLetter(l.ch):
			identifier := l.readIdentifier()
			return l.newToken(token.LookupIdent(identifier), identifier)

		case isDigit(l.ch):
			return l.newToken(token.INT, l.readNumber())

		default:
			tok = l.newTokenByType(token.ILLEGAL)
		}
	}

	l.advance()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: l.line}
}

func (l *Lexer) newTokenByType(tokenType token.TokenType) token.Token {
	return l.newToken(tokenType, string(l.ch))
}

func (l *Lexer) switch2(tok1, tok2 token.TokenType) token.Token {
	if l.peek() == '=' {
		l.advance()
		return l.newToken(tok2, tok2.String())
	}
	return l.newToken(tok1, tok1.String())
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

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.advance()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	// TODO support decimals
	start := l.position
	for isDigit(l.ch) {
		l.advance()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() string {
	l.advance() // consume opening '"'

	start := l.position
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\n' {
			l.line++
		}
		l.advance()
	}
	// TODO handle unterminated strings
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.ch {
		case ' ', '\t', '\r':
			l.advance()

		case '\n':
			l.line++
			l.advance()

		case '/':
			if l.peek() == '/' {
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
