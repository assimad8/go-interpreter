package ast

import (
	"bytes"
	"strings"

	"github.com/assimad8/go-interpreter/internal/token"
)

type Node interface {
	TokenLiteral() string
	String()       string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _,s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
//LET IDEN   EXPRESSION
//let name = "emad"
//let name = 5*5
//let name = add(5,5)
type LetStatement struct {
	Token token.Token // the token.LET token
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral()+ " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}


//identifier like let x = 3; x is the identifier
type Identifier struct {
	Token token.Token //the token.IDENT token
	Value string
}

func (i *Identifier) statementNode() {}
func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {return i.Value}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

//return statement return [expression]
type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral()+ " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

//expressionStatement statement x+4; 5+5; add(3,4)+3;
type ExpressionStatement struct {
	Token token.Token//first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}	
	return ""
}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

//integerLiteral for integer in the expression
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func( il *IntegerLiteral) expressionNode() {}
func( il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func( il *IntegerLiteral) String() string {
	return il.Token.Literal
}

//prefix expression !x, -x...
type PrefixExpression struct {
	Token 		token.Token //! or !
	Operator 	string
	Right 		Expression
}
func( pe *PrefixExpression) expressionNode() {}
func( pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func( pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	
	return out.String()
}

//infix expression x + x, x + 34 , 3 * 4 ....

type InfixExpression struct {
	Token 	 token.Token
	Left	 Expression
	Operator string
	Right	 Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

//Boolean literals
type Boolean struct {
	Token token.Token
	Value bool
}

func(b *Boolean) expressionNode() {}
func(b *Boolean) TokenLiteral() string {return b.Token.Literal}
func(b *Boolean) String() string {return b.Token.Literal}

//If else expression
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}
func(ie *IfExpression) expressionNode() {}
func(ie *IfExpression) TokenLiteral() string {return ie.Token.Literal}
func(ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}
type BlockStatement struct {
	Token token.Token // the { token
	Statements []Statement
}
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

//Function
type FunctionLiteral struct {
	Token 		token.Token // the 'fn' token
	Parameters  []*Identifier
	Body		*BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _,p := range fl.Parameters {
		params = append(params,p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params,", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// call a function in the program
type CallExpression struct {
	Token token.Token // the '(' token
	Function Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func(ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _,a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args,", "))
	out.WriteString(")")

	return out.String()
}

// String datatype

type StringLiteral struct {
	Token token.Token 
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// Array literal
type ArrayLiteral struct {
	Token token.Token // thre '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode(){}
func (al *ArrayLiteral) TokenLiteral() string{
	return al.Token.Literal
}
func (al *ArrayLiteral) String() string{
	var out bytes.Buffer

	elements := []string{}
	for _,e := range al.Elements {
		elements = append(elements,e.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements,", "))
	out.WriteString("]")

	return out.String()
}

// Index expression for array
type IndexExpression struct {
	Token 	token.Token // the [ token
	Left 	Expression // the array or identifier
	Index 	Expression // the index expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

//hash table data type

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}

	for key,value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	
	out.WriteString("{")
	out.WriteString(strings.Join(pairs,", "))
	out.WriteString("}")
	return out.String()
}