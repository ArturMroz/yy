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

	// out.WriteString(ls.TokenLiteral() + " ")
	// out.WriteString(ls.Name.String())
	// out.WriteString(" = ")
	// if ls.Value != nil {
	// 	out.WriteString(ls.Value.String())
	// }
	// out.WriteString(";")
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
