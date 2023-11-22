package lexer

import (
	"yy/token"
)

type Lexer struct {
	Input        string
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           byte   // current char under examination
	numBrackets  int    // depth of string interpolation
	brackets     [5]int // stack of interpolations
}

func New(input string) *Lexer {
	l := &Lexer{Input: input}
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
	case '[':
		tok = l.newToken(token.LBRACKET)
	case ']':
		tok = l.newToken(token.RBRACKET)
	case '\\':
		tok = l.newToken(token.BACKSLASH)

	case '{':
		if l.numBrackets > 0 {
			l.brackets[l.numBrackets-1]++
		}

		tok = l.newToken(token.LBRACE)

	case '}':
		if l.numBrackets > 0 {
			l.brackets[l.numBrackets-1]--

			if l.brackets[l.numBrackets-1] == 0 {
				l.numBrackets--
				tok = l.readString()
				break
			}
		}

		tok = l.newToken(token.RBRACE)

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
		tok = l.readString()

	case 0:
		tok = l.newToken(token.EOF)

	default:
		switch {
		case isLetter(l.ch):
			return l.readIdentifier()

		case isDigit(l.ch):
			return l.readNumber()

		default:
			tok = l.newTokenWithLiteral(token.ERROR, "unexpected character: "+string(l.ch))
		}
	}

	l.advance()
	return tok
}

func (l *Lexer) newToken(tokenType token.Type) token.Token {
	return l.newTokenWithLiteral(tokenType, tokenType.String())
}

func (l *Lexer) newTokenWithLiteral(tokenType token.Type, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Offset: l.position - len(literal) + 1}
}

func (l *Lexer) switch2(tok1, tok2 token.Type, expected byte) token.Token {
	if l.peek() == expected {
		l.advance()
		return l.newToken(tok2)
	}
	return l.newToken(tok1)
}

func (l *Lexer) switchEq(tok1, tok2 token.Type) token.Token {
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

func (l *Lexer) readIdentifier() token.Token {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.advance()
	}

	ident := l.Input[start:l.position]
	return token.Token{Type: token.LookupIdent(ident), Literal: ident, Offset: start}
}

func (l *Lexer) readNumber() token.Token {
	start := l.position
	for isDigit(l.ch) {
		l.advance()
	}

	if l.ch == '.' && isDigit(l.peek()) {
		l.advance() // dot
		for isDigit(l.ch) {
			l.advance()
		}

		return token.Token{Type: token.NUMBER, Literal: l.Input[start:l.position], Offset: start}
	}

	return token.Token{Type: token.INT, Literal: l.Input[start:l.position], Offset: start}
}

func (l *Lexer) readString() token.Token {
	l.advance() // consume opening '"' or '}'

	start := l.position
	escapePositions := []int{}

	for {
		switch l.ch {
		case 0:
			return token.Token{
				Type:    token.ERROR,
				Literal: "unterminated string",
				Offset:  start,
			}

		case '"':
			return token.Token{
				Type: token.STRING,
				// Literal: escapeBrackets(l.Input[start:l.position]),
				Literal: l.escapeString(start, escapePositions),
				Offset:  start,
			}

		case '}':
			// double brackets is an escape sequence, ie {{name}}
			if l.peek() == '}' {
				l.advance()
				l.advance()
				escapePositions = append(escapePositions, l.position)
			}

		case '{':
			// double brackets is an escape sequence, ie {{name}}
			if l.peek() == '{' {
				l.advance()
				l.advance()
				escapePositions = append(escapePositions, l.position)
				break
			}

			l.brackets[l.numBrackets] = 1
			l.numBrackets++

			return token.Token{
				Type:    token.TEMPL_STRING,
				Literal: l.escapeString(start, escapePositions),
				Offset:  start,
			}

		default:
			l.advance()
		}
	}
}

func (l *Lexer) escapeString(start int, escapePositions []int) string {
	if len(escapePositions) == 0 {
		return l.Input[start:l.position]
	}

	result := make([]byte, 0, l.position-start-len(escapePositions))
	lastIdx := start

	for _, v := range escapePositions {
		result = append(result, l.Input[lastIdx:v-1]...)
		lastIdx = v
	}

	result = append(result, l.Input[lastIdx:l.position]...)

	return string(result)
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.ch {
		case ' ', '\t', '\r', '\n':
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
