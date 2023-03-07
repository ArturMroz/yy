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

	tok := token.Token{Line: l.line}

	switch l.ch {
	case '<':
		tok = l.newToken(token.LT)
	case '>':
		tok = l.newToken(token.GT)
	case ';':
		tok = l.newToken(token.SEMICOLON)
	case '(':
		tok = l.newToken(token.LPAREN)
	case ')':
		tok = l.newToken(token.RPAREN)
	case ',':
		tok = l.newToken(token.COMMA)
	case '{':
		tok = l.newToken(token.LBRACE)
	case '}':
		tok = l.newToken(token.RBRACE)
	case '[':
		tok = l.newToken(token.LBRACKET)
	case ']':
		tok = l.newToken(token.RBRACKET)
	case '\\':
		tok = l.newToken(token.BACKSLASH)

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

	case '%':
		switch l.peek() {
		case '{':
			l.advance()
			tok = l.newToken(token.HASHMAP)
		case '=':
			l.advance()
			tok = l.newToken(token.MOD_ASSIGN)
		default:
			tok = l.newToken(token.PERCENT)
		}

	case '.':
		if l.peek() == '.' {
			l.advance()
			tok = l.newToken(token.RANGE)
		}

	case '&':
		if l.peek() == '&' {
			l.advance()
			tok = l.newToken(token.AND)
		}

	case '|':
		if l.peek() == '|' {
			l.advance()
			tok = l.newToken(token.OR)
		}

	case '@':
		if l.peek() == '\\' {
			l.advance()
			tok = l.newToken(token.MACRO)
		}

	case '"':
		tok = l.newTokenWithLiteral(token.STRING, l.readString())

	case 0:
		tok = l.newToken(token.EOF)

	default:
		switch {
		case isLetter(l.ch):
			return l.identifier()

		case isDigit(l.ch):
			return l.number()

		default:
			tok = l.newTokenWithLiteral(token.ILLEGAL, string(l.ch))
		}
	}

	l.advance()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType) token.Token {
	return token.Token{Type: tokenType, Literal: tokenType.String(), Line: l.line}
}

func (l *Lexer) newTokenWithLiteral(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: l.line}
}

func (l *Lexer) switch2(tok1, tok2 token.TokenType) token.Token {
	if l.peek() == '=' {
		l.advance()
		return l.newToken(tok2)
	}
	return l.newToken(tok1)
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

func (l *Lexer) identifier() token.Token {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.advance()
	}

	ident := l.input[start:l.position]
	return l.newTokenWithLiteral(token.LookupIdent(ident), ident)
}

func (l *Lexer) number() token.Token {
	start := l.position
	for isDigit(l.ch) {
		l.advance()
	}

	if l.ch == '.' && isDigit(l.peek()) {
		l.advance() // dot
		for isDigit(l.ch) {
			l.advance()
		}

		return l.newTokenWithLiteral(token.NUMBER, l.input[start:l.position])
	} else {
		return l.newTokenWithLiteral(token.INT, l.input[start:l.position])
	}
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
