package ast

import (
	"fmt"
	"strings"

	"ylang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var b strings.Builder
	for _, s := range p.Statements {
		b.WriteString(s.String())
	}
	return b.String()
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var b strings.Builder

	// fmt.Fprintf(&b, "%s %s = %s;", ls.TokenLiteral(), ls.Name.String(), ls.Value.String())

	fmt.Fprintf(&b, "%s %s = ", ls.TokenLiteral(), ls.Name.String())
	if ls.Value != nil {
		b.WriteString(ls.Value.String())
	}
	b.WriteString(";")

	return b.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type YeetStatement struct {
	Token       token.Token // the 'yeet' token
	ReturnValue Expression
}

func (ys *YeetStatement) statementNode()       {}
func (ys *YeetStatement) TokenLiteral() string { return ys.Token.Literal }
func (ys *YeetStatement) String() string {
	var b strings.Builder
	b.WriteString(ys.TokenLiteral() + " ")
	if ys.ReturnValue != nil {
		b.WriteString(ys.ReturnValue.String())
	}
	b.WriteString(";")
	return b.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type Null struct {
	Token token.Token
}

func (n *Null) expressionNode()      {}
func (n *Null) TokenLiteral() string { return n.Token.Literal }
func (n *Null) String() string       { return n.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var b strings.Builder

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	b.WriteString("[")
	b.WriteString(strings.Join(elements, ", "))
	b.WriteString("]")
	return b.String()
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var b strings.Builder
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	b.WriteString("{")
	b.WriteString(strings.Join(pairs, ", "))
	b.WriteString("}")

	return b.String()
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(ie.Left.String())
	b.WriteString("[")
	b.WriteString(ie.Index.String())
	b.WriteString("])")
	return b.String()
}

type PrefixExpression struct {
	Token    token.Token // prefix token e.g. !
	Operator string      // '-' or '!'
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var b strings.Builder
	b.WriteString("if ")
	b.WriteString(ie.Condition.String())
	b.WriteString(" ")
	b.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		b.WriteString("else ")
		b.WriteString(ie.Alternative.String())
	}
	return b.String()
}

type YoyoExpression struct {
	Token     token.Token // The 'yoyo' token
	Condition Expression
	Body      *BlockStatement
}

func (ye *YoyoExpression) expressionNode()      {}
func (ye *YoyoExpression) TokenLiteral() string { return ye.Token.Literal }
func (ye *YoyoExpression) String() string {
	return fmt.Sprintf("yoyo %s { %s }", ye.Condition.String(), ye.Body.String())
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	stmts := []string{}
	for _, p := range bs.Statements {
		stmts = append(stmts, p.String())
	}

	var b strings.Builder
	b.WriteString("{ ")
	// for _, s := range bs.Statements {
	// 	b.WriteString(s.String())
	// }

	b.WriteString(strings.Join(stmts, "; "))
	b.WriteString(" }")
	return b.String()
}

type FunctionLiteral struct {
	Token      token.Token // 'fun' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var b strings.Builder

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	b.WriteString(fl.TokenLiteral())
	b.WriteString("(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") ")
	b.WriteString(fl.Body.String())

	return b.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var b strings.Builder

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	b.WriteString(ce.Function.String())
	b.WriteString("(")
	b.WriteString(strings.Join(args, ", "))
	b.WriteString(")")

	return b.String()
}
