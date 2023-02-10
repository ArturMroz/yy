package ast

import (
	"testing"

	"ylang/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{Type: token.WALRUS, Literal: ":="},
				Expression: &AssignExpression{
					Token: token.Token{Type: token.WALRUS, Literal: ":="},
					Name: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
						Value: "anotherVar",
					},
				},
			},
		},
	}

	expected := "(myVar := anotherVar)"
	if program.String() != expected {
		t.Errorf("program.String() wrong. want=%q, got=%q", expected, program.String())
	}
}
