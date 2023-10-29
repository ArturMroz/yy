package parser

import (
	"fmt"
	"strconv"

	"yy/ast"
	"yy/lexer"
	"yy/token"
	"yy/yikes"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors    []yikes.YYError
	panicMode bool

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) Errors() []yikes.YYError {
	return p.errors
}

func (p *Parser) advance() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	if p.curToken.Type == token.ERROR {
		p.newError(p.curToken.Literal, p.curToken.Offset)
	}
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
	p.errorAtPeek(t, errMsg)
	return false
}

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	ASSIGNMENT  // = :=
	OR          // ||
	AND         // &&
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
	token.LT_LT:      ASSIGNMENT,
	token.OR:         OR,
	token.AND:        AND,
	token.EQ:         EQUALS,
	token.NOT_EQ:     EQUALS,
	token.LT:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.LT_EQ:      LESSGREATER,
	token.GT_EQ:      LESSGREATER,
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
	p := &Parser{l: l}

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
		token.LBRACE:    p.parseBlockExpression,
		token.YIF:       p.parseYifExpression,
		token.YOLO:      p.parseYoloExpression,
		token.YALL:      p.parseYallExpression,
		token.YOYO:      p.parseYoyoExpression,
		token.BACKSLASH: p.parseLambdaLiteral,
		token.MACRO:     p.parseMacroLiteral,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.OR:         p.parseOrExpression,
		token.AND:        p.parseAndExpression,
		token.PLUS:       p.parseInfixExpression,
		token.MINUS:      p.parseInfixExpression,
		token.ASTERISK:   p.parseInfixExpression,
		token.SLASH:      p.parseInfixExpression,
		token.PERCENT:    p.parseInfixExpression,
		token.EQ:         p.parseInfixExpression,
		token.NOT_EQ:     p.parseInfixExpression,
		token.LT:         p.parseInfixExpression,
		token.GT:         p.parseInfixExpression,
		token.LT_EQ:      p.parseInfixExpression,
		token.GT_EQ:      p.parseInfixExpression,
		token.LT_LT:      p.parseInfixExpression,
		token.RANGE:      p.parseRangeLiteral,
		token.WALRUS:     p.parseDeclareExpression,
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
		if p.panicMode {
			p.sync()
		}
		p.advance()
	}

	return program
}

//
// STATEMENTS
//

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

//
// EXPRESSIONS
//

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errorAtCurrent("unexpected token '%s'", p.curToken.Literal)
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
		p.errorAtCurrent("could not parse %s as integer", p.peekToken.Literal)
		return nil
	}

	return &ast.IntegerLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errorAtCurrent("could not parse %s as float", p.peekToken.Literal)
		return nil
	}

	return &ast.NumberLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := p.curToken.Literal
	matches := []string{}
	result := ""
	lastIdx := 0

	// scan the string to see if it contains any template expressions
	for i := 0; i < len(lit); i++ {
		if lit[i] == '$' && i+1 < len(lit) {
			i++
			if lit[i] == '$' {
				// replace $$ with $
				result += lit[lastIdx : i-1]
				lastIdx = i
			} else if lexer.IsLetter(lit[i]) {
				// replace template expression $var with %s
				start := i
				for i+1 < len(lit) && (lexer.IsLetter(lit[i+1]) || lexer.IsDigit(lit[i+1])) {
					i++
				}
				matches = append(matches, lit[start:i+1])
				result += lit[lastIdx:start-1] + "%s"
				lastIdx = i + 1
			}
		}
	}

	if lastIdx == 0 {
		// there were no template expressions, return regular string literal
		return &ast.StringLiteral{Token: p.curToken, Value: lit}
	}

	result += lit[lastIdx:]

	idents := make([]ast.Expression, 0, len(matches))
	for _, v := range matches {
		idents = append(idents, &ast.Identifier{Value: v})
	}

	return &ast.TemplateStringLiteral{
		Token:    p.curToken,
		Template: result,
		Values:   idents,
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curIs(token.TRUE)}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance()

	expr := p.parseExpression(LOWEST)
	if !p.eat(token.RPAREN, "missing closing ')' in grouped expression") {
		return nil
	}

	return expr
}

func (p *Parser) parseBlockExpression() ast.Expression {
	block := &ast.BlockExpression{Token: p.curToken, Statements: []ast.Statement{}}

	p.advance()
	for !p.curIs(token.RBRACE) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.advance()
	}
	return block
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
			p.errorAtCurrent("expected yif statement or block after 'yels', found '%s'", p.peekToken.Literal)
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

func (p *Parser) parseYoyoExpression() ast.Expression {
	yoyoExpr := &ast.YoyoExpression{Token: p.curToken}
	p.advance()

	yoyoExpr.Condition = p.parseExpression(LOWEST)

	if !p.eat(token.LBRACE, "missing opening '{' after 'yoyo'") {
		return nil
	}

	yoyoExpr.Body = p.parseBlockStatement()

	return yoyoExpr
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

		if p.peekIs(token.COMMA) {
			p.advance()
		} else if !p.peekIs(token.RBRACKET) {
			p.errorAtPeek(token.COMMA, "missing comma after element in array literal")
			return nil
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

		if p.peekIs(token.COMMA) {
			p.advance()
		} else if !p.peekIs(token.RBRACE) {
			p.errorAtPeek(token.COMMA, "missing comma after key-value pair in hashmap literal")
			return nil
		}
	}

	if !p.eat(token.RBRACE, "missing closing '}' in hashmap literal") {
		return nil
	}

	return hashmap
}

func (p *Parser) parseLambdaLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}

	for !p.peekIs(token.LBRACE) && !p.peekIs(token.EOF) {
		p.advance()

		if p.curIs(token.IDENT) {
			param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			fn.Parameters = append(fn.Parameters, param)
		} else {
			p.errorAtCurrent("expected a parameter in lambda declaration, found " + p.curToken.Literal)
			return nil
		}

		if p.peekIs(token.COMMA) { // comma is optional
			p.advance()
		}
	}

	if !p.eat(token.LBRACE, "missing opening '{' before lambda body") {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	fn := &ast.MacroLiteral{Token: p.curToken}

	for !p.peekIs(token.LBRACE) && !p.peekIs(token.EOF) {
		p.advance()

		if p.curIs(token.IDENT) {
			param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			fn.Parameters = append(fn.Parameters, param)
		} else {
			p.errorAtCurrent("expected a parameter in macro declaration, found " + p.curToken.Literal)
			return nil
		}

		if p.peekIs(token.COMMA) { // comma is optional
			p.advance()
		}
	}

	if !p.eat(token.LBRACE, "missing opening '{' before lambda body") {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

//
// INFIX EXPRESSIONS
//

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

func (p *Parser) parseAndExpression(left ast.Expression) ast.Expression {
	expr := &ast.AndExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.advance()
	expr.Right = p.parseExpression(AND)

	return expr
}

func (p *Parser) parseOrExpression(left ast.Expression) ast.Expression {
	expr := &ast.OrExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.advance()
	expr.Right = p.parseExpression(OR)

	return expr
}

func (p *Parser) parseDeclareExpression(maybeIdent ast.Expression) ast.Expression {
	ident, ok := maybeIdent.(*ast.Identifier)
	if !ok {
		p.errorAtCurrent("expected a name when declaring a variable (got '%s')", maybeIdent.TokenLiteral())
		return nil
	}

	declExpr := &ast.DeclareExpression{
		Name:  ident,
		Token: p.curToken,
	}

	p.advance()
	declExpr.Value = p.parseExpression(LOWEST)

	return declExpr
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	assExpr := &ast.AssignExpression{Token: p.curToken}

	switch left.(type) {
	case *ast.Identifier, *ast.IndexExpression:
		assExpr.Left = left
	default:
		p.errorAtCurrent("expected a variable name or index expression when assigning a value (got '%s')", left.TokenLiteral())
		return nil
	}

	p.advance()
	assExpr.Value = p.parseExpression(LOWEST)

	// desugar a += 5 into a = a + 5
	switch assExpr.Token.Type {
	case token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN, token.DIV_ASSIGN, token.MOD_ASSIGN:
		assExpr.Value = &ast.InfixExpression{
			Left:     assExpr.Left,
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
		Token:    p.curToken,
		Function: fn,
	}

	for !p.peekIs(token.RPAREN) && !p.peekIs(token.EOF) {
		p.advance()

		callExpr.Arguments = append(callExpr.Arguments, p.parseExpression(LOWEST))

		if p.peekIs(token.COMMA) {
			p.advance()
		} else if !p.peekIs(token.RPAREN) {
			p.errorAtPeek(token.COMMA, "missing comma after an argument in call expression")
			return nil
		}
	}

	if !p.eat(token.RPAREN, "missing closing ')' in call expression") {
		return nil
	}

	return callExpr
}

//
// ERRORS
//

func (p *Parser) newError(msg string, offset int) {
	if p.panicMode {
		return // don't log cascading errors if we're already panicking
	}

	p.panicMode = true
	p.errors = append(p.errors, yikes.YYError{Msg: msg, Offset: offset})
}

func (p *Parser) errorAtPeek(expected token.TokenType, errMsg string) {
	msg := fmt.Sprintf("%s (expected '%s', found '%s')", errMsg, expected, p.peekToken.Literal)
	p.newError(msg, p.peekToken.Offset)
}

func (p *Parser) errorAtCurrent(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	p.newError(msg, p.curToken.Offset)
}

// sync recovers from panic mode by fastforwarding to the next expr/stmt
func (p *Parser) sync() {
	p.panicMode = false

	p.advance()

	for !p.curIs(token.EOF) {
		if p.curIs(token.SEMICOLON) {
			return
		}

		switch p.peekToken.Type {
		case token.YEET, token.YIF, token.YALL, token.YOYO, token.YOLO, token.BACKSLASH, token.MACRO:
			return

		default:
			p.advance()
		}
	}
}
