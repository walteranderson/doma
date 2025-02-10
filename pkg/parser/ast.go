package parser

import (
	"bytes"
	"doma/pkg/lexer"
	"fmt"
	"strings"
)

type Expression interface {
	TokenLiteral() string
	String() string
}

// ------------------------------
// Atoms
// ------------------------------

type String struct {
	Token lexer.Token
	Value string
}

func (s *String) TokenLiteral() string {
	return s.Token.Literal
}
func (s *String) String() string {
	return s.Value
}

// ---

type Number struct {
	Token lexer.Token
	Value int64
}

func (n *Number) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

// ---

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

// ---

type BuiltinIdentifier struct {
	Token lexer.Token
	Value string
}

func (i *BuiltinIdentifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *BuiltinIdentifier) String() string {
	return i.Value
}

// ---

type Boolean struct {
	Token lexer.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

// ---

type Symbol struct {
	Token lexer.Token
	Value string
}

func (s *Symbol) TokenLiteral() string {
	return s.Token.Literal
}
func (s *Symbol) String() string {
	return s.Token.Literal
}

// ------------------------------
// Lists
// ------------------------------

type Program struct {
	Args []Expression
}

func (p *Program) TokenLiteral() string {
	if len(p.Args) > 0 {
		return p.Args[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, arg := range p.Args {
		out.WriteString(arg.String() + "\n")
	}
	return out.String()
}

// ---

type Form struct {
	Token lexer.Token
	First Expression
	Rest  []Expression
}

func (f *Form) TokenLiteral() string {
	return f.Token.Literal
}
func (f *Form) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	args = append(args, f.First.String())
	for _, r := range f.Rest {
		args = append(args, r.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}

// ---

type List struct {
	Token lexer.Token
	Args  []Expression
}

func (lf *List) TokenLiteral() string {
	return lf.Token.Literal
}
func (lf *List) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, r := range lf.Args {
		args = append(args, r.String())
	}
	out.WriteString("'(")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}
