package evaluator

import (
	"fmt"
	"strings"

	"github.com/assimad8/go-interpreter/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn:func(args ...object.Object) object.Object {
			if len(args)!=1{
				return newError("wrong number of arguments. got=%d, want=1",len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to 'len' not supported. got %s",args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'first' must be ARRAY. got=%s",args[0].Type())
			}

			arg := args[0].(*object.Array)
			if len(arg.Elements) >0 {
				return arg.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'last' must be ARRAY. got=%s",args[0].Type())
			}

			arg := args[0].(*object.Array)
			if len(arg.Elements) >0 {
				return arg.Elements[len(arg.Elements)-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'rest' must be ARRAY. got=%s",args[0].Type())
			}

			arg := args[0].(*object.Array)
			length := len(arg.Elements)
			if length >0 {
				newArray := make([]object.Object,length-1)
				copy(newArray,arg.Elements[:length])
				return &object.Array{Elements: newArray}
			}
			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'push' must be ARRAY. got=%s",args[0].Type())
			}
			arg := args[0].(*object.Array)
			length := len(arg.Elements)
			newArray := make([]object.Object,length+1)
			copy(newArray,arg.Elements)
			newArray[length] = args[1]
			return &object.Array{Elements: newArray}
		},
	},
	"puts": {
		Fn:func(args ...object.Object) object.Object {
			if len(args)<1 {
				return newError("wrong number of arguments. got=%d, want>=1",len(args))
			}
			data := make([]string, 0, len(args))
			for _,arg := range args {
				data = append(data,arg.Inspect())
			}
			fmt.Println(strings.Join(data," "))
			return NULL
		},
	},
}