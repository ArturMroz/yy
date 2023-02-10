package lexer

import (
	"testing"

	"ylang/token"
)

func TestNextToken(t *testing.T) {
	testCases := []struct {
		input    string
		expected []token.Token
	}{
		{
			"+-*/,(){}[]",
			[]token.Token{
				{Type: token.PLUS, Literal: "+"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.ASTERISK, Literal: "*"},
				{Type: token.SLASH, Literal: "/"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.LBRACKET, Literal: "["},
				{Type: token.RBRACKET, Literal: "]"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			`myString := "testy guy"; second = "other string"`,
			[]token.Token{
				{Type: token.IDENT, Literal: "myString"},
				{Type: token.WALRUS, Literal: ":="},
				{Type: token.STRING, Literal: "testy guy"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.IDENT, Literal: "second"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.STRING, Literal: "other string"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			`{"key": "value"}`, // hash map
			[]token.Token{
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.STRING, Literal: "key"},
				{Type: token.COLON, Literal: ":"},
				{Type: token.STRING, Literal: "value"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			`\(x, y) { x + y }`,
			[]token.Token{
				{Type: token.BACKSLASH, Literal: "\\"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tc := range testCases {
		l := New(tc.input)

		for i, exp := range tc.expected {
			tok := l.NextToken()

			if tok.Type != exp.Type {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, exp.Type, tok.Type)
			}
			if tok.Literal != exp.Literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, exp.Literal, tok.Literal)
			}
		}
	}
}
