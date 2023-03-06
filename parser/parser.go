package parser

import (
	"fmt"
	"strconv"
	"strings"

	"yy/ast"
	"yy/lexer"
	"yy/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) Errors() []string {
	return p.errors
}

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

func (p *Parser) eat(t token.TokenType, errMsg string) bool {
	if p.peekIs(t) {
		p.advance()
		return true
	}
	p.peekError(t, errMsg)
	return false
}

func (p *Parser) peekError(t token.TokenType, errMsg string) {
	msg := fmt.Sprintf("[line %d] %s (expected '%s', found '%s')", p.curToken.Line, errMsg, t, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError() {
	msg := fmt.Sprintf("[line %d] unexpected token '%s' near '%s'", p.curToken.Line, p.curToken.Literal, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) newError(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg = fmt.Sprintf("[line %d] %s", p.curToken.Line, msg)
	p.errors = append(p.errors, msg)
}

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	ASSIGNMENT  // = :=
	EQUALS      // == !=
	LESSGREATER // > <
	RANGE       // x..y
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -x !x
	CALL        // function(x)
	INDEX       // array[idx]
)

var precedences = map[token.TokenType]Precedence{
	token.ASSIGN:     ASSIGNMENT,
	token.WALRUS:     ASSIGNMENT,
	token.ADD_ASSIGN: ASSIGNMENT,
	token.SUB_ASSIGN: ASSIGNMENT,
	token.MUL_ASSIGN: ASSIGNMENT,
	token.DIV_ASSIGN: ASSIGNMENT,
	token.MOD_ASSIGN: ASSIGNMENT,
	token.EQ:         EQUALS,
	token.NOT_EQ:     EQUALS,
	token.LT:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.RANGE:      RANGE,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.SLASH:      PRODUCT,
	token.ASTERISK:   PRODUCT,
	token.PERCENT:    PRODUCT,
	token.LPAREN:     CALL,
	token.LBRACKET:   INDEX,
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
		token.NUMBER:    p.parseNumberLiteral,
		token.STRING:    p.parseStringLiteral,
		token.MINUS:     p.parsePrefixExpression,
		token.BANG:      p.parsePrefixExpression,
		token.TRUE:      p.parseBoolean,
		token.FALSE:     p.parseBoolean,
		token.NULL:      p.parseNull,
		token.LPAREN:    p.parseGroupedExpression,
		token.LBRACKET:  p.parseArrayLiteral,
		token.HASHMAP:   p.parseHashmapLiteral,
		token.YIF:       p.parseYifExpression,
		token.YOLO:      p.parseYoloExpression,
		token.YALL:      p.parseYallExpression,
		token.YET:       p.parseYetExpression,
		token.BACKSLASH: p.parseLambdaLiteral,
		token.MACRO:     p.parseMacroLiteral,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:       p.parseInfixExpression,
		token.MINUS:      p.parseInfixExpression,
		token.ASTERISK:   p.parseInfixExpression,
		token.SLASH:      p.parseInfixExpression,
		token.PERCENT:    p.parseInfixExpression,
		token.EQ:         p.parseInfixExpression,
		token.NOT_EQ:     p.parseInfixExpression,
		token.LT:         p.parseInfixExpression,
		token.GT:         p.parseInfixExpression,
		token.RANGE:      p.parseRangeLiteral,
		token.WALRUS:     p.parseAssignExpression,
		token.ASSIGN:     p.parseAssignExpression,
		token.ADD_ASSIGN: p.parseAssignExpression,
		token.SUB_ASSIGN: p.parseAssignExpression,
		token.MUL_ASSIGN: p.parseAssignExpression,
		token.DIV_ASSIGN: p.parseAssignExpression,
		token.MOD_ASSIGN: p.parseAssignExpression,
		token.LPAREN:     p.parseCallExpression,
		token.LBRACKET:   p.parseIndexExpression,
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
		p.noPrefixParseFnError()
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
		p.newError("could not parse %s as integer", p.peekToken.Literal)
		return nil
	}

	return &ast.IntegerLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.newError("could not parse %s as float", p.peekToken.Literal)
		return nil
	}

	return &ast.NumberLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := p.curToken.Literal

	// TODO this handles only the happy path
	matches := []int{}
	for i, ch := range lit {
		if ch == '{' || ch == '}' {
			matches = append(matches, i)
		}
	}

	if len(matches) == 0 {
		return &ast.StringLiteral{Token: p.curToken, Value: lit}
	}

	// TODO support expr in templated strings (currently only Identifiers are supported)
	idents := make([]ast.Expression, 0, len(matches)/2)
	replaced := lit
	offset := 0
	for i := 0; i < len(matches); i += 2 {
		fst, snd := matches[i], matches[i+1]

		ident := strings.TrimSpace(lit[fst+1 : snd])
		replaced = replaced[:fst-offset] + "%s" + replaced[snd+1-offset:]
		offset += snd - fst - 1

		idents = append(idents, &ast.Identifier{Value: ident})
	}

	return &ast.TemplateStringLiteral{
		Token:    p.curToken,
		Template: replaced,
		Values:   idents,
	}
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
	if !p.eat(token.RPAREN, "missing closing ')' in grouped expression") {
		return nil
	}

	return expr
}

func (p *Parser) parseYifExpression() ast.Expression {
	yifExpr := &ast.YifExpression{Token: p.curToken}

	p.advance()
	yifExpr.Condition = p.parseExpression(LOWEST)

	if !p.eat(token.LBRACE, "missing opening '{' after 'yif' condition") {
		return nil
	}

	yifExpr.Consequence = p.parseBlockStatement()

	if p.peekIs(token.YELS) {
		p.advance()

		if p.peekIs(token.LBRACE) {
			p.advance()
			yifExpr.Alternative = p.parseBlockStatement()
		} else if p.peekIs(token.YIF) {
			p.advance()
			// TODO this is a bit clunky
			yifExpr.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: p.parseYifExpression()},
				},
			}
		} else {
			p.newError("expected yif statement or block after 'yels', found '%s'", p.peekToken.Literal)
			return nil
		}
	}

	return yifExpr
}

func (p *Parser) parseYoloExpression() ast.Expression {
	yoloExpr := &ast.YoloExpression{Token: p.curToken}

	if !p.eat(token.LBRACE, "missing opening '{' after 'yolo'") {
		return nil
	}

	yoloExpr.Body = p.parseBlockStatement()

	return yoloExpr
}

func (p *Parser) parseYetExpression() ast.Expression {
	yetExpr := &ast.YetExpression{Token: p.curToken}
	p.advance()

	yetExpr.Condition = p.parseExpression(LOWEST)

	if !p.eat(token.LBRACE, "missing opening '{' after 'yet'") {
		return nil
	}

	yetExpr.Body = p.parseBlockStatement()

	return yetExpr
}

func (p *Parser) parseYallExpression() ast.Expression {
	yallExpr := &ast.YallExpression{Token: p.curToken, KeyName: "yt"}
	p.advance()

	if p.curIs(token.IDENT) && p.peekIs(token.COLON) {
		yallExpr.KeyName = p.curToken.Literal

		p.advance()
		p.advance()
	}

	yallExpr.Iterable = p.parseExpression(LOWEST)

	if !p.eat(token.LBRACE, "missing opening '{' after 'yall'") {
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

		if p.peekIs(token.COMMA) { // TODO make comma nonoptional
			p.advance()
		}
	}

	if !p.eat(token.RBRACKET, "missing closing ']' in array literal") {
		return nil
	}

	return arr
}

func (p *Parser) parseHashmapLiteral() ast.Expression {
	hashmap := &ast.HashmapLiteral{
		Token: p.curToken,
		Pairs: map[ast.Expression]ast.Expression{},
	}

	for !p.peekIs(token.RBRACE) && !p.peekIs(token.EOF) {
		p.advance()

		key := p.parseExpression(LOWEST)

		if !p.eat(token.COLON, "missing ':' in hashmap literal after a key") {
			return nil
		}

		p.advance()
		val := p.parseExpression(LOWEST)

		hashmap.Pairs[key] = val

		if p.peekIs(token.COMMA) { // TODO tighten up after writing more unit tests
			p.advance()
		}
	}

	if !p.eat(token.RBRACE, "missing closing '}' in hashmap literal") {
		return nil
	}

	return hashmap
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

	if !p.eat(token.LBRACE, "missing opening '{' before lambda body") {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	fn := &ast.MacroLiteral{Token: p.curToken}

	if p.peekIs(token.LPAREN) { // parens are optional
		p.advance()
	}

	// TODO dry param parsing
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

	if !p.eat(token.LBRACE, "missing opening '{' before lambda body") {
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
		// TODO support assignment to index expr ie array[2] = 5;
		p.newError("can only assign to an Identifier (got '%s')", maybeIdent)
		return nil
	}

	assExpr := &ast.AssignExpression{
		Name:   ident,
		Token:  p.curToken,
		IsInit: p.curIs(token.WALRUS),
	}

	p.advance()
	assExpr.Value = p.parseExpression(LOWEST)

	// TODO rethink this implementation, maybe move to eval stage instead of desugaring
	switch assExpr.Token.Type {
	case token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN, token.DIV_ASSIGN, token.MOD_ASSIGN:
		assExpr.Value = &ast.InfixExpression{
			Left:     assExpr.Name,
			Right:    assExpr.Value,
			Operator: string(assExpr.Token.Literal[0]),
		}
		assExpr.Token = token.Token{Type: token.ASSIGN, Literal: "="}
	}

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

	if !p.eat(token.RBRACKET, "missing closing ']' when indexing an array") {
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
