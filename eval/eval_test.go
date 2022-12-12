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

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	expected := "Hello World!"
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != expected {
		t.Errorf("String has wrong value. want=%q, got=%q", expected, str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"con" + "cat"`, "concat"},
		{`"" + "cat"`, "cat"},
		{`"" + ""`, ""},
		{`"con" + "cat" + "enation"`, "concatenation"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}
		if str.Value != tt.expected {
			t.Errorf("String has wrong value. want=%q, got=%q", tt.expected, str.Value)
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
		{`"yolo" == "yolo"`, true},
		{`"yolo" == "yeet"`, false},
		{`"yolo" != "yolo"`, false},
		{`"yolo" != "yeet"`, true},
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

func TestArrayLiterals(t *testing.T) {
	// TODO add more test cases
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	if err := testIntegerObject(result.Elements[0], 1); err != nil {
		t.Error(err)
	}
	if err := testIntegerObject(result.Elements[1], 4); err != nil {
		t.Error(err)
	}
	if err := testIntegerObject(result.Elements[2], 6); err != nil {
		t.Error(err)
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		// out of bounds access returns nil
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			if err := testIntegerObject(evaluated, int64(expected)); err != nil {
				t.Error(err)
			}
		case nil:
			if err := testNullObject(evaluated); err != nil {
				t.Error(err)
			}
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
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			"foobar()",
			"identifier not found: foobar",
		},
		{
			"foobar2(x, y)",
			"identifier not found: foobar2",
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
			t.Errorf("wrong error message. want=%q, got=%q", tt.expectedMessage, errObj.Msg)
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

func TestFunctionObject(t *testing.T) {
	input := "fun(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "{ (x + 2) }"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fun(x) { x; }; identity(5);", 5},
		{"let identity = fun(x) { yeet x; }; identity(5);", 5},
		{"let double = fun(x) { x * 2; }; double(5);", 10},
		{"let add = fun(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fun(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fun(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		if err := testIntegerObject(testEval(tt.input), tt.expected); err != nil {
			t.Error(err)
		}
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fun(x) {
fun(y) { x + y };
};
let addTwo = newAdder(2);
addTwo(2);`

	if err := testIntegerObject(testEval(input), 4); err != nil {
		t.Error(err)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`assert(1 == 1)`, nil},
		{`assert(1 == 2)`, "assert failed"},
		{`assert(false)`, "assert failed"},
		{`assert(true)`, nil},
		{`let a = 5; let b = 6; assert(a == b, "expect a to be equal to b")`, "assert failed: expect a to be equal to b"},
		{`assert(1 == 2, "one is different than two")`, "assert failed: one is different than two"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			if err := testIntegerObject(evaluated, int64(expected)); err != nil {
				t.Error(err)
			}
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Msg != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Msg)
			}
		case nil:
			if _, ok := evaluated.(*object.Null); !ok {
				t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
			}

		default:
			t.Errorf("unexpected type, got=%T", expected)
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

var testFiles = []string{
	"first.yeet",
	"fun.yeet",
}

func TestFiles(t *testing.T) {
	for _, f := range testFiles {
		t.Run(f, func(t *testing.T) {
			filename := filepath.Join("../examples", f)
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("couldn't read test file: %s", err)
			}
			sSrc := string(src)
			result := testEval(sSrc)

			if evalError, ok := result.(*object.Error); ok {
				t.Errorf("evaluated to error: %q", evalError.Msg)
			}
		})
	}
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
