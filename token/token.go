package token

type Token struct {
	Type    TokenType
	Literal string
	// TODO add line info for better error reporting
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers + literals
	IDENT
	INT
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	LT
	GT
	EQ
	NOT_EQ
	WALRUS
	RANGE
	BACKSLASH
	ADD_ASSIGN
	SUB_ASSIGN
	MUL_ASSIGN
	DIV_ASSIGN

	// Delimiters
	COMMA
	SEMICOLON
	COLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// Keywords
	TRUE
	FALSE
	NULL
	YO
	YIF
	YELS
	YEET
	YOYO
	YOLO
	YALL
	YET
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	// Identifiers + literals
	IDENT:  "IDENT",
	INT:    "INT",
	STRING: "STRING",

	// Operators
	ASSIGN:     "=",
	PLUS:       "+",
	MINUS:      "-",
	BANG:       "!",
	ASTERISK:   "*",
	SLASH:      "/",
	LT:         "<",
	GT:         ">",
	EQ:         "==",
	NOT_EQ:     "!=",
	WALRUS:     ":=",
	RANGE:      "..",
	BACKSLASH:  "\\",
	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	DIV_ASSIGN: "/=",

	// Delimiters
	COMMA:     ",",
	SEMICOLON: ";",
	COLON:     ":",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",

	// Keywords
	TRUE:  "TRUE",
	FALSE: "FALSE",
	NULL:  "NULL",
	YO:    "YO",
	YIF:   "YIF",
	YELS:  "YELS",
	YEET:  "YEET",
	YOYO:  "YOYO",
	YOLO:  "YOLO",
	YALL:  "YALL",
	YET:   "YET",
}

func (tok TokenType) String() string {
	return tokens[tok]
}

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
