package ast_test

import (
	"testing"

	"yy/ast"
	"yy/token"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: token.Token{Type: token.WALRUS, Literal: ":="},
				Expression: &ast.DeclareExpression{
					Token: token.Token{Type: token.WALRUS, Literal: ":="},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &ast.InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
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
