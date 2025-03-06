package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"yap/ast"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	FLOAT_OBJ        = "FLOAT"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	FOR_OBJ          = "FOR"
)

type ObjectType string

type BuiltinFunction func(args ...Object) Object

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) HashKey() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 0
	}

	return HashKey{Type: b.Type(), Value: val}
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%g", f.Value)
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(f.Value)}
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Enviroment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var msg bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	msg.WriteString("func")
	msg.WriteString("(")
	msg.WriteString(strings.Join(params, ", "))
	msg.WriteString(") {\n")
	msg.WriteString(f.Body.String())
	msg.WriteString("\n}")

	return msg.String()
}

type For struct {
	Identifer *ast.SayStatement
	Condition []ast.Expression
	Body      *ast.BlockStatement
	Env       *Enviroment
}

func (f *For) Type() ObjectType {
	return FOR_OBJ
}

func (f *For) Inspect() string {
	var msg bytes.Buffer

	conditions := []string{}

	for _, condi := range f.Condition {
		conditions = append(conditions, condi.String())
	}

	msg.WriteString("for (")
	if f.Identifer != nil {
		msg.WriteString(f.Identifer.String() + ", ")
	}
	msg.WriteString(strings.Join(conditions, ", "))
	msg.WriteString(") {\n")
	msg.WriteString(f.Body.String())
	msg.WriteString("\n}")

	return msg.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) Inspect() string {
	var msg bytes.Buffer

	elemets := []string{}

	for _, index := range a.Elements {
		elemets = append(elemets, index.Inspect())
	}

	msg.WriteString("[")
	msg.WriteString(strings.Join(elemets, ", "))
	msg.WriteString("]")

	return msg.String()
}

type Hash struct {
	Pairs map[HashKey]HashPair
	Keys  []Object
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var msg bytes.Buffer

	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	msg.WriteString("{")
	msg.WriteString(strings.Join(pairs, ", "))
	msg.WriteString("}")

	return msg.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hashable interface {
	HashKey() HashKey
}
