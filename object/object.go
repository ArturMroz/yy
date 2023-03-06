package object

import (
	"fmt"
	"strings"

	"yy/ast"
)

type Object interface {
	Type() ObjectType
	String() string
}

type ObjectType int

const (
	INTEGER_OBJ ObjectType = iota
	NUMBER_OBJ
	BOOLEAN_OBJ
	STRING_OBJ
	NULL_OBJ

	ARRAY_OBJ
	HASHMAP_OBJ
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
	NUMBER_OBJ:  "NUMBER",
	BOOLEAN_OBJ: "BOOLEAN",
	STRING_OBJ:  "STRING",
	NULL_OBJ:    "NULL",

	ARRAY_OBJ:   "ARRAY",
	HASHMAP_OBJ: "HASHMAP",
	RANGE_OBJ:   "RANGE",

	FUNCTION_OBJ: "FUNCTION",
	BUILTIN_OBJ:  "BUILTIN",
	QUOTE_OBJ:    "QUOTE",
	MACRO_OBJ:    "MACRO",

	ERROR_OBJ:        "ERROR",
	RETURN_VALUE_OBJ: "RETURN_VALUE",
}

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	ABYSS = &String{Value: "Stare at the abyss long enough, and it starts to stare back at you."}
)

func (ot ObjectType) String() string {
	return objectTypes[ot]
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) String() string   { return fmt.Sprintf("%d", i.Value) }

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) String() string   { return fmt.Sprintf("%g", n.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) String() string   { return s.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) String() string   { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) String() string   { return "null" }

type Range struct {
	Start int64
	End   int64
}

func (r *Range) Type() ObjectType { return RANGE_OBJ }
func (r *Range) String() string   { return fmt.Sprintf("%d..%d", r.Start, r.End) }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) String() string {
	var b strings.Builder
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.String())
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

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	return HashKey{Type: s.Type(), Value: hashString(s.Value)}
}

func hashString(key string) uint64 {
	// FNV-1a algorithm
	hash := uint64(2166136261)
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= 16777619
	}
	return hash
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hashmap struct {
	Pairs map[HashKey]HashPair
}

func (h *Hashmap) Type() ObjectType { return HASHMAP_OBJ }
func (h *Hashmap) String() string {
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.String(), pair.Value.String()))
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) String() string   { return rv.Value.String() }

type Error struct {
	Msg string
	// TODO add stack trace, column, line etc
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) String() string   { return e.Msg }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) String() string {
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
func (m *Macro) String() string {
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
func (b *Builtin) String() string   { return "builtin function" }

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) String() string   { return fmt.Sprintf("QUOTE(%s)", q.Node) }
