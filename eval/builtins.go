package eval

import (
	"fmt"
	"math/rand"
	"strings"

	"yy/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args for len (got %d, want 1)", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Hashmap:
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
			return object.NULL
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
			return object.NULL
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of args for push (got %d, want 2)", len(args))
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

	"swap": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("wrong number of args for swap (got %d, want 3)", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to swap must be ARRAY, got %s", args[0].Type())
			}
			if args[1].Type() != object.INTEGER_OBJ {
				return newError("second argument to swap must be INTEGER, got %s", args[1].Type())
			}
			if args[2].Type() != object.INTEGER_OBJ {
				return newError("third argument to swap must be INTEGER, got %s", args[2].Type())
			}

			arr := args[0].(*object.Array)
			i := args[1].(*object.Integer).Value
			j := args[2].(*object.Integer).Value
			length := len(arr.Elements)

			if i < 0 || i >= int64(length) || j < 0 || j >= int64(length) {
				return arr
			}

			newElements := make([]object.Object, length)
			copy(newElements, arr.Elements)
			newElements[i], newElements[j] = newElements[j], newElements[i]

			return &object.Array{Elements: newElements}
		},
	},

	"yoink": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 && len(args) != 2 {
				return newError("wrong number of args for yoink (got %d, want 1 or 2)", len(args))
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
					return object.NULL
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
					return object.NULL
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

	"yahtzee": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args for yahtzee (got %d, want 1)", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				if arg.Value <= 0 {
					return newError("negative integer not supported by yahtzee")
				}
				return &object.Integer{Value: rand.Int63n(arg.Value)}

			case *object.Array:
				max := len(arg.Elements) - 1
				i := rand.Intn(max)
				return arg.Elements[i]

			case *object.String:
				max := len(arg.Value) - 1
				i := rand.Intn(max)
				return &object.String{Value: string(arg.Value[i])}

			case *object.Range:
				min, max := arg.Start, arg.End
				if min > max {
					min, max = max, min
				}
				v := min + rand.Int63n(max-min)
				return &object.Integer{Value: v}

			default:
				return newError("argument passed to yahtzee not supported, got %s", args[0].Type())
			}
		},
	},

	// PRINT

	"yowl": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(strings.ToUpper(arg.String()))
			}
			fmt.Println()
			return object.NULL
		},
	},

	"yap": {
		Fn: func(args ...object.Object) object.Object {
			msg := spaceSeparatedArgs(args...)
			fmt.Println(msg)
			return object.NULL
		},
	},

	// converts an object to string
	"yarn": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args for yarn (got %d, want 1)", len(args))
			}
			return &object.String{Value: args[0].String()}
		},
	},

	// ERROR THROWING

	// throws an error, effectively terminating the program
	"yikes": {
		Fn: func(args ...object.Object) object.Object {
			msg := spaceSeparatedArgs(args...)
			if msg == "" {
				msg = "yikes!"
			}
			return newError(msg)
		},
	},

	"yassert": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 && len(args) != 2 {
				return newError("wrong number of args for yassert (got %d, want 1 or 2)", len(args))
			}

			if isTruthy(args[0]) {
				return object.NULL // all good, nothing to see here
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
		return nil, fmt.Errorf("wrong number of args for `%s` (got %d, want 1)", fnName, len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return nil, fmt.Errorf("argument to `%s` must be ARRAY, got %s", fnName, args[0].Type())
	}

	return args[0].(*object.Array), nil
}

func spaceSeparatedArgs(args ...object.Object) string {
	// need to convert []Object to []any for Sprint to work
	s := make([]any, len(args))
	for i, v := range args {
		s[i] = v
	}
	return fmt.Sprint(s...)
}
