package ast

import (
	"fmt"
	"strings"

	"yy/token"
)

type Node interface {
	Pos() int
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

func (p *Program) Pos() int { return 0 }

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
		b.WriteString(s.String() + ";")
	}
	return b.String()
}

type DeclareExpression struct {
	Token token.Token // the ':=' token
	Name  *Identifier
	Value Expression
}

func (de *DeclareExpression) expressionNode()      {}
func (de *DeclareExpression) Pos() int             { return de.Token.Offset }
func (de *DeclareExpression) TokenLiteral() string { return de.Token.Literal }
func (de *DeclareExpression) String() string {
	return fmt.Sprintf("(%s := %s)", de.Name.String(), de.Value.String())
}

type AssignExpression struct {
	Token token.Token // '=' token
	Left  Expression
	Value Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) Pos() int             { return ae.Token.Offset }
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) String() string {
	return fmt.Sprintf("(%s = %s)", ae.Left.String(), ae.Value.String())
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) Pos() int             { return i.Token.Offset }
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type YeetStatement struct {
	Token       token.Token // the 'yeet' token
	ReturnValue Expression
}

func (ys *YeetStatement) statementNode()       {}
func (ys *YeetStatement) Pos() int             { return ys.Token.Offset }
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
func (de *ExpressionStatement) Pos() int             { return de.Token.Offset }
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

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) Pos() int             { return i.Token.Offset }
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) expressionNode()      {}
func (n *NumberLiteral) Pos() int             { return n.Token.Offset }
func (n *NumberLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NumberLiteral) String() string       { return n.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) Pos() int             { return b.Token.Offset }
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }

type NullLiteral struct {
	Token token.Token
}

func (n *NullLiteral) expressionNode()      {}
func (n *NullLiteral) Pos() int             { return n.Token.Offset }
func (n *NullLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NullLiteral) String() string       { return n.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) Pos() int             { return s.Token.Offset }
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return `"` + s.Token.Literal + `"` }

type TemplateStringLiteral struct {
	Token    token.Token
	Template string
	Values   []Expression
}

func (ts *TemplateStringLiteral) expressionNode()      {}
func (ts *TemplateStringLiteral) Pos() int             { return ts.Token.Offset }
func (ts *TemplateStringLiteral) TokenLiteral() string { return ts.Token.Literal }
func (ts *TemplateStringLiteral) String() string       { return `"` + ts.Token.Literal + `"` }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) Pos() int             { return a.Token.Offset }
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	var b strings.Builder

	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}

	b.WriteString("[")
	b.WriteString(strings.Join(elements, ", "))
	b.WriteString("]")
	return b.String()
}

type RangeLiteral struct {
	Token token.Token // the '..' token
	Start Expression
	End   Expression
}

func (rl *RangeLiteral) expressionNode()      {}
func (rl *RangeLiteral) Pos() int             { return rl.Token.Offset }
func (rl *RangeLiteral) TokenLiteral() string { return rl.Token.Literal }
func (rl *RangeLiteral) String() string {
	return fmt.Sprintf("(%s..%s)", rl.Start, rl.End)
}

type HashmapLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashmapLiteral) expressionNode()      {}
func (hl *HashmapLiteral) Pos() int             { return hl.Token.Offset }
func (hl *HashmapLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashmapLiteral) String() string {
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

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) Pos() int             { return ie.Left.Pos() }
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
func (pe *PrefixExpression) Pos() int             { return pe.Token.Offset }
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
func (ie *InfixExpression) Pos() int             { return ie.Token.Offset }
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

type AndExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ae *AndExpression) expressionNode()      {}
func (ae *AndExpression) Pos() int             { return ae.Token.Offset }
func (ae *AndExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AndExpression) String() string {
	return fmt.Sprintf("(%s && %s)", ae.Left.String(), ae.Right.String())
}

type OrExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ae *OrExpression) expressionNode()      {}
func (de *OrExpression) Pos() int             { return de.Token.Offset }
func (oe *OrExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *OrExpression) String() string {
	return fmt.Sprintf("(%s || %s)", oe.Left.String(), oe.Right.String())
}

type YifExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ye *YifExpression) expressionNode()      {}
func (ye *YifExpression) Pos() int             { return ye.Token.Offset }
func (ye *YifExpression) TokenLiteral() string { return ye.Token.Literal }
func (ye *YifExpression) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "yif %s %s", ye.Condition, ye.Consequence)
	if ye.Alternative != nil {
		b.WriteString(" yels ")
		b.WriteString(ye.Alternative.String())
	}
	return b.String()
}

type YoloExpression struct {
	Token token.Token
	Body  *BlockStatement
}

func (ye *YoloExpression) expressionNode()      {}
func (ye *YoloExpression) Pos() int             { return ye.Token.Offset }
func (ye *YoloExpression) TokenLiteral() string { return ye.Token.Literal }
func (ye *YoloExpression) String() string {
	return fmt.Sprintf("yolo { %s }", ye.Body.String())
}

type YoyoExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ye *YoyoExpression) expressionNode()      {}
func (de *YoyoExpression) Pos() int             { return de.Token.Offset }
func (ye *YoyoExpression) TokenLiteral() string { return ye.Token.Literal }
func (ye *YoyoExpression) String() string {
	return fmt.Sprintf("yoyo %s { %s }", ye.Condition.String(), ye.Body.String())
}

type YallExpression struct {
	Token    token.Token
	Iterable Expression
	KeyName  string
	Body     *BlockStatement
}

func (ye *YallExpression) expressionNode()      {}
func (ye *YallExpression) Pos() int             { return ye.Token.Offset }
func (ye *YallExpression) TokenLiteral() string { return ye.Token.Literal }
func (ye *YallExpression) String() string {
	return fmt.Sprintf("yall %s: %s { %s }", ye.KeyName, ye.Iterable.String(), ye.Body.String())
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) Pos() int             { return bs.Token.Offset }
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	stmts := []string{}
	for _, p := range bs.Statements {
		stmts = append(stmts, p.String())
	}

	var b strings.Builder
	b.WriteString("{ ")

	b.WriteString(strings.Join(stmts, "; "))
	b.WriteString(" }")
	return b.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) Pos() int             { return fl.Token.Offset }
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

type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml *MacroLiteral) expressionNode()      {}
func (ml *MacroLiteral) Pos() int             { return ml.Token.Offset }
func (ml *MacroLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MacroLiteral) String() string {
	var b strings.Builder

	params := []string{}
	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}

	b.WriteString(ml.TokenLiteral())
	b.WriteString("(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") ")
	b.WriteString(ml.Body.String())

	return b.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) Pos() int             { return ce.Token.Offset }
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
