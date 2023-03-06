package token

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers + literals
	IDENT
	INT
	NUMBER
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	PERCENT
	LT
	GT
	EQ
	NOT_EQ
	WALRUS
	RANGE
	BACKSLASH
	MACRO
	ADD_ASSIGN
	SUB_ASSIGN
	MUL_ASSIGN
	DIV_ASSIGN
	MOD_ASSIGN
	HASHMAP

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
	NUMBER: "NUMBER",
	STRING: "STRING",

	// Operators
	ASSIGN:     "=",
	PLUS:       "+",
	MINUS:      "-",
	BANG:       "!",
	ASTERISK:   "*",
	SLASH:      "/",
	PERCENT:    "%",
	LT:         "<",
	GT:         ">",
	EQ:         "==",
	NOT_EQ:     "!=",
	WALRUS:     ":=",
	RANGE:      "..",
	BACKSLASH:  `\`,
	MACRO:      `@\`,
	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	DIV_ASSIGN: "/=",
	MOD_ASSIGN: "%=",
	HASHMAP:    "%{",

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
	YOLO:  "YOLO",
	YALL:  "YALL",
	YET:   "YET",
}

func (tok TokenType) String() string {
	return tokens[tok]
}

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
	"yif":   YIF,
	"yels":  YELS,
	"yeet":  YEET,
	"yolo":  YOLO,
	"yet":   YET,
	"yall":  YALL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
