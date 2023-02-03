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
			"+,(){}[]",
			[]token.Token{
				{token.PLUS, "+"},
				{token.COMMA, ","},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RBRACE, "}"},
				{token.LBRACKET, "["},
				{token.RBRACKET, "]"},
				{token.EOF, ""},
			},
		},
		{
			"yo dawg = 5;",
			[]token.Token{
				{token.YO, "yo"},
				{token.IDENT, "dawg"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			`yo myString = "testy guy"; "other string"`,
			[]token.Token{
				{token.YO, "yo"},
				{token.IDENT, "myString"},
				{token.ASSIGN, "="},
				{token.STRING, "testy guy"},
				{token.SEMICOLON, ";"},
				{token.STRING, "other string"},
				{token.EOF, ""},
			},
		},
		{
			`{"key": "value"}`, // hash map
			[]token.Token{
				{token.LBRACE, "{"},
				{token.STRING, "key"},
				{token.COLON, ":"},
				{token.STRING, "value"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			`\(x, y) { x + y }`,
			[]token.Token{
				{token.BACKSLASH, "\\"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.PLUS, "+"},
				{token.IDENT, "y"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
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
