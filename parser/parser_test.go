package parser_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"yy/ast"
	"yy/lexer"
	"yy/parser"
)

func TestDeclareExpression(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"x := 69;", "x", 69},
		{"y := true;", "y", true},
		{"z := y;", "z", "y"},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)
		if err := testLiteralExpression(expr.(*ast.DeclareExpression).Value, tt.expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestAssignExpression(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"x = 69;", "x", 69},
		{"y = true;", "y", true},
		{"z = y;", "z", "y"},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)
		if err := testLiteralExpression(expr.(*ast.AssignExpression).Value, tt.expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestYeetExpressions(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"yeet 8;", 8},
		{"yeet true;", true},
		{"yeet BoatyMcBoatface;", "BoatyMcBoatface"},
	}

	for _, tt := range tests {
		program := parse(t, tt.input)

		if len(program.Expressions) != 1 {
			t.Errorf("program.Expressions does not contain 1 statements. got=%d", len(program.Expressions))
		}

		expr := program.Expressions[0]

		yeetExpr, ok := expr.(*ast.YeetExpression)
		if !ok {
			t.Errorf("expr not *ast.YeetExpression. got=%T", expr)
		}
		if yeetExpr.TokenLiteral() != "yeet" {
			t.Errorf("Yeetexpr.TokenLiteral wrong, got %q", yeetExpr.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	expected := "foobar"

	expr := parseSingleExpr(t, input)
	if err := testIdentifier(expr, expected); err != nil {
		t.Error(err)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	expected := int64(5)

	expr := parseSingleExpr(t, input)
	if err := testIntegerLiteral(expr, expected); err != nil {
		t.Error(err)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{`"Yo, world!";`, "Yo, world!"},
		{`"42 is the answer"`, "42 is the answer"},
	}

	for _, tt := range testCases {
		expr := parseSingleExpr(t, tt.input)
		literal, ok := expr.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("exp not *ast.StringLiteral. got=%T", expr)
		}
		if literal.Value != tt.expected {
			t.Errorf("literal.Value not %q. got=%q", tt.expected, literal.Value)
		}
	}
}

func TestTemplStringLiteralExpression(t *testing.T) {
	// TODO: either add more test cases or collapse
	testCases := []struct {
		input string
		templ string
	}{
		{
			`"i'm {name} and i'm { 30 + 5 } years old"`,
			"i'm %s and i'm %s years old",
		},
	}

	for _, tt := range testCases {
		expr := parseSingleExpr(t, tt.input)

		literal, ok := expr.(*ast.TemplateStringLiteral)
		if !ok {
			t.Fatalf("exp not *ast.TemplateStringLiteral. got=%T", expr)
		}

		if literal.Template != tt.templ {
			t.Errorf("literal.Template not %q. got=%q", tt.templ, literal.Template)
		}

		if err := testIdentifier(literal.Values[0], "name"); err != nil {
			t.Error(err)
		}

		if err := testInfixExpression(literal.Values[1], 30, "+", 5); err != nil {
			t.Error(err)
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	expr := parseSingleExpr(t, input)
	array, ok := expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", expr)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	if err := testIntegerLiteral(array.Elements[0], 1); err != nil {
		t.Error(err)
	}
	if err := testInfixExpression(array.Elements[1], 2, "*", 2); err != nil {
		t.Error(err)
	}
	if err := testInfixExpression(array.Elements[2], 3, "+", 3); err != nil {
		t.Error(err)
	}
}

func TestParsingArrayLiterals3(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"[]", []int64{}},
		{"[1]", []int64{1}},
		{"[1,2]", []int64{1, 2}},
		{"[1,2,]", []int64{1, 2}},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)
		array, ok := expr.(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf("exp not ast.ArrayLiteral. got=%T", expr)
		}
		if len(array.Elements) != len(tt.expected) {
			t.Fatalf("len(array.Elements) not %d. got=%d", len(tt.expected), len(array.Elements))
		}
		for i := range tt.expected {
			if err := testIntegerLiteral(array.Elements[i], tt.expected[i]); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestRangeLiteral(t *testing.T) {
	tests := []struct {
		input string
		start string
		end   string
	}{
		{"a..b", "a", "b"},
		{"(a+1)..(b+5)", "(a + 1)", "(b + 5)"},
		{"1..5", "1", "5"},
		{"(1+1)..5", "(1 + 1)", "5"},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)

		rangeLit, ok := expr.(*ast.RangeLiteral)
		if !ok {
			t.Fatalf("expr.Expression is not ast.RangeLiteral. got=%T", expr)
		}

		if rangeLit.Start.String() != tt.start {
			t.Errorf("Iterable is not %q. got=%q", tt.start, rangeLit.Start.String())
		}
		if rangeLit.End.String() != tt.end {
			t.Errorf("End is not %q. got=%q", tt.end, rangeLit.End.String())
		}
	}
}

func TestParsingHashLiterals(t *testing.T) {
	tests := []struct {
		input string
		pairs int
	}{
		{"%{}", 0},
		{`%{"a": 1 }`, 1},
		{`%{"a": 1, "b": 2 }`, 2},
		{`%{"a": 1, "b": 2, 3: 3 }`, 3},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)
		hash, ok := expr.(*ast.HashmapLiteral)
		if !ok {
			t.Fatalf("exp is not ast.HashLiteral. got=%T", expr)
		}
		if len(hash.Pairs) != tt.pairs {
			t.Errorf("hash.Pairs has wrong length. want=%d, got=%d", tt.pairs, len(hash.Pairs))
		}
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `%{"one": 1, "two": 2, "three": 3}`
	expected := map[string]int64{"one": 1, "two": 2, "three": 3}

	expr := parseSingleExpr(t, input)
	hash, ok := expr.(*ast.HashmapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", expr)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.Value]
		if err := testIntegerLiteral(value, expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `%{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	expr := parseSingleExpr(t, input)
	hash, ok := expr.(*ast.HashmapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", expr)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string](func(ast.Expression) error){
		"one": func(e ast.Expression) error {
			return testInfixExpression(e, 0, "+", 1)
		},
		"two": func(e ast.Expression) error {
			return testInfixExpression(e, 10, "-", 8)
		},
		"three": func(e ast.Expression) error {
			return testInfixExpression(e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.Value]
		if !ok {
			t.Errorf("No test function for key %q found", literal.Value)
			continue
		}
		if err := testFunc(value); err != nil {
			t.Error(err)
		}
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	// TODO add more test cases
	input := "myArray[1 + 1]"
	expr := parseSingleExpr(t, input)
	indexExp, ok := expr.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", expr)
	}
	if err := testIdentifier(indexExp.Left, "myArray"); err != nil {
		t.Error(err)
	}
	if err := testInfixExpression(indexExp.Index, 1, "+", 1); err != nil {
		t.Error(err)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		expected any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		expr := parseSingleExpr(t, tt.input)
		if err := testPrefixExpression(expr, tt.expected, tt.operator); err != nil {
			t.Error(err)
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		expr := parseSingleExpr(t, tt.input)
		if err := testInfixExpression(expr, tt.leftValue, tt.operator, tt.rightValue); err != nil {
			t.Error(err)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b);",
		},
		{
			"!-a",
			"(!(-a));",
		},
		{
			"a + b + c",
			"((a + b) + c);",
		},
		{
			"a + b - c",
			"((a + b) - c);",
		},
		{
			"a * b * c",
			"((a * b) * c);",
		},
		{
			"a * b / c",
			"((a * b) / c);",
		},
		{
			"a / b * c",
			"((a / b) * c);",
		},
		{
			"a + b / c",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4);((-5) * 5);",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4));",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		{
			"true",
			"true;",
		},
		{
			"false",
			"false;",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false);",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true);",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4);",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2);",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5));",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5));",
		},
		{
			"!(true == true)",
			"(!(true == true));",
		},
		{
			"3 && 4 * 5",
			"(3 && (4 * 5));",
		},
		{
			"3 && 4 || 5",
			"((3 && 4) || 5);",
		},
		{
			"3 || 4 && 5",
			"(3 || (4 && 5));",
		},
		{
			"3 && 4 * 5 && 6 + 7",
			"((3 && (4 * 5)) && (6 + 7));",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d);",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g));",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d);",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])));",
		},
		{
			`add := \a b { a + b }; add(5, 10) == 5 + 2 * 10`,
			`(add := \(a, b) { (a + b) });(add(5, 10) == (5 + (2 * 10)));`,
		},
		{
			`a := 5 + 2 * 10 / 8 - 15;`,
			`(a := ((5 + ((2 * 10) / 8)) - 15));`,
		},
		{
			"a := 3 + 4 * 5 == 3 * 1 + 4 * 5",
			"(a := ((3 + (4 * 5)) == ((3 * 1) + (4 * 5))));",
		},
		{
			"a = b = c = 8",
			"(a = (b = (c = 8)));",
		},
		{
			"a := b := c := 8",
			"(a := (b := (c := 8)));",
		},
		{
			"f := 6 + 2 * 3 g := 3 * 3 + 1 h := f + g",
			"(f := (6 + (2 * 3)));(g := ((3 * 3) + 1));(h := (f + g));",
		},
		{
			"(1+2) .. (8*2)",
			"((1 + 2)..(8 * 2));",
		},
		{
			"1 + 2 .. 8 * 2",
			"((1 + 2)..(8 * 2));",
		},
		{
			"r := 1 + 2 .. 8 * 2",
			"(r := ((1 + 2)..(8 * 2)));",
		},
	}

	for _, tt := range tests {
		program := parse(t, tt.input)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("\nwant %q \ngot  %q", tt.expected, actual)
		}
	}
}

func TestIfExpression2(t *testing.T) {
	input := `yif (x < y) { x } yels { y }`
	expr := parseSingleExpr(t, input)

	yifExpr, ok := expr.(*ast.YifExpression)
	if !ok {
		t.Fatalf("expr.Expression is not ast.IfExpression. got=%T", expr)
	}

	if err := testInfixExpression(yifExpr.Condition, "x", "<", "y"); err != nil {
		t.Error(err)
	}

	if len(yifExpr.Consequence.Expressions) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(yifExpr.Consequence.Expressions))
	}
	consequence := yifExpr.Consequence.Expressions[0]
	if err := testIdentifier(consequence, "x"); err != nil {
		t.Error(err)
	}

	if len(yifExpr.Alternative.Expressions) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(yifExpr.Alternative.Expressions))
	}
	alternative := yifExpr.Alternative.Expressions[0]
	if !ok {
		t.Fatalf("Expressions[0] is not ast.ExpressionStatement. got=%T", yifExpr.Alternative.Expressions[0])
	}
	if err := testIdentifier(alternative, "y"); err != nil {
		t.Error(err)
	}
}

func TestYoloExpression(t *testing.T) {
	tests := []struct {
		input string
		body  string
	}{
		{
			"yolo { i = i + 1 }",
			"{ (i = (i + 1)) }",
		},
		{
			"yolo { yowl() }",
			"{ yowl() }",
		},
		{
			"yolo { yowl(); 5; yap(); 2 + 2; }",
			"{ yowl(); 5; yap(); (2 + 2) }",
		},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)

		yoloExpr, ok := expr.(*ast.YoloExpression)
		if !ok {
			t.Fatalf("expr.Expression is not ast.YoloExpression. got=%T", expr)
		}

		if yoloExpr.Body.String() != tt.body {
			t.Errorf("condition.String() is not %s. got=%s", tt.body, yoloExpr.Body)
		}
	}
}

func TestYoyoExpression(t *testing.T) {
	tests := []struct {
		input     string
		condition string
		body      string
	}{
		{
			"yoyo i < 5 { i = i + 1 }",
			"(i < 5)",
			"{ (i = (i + 1)) }",
		},
		{
			"yoyo true { yowl() }",
			"true",
			"{ yowl() }",
		},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)

		yoyoExpr, ok := expr.(*ast.YoyoExpression)
		if !ok {
			t.Fatalf("expr is not ast.YoyoExpression. got=%T", expr)
		}

		if yoyoExpr.Condition.String() != tt.condition {
			t.Errorf("condition is not %s. got=%s", tt.condition, yoyoExpr.Condition)
		}

		if yoyoExpr.Body.String() != tt.body {
			t.Errorf("condition.String() is not %s. got=%s", tt.body, yoyoExpr.Body)
		}
	}
}

func TestYallExpression(t *testing.T) {
	tests := []struct {
		input    string
		name     string
		iterable string
		body     string
	}{
		{
			"yall array { yt }",
			"yt",
			"array",
			"{ yt }",
		},
		{
			"yall i: array { i }",
			"i",
			"array",
			"{ i }",
		},
		{
			"yall yt: array { yt }",
			"yt",
			"array",
			"{ yt }",
		},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)

		yallExpr, ok := expr.(*ast.YallExpression)
		if !ok {
			t.Fatalf("expr is not ast.YallExpression. got=%T", expr)
		}

		if yallExpr.KeyName != tt.name {
			t.Errorf("KeyName is not %s. got=%s", tt.name, yallExpr.KeyName)
		}

		if yallExpr.Iterable.String() != tt.iterable {
			t.Errorf("Iterable is not %s. got=%s", tt.iterable, yallExpr.Iterable)
		}

		if yallExpr.Body.String() != tt.body {
			t.Errorf("Body is not %s. got=%s", tt.body, yallExpr.Body)
		}
	}
}

func TestYifYelsExpression(t *testing.T) {
	input := `yif (x < y) { x } yels { y }`
	expr := parseSingleExpr(t, input)

	yifExpr, ok := expr.(*ast.YifExpression)
	if !ok {
		t.Fatalf("expr is not ast.IfExpression. got=%T", yifExpr)
	}

	if err := testInfixExpression(yifExpr.Condition, "x", "<", "y"); err != nil {
		t.Error(err)
	}
	if len(yifExpr.Consequence.Expressions) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(yifExpr.Consequence.Expressions))
	}

	consequence := yifExpr.Consequence.Expressions[0]
	if err := testIdentifier(consequence, "x"); err != nil {
		t.Error(err)
	}

	alternative := yifExpr.Alternative.Expressions[0]
	if err := testIdentifier(alternative, "y"); err != nil {
		t.Errorf("alternative wrong: %s", err)
	}
}

func TestParseIfExpressions(t *testing.T) {
	tests := []struct {
		input               string
		expectedCondition   string
		expectedConsequence string
		expectedAlternative string
	}{
		{
			"yif x < y { x }",
			"(x < y)",
			"{ x }",
			"",
		},
		{
			"yif (x < y) { x }",
			"(x < y)",
			"{ x }",
			"",
		},
		{
			"yif x < y { x } yels { y }",
			"(x < y)",
			"{ x }",
			"{ y }",
		},
		{
			"yif null { x } yels { y }",
			"null",
			"{ x }",
			"{ y }",
		},
		{
			"yif (x < y) { yif (x > y) { x } }",
			"(x < y)",
			"{ yif (x > y) { x } }",
			"",
		},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)

		yifExpr, ok := expr.(*ast.YifExpression)
		if !ok {
			t.Fatalf("expr.Expression is not ast.IfExpression. got=%T", expr)
		}

		if yifExpr.Condition.String() != tt.expectedCondition {
			t.Errorf("condition.String() is not %q. got=%q", tt.expectedCondition, yifExpr.Condition.String())
		}

		if yifExpr.Consequence.String() != tt.expectedConsequence {
			t.Errorf("consequence.String() is not %q. got=%q", tt.expectedConsequence, yifExpr.Consequence.String())
		}

		if yifExpr.Alternative == nil {
			if tt.expectedAlternative != "" {
				t.Errorf("exp.Alternative is nil. want=%q", tt.expectedAlternative)
			}
		} else {
			if yifExpr.Alternative.String() != tt.expectedAlternative {
				t.Errorf("exp.Alternative.String() is not %q. got=%q", tt.expectedAlternative, yifExpr.Alternative.String())
			}
		}
	}
}

func TestLambdaLiteralParsing(t *testing.T) {
	input := `\x, y { x + y }`
	expr := parseSingleExpr(t, input)

	function, ok := expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not ast.FunctionLiteral. got=%T", expr)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	if err := testLiteralExpression(function.Parameters[0], "x"); err != nil {
		t.Error(err)
	}
	if err := testLiteralExpression(function.Parameters[1], "y"); err != nil {
		t.Error(err)
	}

	if len(function.Body.Expressions) != 1 {
		t.Fatalf("function.Body.Expressions has not 1 statements. got=%d\n", len(function.Body.Expressions))
	}
	if err := testInfixExpression(function.Body.Expressions[0], "x", "+", "y"); err != nil {
		t.Error(err)
	}
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `@\x, y { x + y; }`
	expr := parseSingleExpr(t, input)

	macro, ok := expr.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("expr.Expression is not ast.MacroLiteral. got=%T", expr)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("macro literal parameters wrong. want 2, got=%d\n", len(macro.Parameters))
	}

	if err := testLiteralExpression(macro.Parameters[0], "x"); err != nil {
		t.Error(err)
	}
	if err := testLiteralExpression(macro.Parameters[1], "y"); err != nil {
		t.Error(err)
	}

	if len(macro.Body.Expressions) != 1 {
		t.Fatalf("macro.Body.Expressions has not 1 statements. got=%d\n", len(macro.Body.Expressions))
	}

	bodyexpr := macro.Body.Expressions[0]

	if err := testInfixExpression(bodyexpr, "x", "+", "y"); err != nil {
		t.Error(err)
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{`\ {};`, []string{}},
		{`\x {};`, []string{"x"}},
		{`\x, y, z {};`, []string{"x", "y", "z"}},
		{`\x y z {};`, []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)
		fn := expr.(*ast.FunctionLiteral)
		if len(fn.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(fn.Parameters))
		}
		for i, ident := range tt.expectedParams {
			if err := testLiteralExpression(fn.Parameters[i], ident); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "myFunction(1, 2 * 3, 4 + 5);"
	expr := parseSingleExpr(t, input)

	callExpr, ok := expr.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expr.Expression is not ast.CallExpression. got=%T", callExpr)
	}
	if err := testIdentifier(callExpr.Function, "myFunction"); err != nil {
		t.Fatal(err)
	}
	if len(callExpr.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(callExpr.Arguments))
	}
	if err := testLiteralExpression(callExpr.Arguments[0], 1); err != nil {
		t.Fatal(err)
	}
	if err := testInfixExpression(callExpr.Arguments[1], 2, "*", 3); err != nil {
		t.Fatal(err)
	}
	if err := testInfixExpression(callExpr.Arguments[2], 4, "+", 5); err != nil {
		t.Fatal(err)
	}
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input        string
		expectedArgs []string
	}{
		{"myFunction();", []string{}},
		{"myFunction(x);", []string{"x"}},
		{"myFunction(x, y, z);", []string{"x", "y", "z"}},
		{"myFunction(x, y, z,);", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		expr := parseSingleExpr(t, tt.input)
		callExpr := expr.(*ast.CallExpression)
		if len(callExpr.Arguments) != len(tt.expectedArgs) {
			t.Errorf("length args wrong. want %d, got=%d\n", len(tt.expectedArgs), len(callExpr.Arguments))
		}
		for i, ident := range tt.expectedArgs {
			if err := testLiteralExpression(callExpr.Arguments[i], ident); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestParsingUnterminatedStrings(t *testing.T) {
	tests := []string{
		`x := "5`,
		`x := "5""`,
		`x := "5"";`,
		`x := "5; z := 2;`,
		`x := "5; 
z := 2;`,
		`x := "5"; z := "2;`,
	}

	for _, tt := range tests {
		parser := parser.New(lexer.New(tt))
		_ = parser.ParseProgram()
		errors := parser.Errors()

		if len(errors) != 1 {
			t.Errorf("expected 1 parsing error, got %d", len(errors))
		}

		expectedErrMsg := "unterminated string"
		if errors[0].Msg != expectedErrMsg {
			t.Errorf("Wrong error msg, want `%s`, got `%s`", expectedErrMsg, errors[0].Msg)
		}
	}
}

//
// HELPERS
//

func parse(t *testing.T, input string) *ast.Program {
	t.Helper()

	l := lexer.New(input)
	parser := parser.New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	return program
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()

	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d error(s)", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func parseSingleExpr(t *testing.T, input string) ast.Expression {
	t.Helper()

	program := parse(t, input)
	if len(program.Expressions) != 1 {
		t.Fatalf("program.Expressions does not contain a single statement. got=%d\n", len(program.Expressions))
	}
	return program.Expressions[0]
}

func testLiteralExpression(expr ast.Expression, expected any) error {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(expr, int64(v))
	case float64:
		return testNumberLiteral(expr, v)
	case string:
		return testIdentifier(expr, v)
	case bool:
		return testBooleanLiteral(expr, v)
	default:
		return fmt.Errorf("type of expr not handled. got=%T", expr)
	}
}

func testIntegerLiteral(expr ast.Expression, value int64) error {
	integ, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		return fmt.Errorf("expr not *ast.IntegerLiteral. got=%T", expr)
	}
	if integ.Value != value {
		return fmt.Errorf("integ.Value not %d. got=%d", value, integ.Value)
	}
	if integ.TokenLiteral() != strconv.FormatInt(value, 10) {
		return fmt.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
	}
	return nil
}

func testNumberLiteral(expr ast.Expression, value float64) error {
	number, ok := expr.(*ast.NumberLiteral)
	if !ok {
		return fmt.Errorf("expr not *ast.NumberLiteral. got=%T", expr)
	}
	if number.Value != value {
		return fmt.Errorf("integ.Value not %f. got=%f", value, number.Value)
	}
	if number.TokenLiteral() != fmt.Sprintf("%g", value) {
		return fmt.Errorf("number.TokenLiteral not %g. got=%s", value, number.TokenLiteral())
	}
	return nil
}

func testIdentifier(expr ast.Expression, value string) error {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("expr not *ast.Identifier. got=%T", expr)
	}
	if ident.Value != value {
		return fmt.Errorf("ident.Value not %s. got=%s", value, ident.Value)
	}
	if ident.TokenLiteral() != value {
		return fmt.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
	}
	return nil
}

func testBooleanLiteral(expr ast.Expression, value bool) error {
	b, ok := expr.(*ast.BooleanLiteral)
	if !ok {
		return fmt.Errorf("expr not *ast.Boolean. got=%T", expr)
	}
	if b.Value != value {
		return fmt.Errorf("bool value not %t. got=%t", value, b.Value)
	}
	if b.TokenLiteral() != strconv.FormatBool(value) {
		return fmt.Errorf("bool TokenLiteral not %t. got=%s", value, b.TokenLiteral())
	}
	return nil
}

func testInfixExpression(expr ast.Expression, left any, operator string, right any) error {
	infixExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		return fmt.Errorf("expr not ast.InfixExpression. got=%T(%s)", expr, expr)
	}
	if err := testLiteralExpression(infixExpr.Left, left); err != nil {
		return err
	}
	if infixExpr.Operator != operator {
		return fmt.Errorf("infix operator not '%s'. got=%q", operator, infixExpr.Operator)
	}

	return testLiteralExpression(infixExpr.Right, right)
}

func testPrefixExpression(expr ast.Expression, expectedRight any, operator string) error {
	prefixExpr, ok := expr.(*ast.PrefixExpression)
	if !ok {
		return fmt.Errorf("expr not ast.PrefixExpression. got=%T(%s)", expr, expr)
	}
	if prefixExpr.Operator != operator {
		return fmt.Errorf("prefix operator not '%s'. got=%q", operator, prefixExpr.Operator)
	}

	return testLiteralExpression(prefixExpr.Right, expectedRight)
}

//
// PROGRAMS FROM EXAMPLES/
//

const examplesDir = "../examples"

func TestParseFiles(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testFiles, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("couldn't read example files dir: %s", err)
	}

	for _, f := range testFiles {
		t.Run(f.Name(), func(t *testing.T) {
			filename := filepath.Join("..", "examples", f.Name())
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("couldn't read test file: %s", err)
			}

			_ = parse(t, string(src))
		})
	}
}

func BenchmarkParse(b *testing.B) {
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

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				l := lexer.New(sSrc)
				parser := parser.New(l)
				_ = parser.ParseProgram()
			}
		})
	}
}
