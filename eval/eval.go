package eval

import (
	"fmt"
	"reflect"

	"ylang/ast"
	"ylang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)

	case *ast.YeetStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		var args []object.Object
		for _, a := range node.Arguments {
			evaluated := Eval(a, env)
			if isError(evaluated) {
				return evaluated
			}
			args = append(args, evaluated)
		}

		return applyFunction(fn, args)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *ast.AssignExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if node.IsInit {
			env.Set(node.Name.Value, val)
			return val
		}
		if ok := env.Update(node.Name.Value, val); ok {
			return val
		}
		return newError("identifier not found: " + node.Name.Value)

	case *ast.Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}
		if builtin, ok := builtins[node.Value]; ok {
			return builtin
		}
		return newError("identifier not found: " + node.Value)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(node.Consequence, env)
		} else if node.Alternative != nil {
			return Eval(node.Alternative, env)
		} else {
			return NULL
		}

	case *ast.YoyoExpression:
		extendedEnv := object.NewEnclosedEnvironment(env)
		if node.Initialiser != nil {
			init := Eval(node.Initialiser, extendedEnv)
			if isError(init) {
				return init
			}
		}

		var result object.Object
		for {
			condition := Eval(node.Condition, extendedEnv)
			if isError(condition) {
				return condition
			}

			if !isTruthy(condition) {
				return result
			}

			result = Eval(node.Body, extendedEnv)

			if node.Post != nil {
				post := Eval(node.Post, extendedEnv)
				if isError(post) {
					return post
				}
			}
		}

	case *ast.YallExpression:
		var result object.Object
		iter := Eval(node.Iterable, env)
		extendedEnv := object.NewEnclosedEnvironment(env)

		switch iter := iter.(type) {
		case *object.Array:
			for _, v := range iter.Elements {
				extendedEnv.Set(node.KeyName, v)
				result = Eval(node.Body, extendedEnv)
			}

		case *object.String:
			for _, v := range iter.Value {
				extendedEnv.Set(node.KeyName, &object.String{Value: string(v)})
				result = Eval(node.Body, extendedEnv)
			}

		case *object.Range:
			if iter.Start <= iter.End {
				for i := iter.Start; i <= iter.End; i++ {
					extendedEnv.Set(node.KeyName, &object.Integer{Value: i})
					result = Eval(node.Body, extendedEnv)
				}
			} else {
				for i := iter.Start; i >= iter.End; i-- {
					extendedEnv.Set(node.KeyName, &object.Integer{Value: i})
					result = Eval(node.Body, extendedEnv)
				}
			}

		default:
			return newError("cannot iterate over %s, type of %T", iter, iter)
		}

		return result

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return toYeetBool(node.Value)

	case *ast.RangeLiteral:
		start := Eval(node.Start, env)
		if isError(start) {
			return start
		}
		end := Eval(node.End, env)
		if isError(end) {
			return end
		}
		if start.Type() != object.INTEGER_OBJ || end.Type() != object.INTEGER_OBJ {
			return newError("only integers can be used to create a range")
		}

		return &object.Range{
			Start: start.(*object.Integer).Value,
			End:   end.(*object.Integer).Value,
		}

	case *ast.ArrayLiteral:
		var elts []object.Object
		for _, elt := range node.Elements {
			evaluated := Eval(elt, env)
			if isError(evaluated) {
				return evaluated
			}
			elts = append(elts, evaluated)
		}
		return &object.Array{Elements: elts}

	case *ast.HashLiteral:
		hash := &object.Hash{Pairs: map[object.HashKey]object.HashPair{}}
		for k, v := range node.Pairs {
			key := Eval(k, env)
			if isError(key) {
				return key
			}
			hashKey, ok := key.(object.Hashable)
			if !ok {
				return newError("key not hashable: %s", key.Type())
			}

			val := Eval(v, env)
			if isError(val) {
				return val
			}

			pair := object.HashPair{Key: key, Value: val}
			hash.Pairs[hashKey.HashKey()] = pair
		}
		return hash

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		idx := Eval(node.Index, env)
		if isError(idx) {
			return idx
		}

		switch {
		case left.Type() == object.ARRAY_OBJ && idx.Type() == object.INTEGER_OBJ:
			i := idx.(*object.Integer).Value
			a := left.(*object.Array)
			if i < 0 || i >= int64(len(a.Elements)) {
				return NULL
			}
			return a.Elements[i]

		case left.Type() == object.HASH_OBJ:
			hashMap := left.(*object.Hash)
			key, ok := idx.(object.Hashable)
			if !ok {
				return newError("key not hashable: %s", idx.Type())
			}

			pair, ok := hashMap.Pairs[key.HashKey()]
			if !ok {
				return NULL
			}
			return pair.Value

		default:
			return newError("index operator not supported: %s", idx.Type())
		}

	case *ast.Null:
		return NULL

	default:
		return newError("ast object not supported %q %T", node, node)
	}
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range statements {
		result = Eval(stmt, env)

		if result != nil {
			rtype := result.Type()
			if rtype == object.RETURN_VALUE_OBJ || rtype == object.ERROR_OBJ {
				// don't unwrap return value and let it bubble so it stops execution in outer block statement
				return result
			}
		}
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)
		for paramIdx, param := range fn.Parameters {
			extendedEnv.Set(param.Value, args[paramIdx])
		}

		evaluated := Eval(fn.Body, extendedEnv)

		// unwrap return value so it doesn't stop eval in outer scope
		if returnValue, ok := evaluated.(*object.ReturnValue); ok {
			return returnValue.Value
		}
		return evaluated

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return toYeetBool(!isTruthy(right))

	case "-":
		if right.Type() != object.INTEGER_OBJ {
			return newError("unknown operator: %s%s", op, right.Type())
		}
		right := right.(*object.Integer)
		return &object.Integer{Value: -right.Value}

	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		right := right.(*object.Integer)
		left := left.(*object.Integer)

		switch operator {
		case "+":
			return &object.Integer{Value: left.Value + right.Value}
		case "-":
			return &object.Integer{Value: left.Value - right.Value}
		case "*":
			return &object.Integer{Value: left.Value * right.Value}
		case "/":
			return &object.Integer{Value: left.Value / right.Value}
		case "<":
			return toYeetBool(left.Value < right.Value)
		case ">":
			return toYeetBool(left.Value > right.Value)
		case "==":
			return toYeetBool(left.Value == right.Value)
		case "!=":
			return toYeetBool(left.Value != right.Value)
		}

	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		right := right.(*object.Boolean)
		left := left.(*object.Boolean)

		switch operator {
		case "==":
			return toYeetBool(left.Value == right.Value)
		case "!=":
			return toYeetBool(left.Value != right.Value)
		}

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		right := right.(*object.String)
		left := left.(*object.String)

		switch operator {
		case "+":
			return &object.String{Value: left.Value + right.Value}
		case "==":
			return toYeetBool(left.Value == right.Value)
		case "!=":
			return toYeetBool(left.Value != right.Value)
		}

	case left.Type() == object.NULL_OBJ && right.Type() == object.NULL_OBJ:
		return toYeetBool(operator == "==")

	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		right := right.(*object.Array)
		left := left.(*object.Array)

		switch operator {
		case "+":
			return &object.Array{Elements: append(left.Elements, right.Elements...)}
		case "==":
			// TODO DeepEqual performance isn't great, replace it
			return toYeetBool(reflect.DeepEqual(left.Elements, right.Elements))
		case "!=":
			return toYeetBool(!reflect.DeepEqual(left.Elements, right.Elements))
		}
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func isTruthy(o object.Object) bool {
	// Ruby's truthiness rules: nil & false are falsy, everything else is truthy
	switch o {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}

func toYeetBool(v bool) object.Object {
	if v {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func newError(format string, args ...any) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, args...)}
}
