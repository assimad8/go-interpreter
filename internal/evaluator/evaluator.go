package evaluator

import (
	"fmt"
	"strings"

	"github.com/assimad8/go-interpreter/internal/ast"
	"github.com/assimad8/go-interpreter/internal/object"
)

var (

	NULL  = &object.NULL{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)


func Eval(node ast.Node,env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements,env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression,env)
	//Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value:node.Value}
	case *ast.PrefixExpression:
		right := Eval(node.Right,env)
		if isError(right){
			return right
		}
		return evalPrefixExpression(node.Operator,right)
	case *ast.InfixExpression:
		left := Eval(node.Left,env)
		if isError(left){
			return left
		}
		right := Eval(node.Right,env)
		if isError(right){
			return right
		}
		return evalInfixExpression(node.Operator,left,right)
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements,env)
	case *ast.IfExpression:
		return evalIfExpression(node,env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue,env)
		if isError(val){
			return val
		}
		return &object.ReturnVALUE{Value:val}
	case *ast.LetStatement:
		val := Eval(node.Value,env)
		if isError(val){
			return val
		}
		env.Set(node.Name.Value,val)
	case *ast.Identifier:
		return evalIdentifier(node,env)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		return &object.Function{Parameters: parameters,Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function,env)
		if isError(function){
			return function
		}
		args := evalExpressions(node.Arguments,env)
		if len(args) == 1 && isError(args[0]){
			return args[0]
		}
		return applyFunction(function,args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements,env)
		if len(elements)==1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left,env)
		if isError(left){
			return left
		}
		index := Eval(node.Index,env)
		if isError(index){
			return index
		}
		return evalIndexExpression(left,index)
	case *ast.HashLiteral:
		return evalHashLiteral(node,env)

	}
	return nil
}

func evalHashLiteral(node *ast.HashLiteral,env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode,valueNode := range node.Pairs{
		key := Eval(keyNode,env)
		if isError(key){
			return key
		}
		hashKey,ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s",key.Type())
		}
		value := Eval(valueNode,env)
		if isError(value){
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed]=object.HashPair{Key: key,Value:value}
	}
	return &object.Hash{Pairs: pairs}
}
func evalIndexExpression(left ,index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type()==object.INTEGER_OBJ :
		return evalArrayIndexExpression(left,index)
	case left.Type()==object.HASH_OBJ:
		return evalHashIndexExpression(left,index)
	default:
		return newError("index operator not supported: %s",left.Type())
	}
}
func evalHashIndexExpression(hash,index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key ,ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s",index.Type())
	}

	pair,ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}
func evalArrayIndexExpression(array,index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements)-1)

	if idx>max || idx<0 {
		return NULL
	}
	return arrayObj.Elements[idx]
}

func applyFunction(fn object.Object,args []object.Object) object.Object {
	
	switch fn:= fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn,args)
		evaluated := Eval(fn.Body,extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s",fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function,args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx ,param := range fn.Parameters {
		env.Set(param.Value,args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue,ok := obj.(*object.ReturnVALUE); ok {
		return returnValue
	}
	return obj
}

func evalExpressions(exps []ast.Expression,env *object.Environment) []object.Object {

	var result []object.Object

	for _,e := range exps {
		evaluated := Eval(e,env)
		if isError(evaluated){
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}
func evalIdentifier(node *ast.Identifier ,env *object.Environment) object.Object{
	if val,ok := env.Get(node.Value);ok {
		return val
	}
	if builtin,ok := builtins[node.Value];ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}
func evalIfExpression(ie *ast.IfExpression,env *object.Environment) object.Object{
	condition := Eval(ie.Condition,env)

	if isError(condition){
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence,env)
	}else if ie.Alternative != nil {
		return Eval(ie.Alternative,env)
	}else {
		return NULL
	}
}
func isTruthy(obj object.Object) bool {
	if obj == NULL || obj == FALSE {
		return false
	}
	return true
}

func evalBlockStatement(stmts []ast.Statement,env *object.Environment) object.Object {
	var result object.Object

	for _,statement := range stmts {
		result = Eval(statement,env)

		if result != nil{ 
			rt := result.Type() 
			
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ{
				return result
			}
		}
	}
	return result
}

func evalProgram(stmts []ast.Statement,env *object.Environment) object.Object {
	var result object.Object

	for _,statement := range stmts {
		result = Eval(statement,env) 

		switch result := result.(type) {
		case *object.ReturnVALUE :
			return result.Value
		case *object.Error :
			return result
		}
	}

	return result
}

func evalInfixExpression(op string,left object.Object,right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ :
		return evalIntegerInfixExpression(op,left,right)
	
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ || left.Type() == object.STRING_OBJ &&right.Type() == object.INTEGER_OBJ ||right.Type() == object.STRING_OBJ &&left.Type() == object.INTEGER_OBJ:
		return evalStringInfixExpression(op,left,right)

	case op == "==":
		return nativeBoolToBooleanObject(left == right)

	case op == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",left.Type(),op,right.Type())
	default:
		return newError("unknown operator: %s %s %s",left.Type(),op,right.Type())
	}
}
func evalStringInfixExpression(op string,left object.Object,right object.Object) object.Object {
	switch op {
	case "+":
		leftVal ,ok:= left.(*object.String)
		if ok {
			leftVal := leftVal.Value
			rightVal := right.(*object.Integer).Value
			return &object.String{Value:fmt.Sprintf("%s%d",leftVal,rightVal)}
		}else {
				leftVal := left.(*object.Integer).Value
				rightVal := right.(*object.String).Value
				return &object.String{Value:fmt.Sprintf("%d%s",leftVal,rightVal)}
			}
	case "*":
		leftVal ,ok:= left.(*object.String)
		if ok {
			leftVal := leftVal.Value
			rightVal := right.(*object.Integer).Value
			return &object.String{Value:strings.Repeat(leftVal,int(rightVal))}
			}else {
				leftVal := left.(*object.Integer).Value
				rightVal := right.(*object.String).Value
				return &object.String{Value:strings.Repeat(rightVal,int(leftVal))}
		}
	

	default:
		return newError("unknown operator: %s %s %s",left.Type(),op,right.Type())
	}
}

func evalIntegerInfixExpression(op string,left object.Object,right object.Object) object.Object {
	leftVal   := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	
	switch op {
	case "+":
		return &object.Integer{Value:leftVal+rightVal}
	case "-":
		return &object.Integer{Value:leftVal-rightVal}
	case "*":
		return &object.Integer{Value:leftVal*rightVal}
	case "/":
		return &object.Integer{Value:leftVal/rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal>rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal<rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal==rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal!=rightVal)
	default:
		return newError("unknown operator: %s %s %s",left.Type(),op,right.Type())
		
	}
}

func evalPrefixExpression(op string,right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default :
		return newError("unknown operator: %s%s",op,right.Type())
	}
}


func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default :
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {

	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s",right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value:-value}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format,a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}