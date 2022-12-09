package eval

import (
	"fmt"
	"testing"

	"ylang/lexer"
	"ylang/object"
	"ylang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if err := testIntegerObject(evaluated, tt.expected); err != nil {
			t.Error(err)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if err := testBooleanObject(evaluated, tt.expected); err != nil {
			t.Error(err)
		}
	}
}

//
// HELPERS
//

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObject(obj object.Object, expected int64) error {
	result, ok := obj.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testBooleanObject(obj object.Object, expected bool) error {
	result, ok := obj.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}
	return nil
}
