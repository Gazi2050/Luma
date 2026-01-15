package token

type TokenType string

const (
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"

	// Separators
	SEMICOLON = ";"
	COMMA     = ","
	COLON     = ":"
	DOT       = "."
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	// Comparison
	EQ  = "=="
	NEQ = "!="
	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	// Keywords
	LET    = "LET"
	CONST  = "CONST"
	FN     = "FN"
	RETURN = "RETURN"
	IF     = "IF"
	ELSE   = "ELSE"
	LOOP   = "LOOP"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	// Built-in placeholder
	LOG = "LOG"
	LEN = "LEN"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
}

func LookupIdent(ident string) TokenType {
	switch ident {
	case "let":
		return LET
	case "const":
		return CONST
	case "fn":
		return FN
	case "return":
		return RETURN
	case "if":
		return IF
	case "else":
		return ELSE
	case "loop":
		return LOOP
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "log":
		return LOG
	case "len":
		return LEN
	default:
		return IDENT
	}
}
