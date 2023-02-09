package parser

import (
	"fmt"
	"strconv"

	"ylang/ast"
	"ylang/lexer"
	"ylang/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	EQUALS      // == !=
	LESSGREATER // > <
	RANGE       // x..y
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -x !x
	ASSIGNMENT  // = :=
	CALL        // my_function(x)
	INDEX       // my_array[idx]
)

var precedences = map[token.TokenType]Precedence{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.RANGE:    RANGE,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.ASSIGN:   ASSIGNMENT,
	token.WALRUS:   ASSIGNMENT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func getPrecedence(tok token.Token) Precedence {
	if p, ok := precedences[tok.Type]; ok {
		return p
	}
	return LOWEST
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT:     p.parseIdentifier,
		token.INT:       p.parseIntegerLiteral,
		token.STRING:    p.parseStringLiteral,
		token.MINUS:     p.parsePrefixExpression,
		token.BANG:      p.parsePrefixExpression,
		token.TRUE:      p.parseBoolean,
		token.FALSE:     p.parseBoolean,
		token.NULL:      p.parseNull,
		token.LPAREN:    p.parseGroupedExpression,
		token.LBRACKET:  p.parseArrayLiteral,
		token.LBRACE:    p.parseHashLiteral,
		token.IF:        p.parseIfExpression,
		token.YOYO:      p.parseYoyoExpression,
		token.YALL:      p.parseYallExpression,
		token.FUNCTION:  p.parseFunctionLiteral,
		token.BACKSLASH: p.parseLambdaLiteral,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NOT_EQ:   p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.RANGE:    p.parseRangeLiteral,
		token.WALRUS:   p.parseAssignExpression,
		token.ASSIGN:   p.parseAssignExpression,
		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
	}

	// read two tokens, so curToken and peekToken are both set
	p.advance()
	p.advance()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.advance()
	}

	return program
}

// STATEMENTS

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.YEET:
		return p.parseYeetStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.advance()
	for !p.curIs(token.RBRACE) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.advance()
	}
	return block
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.curToken}
	if !p.consume(token.IDENT, "missing identfier after 'let' keyword") {
		return nil // TODO return error
	}

	letStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.consume(token.ASSIGN, "missing '=' after identifier in assignment") {
		return nil
	}

	p.advance()

	letStmt.Value = p.parseExpression(LOWEST)

	if p.peekIs(token.SEMICOLON) { // optional semicolon
		p.advance()
	}

	return letStmt
}

func (p *Parser) parseYeetStatement() *ast.YeetStatement {
	yeetStmt := &ast.YeetStatement{Token: p.curToken}
	p.advance()

	yeetStmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekIs(token.SEMICOLON) { // optional semicolon
		p.advance()
	}

	return yeetStmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	exprStmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekIs(token.SEMICOLON) { // optional semicolon
		p.advance()
	}

	return exprStmt
}

// EXPRESSIONS

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	for !p.peekIs(token.SEMICOLON) && precedence < getPrecedence(p.peekToken) {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.advance()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.advance()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// PREFIX EXPRESSIONS

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curIs(token.TRUE)}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.Null{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance()

	expr := p.parseExpression(LOWEST)
	if !p.consume(token.RPAREN, "missing closing ')' in grouped expression") {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpr := &ast.IfExpression{Token: p.curToken}

	p.advance()
	ifExpr.Condition = p.parseExpression(LOWEST)

	if !p.consume(token.LBRACE, "missing opening '{' after condition") {
		return nil
	}

	ifExpr.Consequence = p.parseBlockStatement()

	if p.peekIs(token.ELSE) {
		p.advance()
		if !p.consume(token.LBRACE, "missing opening '{' after 'else'") {
			return nil
		}
		ifExpr.Alternative = p.parseBlockStatement()
	}

	return ifExpr
}

func (p *Parser) parseYoyoExpression() ast.Expression {
	yoyoExpr := &ast.YoyoExpression{Token: p.curToken}
	p.advance()

	if !p.curIs(token.SEMICOLON) {
		yoyoExpr.Initialiser = p.parseExpression(LOWEST)
		p.advance()
	}

	p.advance()

	if !p.curIs(token.SEMICOLON) {
		yoyoExpr.Condition = p.parseExpression(LOWEST)
		p.advance()
	}

	if !p.peekIs(token.LBRACE) {
		p.advance()
		yoyoExpr.Post = p.parseExpression(LOWEST)
	}

	if !p.consume(token.LBRACE, "missing opening '{' after 'yoyo'") {
		return nil
	}

	yoyoExpr.Body = p.parseBlockStatement()

	return yoyoExpr
}

func (p *Parser) parseYallExpression() ast.Expression {
	yallExpr := &ast.YallExpression{Token: p.curToken}
	p.advance()

	if p.curIs(token.IDENT) && p.peekIs(token.COLON) {
		yallExpr.KeyName = p.curToken.Literal

		p.advance()
		p.advance()
	} else {
		yallExpr.KeyName = "yt" // yt stands for yeeterator
	}

	yallExpr.Iterable = p.parseExpression(LOWEST)

	if !p.consume(token.LBRACE, "missing opening '{' after 'yall'") {
		return nil
	}

	yallExpr.Body = p.parseBlockStatement()

	return yallExpr
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.curToken}

	for !p.peekIs(token.RBRACKET) && !p.peekIs(token.EOF) {
		p.advance()

		arr.Elements = append(arr.Elements, p.parseExpression(LOWEST))

		if p.peekIs(token.COMMA) { // TODO tighten up after writing more unit tests
			p.advance()
		}
	}

	if !p.consume(token.RBRACKET, "missing closing ']' in array literal") {
		return nil
	}

	return arr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken, Pairs: map[ast.Expression]ast.Expression{}}

	for !p.peekIs(token.RBRACE) && !p.peekIs(token.EOF) {
		p.advance()

		key := p.parseExpression(LOWEST)

		if !p.consume(token.COLON, "missing ':' in hash literal after a key") {
			return nil
		}

		p.advance()
		val := p.parseExpression(LOWEST)

		hash.Pairs[key] = val

		if p.peekIs(token.COMMA) { // TODO tighten up after writing more unit tests
			p.advance()
		}
	}

	if !p.consume(token.RBRACE, "missing closing '}' in hash literal") {
		return nil
	}

	return hash
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken, Parameters: []*ast.Identifier{}}
	if !p.consume(token.LPAREN, "missing opening '(' after function") {
		return nil
	}

	for !p.peekIs(token.RPAREN) && !p.peekIs(token.EOF) {
		p.advance()

		param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		fn.Parameters = append(fn.Parameters, param)

		if p.peekIs(token.COMMA) {
			p.advance()
		}
	}

	if !p.consume(token.RPAREN, "missing closing ')' after function parameters") {
		return nil
	}

	if !p.consume(token.LBRACE, "missing opening '{' before function body") {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseLambdaLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekIs(token.LPAREN) { // parens are optional
		p.advance()
	}

	for !p.peekIs(token.RPAREN) && !p.peekIs(token.LBRACE) && !p.peekIs(token.EOF) {
		p.advance()

		param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		fn.Parameters = append(fn.Parameters, param)

		if p.peekIs(token.COMMA) { // comma is optional
			p.advance()
		}
	}

	if p.peekIs(token.RPAREN) { // parens are optional
		p.advance()
	}

	if !p.consume(token.LBRACE, "missing opening '{' before lambda body") {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

// INFIX EXPRESSIONS

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := getPrecedence(p.curToken)
	p.advance()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseAssignExpression(maybeIdent ast.Expression) ast.Expression {
	ident, ok := maybeIdent.(*ast.Identifier)
	if !ok {
		errMsg := fmt.Sprintf("can only assign to an Identifier (got '%s', type of %T)", maybeIdent, maybeIdent)
		p.errors = append(p.errors, errMsg)
		return nil
	}

	assExpr := &ast.AssignExpression{
		Name:   ident,
		Token:  p.curToken,
		IsInit: p.curIs(token.WALRUS),
	}

	p.advance()
	assExpr.Value = p.parseExpression(LOWEST)

	return assExpr
}

func (p *Parser) parseRangeLiteral(left ast.Expression) ast.Expression {
	rangeLit := &ast.RangeLiteral{
		Token: p.curToken,
		Start: left,
	}

	p.advance()
	rangeLit.End = p.parseExpression(LOWEST)

	return rangeLit
}

func (p *Parser) parseIndexExpression(array ast.Expression) ast.Expression {
	indexExpr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  array,
	}

	p.advance()
	indexExpr.Index = p.parseExpression(LOWEST)

	if !p.consume(token.RBRACKET, "missing closing ']' when indexing an array") {
		return nil
	}

	return indexExpr
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	callExpr := &ast.CallExpression{
		Token:     p.curToken,
		Function:  fn,
		Arguments: []ast.Expression{},
	}

	p.advance()

	for !p.curIs(token.RPAREN) && !p.curIs(token.EOF) {
		callExpr.Arguments = append(callExpr.Arguments, p.parseExpression(LOWEST))

		p.advance()

		if p.curIs(token.COMMA) { // TODO has to be more strict
			p.advance()
		}
	}

	// TODO check this
	// if !p.consume(token.RPAREN, "missing closing ')' after call expression") {
	// 	return nil
	// }

	return callExpr
}

// utils

func (p *Parser) advance() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) consume(t token.TokenType, errMsg string) bool {
	// TODO return error instead of bool?
	if p.peekIs(t) {
		p.advance()
		return true
	} else {
		p.peekError(t, errMsg) // TODO move this up?
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType, errMsg string) {
	msg := fmt.Sprintf("%s (expected '%s', got '%s')", errMsg, t, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for '%s' found", t)
	p.errors = append(p.errors, msg)
}
