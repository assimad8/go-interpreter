package lexer

import "github.com/assimad8/go-interpreter/internal/token"

type Lexer struct {
	input        string
	position     int  //current position in input
	readPosition int  // current position in input after current char
	ch           byte //current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.ch = 0
	} else {
		lex.ch = lex.input[lex.readPosition]
	}
	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) NextToken() token.Token {
	for lex.isSkipped(){
		lex.readChar()
	}
	var tk token.Token

	switch lex.ch {
	case '=':
		if (lex.peekChar() == '='){
			ch := lex.ch
			lex.readChar()
			tk = token.Token{Type: token.EQ, Literal: string(ch) + string(lex.ch)}
			}else{
				tk = newToken(token.ASSIGN,lex.ch)
			}
		case '+':
			if(lex.peekChar() == '+'){
			ch := lex.ch
			lex.readChar()
			tk = token.Token{Type: token.PLUS_PLUS, Literal: string(ch) + string(lex.ch)}
		}else{
			tk = newToken(token.PLUS,lex.ch)
		}
	case '-':
		if(lex.peekChar() == '-'){
			ch := lex.ch
			lex.readChar()
			tk = token.Token{Type: token.MINUS_MINUS, Literal: string(ch) + string(lex.ch)}
		}else{
			tk = newToken(token.MINUS,lex.ch)
		}
	case '*':
		tk = newToken(token.ASTERISK,lex.ch)
	case '/':
		tk = newToken(token.SLASH,lex.ch)
	case '>':
		if(lex.peekChar()=='='){
			ch := lex.ch
			lex.readChar()
			tk = token.Token{Type: token.GT_EQ,Literal: string(ch)+string(lex.ch)}
		}else{
			tk = newToken(token.GT,lex.ch)
		}
	case '<':
		if(lex.peekChar()=='='){
			ch := lex.ch
			lex.readChar()
			tk = token.Token{Type:token.LT_EQ,Literal: string(ch)+string(lex.ch)}
		}else{
			tk = newToken(token.LT,lex.ch)
		}
	case '!':
		if(lex.peekChar()=='='){
			ch := lex.ch
			lex.readChar()
			tk = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(lex.ch)}
		}else{
			tk = newToken(token.BANG,lex.ch)
		}
	case ';':
		tk = newToken(token.SEMICOLON,lex.ch)
	case ',':
		tk = newToken(token.COMMA,lex.ch)
	case '(':
		tk = newToken(token.LPAREN,lex.ch)
	case ')':
		tk = newToken(token.RPAREN,lex.ch)
	case '{':
		tk = newToken(token.LBRACE,lex.ch)
	case '}':
		tk = newToken(token.RBRACE,lex.ch)
	case 0:
		tk.Literal = ""
		tk.Type = token.EOF
	case '"':
		tk.Type 	= token.STRING
		tk.Literal	= lex.readString()
	case '\'':
		tk.Type 	= token.STRING
		tk.Literal	= lex.readString()
	default:
		if isLetter(lex.ch) {
			tk.Literal = lex.readIdentifier()
			tk.Type = token.LookupIden(tk.Literal)
			return tk
		}else if isNumber(lex.ch){
			tk.Type = token.INT
			tk.Literal = lex.readNumber()
			return tk
		}else{
			tk = newToken(token.ILLEGAL,lex.ch)
		}
	}
	lex.readChar()
	return tk
}

func newToken(tokenType token.TokenType,ch byte) token.Token {
	return token.Token{Type:tokenType,Literal:string(ch)}
}

func (lex *Lexer) readString() string {
	position := lex.position+1
	char := lex.ch
	for {
		lex.readChar()
		if lex.ch == char || lex.ch == 0 {
			break
		}
	}
	return lex.input[position:lex.position]
}

func (lex *Lexer) readIdentifier() string {
	position := lex.position
	for isLetter(lex.ch){
		lex.readChar()
	}
	return lex.input[position:lex.position]
}
func (lex *Lexer) readNumber() string {
	position := lex.position
	for isNumber(lex.ch){
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

func (lex *Lexer) peekChar() byte {
	if(lex.readPosition>=len(lex.input)){
		return 0
	}
	return lex.input[lex.readPosition]
}

func (lex *Lexer) isSkipped() bool {
	if(lex.ch==' ' || lex.ch=='\t' || lex.ch=='\n'||lex.ch=='\r'){
		return true
	}
	return false
}

func isLetter(ch byte) bool {
	return 'a'<= ch && ch <= 'z' || 'A'<= ch && ch <= 'Z' || ch == '_'
}
func isNumber(ch byte) bool {
	return '0'<= ch && ch <= '9'
}