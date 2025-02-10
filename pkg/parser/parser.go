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
		lexer.STRINGREF,
		lexer.FIRST,
		lexer.REST,
		lexer.EQ,
		lexer.LT,
		lexer.LTE,
		lexer.GT,
		lexer.GTE:
		return p.parseBuiltinIdentifier()
	case lexer.TICK:
		return p.parseListShorthand()
	case lexer.LPAREN:
		return p.parseForm()
	default:
		p.errors = append(p.errors, fmt.Sprintf("illegal character: %s", p.cur.Literal))
		return nil
	}
}

func (p *Parser) parseBuiltinIdentifier() Expression {
	return &BuiltinIdentifier{
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
	if p.peek.Type == lexer.LIST {
		return p.parseList()
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
