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
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // myFunction(x)
	INDEX       // myArray[idx]
)

var precedences = map[token.TokenType]Precedence{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func (p *Parser) peekPrecedence() Precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() Precedence {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT:    p.parseIdentifier,
		token.INT:      p.parseIntegerLiteral,
		token.STRING:   p.parseStringLiteral,
		token.MINUS:    p.parsePrefixExpression,
		token.BANG:     p.parsePrefixExpression,
		token.TRUE:     p.parseBoolean,
		token.FALSE:    p.parseBoolean,
		token.NULL:     p.parseNull,
		token.LPAREN:   p.parseGroupedExpression,
		token.LBRACKET: p.parseArrayLiteral,
		token.LBRACE:   p.parseHashLiteral,
		token.IF:       p.parseIfExpression,
		token.FUNCTION: p.parseFunctionLiteral,
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
		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
	}

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

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

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.curToken}
	if !p.consume(token.IDENT, "missing identfier after 'let' keyword") {
		return nil // TODO return error
	}

	letStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.consume(token.ASSIGN, "missing '=' after identifier in assignment") {
		return nil
	}

	p.nextToken()

	letStmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) { // optional semicolon
		p.nextToken()
	}

	return letStmt
}

func (p *Parser) parseYeetStatement() *ast.YeetStatement {
	yeetStmt := &ast.YeetStatement{Token: p.curToken}
	p.nextToken()

	yeetStmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) { // optional semicolon
		p.nextToken()
	}

	return yeetStmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	exprStmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekTokenIs(token.SEMICOLON) { // optional semicolon
		p.nextToken()
	}

	return exprStmt
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)
	return expr
}

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
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.Null{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if !p.consume(token.RPAREN, "missing closing ')' in grouped expression") {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpr := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	ifExpr.Condition = p.parseExpression(LOWEST)

	if !p.consume(token.LBRACE, "missing opening '{' after condition") {
		return nil
	}

	ifExpr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.consume(token.LBRACE, "missing opening '{' after 'else'") {
			return nil
		}
		ifExpr.Alternative = p.parseBlockStatement()
	}

	return ifExpr
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.curToken}

	for !p.peekTokenIs(token.RBRACKET) && !p.peekTokenIs(token.EOF) {
		p.nextToken()

		arr.Elements = append(arr.Elements, p.parseExpression(LOWEST))

		if p.peekTokenIs(token.COMMA) { // TODO tighten up after writing more unit tests
			p.nextToken()
		}
	}

	if !p.consume(token.RBRACKET, "missing closing ']' in array literal") {
		return nil
	}

	return arr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken, Pairs: map[ast.Expression]ast.Expression{}}

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()

		key := p.parseExpression(LOWEST)

		if !p.consume(token.COLON, "missing ':' in hash literal after a key") {
			return nil
		}

		p.nextToken()
		val := p.parseExpression(LOWEST)

		hash.Pairs[key] = val

		if p.peekTokenIs(token.COMMA) { // TODO tighten up after writing more unit tests
			p.nextToken()
		}
	}

	if !p.consume(token.RBRACE, "missing closing '}' in hash literal") {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(array ast.Expression) ast.Expression {
	indexExpr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  array,
	}

	p.nextToken()
	indexExpr.Index = p.parseExpression(LOWEST)

	if !p.consume(token.RBRACKET, "missing closing ']' when indexing an array") {
		return nil
	}

	return indexExpr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken, Parameters: []*ast.Identifier{}}
	if !p.consume(token.LPAREN, "missing opening '(' after function") {
		return nil
	}

	for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
		p.nextToken()

		param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		fn.Parameters = append(fn.Parameters, param)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
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

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	callExpr := &ast.CallExpression{
		Token:     p.curToken,
		Function:  fn,
		Arguments: []ast.Expression{},
	}

	p.nextToken()

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		callExpr.Arguments = append(callExpr.Arguments, p.parseExpression(LOWEST))

		p.nextToken()

		if p.curTokenIs(token.COMMA) { // TODO has to be more strict
			p.nextToken()
		}
	}

	// TODO check this
	// if !p.consume(token.RPAREN, "missing closing ')' after call expression") {
	// 	return nil
	// }

	return callExpr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}
	return block
}

// utils

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) consume(t token.TokenType, errMsg string) bool {
	// TODO return error instead of bool?
	if p.peekTokenIs(t) {
		p.nextToken()
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
