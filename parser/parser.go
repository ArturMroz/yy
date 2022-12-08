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
	PREFIX      // -X or !X
	CALL        // myFunction(X)
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
		token.MINUS:    p.parsePrefixExpression,
		token.BANG:     p.parsePrefixExpression,
		token.TRUE:     p.parseBoolean,
		token.FALSE:    p.parseBoolean,
		token.LPAREN:   p.parseGroupedExpression,
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
	if !p.match(token.IDENT) {
		return nil // TODO return error
	}

	letStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.match(token.ASSIGN) {
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if !p.match(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpr := &ast.IfExpression{Token: p.curToken}
	if !p.match(token.LPAREN) {
		return nil
	}

	p.nextToken()
	ifExpr.Condition = p.parseExpression(LOWEST)
	if !p.match(token.RPAREN) {
		return nil
	}

	if !p.match(token.LBRACE) {
		return nil
	}

	ifExpr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.match(token.LBRACE) {
			return nil
		}
		ifExpr.Alternative = p.parseBlockStatement()
	}

	return ifExpr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken, Parameters: []*ast.Identifier{}}
	if !p.match(token.LPAREN) {
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

	if !p.match(token.RPAREN) {
		return nil
	}

	if !p.match(token.LBRACE) {
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

	return callExpr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		// if stmt != nil { // TODO this will always be true
		// 	block.Statements = append(block.Statements, stmt)
		// }
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

func (p *Parser) match(t token.TokenType) bool {
	// TODO return error instead of bool?
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t) // TODO move this up?
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	// TODO allow custom error msg to be passed
	msg := fmt.Sprintf("parse error: expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("[!] no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
