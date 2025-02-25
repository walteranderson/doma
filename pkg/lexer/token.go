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

	FIRST    = "FIRST"
	REST     = "REST"
	CONS     = "CONS"
	LENGTH   = "LENGTH"
	LIST_REF = "LIST_REF"
)

var keywords = map[string]TokenType{
	"true":     TRUE,
	"false":    FALSE,
	"lambda":   LAMBDA,
	"if":       IF,
	"define":   DEFINE,
	"display":  DISPLAY,
	"list":     LIST,
	"eq":       EQ,
	"first":    FIRST,
	"rest":     REST,
	"length":   LENGTH,
	"cons":     CONS,
	"list-ref": LIST_REF,
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
	EQ,
	LT,
	LTE,
	GT,
	GTE,
	FIRST,
	REST,
	LENGTH,
	CONS,
	LIST_REF,
}

func IsBuiltinToken(token TokenType) bool {
	for _, t := range builtins {
		if t == token {
			return true
		}
	}
	return false
}
