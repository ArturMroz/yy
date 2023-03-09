package ast

import (
	"testing"

	"yy/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{Type: token.WALRUS, Literal: ":="},
				Expression: &DeclareExpression{
					Token: token.Token{Type: token.WALRUS, Literal: ":="},
					Name: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "18"},
							Value: 18,
						},
						Operator: "+",
					},
				},
			},
		},
	}

	expected := "(myVar := (5 + 18));"
	if program.String() != expected {
		t.Errorf("program.String() wrong. want %q, got %q", expected, program.String())
	}
}
