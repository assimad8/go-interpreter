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
	ASSIGN 	= "="
	PLUS 	= "+"

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
)

var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
}

func LookupIden(ident string) TokenType {
	if tokentype,ok := keywords[ident];ok {
		return tokentype
	}
	return IDENT
}