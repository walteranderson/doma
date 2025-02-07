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

// ---

type IfProcedure struct {
	Token       lexer.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (i *IfProcedure) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IfProcedure) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	args = append(args, i.Condition.String())
	args = append(args, i.Consequence.String())
	args = append(args, i.Alternative.String())
	out.WriteString("(if ")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}

// ---

type LambdaProcedure struct {
	Token  lexer.Token
	Params []*Identifier
	Body   []Expression
}

func (l *LambdaProcedure) TokenLiteral() string {
	return l.Token.Literal
}
func (l *LambdaProcedure) String() string {
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

// ---

type DefineProcedure struct {
	Token lexer.Token
	Name  string
	Value Expression
}

func (d *DefineProcedure) TokenLiteral() string {
	return d.Token.Literal
}
func (d *DefineProcedure) String() string {
	var out bytes.Buffer
	out.WriteString("(define ")
	out.WriteString(d.Name)
	out.WriteString(" " + d.Value.String())
	out.WriteString(")")
	return out.String()
}

// ---

type DisplayProcedure struct {
	Token lexer.Token
	Args  []Expression
}

func (d *DisplayProcedure) TokenLiteral() string {
	return d.Token.Literal
}
func (d *DisplayProcedure) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, arg := range d.Args {
		args = append(args, arg.String())
	}
	out.WriteString("(display ")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}

// ---

type MathProcedure struct {
	Token    lexer.Token
	Operator string
	Args     []Expression
}

func (m *MathProcedure) TokenLiteral() string {
	return m.Token.Literal
}
func (m *MathProcedure) String() string {
	var out bytes.Buffer
	args := make([]string, 0)
	for _, arg := range m.Args {
		args = append(args, arg.String())
	}
	out.WriteString("(")
	out.WriteString(m.Operator + " ")
	out.WriteString(strings.Join(args, " "))
	out.WriteString(")")
	return out.String()
}

// ---

type EqProcedure struct {
	Token lexer.Token
	Left  Expression
	Right Expression
}

func (e *EqProcedure) TokenLiteral() string {
	return e.Token.Literal
}
func (e *EqProcedure) String() string {
	var out bytes.Buffer
	out.WriteString("(eq ")
	out.WriteString(e.Left.String() + " ")
	out.WriteString(e.Right.String())
	out.WriteString(")")
	return out.String()
}

// ---

type ComparisonProcedure struct {
	Token    lexer.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (c *ComparisonProcedure) TokenLiteral() string {
	return c.Token.Literal
}
func (c *ComparisonProcedure) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(c.Operator + " ")
	out.WriteString(c.Left.String() + " ")
	out.WriteString(c.Right.String())
	out.WriteString(")")
	return out.String()
}
