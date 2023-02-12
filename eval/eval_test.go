package eval

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"yy/lexer"
	"yy/object"
	"yy/parser"
)

type evalTestCase struct {
	input    string
	expected any
}

type errmsg struct {
	msg string
}

func TestEvalIntegerExpression(t *testing.T) {
	runEvalTests(t, []evalTestCase{
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
	})
}

func TestStringLiteral(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`"piece of yarn"`, "piece of yarn"},
		{`"Greetings, Earth!"`, "Greetings, Earth!"},
	})
}

func TestStringConcatenation(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`"con" + "cat"`, "concat"},
		{`"" + "cat"`, "cat"},
		{`"" + ""`, ""},
		{`"con" + "cat" + "enation"`, "concatenation"},
	})
}

func TestEvalBooleanExpression(t *testing.T) {
	runEvalTests(t, []evalTestCase{
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
		{`"[1, 2, 3]" == "[4, 5, 6]"`, false},
		{`"[1, 2, 3]" != "[1, 2, 3]"`, false},
		{`"[1, 2, 3]" != "[4, 5, 6]"`, true},
		{"null == null", true},
		{"null != null", false},

		// mixed types
		{"true != null", true},
		{"true == null", false},
		{"2 != null", true},
		{"2 == null", false},
		{"[] != null", true},
		{"[] == null", false},
		{"[1, 2, 3] != null", true},
		{"[1, 2, 3] == null", false},
		{`"Testy McTestface" != null`, true},
		{`"Testy McTestface" == null`, false},
		{`"" != null`, true},
		{`"" == null`, false},
		{`"1" == 1`, false},
		{`"1" != 1`, true},
		{`"1" == [1]`, false},
		{`"1" != [1]`, true},
	})
}

func TestBangOperator(t *testing.T) {
	runEvalTests(t, []evalTestCase{
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
	})
}

func TestIntegerArrayLiterals(t *testing.T) {
	runEvalTests(t, []evalTestCase{
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
	})
}

func TestArrayIndexExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
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
	})
}

func TestHashLiterals(t *testing.T) {
	input := `
two := "two";
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

	evaluated := testEval(t, input)
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
	runEvalTests(t, []evalTestCase{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`key := "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},

		{
			`{"name": "Vars McVariable"}[\x { x }];`,
			errmsg{"key not hashable: FUNCTION"},
		},
	})
}

func TestYifYelsExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"yif true { 10 }", 10},
		{"yif false { 10 }", nil},
		{"yif null { 10 }", nil},
		{"yif 1 { 10 }", 10},
		{"yif 1 < 2 { 10 }", 10},
		{"yif 1 > 2 { 10 }", nil},
		{"yif 1 > 2 { 10 } yels { 20 }", 20},
		{`yif 1 > 2 { "nope" } yels { yif 2 > 5 { "nope" } yels { 20 } }`, 20},
		// {`yif 1 > 2 { "nope } yels yif 2 > 5 { "nope" } yels { 20 }`, 20}, // TODO fix
		{"yif 1 < 2 { 10 } yels { 20 }", 10},
		{"yif null { 10 } yels { 20 }", 20},
		{"result := yif null { 10 } yels { 20 }; result", 20},
		{"5 + yif null { 10 } yels { 20 }", 25},
		{"yif null { 10 } yels { 20 } * 2", 40},
		{"5 + yif null { 10 } yels { 20 } * 2", 45},
		{"result := 3 * yif null { 10 } yels { 20 } + 9; result", 69},
	})
}

func TestYoloExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		// normal operations in yolo mode are still ok
		{"yolo { 2 + 2 }", 4},
		{"yolo { 2 + 2; 8 }", 8},
		{"yolo { a := 1; a }", 1},
		{"yolo { a := 1; b := 2; a + b }", 3},
		{"result := yolo { a := 1; b := 2; a + b }; result", 3},

		// arrays
		{"yolo { 3 + [1, 2, 3] }", []int64{4, 5, 6}},
		{"yolo { [1, 2, 3] + 4 }", []int64{5, 6, 7}},
		{"yolo { 3 * [1, 2, 3] }", []int64{3, 6, 9}},
		{"yolo { [1, 2, 3] * 3 }", []int64{3, 6, 9}},

		// strings
		{`2 + "troll"`, errmsg{"type mismatch: INTEGER + STRING"}},
		{`yolo { 3 * "22" }`, 66},
		{`yolo { "22" * 3 }`, 66},
		{`yolo { 3 * "test" }`, "testtesttest"},
		{`yolo { "test" * 3 }`, "testtesttest"},
		{`yolo { "test" * 0 }`, ""},
		{`yolo { "test" * -5 }`, ABYSS.Value},
		{`yolo { "test" / 0 }`, ABYSS.Value},
		{`yolo { 2 + "troll" }`, "2troll"},

		// what happens in yolo, stays in yolo
		{"yolo { a := 1; a }; a", errmsg{"identifier not found: a"}},
	})
}

func TestYoyoExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"i := 0; yoyo ; i < 5; { i = i + 1 }", 5},
		{"yoyo i := 0; i < 5; { i = i + 1 }", 5},
		{"yoyo i := 0; i < 5; i = i + 1 { i }", 4},
		{"i := 69; yoyo i := 0; i < 5; i = i + 1 { i }; i", 69},
		{"i := 69; yoyo i = 0; i < 5; i = i + 1 { i }; i", 5},
		{"result := (yoyo i := 0; i < 5; i = i + 1 { i }); result", 4},
		{"result := yoyo i := 0; i < 5; i = i + 1 { i }; result", 4},

		{"yoyo i = 0; i < 5; i = i + 1 { i }", errmsg{"identifier not found: i"}},
		{"yoyo i := 0; i < 5; i = i + 1 { i }; i", errmsg{"identifier not found: i"}},
	})
}

func TestYetExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"i := 0; yet i < 5 { i = i + 1 }", 5},
		{"sum := 0; i := 1; yet i < 5 { sum = sum + i; i = i + 1 }; sum", 10},
		{"i := 1; yet false { i = 69 }; i", 1},
	})
}

func TestYallExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"yall [1, 2, 3] { yt }", 3},
		{"arr := [1, 2, 3]; yall arr { yt }", 3},
		{`yall "testme" { yt }`, "e"},
		{`s := ""; yall "test" { s = s + yt + "-" }; s`, "t-e-s-t-"},
		{"sum := 0; yall [1, 2, 3] { sum = sum + yt }; sum", 6},
		{`my_str := "swag"; yall my_str { yt }`, "g"},
		{`yall 0..5 { yt }`, 5},
		{`yall 4..4 { yt }`, 4},
		{`sum := 0; yall 1..4 { sum = sum + yt }; sum`, 10},
		{`yall i: 0..5 { i }`, 5},
		{`sum := 0; yall j: 1..4 { sum = sum + j }; sum`, 10},

		{`yall 0..5 { x }`, errmsg{"identifier not found: x"}},
		{`yall i: 0..5 { yt }`, errmsg{"identifier not found: yt"}},
	})
}

func TestYeetStatements(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"yeet 10;", 10},
		{"yeet 10; 9;", 10},
		{"yeet 2 * 5; 9;", 10},
		{"9; yeet 2 * 5; 9;", 10},
		{
			`
yif 10 > 1 {
	yif 10 > 1 {
		yeet 10;
	}
	yeet 1;
}`,
			10,
		},
	})
}

func TestRangeLiterals(t *testing.T) {
	tests := []struct {
		input string
		start int64
		end   int64
	}{
		{"0..5", 0, 5},
		{"5..0", 5, 0},
		{"(1+2)..(5*2)", 3, 10},
		{"3 + 2 * 2 .. 5 - 2 * 2", 7, 1},
		{"-5..-2", -5, -2},
		{"a := 1; b := 8; a..b", 1, 8},
		{"a := 1; b := 8; (a+3)..(b-1)", 4, 7},
		{"a := 1; b := 8; a+3 .. b-1", 4, 7},
		{"r := 1+3 .. 9-1; r", 4, 8},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		rangeObj, ok := evaluated.(*object.Range)
		if !ok {
			t.Errorf("obj not *object.Range. got=%q, type of %T", evaluated, evaluated)
			continue
		}
		if rangeObj.Start != tt.start {
			t.Errorf("start wrong: want %d, got %d", tt.start, rangeObj.Start)
		}
		if rangeObj.End != tt.end {
			t.Errorf("end wrong: want %d, got %d", tt.end, rangeObj.End)
		}
	}
}

func TestAssignExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
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
		{"a := b := c := 8; a + b + c", 24},

		{"a = 8", errmsg{"identifier not found: a"}},
	})
}

func TestLambdaApplication(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`nope := \() { 69 }; nope();`, 69},
		{`identity := \(x) { x; }; identity(5);`, 5},
		{`identity := \(x) { yeet x; }; identity(5);`, 5},
		{`double := \(x) { x * 2; }; double(5);`, 10},
		{`add := \(x, y) { x + y; }; add(5, 5);`, 10},
		{`add := \(x, y) { x + y; }; add(5 + 5, add(5, 5));`, 20},
		{`add := \(x, y, z) { x + y + z; }; add(1, 2, 3);`, 6},

		{`nope := \ { 69 }; nope();`, 69},
		{`add := \x, y { x + y; }; add(5, 5);`, 10},
		{`add := \x, y { x + y; }; add(5 + 5, add(5, 5));`, 20},
		{`add := \x, y, z { x + y + z; }; add(1, 2, 3);`, 6},

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
	})
}

func TestClosures(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{
			`
newAdder := \x { 
    \n { x + n } 
}
addTwo := newAdder(2)
addTwo(2)`,
			4,
		},
		{
			`
newGenerator := \ {
    n := 0
    \ { n = n + 2 }
}
gen := newGenerator()
gen() + gen() + gen()`,
			12, // 2 + 4 + 6
		},
	})
}

func TestBuiltinFunctions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, errmsg{"argument to `len` not supported, got INTEGER"}},
		{`len("one", "two")`, errmsg{"wrong number of arguments. got=2, want=1"}},
		{`yassert(1 == 1)`, nil},
		{`yassert(1 == 2)`, errmsg{"yassert failed"}},
		{`yassert(false)`, errmsg{"yassert failed"}},
		{`a := 5; b := 6; yassert(a == b)`, errmsg{"yassert failed"}},
		{`yassert(true)`, nil},
		{`yassert(1 == 2, "one isn't two")`, errmsg{"yassert failed: one isn't two"}},
	})
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
			"yif (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
yif (10 > 1) {
	yif (10 > 1) {
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
		evaluated := testEval(t, tt.input)
		if err := testErrorObject(evaluated, tt.expectedMessage); err != nil {
			t.Error(err)
		}
	}
}

//
// HELPERS
//

func runEvalTests(t *testing.T, tests []evalTestCase) {
	t.Helper()

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch expected := tt.expected.(type) {
		case int:
			if err := testIntegerObject(evaluated, int64(expected)); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case []int64:
			if err := testIntegerArray(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case bool:
			if err := testBooleanObject(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case string:
			if err := testStringObject(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case errmsg:
			if err := testErrorObject(evaluated, expected.msg); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case nil:
			if err := testNullObject(evaluated); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		default:
			t.Errorf("unexpected type, got=%T", expected)
		}
	}
}

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Error(err)
		}
		t.Fatalf("parsing errors, bailing")
	}

	env := object.NewEnvironment()

	return Eval(program, env)
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

func testIntegerArray(obj object.Object, expected []int64) error {
	result, ok := obj.(*object.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got=%T (%+v)", obj, obj)
	}
	for i := range expected {
		if err := testIntegerObject(result.Elements[i], expected[i]); err != nil {
			return err
		}
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
		return fmt.Errorf("object has wrong value. want=%t, got=%t", expected, result.Value)
	}
	return nil
}

func testNullObject(obj object.Object) error {
	if obj != NULL {
		return fmt.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
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
	if testing.Short() {
		t.Skip("skipping testing files in short mode")
	}

	testFiles, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("couldn't read example files dir: %s", err)
	}

	for _, f := range testFiles {
		t.Run(f.Name(), func(t *testing.T) {
			filename := filepath.Join(examplesDir, f.Name())
			// filename := filepath.Join(examplesDir, "builtins.yeet")
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("couldn't read test file: %s", err)
			}

			result := testEval(t, string(src))
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
