package eval

import (
	"bytes"
	"doma/pkg/lexer"
	"doma/pkg/parser"
	"fmt"
	"strings"
)

type ObjectType string

const (
	ERROR_OBJ   = "ERROR"
	NUMBER_OBJ  = "NUMBER"
	STRING_OBJ  = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"
	LIST_OBJ    = "LIST"
	LAMBDA_OBJ  = "LAMBDA"
	BUILTIN_OBJ = "BUILTIN"
	SYMBOL_OBJ  = "SYMBOL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// ---

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// ---

type Number struct {
	Value int64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string {
	return fmt.Sprintf("%d", n.Value)
}

// ---

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string {
	return s.Value
}

// ---

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "#t"
	} else {
		return "#f"
	}
}

// ---

type Symbol struct {
	Value string
}

func (s *Symbol) Type() ObjectType { return SYMBOL_OBJ }
func (s *Symbol) Inspect() string {
	return fmt.Sprintf("'%s", s.Value)
}

// ---

type List struct {
	Args []parser.Expression
}

func (l *List) Type() ObjectType { return LIST_OBJ }
func (l *List) Inspect() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, arg := range l.Args {
		args = append(args, arg.String())
	}
	out.WriteString("'(")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}

type Lambda struct {
	Params []*parser.Identifier
	Body   []parser.Expression
	Env    *Env
}

func (l *Lambda) Type() ObjectType { return LAMBDA_OBJ }
func (l *Lambda) Inspect() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, p := range l.Params {
		args = append(args, p.String())
	}
	for _, b := range l.Body {
		args = append(args, b.String())
	}
	out.WriteString("(lambda ")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}

type Builtin struct {
	Value lexer.TokenType
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string {
	return fmt.Sprintf("<procedure:%s>", b.Value)
}
