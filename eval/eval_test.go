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

	if err := testStringObject(evaluated, expected); err != nil {
		t.Error(err)
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
		{`"[1, 2, 3]" == "[1, 2, 3]"`, true},
		{`"[1, 2, 3]" == "[1, 2, 9]"`, false},
		{`"[1, 2, 3]" != "[1, 2, 3]"`, false},
		{`"[1, 2, 3]" != "[1, 2, 9]"`, true},
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

func TestIntegerArrayLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{
			"[1, 2, 3, 4, 5]",
			[]int64{1, 2, 3, 4, 5},
		},
		{
			"[1, 2 * 2, 3 + 3]",
			[]int64{1, 4, 6},
		},
		{
			"[4 / 2, 5 - 1, 8 * 4]",
			[]int64{2, 4, 32},
		},
		{
			"[1 + 1, 2 + 2, 3 + 3]",
			[]int64{2, 4, 6},
		},
		{
			"[1 + 2 * 3 2 + 2 / 2 3 + 3]",
			[]int64{7, 3, 6},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		array, ok := evaluated.(*object.Array)
		if !ok {
			t.Errorf("exp not *object.Array. got=%T", array)
			continue
		}

		for i := range tt.expected {
			if err := testIntegerObject(array.Elements[i], tt.expected[i]); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"[1 2 3][2]", 3},
		{"i := 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{
			"myArray := [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"myArray := [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"myArray := [1, 2, 3]; i := myArray[0]; myArray[i]",
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

func TestHashLiterals(t *testing.T) {
	input := `
let two = "two";
{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
}`
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key %q in Pairs", expectedKey)
		}
		if err := testIntegerObject(pair.Value, expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
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
		{"result := if null { 10 } else { 20 }; result", 20},
		{"5 + if null { 10 } else { 20 }", 25},
		{"if null { 10 } else { 20 } * 2", 40},
		{"5 + if null { 10 } else { 20 } * 2", 45},
		{"result := 3 * if null { 10 } else { 20 } + 9; result", 69},
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

func TestYoyoExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"i := 0; yoyo ; i < 5; { i = i + 1 }", 5},
		{"yoyo i := 0; i < 5; { i = i + 1 }", 5},
		{"yoyo i := 0; i < 5; i = i + 1 { i }", 4},
		{"i := 69; yoyo i := 0; i < 5; i = i + 1 { i }; i", 69},
		// {"i := 69; yoyo i = 0; i < 5; i = i + 1 { i }; i", 5}, // TODO fix test
		{"result := (yoyo i := 0; i < 5; i = i + 1 { i }); result", 4},
		{"result := yoyo i := 0; i < 5; i = i + 1 { i }; result", 4},

		// TODO test error handling
		// {"yoyo i = 0; i < 5; i = i + 1 { i }",  "identifier not found: i"},
		// {"yoyo i := 0; i < 5; i = i + 1 { i }; i", "identifier not found: i"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if err := testIntegerObject(evaluated, tt.expected); err != nil {
			t.Error(err)
		}
	}
}

func TestYoniExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"yoni [1 2 3] { yt }", 3},
		{"arr := [1 2 3]; yoni arr { yt }", 3},
		{`yoni "testme" { yt }`, "e"},
		{`my_str := "swag"; yoni my_str { yt }`, "g"},
		// {`yoni 0..5 { yt }`, 5},
		// {`yoni i : 0..5 { i }`, 5},
		// {`yoni v, i : 0..5 { yt }`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			if err := testIntegerObject(evaluated, int64(expected)); err != nil {
				t.Error(err)
			}
		case string:
			if err := testStringObject(evaluated, expected); err != nil {
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
		{
			`{"name": "Monkey"}[fun(x) { x }];`,
			"key not hashable: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if err := testErrorObject(evaluated, tt.expectedMessage); err != nil {
			t.Error(err)
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

func TestAssignExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"a := 8", 8},
		{"a := 8; a", 8},
		{"a := 8 * 5 + 3 / 2 - 2 * (2 + 3) * 3", 11},
		{"a := 8; b := a", 8},
		{"a := 8; b := a + 2", 10},
		{"a := 8; b := 2; c := a + b", 10},
		{"a := 8 == 5", false},
		{"a := 8 != 5", true},
		{"a := 8 > 5", true},

		{"a := 8; a = 15", 15},
		{"a := 8; b := 2; a = b", 2},

		{"a = 8", "identifier not found: a"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			if err := testIntegerObject(evaluated, int64(expected)); err != nil {
				t.Error(err)
			}
		case bool:
			if err := testBooleanObject(evaluated, expected); err != nil {
				t.Error(err)
			}
		case string: // errors
			if err := testErrorObject(evaluated, expected); err != nil {
				t.Error(err)
			}
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

func TestLambdaObject(t *testing.T) {
	tests := []struct {
		input  string
		params []string
		body   string
	}{
		{`\() { 69 };`, []string{}, "{ 69 }"},
		{`\(x) { x + 2; };`, []string{"x"}, "{ (x + 2) }"},
		{`\(x, y) { x + y };`, []string{"x", "y"}, "{ (x + y) }"},
		{`\(x, y, z) { x * y - z };`, []string{"x", "y", "z"}, "{ ((x * y) - z) }"},

		{`\ { 69 };`, []string{}, "{ 69 }"},
		{`\x { x + 2; };`, []string{"x"}, "{ (x + 2) }"},
		{`\x, y { x + y };`, []string{"x", "y"}, "{ (x + y) }"},
		{`\x, y, z { x * y - z };`, []string{"x", "y", "z"}, "{ ((x * y) - z) }"},

		{`\x y { x + y }`, []string{"x", "y"}, "{ (x + y) }"},
		{`\x y z { x * y - z }`, []string{"x", "y", "z"}, "{ ((x * y) - z) }"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		fn, ok := evaluated.(*object.Function)
		if !ok {
			t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
		}
		if len(fn.Parameters) != len(tt.params) {
			t.Fatalf("lambda has wrong parameters. got=%+v, expected=%+v", fn.Parameters, tt.params)
		}
		for i, param := range tt.params {
			if fn.Parameters[i].String() != param {
				t.Fatalf("parameter is not '%s'. got=%q", param, fn.Parameters[i])
			}
		}
		if fn.Body.String() != tt.body {
			t.Fatalf("body is not %q. got=%q", tt.body, fn.Body.String())
		}
	}
}

func TestLambdaApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`let nope = \() { 69 }; nope();`, 69},
		{`let identity = \(x) { x; }; identity(5);`, 5},
		{`let identity = \(x) { yeet x; }; identity(5);`, 5},
		{`let double = \(x) { x * 2; }; double(5);`, 10},
		{`let add = \(x, y) { x + y; }; add(5, 5);`, 10},
		{`let add = \(x, y) { x + y; }; add(5 + 5, add(5, 5));`, 20},
		{`let add = \(x, y, z) { x + y + z; }; add(1, 2, 3);`, 6},

		{`let nope = \ { 69 }; nope();`, 69},
		{`let add = \x, y { x + y; }; add(5, 5);`, 10},
		{`let add = \x, y { x + y; }; add(5 + 5, add(5, 5));`, 20},
		{`let add = \x, y, z { x + y + z; }; add(1, 2, 3);`, 6},

		{`add := \x y { x + y; }; add(5, 5);`, 10},
		{`add := \x y { x + y; }; add(5+5 add(5, 5));`, 20},
		{`add := \x y z { x + y + z }; add(1 2 3);`, 6},

		{`add := \x y { x + y; }; add(5, 5);`, 10},
		{`add := \x y { x + y; }; add(5+5 add(5, 5));`, 20},
		{`add := \x y z { x + y + z }; add(1 2 3);`, 6},

		{`\(x) { x; }(5)`, 5},
		{`\x { x }(5)`, 5},
		{`\x, y { x * y }(3, 5)`, 15},
		{`\x y { x * y }(3 5)`, 15},
		{`\x y { x * y }(3+2 5+1)`, 30},
		{`\x y { x * y }(3 + 2 5 + 1)`, 30},
		{`\x y { x * y }(3 + 2, 5 + 1)`, 30},
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
		return fmt.Errorf("Integer object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testStringObject(obj object.Object, expected string) error {
	result, ok := obj.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("String has wrong value. want=%q, got=%q", expected, result.Value)
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

func testErrorObject(obj object.Object, expectedMsg string) error {
	result, ok := obj.(*object.Error)
	if !ok {
		return fmt.Errorf("object is not Error. got=%T (%+v)", obj, obj)
	}
	if result.Msg != expectedMsg {
		return fmt.Errorf("wrong error message. want=%q, got=%q", expectedMsg, result.Msg)
	}
	return nil
}

//
// PROGRAMS FROM EXAMPLES/
//

const examplesDir = "../examples"

func TestExampleFiles(t *testing.T) {
	testFiles, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("couldn't read example files dir: %s", err)
	}

	for _, f := range testFiles {
		t.Run(f.Name(), func(t *testing.T) {
			filename := filepath.Join(examplesDir, f.Name())
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("couldn't read test file: %s", err)
			}

			result := testEval(string(src))
			if evalError, ok := result.(*object.Error); ok {
				t.Errorf("runtime error: %q", evalError.Msg)
			}
		})
	}
}

//
// BENCHMARKS
//

func BenchmarkEval(b *testing.B) {
	testFiles, err := os.ReadDir(examplesDir)
	if err != nil {
		b.Fatalf("couldn't read example files dir: %s", err)
	}

	for _, f := range testFiles {
		b.Run(f.Name(), func(b *testing.B) {
			b.StopTimer()
			filename := filepath.Join(examplesDir, f.Name())
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
