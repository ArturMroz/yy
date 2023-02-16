package eval

import (
	"fmt"
	"strings"

	"yy/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}

			case *object.Range:
				length := arg.End - arg.Start
				if length < 0 {
					length = -length
				}
				return &object.Integer{Value: length + 1}

			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	// ARRAYS

	"last": {
		Fn: func(args ...object.Object) object.Object {
			arr, err := checkArray("last", args...)
			if err != nil {
				return newError(err.Error())
			}

			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},

	"rest": {
		Fn: func(args ...object.Object) object.Object {
			arr, err := checkArray("rest", args...)
			if err != nil {
				return newError(err.Error())
			}

			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)

			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},

	"yoink": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 && len(args) != 2 {
				return newError("wrong number of arguments for `yoink`. got=%d, want=1 or 2", len(args))
			}
			if len(args) == 2 && args[1].Type() != object.INTEGER_OBJ {
				return newError("second argument to `yoink` must be INTEGER, got %s", args[1].Type())
			}

			switch arg := args[0].(type) {
			case *object.Array:
				pos := len(arg.Elements) - 1
				if len(args) == 2 {
					pos = int(args[1].(*object.Integer).Value)
				}
				if pos >= len(arg.Elements) {
					return NULL
				}

				elt := arg.Elements[pos]
				arg.Elements = append(arg.Elements[:pos], arg.Elements[pos+1:]...)
				return elt

			case *object.String:
				pos := len(arg.Value) - 1
				if len(args) == 2 {
					pos = int(args[1].(*object.Integer).Value)
				}
				if pos >= len(arg.Value) {
					return NULL
				}

				elt := arg.Value[pos]
				arg.Value = arg.Value[:pos] + arg.Value[pos+1:]
				return &object.String{Value: string(elt)}

			default:
				// TODO support more types
				return newError("cannot yoink from %s", args[0].Type())

			}
		},
	},

	// PRINT

	"yell": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(strings.ToUpper(arg.Inspect()))
			}
			fmt.Println()
			return NULL
		},
	},

	"yelp": {
		Fn: func(args ...object.Object) object.Object {
			s := make([]any, len(args))
			for i, v := range args {
				s[i] = v.Inspect()
			}
			fmt.Println(s...)
			return NULL
		},
	},

	// DEBUG

	"yassert": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 && len(args) != 2 {
				return newError("wrong number of arguments for `yassert`. got=%d, want=1 or 2", len(args))
			}

			if isTruthy(args[0]) {
				return NULL // all good, nothing to see here
			}

			msg := "yassert failed"
			if len(args) == 2 {
				if v, ok := args[1].(*object.String); ok {
					msg += ": " + v.Value
				}
			}
			return newError(msg)
		},
	},
}

func checkArray(fnName string, args ...object.Object) (*object.Array, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to `%s`. got=%d, want=1", fnName, len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return nil, fmt.Errorf("argument to `%s` must be ARRAY, got %s", fnName, args[0].Type())
	}

	return args[0].(*object.Array), nil
}
