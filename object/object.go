package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"yy/ast"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type ObjectType int

const (
	INTEGER_OBJ ObjectType = iota
	BOOLEAN_OBJ
	STRING_OBJ
	NULL_OBJ

	ARRAY_OBJ
	HASH_OBJ
	RANGE_OBJ

	FUNCTION_OBJ
	BUILTIN_OBJ
	QUOTE_OBJ
	MACRO_OBJ

	ERROR_OBJ
	RETURN_VALUE_OBJ
)

var objectTypes = [...]string{
	INTEGER_OBJ: "INTEGER",
	BOOLEAN_OBJ: "BOOLEAN",
	STRING_OBJ:  "STRING",
	NULL_OBJ:    "NULL",

	ARRAY_OBJ: "ARRAY",
	HASH_OBJ:  "HASH",
	RANGE_OBJ: "RANGE",

	FUNCTION_OBJ: "FUNCTION",
	BUILTIN_OBJ:  "BUILTIN",
	QUOTE_OBJ:    "QUOTE",
	MACRO_OBJ:    "MACRO",

	ERROR_OBJ:        "ERROR",
	RETURN_VALUE_OBJ: "RETURN_VALUE",
}

func (ot ObjectType) String() string {
	return objectTypes[ot]
}

const YoloKey = "$$yolo"

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type Range struct {
	Start int64
	End   int64
}

func (r *Range) Type() ObjectType { return RANGE_OBJ }
func (r *Range) Inspect() string  { return fmt.Sprintf("%d..%d", r.Start, r.End) }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var b strings.Builder
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	b.WriteString("[")
	b.WriteString(strings.Join(elements, ", "))
	b.WriteString("]")
	return b.String()
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	// TODO handle hash collisions
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Msg string
	// TODO add stack trace, column, line etc
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return e.Msg }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var b strings.Builder

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	b.WriteString("fun")
	b.WriteString("(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") {\n")
	b.WriteString(f.Body.String())
	b.WriteString("\n}")

	return b.String()
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string {
	var b strings.Builder

	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	b.WriteString("@\\")
	b.WriteString("(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") {\n")
	b.WriteString(m.Body.String())
	b.WriteString("\n}")

	return b.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string  { return fmt.Sprintf("QUOTE(%s)", q.Node) }
