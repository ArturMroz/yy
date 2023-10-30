package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},
		{
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: one()},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: two()},
				},
			},
		},
		{
			&InfixExpression{Left: one(), Operator: "+", Right: two()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&InfixExpression{Left: two(), Operator: "+", Right: one()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpression{Operator: "-", Right: one()},
			&PrefixExpression{Operator: "-", Right: two()},
		},
		{
			&IndexExpression{Left: one(), Index: one()},
			&IndexExpression{Left: two(), Index: two()},
		},
		{
			&YifExpression{
				Condition: one(),
				Consequence: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
				Alternative: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&YifExpression{
				Condition: two(),
				Consequence: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
				Alternative: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&YoyoExpression{
				Condition: one(),
				Body: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&YoyoExpression{
				Condition: two(),
				Body: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&YallExpression{
				Iterable: &ArrayLiteral{
					Elements: []Expression{one()},
				},
				Body: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&YallExpression{
				Iterable: &ArrayLiteral{
					Elements: []Expression{two()},
				},
				Body: &BlockExpression{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&YeetStatement{ReturnValue: one()},
			&YeetStatement{ReturnValue: two()},
		},
		{
			&DeclareExpression{Value: one()},
			&DeclareExpression{Value: two()},
		},
		{
			&ArrayLiteral{Elements: []Expression{one(), one(), two()}},
			&ArrayLiteral{Elements: []Expression{two(), two(), two()}},
		},
		{
			&RangeLiteral{Start: one(), End: one()},
			&RangeLiteral{Start: two(), End: two()},
		},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)

		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expected)
		}
	}

	hashLiteral := &HashmapLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}

	Modify(hashLiteral, turnOneIntoTwo)

	for key, val := range hashLiteral.Pairs {
		key, _ := key.(*IntegerLiteral)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, key.Value)
		}
		val, _ := val.(*IntegerLiteral)
		if val.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, val.Value)
		}
	}
}
