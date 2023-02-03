package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"ylang/ast"
	"ylang/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 69;", "x", 69},
		{"let y = true;", "y", true},
		{"let z = y;", "z", "y"},
	}

	for _, tt := range tests {
		program := parse(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		if err := testLetStatement(program.Statements[0], tt.expectedIdentifier); err != nil {
			t.Error(err)
		}
		if err := testLiteralExpression(program.Statements[0].(*ast.LetStatement).Value, tt.expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestYeetStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
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

	stmt := parseSingleStmt(t, input)
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expr not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != expected {
		t.Errorf("ident.Value not %s. got=%s", expected, ident.Value)
	}
	if ident.TokenLiteral() != expected {
		t.Errorf("ident.TokenLiteral not %s. got=%s", expected, ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	expected := int64(5)
	stmt := parseSingleStmt(t, input)
	if err := testIntegerLiteral(stmt.Expression, expected); err != nil {
		t.Error(err)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"Yo, world!";`
	expected := "Yo, world!"

	stmt := parseSingleStmt(t, input)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != expected {
		t.Errorf("literal.Value not %q. got=%q", expected, literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	// TODO add more cases
	input := "[1, 2 * 2, 3 + 3]"

	stmt := parseSingleStmt(t, input)
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
		stmt := parseSingleStmt(t, tt.input)
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

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	stmt := parseSingleStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
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
		expectedValue := expected[literal.String()]
		if err := testIntegerLiteral(value, expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestParsingHashLiterals(t *testing.T) {
	// TODO more test cases
	input := "{}"

	stmt := parseSingleStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	stmt := parseSingleStmt(t, input)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
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
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
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
	stmt := parseSingleStmt(t, input)
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
		stmt := parseSingleStmt(t, tt.input)
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
		stmt := parseSingleStmt(t, tt.input)
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
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		program := parse(t, tt.input)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	// TODO add more test cases
	input := `if (x < y) { x }`
	stmt := parseSingleStmt(t, input)

	expr, ok := stmt.Expression.(*ast.IfExpression)
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

func TestYoyoExpression(t *testing.T) {
	input := "yoyo x < y { i }"
	stmt := parseSingleStmt(t, input)

	expr, ok := stmt.Expression.(*ast.YoyoExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.YoyoExpression. got=%T", stmt.Expression)
	}

	if err := testInfixExpression(expr.Condition, "x", "<", "y"); err != nil {
		t.Error(err)
	}
	if len(expr.Body.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(expr.Body.Statements))
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	stmt := parseSingleStmt(t, input)

	expr, ok := stmt.Expression.(*ast.IfExpression)
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
			"if x < y { x }",
			"(x < y)",
			"{ x }",
			"",
		},
		{
			"if (x < y) { x }",
			"(x < y)",
			"{ x }",
			"",
		},
		{
			"if x < y { x } else { y }",
			"(x < y)",
			"{ x }",
			"{ y }",
		},
		{
			"if null { x } else { y }",
			"null",
			"{ x }",
			"{ y }",
		},
		{
			"if (x < y) { if (x > y) { x } }",
			"(x < y)",
			"{ if (x > y) { x } }",
			"",
		},
	}

	for _, tt := range tests {
		stmt := parseSingleStmt(t, tt.input)

		ifExpr, ok := stmt.Expression.(*ast.IfExpression)
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

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fun(x, y) { x + y; }`
	stmt := parseSingleStmt(t, input)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(function.Parameters[0], "x")
	testLiteralExpression(function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}
	testInfixExpression(bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fun() {};", []string{}},
		{"fun(x) {};", []string{"x"}},
		{"fun(x, y, z) {};", []string{"x", "y", "z"}},
		{"fun(x, y, z,) {};", []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		stmt := parseSingleStmt(t, tt.input)
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

func TestLambdaLiteralParsing(t *testing.T) {
	input := `\(x, y) { x + y; }`
	stmt := parseSingleStmt(t, input)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(function.Parameters[0], "x")
	testLiteralExpression(function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}
	testInfixExpression(bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := "myFunction(1, 2 * 3, 4 + 5);"
	stmt := parseSingleStmt(t, input)

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
		stmt := parseSingleStmt(t, tt.input)
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

func parseSingleStmt(t *testing.T, input string) *ast.ExpressionStatement {
	program := parse(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
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
	case int64:
		return testIntegerLiteral(expr, v)
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

func testLetStatement(stmt ast.Statement, name string) error {
	if stmt.TokenLiteral() != "let" {
		return fmt.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		return fmt.Errorf("s not *ast.LetStatement. got=%T", stmt)
	}
	if letStmt.Name.Value != name {
		return fmt.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
	}
	if letStmt.Name.TokenLiteral() != name {
		return fmt.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
	}

	return nil
}

var testFiles = []string{
	"fun.yeet",
	"vars.yeet",
}

func TestParseFiles(t *testing.T) {
	for _, f := range testFiles {
		t.Run(f, func(t *testing.T) {
			filename := filepath.Join("..", "examples", f)
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("couldn't read test file: %s", err)
			}

			_ = parse(t, string(src))
		})
	}
}

func BenchmarkParse(b *testing.B) {
	for _, f := range testFiles {
		b.Run(f, func(b *testing.B) {
			b.StopTimer()
			filename := filepath.Join("..", "examples", f)
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
