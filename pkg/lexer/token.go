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

	LAMBDA   = "LAMBDA"
	IF       = "IF"
	DEFINE   = "DEFINE"
	DISPLAY  = "DISPLAY"
	PRINTF   = "PRINTF"
	LIST     = "LIST"
	EQ       = "EQ"
	FIRST    = "FIRST"
	REST     = "REST"
	CONS     = "CONS"
	LENGTH   = "LENGTH"
	LIST_REF = "LIST_REF"
	BEGIN    = "BEGIN"
)

var keywords = map[string]TokenType{
	"true":     TRUE,
	"false":    FALSE,
	"lambda":   LAMBDA,
	"if":       IF,
	"define":   DEFINE,
	"display":  DISPLAY,
	"printf":   PRINTF,
	"list":     LIST,
	"eq":       EQ,
	"first":    FIRST,
	"rest":     REST,
	"length":   LENGTH,
	"cons":     CONS,
	"list-ref": LIST_REF,
	"begin":    BEGIN,
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
	PRINTF,
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
	BEGIN,
}

func IsBuiltinToken(token TokenType) bool {
	for _, t := range builtins {
		if t == token {
			return true
		}
	}
	return false
}
