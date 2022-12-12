package eval

import "ylang/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	"assert": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 && len(args) != 2 {
				return newError("wrong number of arguments for `assert`. got=%d, want=1 or 2", len(args))
			}

			if isTruthy(args[0]) {
				return NULL // all good, nothing to see here
			}

			msg := "assert failed"
			if len(args) == 2 {
				if v, ok := args[1].(*object.String); ok {
					msg += ": " + v.Value
				}
			}
			return newError(msg)
		},
	},
}
