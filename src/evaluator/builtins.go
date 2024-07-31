package evaluator

import (
	"fmt"

	"github.com/tjapit/monkey/src/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(
					"wrong number of arguments. want=%d, got =%d",
					1,
					len(args),
				)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got =%s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(
					"wrong number of arguments. want=%d, got =%d",
					1,
					len(args),
				)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError(
					"argument to `first` must be ARRAY, got =%s",
					args[0].Type(),
				)
			}

			arr := args[0].(*object.Array).Elements
			if len(arr) > 0 {
				return arr[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(
					"wrong number of arguments. want=%d, got =%d",
					1,
					len(args),
				)
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(
					"argument to `last` must be ARRAY, got =%s",
					args[0].Type(),
				)
			}

			arr := args[0].(*object.Array).Elements
			if len(arr) > 0 {
				return arr[len(arr)-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(
					"wrong number of arguments. want=%d, got =%d",
					1,
					len(args),
				)
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(
					"argument to `rest` must be ARRAY, got =%s",
					args[0].Type(),
				)
			}

			arr := args[0].(*object.Array).Elements
			if len(arr) > 0 {
				return &object.Array{Elements: arr[1:]}
			}
			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(
					"wrong number of arguments. want=%d, got =%d",
					2,
					len(args),
				)
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(
					"first argument to `push` must be ARRAY, got =%s",
					args[0].Type(),
				)
			}

			arr := args[0].(*object.Array)
			arr.Elements = append(arr.Elements, args[1])

			return arr
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
