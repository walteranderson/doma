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
	ERROR_OBJ     = "ERROR"
	NUMBER_OBJ    = "NUMBER"
	STRING_OBJ    = "STRING"
	BOOLEAN_OBJ   = "BOOLEAN"
	LIST_OBJ      = "LIST"
	LAMBDA_OBJ    = "LAMBDA"
	BUILTIN_OBJ   = "BUILTIN"
	PROCEDURE_OBJ = "PROCEDURE"
	SYMBOL_OBJ    = "SYMBOL"
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
	return "#<procedure>"
}

type Builtin struct {
	Value lexer.TokenType
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string {
	return fmt.Sprintf("#<procedure:%s>", b.Value)
}

type Procedure struct {
	Name  string
	Value *Lambda
}

func (s *Procedure) Type() ObjectType { return PROCEDURE_OBJ }
func (s *Procedure) Inspect() string {
	return fmt.Sprintf("#<procedure:%s>", s.Name)
}
