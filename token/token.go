package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	// TODO add line info for better error reporting
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="
	WALRUS   = ":="
	RANGE    = ".."

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	BACKSLASH = "\\"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	YO       = "YO"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	YEET     = "YEET"
	YOYO     = "YOYO"
	YONI     = "YONI"
	NULL     = "NULL"
)

var keywords = map[string]TokenType{
	"fun":   FUNCTION,
	"let":   LET,
	"yo":    YO,
	"true":  TRUE,
	"false": FALSE,
	"if":    IF,
	"else":  ELSE,
	"yeet":  YEET,
	"yoyo":  YOYO,
	"yoni":  YONI,
	"null":  NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
