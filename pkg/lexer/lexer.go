package lexer

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	l.skipComments()
	switch l.ch {
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	case '(':
		tok.Type = LPAREN
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = RPAREN
		tok.Literal = string(l.ch)
	case '+':
		tok.Type = PLUS
		tok.Literal = string(l.ch)
	case '/':
		tok.Type = SLASH
		tok.Literal = string(l.ch)
	case '*':
		tok.Type = ASTERISK
		tok.Literal = string(l.ch)
	case '\'':
		if isLetter(l.peekChar()) {
			l.readChar()
			tok.Type = SYMBOL
			tok.Literal = l.readIdent()
			return tok
		} else {
			tok.Type = TICK
			tok.Literal = string(l.ch)
		}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case '=':
		tok.Type = EQ
		tok.Literal = string(l.ch)
	case '-':
		if isDigit(l.peekChar()) {
			ch := l.ch
			l.readChar()
			num := l.readInt()
			tok.Type = NUMBER
			tok.Literal = string(ch) + string(num)
			return tok
		} else {
			tok.Type = MINUS
			tok.Literal = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			tok.Type = LTE
			lit := string(l.ch)
			l.readChar()
			lit += string(l.ch)
			tok.Literal = lit
		} else {
			tok.Type = LT
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			tok.Type = GT
			lit := string(l.ch)
			l.readChar()
			lit += string(l.ch)
			tok.Literal = lit
		} else {
			tok.Type = GTE
			tok.Literal = string(l.ch)
		}
	case '#':
		if l.peekChar() == 't' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok.Type = TRUE
			tok.Literal = literal
		} else if l.peekChar() == 'f' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok.Type = FALSE
			tok.Literal = literal
		} else {
			tok.Type = ILLEGAL
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readInt()
			return tok
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readInt() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	for l.ch == ';' {
		for l.ch != 0 && l.ch != '\n' {
			l.readChar()
		}
		l.skipWhitespace()
	}
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_' || ch == '-' || ch == '/'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
