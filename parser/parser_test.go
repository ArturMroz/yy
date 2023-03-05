package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"yy/ast"
	"yy/lexer"
)

func TestAssignExpression(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"x := 69;", "x", 69},
		{"y := true;", "y", true},
		{"z := y;", "z", "y"},
		{"x = 69;", "x", 69},
		{"y = true;", "y", true},
		{"z = y;", "z", "y"},
	}

	for _, tt := range tests {
		stmt := parseSingleExprStmt(t, tt.input)
		if err := testLiteralExpression(stmt.Expression.(*ast.AssignExpression).Value, tt.expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestYeetStatements(t *testing.T) {
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

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]

		yeetStmt, ok := stmt.(*ast.YeetStatement)
		if !ok {
			t.Errorf("stmt not *ast.YeetStatement. got=%T", stmt)
		}
		if yeetStmt.TokenLiteral() != "yeet" {
			t.Errorf("YeetStmt.TokenLiteral wrong, got %q", yeetStmt.TokenLiteral())
		}
		if err := testLiteralExpression(yeetStmt.ReturnValue, tt.expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	expected := "foobar"

	stmt := parseSingleExprStmt(t, input)
	if err := testIdentifier(stmt.Expression, expected); err != nil {
		t.Error(err)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	expected := int64(5)

	stmt := parseSingleExprStmt(t, input)
	if err := testIntegerLiteral(stmt.Expression, expected); err != nil {
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
		stmt := parseSingleExprStmt(t, tt.input)
		literal, ok := stmt.Expression.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
		}
		if literal.Value != tt.expected {
			t.Errorf("literal.Value not %q. got=%q", tt.expected, literal.Value)
		}
	}
}

func TestTemplStringLiteralExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		templ    string
	}{
		{
			`"i'm {age} years old"`,
			"i'm {age} years old",
			"i'm %s years old",
		},
		{
			`"i have {apples} apples and {carrots} carrots"`,
			"i have {apples} apples and {carrots} carrots",
			"i have %s apples and %s carrots",
		},
		{
			`"i have {apples} apples and {carrots} carrots and {kiwi} kiwis`,
			"i have {apples} apples and {carrots} carrots and {kiwi} kiwis",
			"i have %s apples and %s carrots and %s kiwis",
		},
		{
			`"i have fruit: {apples} {carrots} {kiwi}.`,
			"i have fruit: {apples} {carrots} {kiwi}.",
			"i have fruit: %s %s %s.",
		},
	}

	for _, tt := range testCases {
		stmt := parseSingleExprStmt(t, tt.input)
		literal, ok := stmt.Expression.(*ast.TemplateStringLiteral)
		if !ok {
			t.Fatalf("exp not *ast.TemplateStringLiteral. got=%T", stmt.Expression)
		}
		if literal.Token.Literal != tt.expected {
			t.Errorf("literal.Value not %q. got=%q", tt.expected, literal.Token.Literal)
		}
		if literal.Template != tt.templ {
			t.Errorf("literal.Template not %q. got=%q", tt.templ, literal.Template)
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	stmt := parseSingleExprStmt(t, input)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
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
		stmt := parseSingleExprStmt(t, tt.input)
		array, ok := stmt.Expression.(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
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
		stmt := parseSingleExprStmt(t, tt.input)

		rangeLit, ok := stmt.Expression.(*ast.RangeLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.RangeLiteral. got=%T", stmt.Expression)
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
		stmt := parseSingleExprStmt(t, tt.input)
		hash, ok := stmt.Expression.(*ast.HashmapLiteral)
		if !ok {
			t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
		}
		if len(hash.Pairs) != tt.pairs {
			t.Errorf("hash.Pairs has wrong length. want=%d, got=%d", tt.pairs, len(hash.Pairs))
		}
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `%{"one": 1, "two": 2, "three": 3}`
	expected := map[string]int64{"one": 1, "two": 2, "three": 3}

	stmt := parseSingleExprStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashmapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
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
	stmt := parseSingleExprStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashmapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
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
	// TOOD add more test cases
	input := "myArray[1 + 1]"
	stmt := parseSingleExprStmt(t, input)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
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
		stmt := parseSingleExprStmt(t, tt.input)
		if err := testPrefixExpression(stmt.Expression, tt.expected, tt.operator); err != nil {
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
		stmt := parseSingleExprStmt(t, tt.input)
		if err := testInfixExpression(stmt.Expression, tt.leftValue, tt.operator, tt.rightValue); err != nil {
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
			"add(a * b[2] b[1] 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])));",
		},
		{
			`add := \a b { a + b }; add(5 10) == 5 + 2 * 10`,
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

func TestIfExpression(t *testing.T) {
	// TODO add more test cases
	input := `yif (x < y) { x }`
	stmt := parseSingleExprStmt(t, input)

	expr, ok := stmt.Expression.(*ast.YifExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if err := testInfixExpression(expr.Condition, "x", "<", "y"); err != nil {
		t.Error(err)
	}
	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(expr.Consequence.Statements))
	}
	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expr.Consequence.Statements[0])
	}
	if err := testIdentifier(consequence.Expression, "x"); err != nil {
		t.Error(err)
	}
	if expr.Alternative != nil {
		t.Errorf("expr.Alternative.Statements was not nil. got=%+v", expr.Alternative)
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
		stmt := parseSingleExprStmt(t, tt.input)

		yyExpr, ok := stmt.Expression.(*ast.YoloExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.YoloExpression. got=%T", stmt.Expression)
		}

		if yyExpr.Body.String() != tt.body {
			t.Errorf("condition.String() is not %s. got=%s", tt.body, yyExpr.Body)
		}
	}
}

func TestYetExpression(t *testing.T) {
	tests := []struct {
		input     string
		condition string
		body      string
	}{
		{
			"yet i < 5 { i = i + 1 }",
			"(i < 5)",
			"{ (i = (i + 1)) }",
		},
		{
			"yet true { yowl() }",
			"true",
			"{ yowl() }",
		},
	}

	for _, tt := range tests {
		stmt := parseSingleExprStmt(t, tt.input)

		yyExpr, ok := stmt.Expression.(*ast.YetExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.YetExpression. got=%T", stmt.Expression)
		}

		if yyExpr.Condition.String() != tt.condition {
			t.Errorf("condition is not %s. got=%s", tt.condition, yyExpr.Condition)
		}

		if yyExpr.Body.String() != tt.body {
			t.Errorf("condition.String() is not %s. got=%s", tt.body, yyExpr.Body)
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
		stmt := parseSingleExprStmt(t, tt.input)

		yallExpr, ok := stmt.Expression.(*ast.YallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.YallExpression. got=%T", stmt.Expression)
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
	stmt := parseSingleExprStmt(t, input)

	expr, ok := stmt.Expression.(*ast.YifExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if err := testInfixExpression(expr.Condition, "x", "<", "y"); err != nil {
		t.Error(err)
	}
	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence is not ast.ExpressionStatement. got=%T", expr.Consequence.Statements[0])
	}
	if err := testIdentifier(consequence.Expression, "x"); err != nil {
		t.Error(err)
	}

	alternative, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative is not ast.ExpressionStatement. got=%T", expr.Alternative.Statements[0])
	}
	if err := testIdentifier(alternative.Expression, "y"); err != nil {
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
		stmt := parseSingleExprStmt(t, tt.input)

		ifExpr, ok := stmt.Expression.(*ast.YifExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
		}

		if ifExpr.Condition.String() != tt.expectedCondition {
			t.Errorf("condition.String() is not %q. got=%q", tt.expectedCondition, ifExpr.Condition.String())
		}

		if ifExpr.Consequence.String() != tt.expectedConsequence {
			t.Errorf("consequence.String() is not %q. got=%q", tt.expectedConsequence, ifExpr.Consequence.String())
		}

		if ifExpr.Alternative == nil {
			if tt.expectedAlternative != "" {
				t.Errorf("exp.Alternative is nil. want=%q", tt.expectedAlternative)
			}
		} else {
			if ifExpr.Alternative.String() != tt.expectedAlternative {
				t.Errorf("exp.Alternative.String() is not %q. got=%q", tt.expectedAlternative, ifExpr.Alternative.String())
			}
		}
	}
}

func TestLambdaLiteralParsing(t *testing.T) {
	input := `\(x, y) { x + y }`
	stmt := parseSingleExprStmt(t, input)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
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

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	if err := testInfixExpression(bodyStmt.Expression, "x", "+", "y"); err != nil {
		t.Error(err)
	}
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `@\x, y { x + y; }`
	stmt := parseSingleExprStmt(t, input)

	macro, ok := stmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MacroLiteral. got=%T", stmt.Expression)
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

	if len(macro.Body.Statements) != 1 {
		t.Fatalf("macro.Body.Statements has not 1 statements. got=%d\n", len(macro.Body.Statements))
	}

	bodyStmt, ok := macro.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("macro body stmt is not ast.ExpressionStatement. got=%T", macro.Body.Statements[0])
	}

	if err := testInfixExpression(bodyStmt.Expression, "x", "+", "y"); err != nil {
		t.Error(err)
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{`\() {};`, []string{}},
		{`\ {};`, []string{}},
		{`\(x) {};`, []string{"x"}},
		{`\x {};`, []string{"x"}},
		{`\(x, y, z) {};`, []string{"x", "y", "z"}},
		{`\(x, y, z,) {};`, []string{"x", "y", "z"}},
		{`\(x y z) {};`, []string{"x", "y", "z"}},
		{`\x y z {};`, []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		stmt := parseSingleExprStmt(t, tt.input)
		fn := stmt.Expression.(*ast.FunctionLiteral)
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
	stmt := parseSingleExprStmt(t, input)

	expr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}
	if err := testIdentifier(expr.Function, "myFunction"); err != nil {
		t.Fatal(err)
	}
	if len(expr.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(expr.Arguments))
	}
	if err := testLiteralExpression(expr.Arguments[0], 1); err != nil {
		t.Fatal(err)
	}
	if err := testInfixExpression(expr.Arguments[1], 2, "*", 3); err != nil {
		t.Fatal(err)
	}
	if err := testInfixExpression(expr.Arguments[2], 4, "+", 5); err != nil {
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
		stmt := parseSingleExprStmt(t, tt.input)
		callExpr := stmt.Expression.(*ast.CallExpression)
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

//
// HELPERS
//

func parse(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
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

func parseSingleExprStmt(t *testing.T, input string) *ast.ExpressionStatement {
	program := parse(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain a single statement. got=%d\n", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	return stmt
}

func testLiteralExpression(expr ast.Expression, expected any) error {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(expr, int64(v))
	// case int64:
	// 	return testIntegerLiteral(expr, v)
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
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		return fmt.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
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
	b, ok := expr.(*ast.Boolean)
	if !ok {
		return fmt.Errorf("expr not *ast.Boolean. got=%T", expr)
	}
	if b.Value != value {
		return fmt.Errorf("bool value not %t. got=%t", value, b.Value)
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
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

// func testLetStatement(stmt ast.Statement, name string) error {
// 	if stmt.TokenLiteral() != "let" {
// 		return fmt.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
// 	}

// 	letStmt, ok := stmt.(*ast.LetStatement)
// 	if !ok {
// 		return fmt.Errorf("s not *ast.LetStatement. got=%T", stmt)
// 	}
// 	if letStmt.Name.Value != name {
// 		return fmt.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
// 	}
// 	if letStmt.Name.TokenLiteral() != name {
// 		return fmt.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
// 	}

// 	return nil
// }

//
// PROGRAMS FROM EXAMPLES/
//

const examplesDir = "../examples"

func TestParseFiles(t *testing.T) {
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
				parser := New(l)
				_ = parser.ParseProgram()
			}
		})
	}
}
