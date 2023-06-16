package eval

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"yy/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for len (got %d, want 1)", len(args))
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
				return newErrorWithoutPos("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	// ARRAYS

	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for `last` (got %d, want 1)", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					return arg.Elements[length-1]
				}
				return object.NULL

			case *object.String:
				length := len(arg.Value)
				if length > 0 {
					return &object.String{Value: string(arg.Value[(length - 1)])}
				}
				return object.NULL

			default:
				return newErrorWithoutPos("argument of type %s not supported by `last` function", arg.Type())
			}
		},
	},

	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for rest (got %d, want 1)", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1)
					copy(newElements, arg.Elements[1:length])
					return &object.Array{Elements: newElements}
				}
				return object.NULL

			case *object.String:
				if len(arg.Value) > 0 {
					return &object.String{Value: arg.Value[1:]}
				}
				return object.NULL

			default:
				return newErrorWithoutPos("argument of type %s not supported by `rest` function", arg.Type())
			}
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorWithoutPos("wrong number of args for push (got %d, want 2)", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newErrorWithoutPos("argument to `push` must be ARRAY, got %s", args[0].Type())
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
				return newErrorWithoutPos("wrong number of args for yoink (got %d, want 1 or 2)", len(args))
			}
			if len(args) == 2 && args[1].Type() != object.INTEGER_OBJ {
				return newErrorWithoutPos("second argument to `yoink` must be INTEGER, got %s", args[1].Type())
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
				return newErrorWithoutPos("cannot yoink from %s", args[0].Type())

			}
		},
	},

	// RANDOM

	"yahtzee": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 && len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for yahtzee (got %d, want 0 or 1)", len(args))
			}

			if len(args) == 0 {
				return &object.Number{Value: rand.Float64()}
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				if arg.Value <= 0 {
					return newErrorWithoutPos("negative integer not supported by yahtzee")
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
				return newErrorWithoutPos("argument passed to yahtzee not supported, got %s", args[0].Type())
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

	"yelp": {
		Fn: func(args ...object.Object) object.Object {
			msg := spaceSeparatedArgs(args...)
			fmt.Print(msg)
			return object.NULL
		},
	},

	// CONVERT

	"yarn": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for yarn (got %d, want 1)", len(args))
			}
			return &object.String{Value: args[0].String()}
		},
	},

	"chr": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for chr (got %d, want 1)", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.String{Value: string(rune(arg.Value))}
			case *object.String:
				if len(arg.Value) == 0 {
					return object.NULL
				}
				return &object.String{Value: string(arg.Value[0])}
			default:
				return newErrorWithoutPos("unsupported argument type for chr, got %s", arg.Type())
			}
		},
	},

	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for int (got %d, want 1)", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.Number:
				return &object.Integer{Value: int64(arg.Value)}
			case *object.Boolean:
				v := 0
				if arg.Value {
					v = 1
				}
				return &object.Integer{Value: int64(v)}
			case *object.String:
				val, err := strconv.ParseInt(arg.Value, 0, 64)
				if err != nil {
					return newErrorWithoutPos("could not parse %s as integer", arg.Value)
				}
				return &object.Integer{Value: val}
			default:
				return newErrorWithoutPos("unsupported argument type for int, got %s", arg.Type())
			}
		},
	},

	"float": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorWithoutPos("wrong number of args for float (got %d, want 1)", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Number:
				return arg
			case *object.Integer:
				return &object.Number{Value: float64(arg.Value)}
			case *object.String:
				val, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return newErrorWithoutPos("could not parse %s as float", arg.Value)
				}
				return &object.Number{Value: val}
			case *object.Boolean:
				v := 0
				if arg.Value {
					v = 1
				}
				return &object.Number{Value: float64(v)}
			default:
				return newErrorWithoutPos("unsupported argument type for float, got %s", arg.Type())
			}
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
			return newErrorWithoutPos(msg)
		},
	},

	"yassert": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 && len(args) != 2 {
				return newErrorWithoutPos("wrong number of args for yassert (got %d, want 1 or 2)", len(args))
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
			return newErrorWithoutPos(msg)
		},
	},

	"yassert_eq": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 && len(args) != 3 {
				return newErrorWithoutPos("wrong number of args for yassert_eq (got %d, want 2 or 3)", len(args))
			}

			if reflect.DeepEqual(args[0], args[1]) {
				return object.NULL // all good, nothing to see here
			}

			msg := fmt.Sprintf("yassert failed: want %s, got %s", args[0], args[1])
			if len(args) == 3 {
				if v, ok := args[2].(*object.String); ok {
					msg += fmt.Sprintf(" (%s)", v.Value)
				}
			}
			return newErrorWithoutPos(msg)
		},
	},
}

func spaceSeparatedArgs(args ...object.Object) string {
	// need to convert []Object to []any for Sprint to work
	s := make([]any, len(args))
	for i, v := range args {
		s[i] = v
	}
	return fmt.Sprint(s...)
}
