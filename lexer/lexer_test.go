package lexer

import (
	"testing"

	"ylang/token"
)

func TestNextToken(t *testing.T) {
	type lexerTestCase struct {
		input    string
		expected []token.Token
	}

	testCases := []lexerTestCase{
		// {input: "+", result: []testInner{testInner{expectedType: token.PLUS, expectedLiteral: "+"}}},
		// {input: "+", result: []testInner{{expectedType: token.PLUS, expectedLiteral: "+"}}},
		// {"+", []testInner{{expectedType: token.PLUS, expectedLiteral: "+"}}},
		{
			"+,(){}", []token.Token{
				{token.PLUS, "+"},
				{token.COMMA, ","},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			"yo dawg = 5;", []token.Token{
				{token.YO, "yo"},
				{token.IDENT, "dawg"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
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
