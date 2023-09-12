package lexer

import "yy/token"

type Lexer struct {
	Input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	line         int  // current line
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{Input: input, line: 1}
	l.advance()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.ch {
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
		tok = l.switchEq(token.PLUS, token.ADD_ASSIGN)
	case '-':
		tok = l.switchEq(token.MINUS, token.SUB_ASSIGN)
	case '*':
		tok = l.switchEq(token.ASTERISK, token.MUL_ASSIGN)
	case '/':
		tok = l.switchEq(token.SLASH, token.DIV_ASSIGN)
	case '=':
		tok = l.switchEq(token.ASSIGN, token.EQ)
	case '!':
		tok = l.switchEq(token.BANG, token.NOT_EQ)
	case ':':
		tok = l.switchEq(token.COLON, token.WALRUS)
	case '>':
		tok = l.switchEq(token.GT, token.GT_EQ)

	case '.':
		tok = l.switch2(token.DOT, token.RANGE, '.')
	case '&':
		tok = l.switch2(token.AMPERSAND, token.AND, '&')
	case '|':
		tok = l.switch2(token.PIPE, token.OR, '|')
	case '@':
		tok = l.switch2(token.AT, token.MACRO, '\\')

	case '<':
		switch l.peek() {
		case '=':
			l.advance()
			tok = l.newToken(token.LT_EQ)
		case '<':
			l.advance()
			tok = l.newToken(token.LT_LT)
		default:
			tok = l.newToken(token.LT)
		}

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

	case '"':
		tok = l.string()

	case 0:
		tok = l.newToken(token.EOF)

	default:
		switch {
		case IsLetter(l.ch):
			return l.identifier()

		case IsDigit(l.ch):
			return l.number()

		default:
			tok = l.newTokenWithLiteral(token.ILLEGAL, string(l.ch))
		}
	}

	l.advance()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType) token.Token {
	return l.newTokenWithLiteral(tokenType, tokenType.String())
}

func (l *Lexer) newTokenWithLiteral(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: l.line, Offset: l.position - len(literal) + 1}
}

func (l *Lexer) switch2(tok1, tok2 token.TokenType, expected byte) token.Token {
	if l.peek() == expected {
		l.advance()
		return l.newToken(tok2)
	}
	return l.newToken(tok1)
}

func (l *Lexer) switchEq(tok1, tok2 token.TokenType) token.Token {
	return l.switch2(tok1, tok2, '=')
}

func (l *Lexer) advance() {
	if l.readPosition >= len(l.Input) {
		l.ch = 0
	} else {
		l.ch = l.Input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.Input) {
		return 0
	}
	return l.Input[l.readPosition]
}

func (l *Lexer) identifier() token.Token {
	start := l.position
	for IsLetter(l.ch) || IsDigit(l.ch) {
		l.advance()
	}

	ident := l.Input[start:l.position]
	return token.Token{Type: token.LookupIdent(ident), Literal: ident, Offset: start}
}

func (l *Lexer) number() token.Token {
	start := l.position
	for IsDigit(l.ch) {
		l.advance()
	}

	if l.ch == '.' && IsDigit(l.peek()) {
		l.advance() // dot
		for IsDigit(l.ch) {
			l.advance()
		}

		return token.Token{Type: token.NUMBER, Literal: l.Input[start:l.position], Offset: start}
	} else {
		return token.Token{Type: token.INT, Literal: l.Input[start:l.position], Offset: start}
	}
}

func (l *Lexer) string() token.Token {
	l.advance() // consume opening '"'

	start := l.position
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\n' {
			l.line++
		}
		l.advance()
	}

	if l.ch == 0 {
		return token.Token{Type: token.ERROR, Literal: "unterminated string", Offset: start}
	}

	return token.Token{Type: token.STRING, Literal: l.Input[start:l.position], Offset: start}
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

func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
