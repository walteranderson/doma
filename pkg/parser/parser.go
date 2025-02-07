package parser

import (
	"doma/pkg/lexer"
	"fmt"
	"strconv"
)

type Parser struct {
	l      *lexer.Lexer
	cur    lexer.Token
	peek   lexer.Token
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{Args: []Expression{}}

	for !p.curTokenIs(lexer.EOF) {
		expr := p.parseExpression()
		if expr != nil {
			program.Args = append(program.Args, expr)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseExpression() Expression {
	switch p.cur.Type {
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.STRING:
		return p.parseString()
	case lexer.IDENT:
		return p.parseIdent()
	case lexer.TRUE, lexer.FALSE:
		return p.parseBoolean()
	case lexer.PLUS,
		lexer.MINUS,
		lexer.ASTERISK,
		lexer.SLASH,
		lexer.LAMBDA,
		lexer.IF,
		lexer.DEFINE,
		lexer.DISPLAY,
		lexer.EQ,
		lexer.LT,
		lexer.LTE,
		lexer.GT,
		lexer.GTE:
		return p.parseProcedureIdent()
	case lexer.TICK:
		return p.parseListShorthand()
	case lexer.LPAREN:
		return p.parseForm()
	default:
		return nil
	}
}

func (p *Parser) parseProcedureIdent() Expression {
	return &ProcedureIdentifier{
		Token: p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseListShorthand() Expression {
	p.nextToken()
	lf := &List{
		Token: p.cur, // (
		Args:  make([]Expression, 0),
	}
	p.nextToken()
	for !p.curTokenIs(lexer.RPAREN) {
		expr := p.parseExpression()
		lf.Args = append(lf.Args, expr)
		p.nextToken()
	}
	return lf
}

func (p *Parser) parseList() Expression {
	lf := &List{
		Token: p.cur, // (
		Args:  make([]Expression, 0),
	}
	p.nextToken()
	p.nextToken()
	for !p.curTokenIs(lexer.RPAREN) {
		expr := p.parseExpression()
		lf.Args = append(lf.Args, expr)
		p.nextToken()
	}
	return lf
}

func (p *Parser) parseBoolean() Expression {
	return &Boolean{
		Token: p.cur, // true, false, #t, #f
		Value: p.curTokenIs(lexer.TRUE),
	}
}

func (p *Parser) parseIdent() Expression {
	return &Identifier{
		Token: p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseForm() Expression {
	switch p.peek.Type {
	case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH:
		return p.parseMathProc()
	case lexer.LT, lexer.LTE, lexer.GT, lexer.GTE:
		return p.parseComparisonProc()
	case lexer.LAMBDA:
		return p.parseLambdaProc()
	case lexer.IF:
		return p.parseIfProc()
	case lexer.DEFINE:
		return p.parseDefineProc()
	case lexer.DISPLAY:
		return p.parseDisplayProc()
	case lexer.LIST:
		return p.parseList()
	case lexer.EQ:
		return p.parseEqProc()
	}

	form := &Form{
		Token: p.cur, // (
	}
	p.nextToken()
	form.First = p.parseExpression()
	p.nextToken()
	form.Rest = make([]Expression, 0)
	for !p.curTokenIs(lexer.RPAREN) {
		expr := p.parseExpression()
		form.Rest = append(form.Rest, expr)
		p.nextToken()
	}
	return form
}

func (p *Parser) parseComparisonProc() Expression {
	cmp := &ComparisonProcedure{
		Token:    p.cur,          // (
		Operator: p.peek.Literal, // < > <= >=
	}
	p.nextToken() // < > <= >=
	p.nextToken()
	cmp.Left = p.parseExpression()
	if cmp.Left == nil {
		p.errors = append(p.errors, "expected 2 arguments for comparison")
		return nil
	}
	p.nextToken()
	cmp.Right = p.parseExpression()
	if cmp.Right == nil {
		p.errors = append(p.errors, "expected 2 arguments for comparison")
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(lexer.RPAREN) {
		p.errors = append(p.errors, "expected 2 arguments for comparison")
		return nil
	}
	return cmp
}

func (p *Parser) parseEqProc() Expression {
	eq := &EqProcedure{
		Token: p.cur, // (
	}
	p.nextToken() // eq
	if p.peekTokenIs(lexer.RPAREN) {
		p.errors = append(p.errors, fmt.Sprintf("eq missing arguments"))
		return nil
	}
	p.nextToken()
	eq.Left = p.parseExpression()
	p.nextToken()
	eq.Right = p.parseExpression()
	p.nextToken()

	return eq
}

func (p *Parser) parseDisplayProc() Expression {
	disp := &DisplayProcedure{
		Token: p.cur, // (
	}
	p.nextToken() // display
	if p.peekTokenIs(lexer.RPAREN) {
		p.errors = append(p.errors, fmt.Sprintf("display missing argument"))
		return nil
	}
	p.nextToken()
	for !p.curTokenIs(lexer.RPAREN) {
		expr := p.parseExpression()
		disp.Args = append(disp.Args, expr)
		p.nextToken()
	}
	return disp
}

func (p *Parser) parseDefineProc() Expression {
	def := &DefineProcedure{
		Token: p.cur, // (
	}
	p.nextToken() // define
	if !p.peekTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("define expects first argument to be an identifier, got %s", p.cur.Type))
		return nil
	}
	p.nextToken()
	def.Name = p.cur.Literal
	p.nextToken()
	def.Value = p.parseExpression()
	return def
}

func (p *Parser) parseIfProc() Expression {
	i := &IfProcedure{
		Token: p.cur, // (
	}
	p.nextToken() // if
	p.nextToken()
	i.Condition = p.parseExpression()
	p.nextToken()
	i.Consequence = p.parseExpression()
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return i
	}
	p.nextToken()
	i.Alternative = p.parseExpression()
	p.nextToken()
	return i
}

func (p *Parser) parseLambdaProc() Expression {
	l := &LambdaProcedure{
		Token: p.cur, // (
	}
	p.nextToken() // lambda
	p.nextToken()
	expr := p.parseExpression()
	lst, ok := expr.(*List)
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("lambda expects first argument to be a list"))
		return nil
	}
	params := make([]*Identifier, 0)
	for _, arg := range lst.Args {
		ident, ok := arg.(*Identifier)
		if !ok {
			p.errors = append(p.errors, fmt.Sprintf("all parameters should be identifiers, got %s", arg.String()))
			return nil
		}
		params = append(params, ident)
	}
	l.Params = params
	p.nextToken()
	for !p.curTokenIs(lexer.RPAREN) {
		expr := p.parseExpression()
		l.Body = append(l.Body, expr)
		p.nextToken()
	}

	return l
}

func (p *Parser) parseMathProc() Expression {
	proc := &MathProcedure{
		Token:    p.cur,          // (
		Operator: p.peek.Literal, // + - / *
		Args:     make([]Expression, 0),
	}
	p.nextToken() // + - / *
	p.nextToken()
	for !p.curTokenIs(lexer.RPAREN) {
		expr := p.parseExpression()
		proc.Args = append(proc.Args, expr)
		p.nextToken()
	}
	return proc
}

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.cur.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peek.Type == t
}

func (p *Parser) parseString() Expression {
	return &String{
		Token: p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseNumber() Expression {
	value, err := strconv.ParseInt(p.cur.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("Could not parse %s as integer", p.cur.Literal))
		return nil
	}
	return &Number{
		Token: p.cur,
		Value: value,
	}
}
