package eval

import (
	"fmt"
	"os"
	"path/filepath"
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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if err := testBooleanObject(evaluated, tt.expected); err != nil {
			t.Error(err)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!0", false},
		{"!null", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!0", true},
		{"!!null", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if err := testBooleanObject(evaluated, tt.expected); err != nil {
			t.Error(err)
		}
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if null { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
		{"if null { 10 } else { 20 }", 20},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			if err := testIntegerObject(evaluated, int64(integer)); err != nil {
				t.Error(err)
			}
		} else {
			if err := testNullObject(evaluated); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestYeetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"yeet 10;", 10},
		{"yeet 10; 9;", 10},
		{"yeet 2 * 5; 9;", 10},
		{"9; yeet 2 * 5; 9;", 10},
		{
			`
if 10 > 1 {
	if 10 > 1 {
		yeet 10;
	}
	yeet 1;
}`,
			10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if err := testIntegerObject(evaluated, tt.expected); err != nil {
			t.Error(err)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
	if (10 > 1) {
		yeet true + false;
	}
	yeet 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Msg != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Msg)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		if err := testIntegerObject(testEval(tt.input), tt.expected); err != nil {
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
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNullObject(obj object.Object) error {
	if obj != NULL {
		return fmt.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
	}
	return nil
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

//
// BENCHMARKS
//

func TestFiles(t *testing.T) {
	for _, f := range []string{
		"first.yeet",
		// "fib.yeet",
	} {
		t.Run(f, func(t *testing.T) {
			filename := filepath.Join("../examples", f)
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("couldn't read test file: %s", err)
			}
			sSrc := string(src)
			// l := lexer.New(sSrc)
			// p := parser.New(l)
			// program := p.ParseProgram()
			// checkParserErrors(t, p)
			result := testEval(sSrc)

			if evalError, ok := result.(*object.Error); ok {
				t.Errorf("evaluated to error: %q", evalError.Msg)
			}
		})
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()

	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func BenchmarkEval(b *testing.B) {
	for _, f := range []string{
		"fun.yeet",
		"first.yeet",
	} {
		b.Run(f, func(b *testing.B) {
			b.StopTimer()
			filename := filepath.Join("../examples", f)
			src, err := os.ReadFile(filename)
			if err != nil {
				b.Fatalf("couldn't read test file: %s", err)
			}
			sSrc := string(src)
			l := lexer.New(sSrc)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				_ = Eval(program, env)
			}
		})
	}
}
