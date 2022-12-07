package parser

import (
	"fmt"
	"testing"

	"ylang/ast"
	"ylang/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let z = 69;
`
	expected := []string{"x", "y", "z"}

	program := parse(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for i, exp := range expected {
		if err := testLetStatement(program.Statements[i], exp); err != nil {
			t.Error(err)
		}
	}
}

// func TestLetStatements2(t *testing.T) {
// 	tests := []struct {
// 		input              string
// 		expectedIdentifier string
// 		expectedValue      any
// 	}{
// 		{"let x = 5;", "x", 5},
// 		{"let y = true;", "y", true},
// 		{"let foobar = y;", "foobar", "y"},
// 	}

// 	for _, tt := range tests {
// 		l := lexer.New(tt.input)
// 		p := New(l)
// 		program := p.ParseProgram()
// 		checkParserErrors(t, p)
// 		if len(program.Statements) != 1 {
// 			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
// 		}
// 		if err := testLetStatement(program.Statements[0], tt.expectedIdentifier); err != nil {
// 			t.Error(err)
// 		}
// 		if err := testLiteralExpression(program.Statements[0].(*ast.LetStatement).Value, tt.expectedValue); err != nil {
// 			t.Error(err)
// 		}
// 	}
// }

func TestYeetStatements(t *testing.T) {
	input := `
yeet 6;
yeet 9;
yeet 69;
`
	program := parse(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		yeetStmt, ok := stmt.(*ast.YeetStatement)

		if !ok {
			t.Errorf("stmt not *ast.YeetStatement. got=%T", stmt)
		}
		if yeetStmt.TokenLiteral() != "yeet" {
			t.Errorf("YeetStmt.TokenLiteral wrong, got %q", yeetStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	expected := "foobar"

	stmt := singleStmtProgram(t, input)

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
	stmt := singleStmtProgram(t, input)
	if err := testIntegerLiteral(stmt.Expression, expected); err != nil {
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
		stmt := singleStmtProgram(t, tt.input)
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
		stmt := singleStmtProgram(t, tt.input)
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
	input := `if (x < y) { x }`
	stmt := singleStmtProgram(t, input)

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

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	stmt := singleStmtProgram(t, input)

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

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fun(x, y) { x + y; }`
	stmt := singleStmtProgram(t, input)

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
		stmt := singleStmtProgram(t, tt.input)
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

// helpers

func parse(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
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

func singleStmtProgram(t *testing.T, input string) *ast.ExpressionStatement {
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

// type YeetValue interface {
// 	int | int64 | string | bool
// }

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
