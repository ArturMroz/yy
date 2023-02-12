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
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NOT_EQ    = "!="
	WALRUS    = ":="
	RANGE     = ".."
	BACKSLASH = "\\"

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

	// Keywords
	TRUE  = "TRUE"
	FALSE = "FALSE"
	NULL  = "NULL"
	YO    = "YO"
	YIF   = "YIF"
	YELS  = "YELS"
	YEET  = "YEET"
	YOYO  = "YOYO"
	YOLO  = "YOLO"
	YALL  = "YALL"
	YET   = "YET"
)

var keywords = map[string]TokenType{
	"yo":    YO,
	"true":  TRUE,
	"false": FALSE,
	"yif":   YIF,
	"yels":  YELS,
	"yeet":  YEET,
	"yoyo":  YOYO,
	"yolo":  YOLO,
	"yet":   YET,
	"yall":  YALL,
	"null":  NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
