package lexer

import (
	"testing"

	"ylang/token"
)

func TestNextToken(t *testing.T) {
	// = `=+()){},;`

	// type testInner struct {
	// 	expectedType    token.TokenType
	// 	expectedLiteral string
	// }

	// type myTest struct {
	// 	input  string
	// 	result []struct {
	// 		expectedType    token.TokenType
	// 		expectedLiteral string
	// 	}
	// }

	// tests := []myTest{
	// 	// {input: "+", result: []testInner{testInner{expectedType: token.PLUS, expectedLiteral: "+"}}},
	// 	// {input: "+", result: []testInner{{expectedType: token.PLUS, expectedLiteral: "+"}}},
	// 	// {"+", []testInner{{expectedType: token.PLUS, expectedLiteral: "+"}}},

	// 	{input: "+", result: {{expectedType: token.PLUS, expectedLiteral: "+"}}},
	// 	// {
	// 	// 	"+", result{
	// 	// 		{expectedType: token.PLUS, expectedLiteral: "+"},
	// 	// 	},
	// 	// },
	// }

	type expected struct {
		ttype   token.TokenType
		literal string
	}

	type myTest struct {
		input  string
		result []expected
	}

	testCases := []myTest{
		// {input: "+", result: []testInner{testInner{expectedType: token.PLUS, expectedLiteral: "+"}}},
		// {input: "+", result: []testInner{{expectedType: token.PLUS, expectedLiteral: "+"}}},
		// {"+", []testInner{{expectedType: token.PLUS, expectedLiteral: "+"}}},
		{
			"+,(){}", []expected{
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
			"yo dawg = 5;", []expected{
				{token.YO, "yo"},
				{token.IDENT, "dawg"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
	}

	// tests := []struct {
	// 	input  string
	// 	result []struct {
	// 		expectedType    token.TokenType
	// 		expectedLiteral string
	// 	}
	// }{

	// tests := []myTest{
	// 	// tests := []struct {
	// 	// 	input  string
	// 	// 	result []struct {
	// 	// 		expectedType    token.TokenType
	// 	// 		expectedLiteral string
	// 	// 	}
	// 	// }{
	// 	{
	// 		input: `=+()){},;`, {
	// 			// {expectedType: token.PLUS, expectedLiteral: "+"},

	// 			{token.ASSIGN, "="},
	// 			// {token.PLUS, "+"},
	// 			// {token.LPAREN, "("},
	// 			// {token.RPAREN, ")"},
	// 			// {token.LBRACE, "{"},
	// 			// {token.RBRACE, "}"},
	// 			// {token.COMMA, ","},
	// 			// {token.SEMICOLON, ";"},
	// 			// {token.EOF, ""},
	// 		},
	// 	},
	// }

	for _, tc := range testCases {
		l := New(tc.input)

		for i, expected := range tc.result {
			tok := l.NextToken()

			if tok.Type != expected.ttype {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, expected.ttype, tok.Type)
			}
			if tok.Literal != expected.literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, expected.literal, tok.Literal)
			}
		}
	}
}
