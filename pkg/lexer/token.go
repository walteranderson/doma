package lexer

type Token struct {
	Type    TokenType
	Literal string
}

type TokenType string

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
	LPAREN  = "LPAREN"
	RPAREN  = "RPAREN"
	STRING  = "STRING"
	NUMBER  = "NUMBER"
	IDENT   = "IDENT"
	TRUE    = "TRUE"
	FALSE   = "FALSE"
	SYMBOL  = "SYMBOL"

	PLUS     = "PLUS"
	MINUS    = "MINUS"
	ASTERISK = "ASTERISK"
	SLASH    = "SLASH"
	TICK     = "TICK"
	LT       = "LT"
	LTE      = "LTE"
	GT       = "GT"
	GTE      = "GTE"

	LAMBDA  = "LAMBDA"
	IF      = "IF"
	DEFINE  = "DEFINE"
	DISPLAY = "DISPLAY"
	LIST    = "LIST"
	EQ      = "EQ"

	FIRST     = "FIRST"
	REST      = "REST"
	STRINGREF = "STRINGREF"
)

var keywords = map[string]TokenType{
	"true":       TRUE,
	"false":      FALSE,
	"lambda":     LAMBDA,
	"if":         IF,
	"define":     DEFINE,
	"display":    DISPLAY,
	"list":       LIST,
	"eq":         EQ,
	"first":      FIRST,
	"rest":       REST,
	"string-ref": STRINGREF,
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

var builtins = []TokenType{
	PLUS,
	MINUS,
	ASTERISK,
	SLASH,
	LAMBDA,
	IF,
	DEFINE,
	DISPLAY,
	STRINGREF,
	EQ,
	LT,
	LTE,
	GT,
	GTE,
}

func IsBuiltinToken(token TokenType) bool {
	for _, t := range builtins {
		if t == token {
			return true
		}
	}
	return false
}
