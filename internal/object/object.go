package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/assimad8/go-interpreter/internal/ast"
)

type ObjectType string

const (
	STRING_OBJ 		= "STRING"
	INTEGER_OBJ 		= "INTEGER"
	BOOLEAN_OBJ 		= "BOOLEAN"
	RETURN_VALUE_OBJ	= "RETURN"
	ERROR_OBJ			= "ERROR"
	FUNCTION_OBJ		= "FUNCTION"
	NULL_OBJ 			= "NULL"
)

type Object interface {
	Type() 		ObjectType
	Inspect() 	string 
}


//Integer data type
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d",i.Value)
}
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

//Boolean data type
type Boolean struct {
	Value bool
}

func(b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ 
}
func(b *Boolean) Inspect() string {
	return fmt.Sprintf("%t",b.Value)
}

//NULL data type
type NULL struct {}

func(n *NULL) Type() ObjectType {
	return NULL_OBJ
}
func(n *NULL) Inspect() string {
	return "null"
}

//Return statement
type ReturnVALUE struct {
	Value Object
}
func(r *ReturnVALUE) Type() ObjectType {
	return RETURN_VALUE_OBJ
}
func(r *ReturnVALUE) Inspect() string {
	return r.Value.Inspect()
}

// Error object
type Error struct {
	Message string
}
func(e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func(e *Error) Inspect() string {
	return "ERROR: "+e.Message
}



//Function object
type Function struct {
	Parameters []*ast.Identifier
	Body		*ast.BlockStatement
	Env         *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _,p := range f.Parameters {
		params = append(params,p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params,", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

//String datatype
type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}
func (s *String) Inspect() string {
	return s.Value
}



