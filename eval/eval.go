package eval

import (
	"fmt"
	"reflect"

	"yy/ast"
	"yy/object"
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

	case *ast.CallExpression:
		return evalCallExpr(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}

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
		if env.IsYoloMode() {
			env.Set(node.Name.Value, val)
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
		return evalPrefixExpression(node.Operator, right, env.IsYoloMode())

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, env.IsYoloMode())

	case *ast.YifExpression:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(node.Consequence, env)
		} else if node.Alternative != nil {
			return Eval(node.Alternative, env)
		} else {
			return object.NULL
		}

	case *ast.YoloExpression:
		extendedEnv := object.NewEnclosedEnvironment(env)
		env.SetYoloMode()
		return Eval(node.Body, extendedEnv)

	case *ast.YetExpression:
		extendedEnv := object.NewEnclosedEnvironment(env)

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
			if isErrorOrReturn(result) {
				return result
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
				if isErrorOrReturn(result) {
					return result
				}
			}

		case *object.String:
			for _, v := range iter.Value {
				extendedEnv.Set(node.KeyName, &object.String{Value: string(v)})
				result = Eval(node.Body, extendedEnv)
				if isErrorOrReturn(result) {
					return result
				}
			}

		case *object.Range:
			if iter.Start <= iter.End {
				for i := iter.Start; i <= iter.End; i++ {
					extendedEnv.Set(node.KeyName, &object.Integer{Value: i})
					result = Eval(node.Body, extendedEnv)
					if isErrorOrReturn(result) {
						return result
					}
				}
			} else {
				for i := iter.Start; i >= iter.End; i-- {
					extendedEnv.Set(node.KeyName, &object.Integer{Value: i})
					result = Eval(node.Body, extendedEnv)
					if isErrorOrReturn(result) {
						return result
					}
				}
			}

		default:
			return newError("cannot iterate over %s, type of %s", iter, iter.Type())
		}

		return result

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.TemplateStringLiteral:
		vals := []any{}
		for _, v := range node.Values {
			cur := Eval(v, env)
			if isError(cur) {
				return cur
			}
			vals = append(vals, cur)
		}
		value := fmt.Sprintf(node.Template, vals...)

		return &object.String{Value: value}

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
		elts := []object.Object{}
		for _, elt := range node.Elements {
			evaluated := Eval(elt, env)
			if isError(evaluated) {
				return evaluated
			}
			elts = append(elts, evaluated)
		}
		return &object.Array{Elements: elts}

	case *ast.HashmapLiteral:
		hashmap := &object.Hashmap{Pairs: map[object.HashKey]object.HashPair{}}
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
			hashmap.Pairs[hashKey.HashKey()] = pair
		}
		return hashmap

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
			arr := left.(*object.Array)
			if i < 0 || i >= int64(len(arr.Elements)) {
				return object.NULL
			}
			return arr.Elements[i]

		case left.Type() == object.STRING_OBJ && idx.Type() == object.INTEGER_OBJ:
			i := idx.(*object.Integer).Value
			str := left.(*object.String)
			if i < 0 || i >= int64(len(str.Value)) {
				return object.NULL
			}
			return &object.String{Value: string(str.Value[i])}

		case left.Type() == object.HASHMAP_OBJ:
			hashmap := left.(*object.Hashmap)
			key, ok := idx.(object.Hashable)
			if !ok {
				return newError("key not hashable: %s", idx.Type())
			}

			pair, ok := hashmap.Pairs[key.HashKey()]
			if !ok {
				return object.NULL
			}
			return pair.Value

		default:
			return newError("index operator not supported: %s", idx.Type())
		}

	case *ast.Null:
		return object.NULL

	case nil:
		return newError("unexpected error: something most likely went wrong at the parsing stage")

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

func evalCallExpr(callExpr *ast.CallExpression, env *object.Environment) object.Object {
	if callExpr.Function.TokenLiteral() == "quote" { // TODO this is ugly
		return quote(callExpr.Arguments[0], env) // quote only supports 1 arg
	}

	fn := Eval(callExpr.Function, env)
	if isError(fn) {
		return fn
	}

	var args []object.Object
	for _, a := range callExpr.Arguments {
		evaluated := Eval(a, env)
		if isError(evaluated) {
			return evaluated
		}
		args = append(args, evaluated)
	}

	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return newError("wrong number of args for %s (got %d, want %d)",
				callExpr.Function.TokenLiteral(), len(args), len(fn.Parameters))
		}

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

func evalPrefixExpression(op string, right object.Object, yoloOK bool) object.Object {
	switch op {
	case "!":
		return toYeetBool(!isTruthy(right))

	case "-":
		switch {
		case right.Type() == object.INTEGER_OBJ:
			rightVal := right.(*object.Integer).Value
			return &object.Integer{Value: -rightVal}

		case right.Type() == object.NUMBER_OBJ:
			rightVal := right.(*object.Number).Value
			return &object.Number{Value: -rightVal}

		case yoloOK:
			return yoloPrefixExpression(op, right)
		}
	}

	return newError("unknown operator: %s%s", op, right.Type())
}

func evalInfixExpression(op string, left, right object.Object, yoloOK bool) object.Object {
	switch {
	// as a special case, NUMBER & INTEGER types can be mixed together outside of yolo mode
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.INTEGER_OBJ:
		left := left.(*object.Number)
		right := right.(*object.Integer)
		rVal := float64(right.Value)

		switch op {
		case "+":
			return &object.Number{Value: left.Value + rVal}
		case "-":
			return &object.Number{Value: left.Value - rVal}
		case "*":
			return &object.Number{Value: left.Value * rVal}
		case "/":
			return &object.Number{Value: left.Value / rVal}
		case "%":
			return &object.Number{Value: float64(int64(left.Value) % int64(right.Value))}
		case "<":
			return toYeetBool(left.Value < rVal)
		case ">":
			return toYeetBool(left.Value > rVal)
		case "==":
			return toYeetBool(left.Value == rVal)
		case "!=":
			return toYeetBool(left.Value != rVal)
		}

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.NUMBER_OBJ:
		left := left.(*object.Integer)
		right := right.(*object.Number)
		lVal := float64(left.Value)

		switch op {
		case "+":
			return &object.Number{Value: lVal + right.Value}
		case "-":
			return &object.Number{Value: lVal - right.Value}
		case "*":
			return &object.Number{Value: lVal * right.Value}
		case "/":
			return &object.Number{Value: lVal / right.Value}
		case "%":
			return &object.Number{Value: float64(int64(left.Value) % int64(right.Value))}
		case "<":
			return toYeetBool(lVal < right.Value)
		case ">":
			return toYeetBool(lVal > right.Value)
		case "==":
			return toYeetBool(lVal == right.Value)
		case "!=":
			return toYeetBool(lVal != right.Value)
		}

		// mixing of all the other types is allowed only in yolo mode
	case left.Type() != right.Type():
		switch op {
		case "==":
			return toYeetBool(false)
		case "!=":
			return toYeetBool(true)
		}

		if yoloOK {
			return yoloInfixExpression(op, left, right)
		}
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())

	case left.Type() == object.NULL_OBJ && right.Type() == object.NULL_OBJ:
		switch op {
		case "==":
			return toYeetBool(true)
		case "!=":
			return toYeetBool(false)
		}

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		right := right.(*object.Integer)
		left := left.(*object.Integer)

		switch op {
		case "+":
			return &object.Integer{Value: left.Value + right.Value}
		case "-":
			return &object.Integer{Value: left.Value - right.Value}
		case "*":
			return &object.Integer{Value: left.Value * right.Value}
		case "/":
			return &object.Integer{Value: left.Value / right.Value}
		case "%":
			return &object.Integer{Value: left.Value % right.Value}
		case "<":
			return toYeetBool(left.Value < right.Value)
		case ">":
			return toYeetBool(left.Value > right.Value)
		case "==":
			return toYeetBool(left.Value == right.Value)
		case "!=":
			return toYeetBool(left.Value != right.Value)
		}

	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		right := right.(*object.Number)
		left := left.(*object.Number)

		switch op {
		case "+":
			return &object.Number{Value: left.Value + right.Value}
		case "-":
			return &object.Number{Value: left.Value - right.Value}
		case "*":
			return &object.Number{Value: left.Value * right.Value}
		case "/":
			return &object.Number{Value: left.Value / right.Value}
		case "%":
			return &object.Number{Value: float64(int64(left.Value) % int64(right.Value))}
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

		switch op {
		case "==":
			return toYeetBool(left.Value == right.Value)
		case "!=":
			return toYeetBool(left.Value != right.Value)
		}

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		right := right.(*object.String)
		left := left.(*object.String)

		switch op {
		case "+":
			return &object.String{Value: left.Value + right.Value}
		case "==":
			return toYeetBool(left.Value == right.Value)
		case "!=":
			return toYeetBool(left.Value != right.Value)
		}

	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		right := right.(*object.Array)
		left := left.(*object.Array)

		switch op {
		case "+":
			return &object.Array{Elements: append(left.Elements, right.Elements...)}
		case "==":
			return toYeetBool(reflect.DeepEqual(left.Elements, right.Elements))
		case "!=":
			return toYeetBool(!reflect.DeepEqual(left.Elements, right.Elements))
		}
	}

	return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func isTruthy(obj object.Object) bool {
	// Ruby's truthiness rule: nil & false are falsy, everything else is truthy
	switch obj {
	case object.NULL, object.FALSE:
		return false
	default:
		return true
	}
}

func toYeetBool(b bool) object.Object {
	if b {
		return object.TRUE
	}
	return object.FALSE
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func isErrorOrReturn(obj object.Object) bool {
	return obj != nil && (obj.Type() == object.ERROR_OBJ || obj.Type() == object.RETURN_VALUE_OBJ)
}

func newError(format string, args ...any) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, args...)}
}
