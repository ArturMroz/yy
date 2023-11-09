package token

type Token struct {
	Type    Type
	Literal string
	Offset  int
}

type Type int

const (
	ILLEGAL Type = iota
	EOF
	ERROR

	// Identifiers + literals.

	IDENT
	INT
	NUMBER
	STRING

	// Operators.

	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	PERCENT
	LT
	GT
	DOT
	AT
	AMPERSAND
	PIPE
	BACKSLASH
	ASSIGN
	ADD_ASSIGN
	SUB_ASSIGN
	MUL_ASSIGN
	DIV_ASSIGN
	MOD_ASSIGN
	OR
	AND
	EQ
	NOT_EQ
	LT_EQ
	GT_EQ
	LT_LT
	WALRUS
	RANGE
	MACRO
	HASHMAP

	// Delimiters.

	COMMA
	SEMICOLON
	COLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// Keywords.

	TRUE
	FALSE
	NULL
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
	ERROR:   "ERROR",

	// Identifiers + literals

	IDENT:  "IDENT",
	INT:    "INT",
	NUMBER: "NUMBER",
	STRING: "STRING",

	// Operators

	PLUS:      "+",
	MINUS:     "-",
	BANG:      "!",
	ASTERISK:  "*",
	SLASH:     "/",
	PERCENT:   "%",
	LT:        "<",
	GT:        ">",
	DOT:       ".",
	AT:        "@",
	AMPERSAND: "&",
	PIPE:      "|",
	BACKSLASH: `\`,
	ASSIGN:    "=",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	DIV_ASSIGN: "/=",
	MOD_ASSIGN: "%=",
	OR:         "||",
	AND:        "&&",
	EQ:         "==",
	NOT_EQ:     "!=",
	LT_EQ:      "<=",
	GT_EQ:      ">=",
	LT_LT:      "<<",
	WALRUS:     ":=",
	RANGE:      "..",
	MACRO:      `@\`,
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
	YIF:   "YIF",
	YELS:  "YELS",
	YEET:  "YEET",
	YOLO:  "YOLO",
	YALL:  "YALL",
	YET:   "YET",
}

func (tok Type) String() string {
	return tokens[tok]
}

var keywords = map[string]Type{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
	"yif":   YIF,
	"yels":  YELS,
	"yeet":  YEET,
	"yolo":  YOLO,
	"yoyo":  YOYO,
	"yall":  YALL,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
