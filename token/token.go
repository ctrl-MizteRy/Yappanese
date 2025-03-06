package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	DASH      = "/"
	BANG      = "!"
	POWER     = "**"
	MOD       = "%"
	DECREMENT = "--"
	INCREMENT = "++"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	FLOAT     = "FLOAT"

	LT      = "<"
	GT      = ">"
	TERNARY = "?"
	EQ      = "=="
	NEQ     = "!="

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	ELIF     = "ELIF"
	RETURN   = "RETURN"
	GLOBAL   = "GLOBAL"
	FOR      = "FOR"
)

var keywords = map[string]TokenType{
	"func":        FUNCTION,
	"propose":     LET,
	"true":        TRUE,
	"false":       FALSE,
	"perhaps":     IF,
	"otherwise":   ELSE,
	"perchance":   ELIF,
	"sayless":     RETURN,
	"nocap":       TRUE,
	"cap":         FALSE,
	"ackchyually": CONST,
	"worldwide":   GLOBAL,
	".":           FLOAT,
	"for":         FOR,
}

func LookupIdent(indent string) TokenType {
	if tok, ok := keywords[indent]; ok {
		return tok
	}
	return IDENT
}
