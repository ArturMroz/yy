package lexer_test

import (
	"testing"

	"yy/lexer"
	"yy/token"
)

type lexerTestCase struct {
	input    string
	expected []token.Token
}

func TestNextToken(t *testing.T) {
	runLexerTests(t, []lexerTestCase{
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
				{Type: token.EOF, Literal: "EOF"},
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
				{Type: token.EOF, Literal: "EOF"},
			},
		},
		{
			`35; 38.23; -4`,
			[]token.Token{
				{Type: token.INT, Literal: "35"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.NUMBER, Literal: "38.23"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.INT, Literal: "4"},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
		{
			`%{"key": "value"}`, // hash map
			[]token.Token{
				{Type: token.HASHMAP, Literal: "%{"},
				{Type: token.STRING, Literal: "key"},
				{Type: token.COLON, Literal: ":"},
				{Type: token.STRING, Literal: "value"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: "EOF"},
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
				{Type: token.EOF, Literal: "EOF"},
			},
		},
		{
			`@\x y { x + y }`,
			[]token.Token{
				{Type: token.MACRO, Literal: "@\\"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
	})
}

func TestLexingYifExpression(t *testing.T) {
	runLexerTests(t, []lexerTestCase{
		{
			"yif 5 > 8 { 5 } yels { 8 }",
			[]token.Token{
				{Type: token.YIF, Literal: "yif"},
				{Type: token.INT, Literal: "5"},
				{Type: token.GT, Literal: ">"},
				{Type: token.INT, Literal: "8"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.INT, Literal: "5"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.YELS, Literal: "yels"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.INT, Literal: "8"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
	})
}

func TestLexingInterpolatedStrings(t *testing.T) {
	runLexerTests(t, []lexerTestCase{
		{
			"`interpolate { 1 + 2 } right now { foo } please`",
			[]token.Token{
				{Type: token.TEMPL_STRING, Literal: "interpolate "},
				{Type: token.INT, Literal: "1"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.INT, Literal: "2"},
				{Type: token.TEMPL_STRING, Literal: " right now "},
				{Type: token.IDENT, Literal: "foo"},
				{Type: token.STRING, Literal: " please"},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
		{
			"`{ 1 + 2 } { foo }{banana}`",
			[]token.Token{
				{Type: token.TEMPL_STRING, Literal: ""},
				{Type: token.INT, Literal: "1"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.INT, Literal: "2"},
				{Type: token.TEMPL_STRING, Literal: " "},
				{Type: token.IDENT, Literal: "foo"},
				{Type: token.TEMPL_STRING, Literal: ""},
				{Type: token.IDENT, Literal: "banana"},
				{Type: token.STRING, Literal: ""},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
		{
			"`i'm {{age}} yr old`",
			[]token.Token{
				{Type: token.STRING, Literal: "i'm {age} yr old"},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
		{
			"`var {{age}} = {age}. And that's a bracket }}.`",
			[]token.Token{
				{Type: token.TEMPL_STRING, Literal: "var {age} = "},
				{Type: token.IDENT, Literal: "age"},
				{Type: token.STRING, Literal: ". And that's a bracket }."},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
		{
			"`i'm { yif 5 > 8 { 5 } yels { 8 } } yr old`",
			[]token.Token{
				{Type: token.TEMPL_STRING, Literal: "i'm "},
				{Type: token.YIF, Literal: "yif"},
				{Type: token.INT, Literal: "5"},
				{Type: token.GT, Literal: ">"},
				{Type: token.INT, Literal: "8"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.INT, Literal: "5"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.YELS, Literal: "yels"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.INT, Literal: "8"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.STRING, Literal: " yr old"},
				{Type: token.EOF, Literal: "EOF"},
			},
		},
	})
}

func runLexerTests(t *testing.T, testCases []lexerTestCase) {
	for _, tc := range testCases {
		l := lexer.New(tc.input)

		for i, exp := range tc.expected {
			tok := l.NextToken()

			if tok.Type != exp.Type {
				t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, exp.Type, tok.Type)
			}
			if tok.Literal != exp.Literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, exp.Literal, tok.Literal)
			}
		}
	}
}
