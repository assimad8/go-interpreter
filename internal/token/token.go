package token

type TokenType string

type Token struct {
	Type TokenType 
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF		= "OEF"

	//identifiers + literals
	IDENT 	= "IDENT"
	INT		= "INT"

	//Operators
	ASSIGN 		= "="
	PLUS 		= "+"
	MINUS 		= "-"
	BANG 		= "!"
	ASTERISK 	= "*"
	SLASH 		= "/"

	GT	= ">"
	GT_EQ	= ">="
	LT	= "<"
	LT_EQ	= "<="

	EQ		= "=="
	NOT_EQ	= "!="

	PLUS_PLUS = "++"
	MINUS_MINUS = "--"

	//Delimeters
	COMMA		= ","
	SEMICOLON	= ";"

	LPAREN	= "("
	RPAREN	= ")"
	LBRACE	= "{"
	RBRACE	= "}"

	//Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
)

var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
	"if": IF,
	"true": TRUE,
	"false": FALSE,
	"else": ELSE,
	"return": RETURN,
}

func LookupIden(ident string) TokenType {
	if tokentype,ok := keywords[ident];ok {
		return tokentype
	}
	return IDENT
}