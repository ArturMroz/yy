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

	case *ast.DeclareExpression:
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

		switch node := node.Left.(type) {
		case *ast.Identifier:
			if ok := env.Update(node.Value, val); ok {
				return val
			}
			if env.IsYoloMode() {
				env.Set(node.Value, val)
				return val
			}
			return newError(node.Pos(), "identifier not found: "+node.Value)

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
					return newError(
						node.Left.Pos(),
						"assign out of bounds for array %s",
						node.Left.(*ast.Identifier).Value)
				}
				arr.Elements[i] = val
				return val

			case left.Type() == object.STRING_OBJ && idx.Type() == object.INTEGER_OBJ:
				i := idx.(*object.Integer).Value
				str := left.(*object.String)
				if i < 0 || i >= int64(len(str.Value)) {
					return newError(
						node.Left.Pos(),
						"assign out of bounds for string %s",
						node.Left.(*ast.Identifier).Value)
				}
				str.Value = str.Value[:i] + val.String() + str.Value[i+1:]
				return val

			case left.Type() == object.HASHMAP_OBJ:
				hashmap := left.(*object.Hashmap)
				key, ok := idx.(object.Hashable)
				if !ok {
					return newError(node.Index.Pos(), "key not hashable: %s", idx.Type())
				}

				hashmap.Pairs[key.HashKey()] = object.HashPair{Key: idx, Value: val}
				return val

			default:
				return newError(node.Index.Pos(), "index operator not supported: %s, type of %s", idx.String(), idx.Type())
			}
		}

		return newError(node.Left.Pos(), "identifier not found: "+node.Left.String())

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		idx := Eval(node.Index, env)
		if isError(idx) {
			return idx
		}

		switch left := left.(type) {
		case *object.Array:
			switch idx := idx.(type) {
			case *object.Integer:
				i := idx.Value
				if i < 0 || i >= int64(len(left.Elements)) {
					return object.NULL
				}
				return left.Elements[i]

			case *object.Range:
				start := idx.Start
				end := idx.End
				if start < 0 {
					start = 0
				}
				if end > int64(len(left.Elements)) {
					end = int64(len(left.Elements))
				}

				// copy the array so modyfing a value in the original array doesn't affect copied array
				return &object.Array{Elements: append([]object.Object{}, left.Elements[start:end]...)}
			}

		case *object.String:
			switch idx := idx.(type) {
			case *object.Integer:
				i := idx.Value
				if i < 0 || i >= int64(len(left.Value)) {
					return object.NULL
				}
				return &object.String{Value: string(left.Value[i])}

			case *object.Range:
				start := idx.Start
				end := idx.End
				if start < 0 {
					start = 0
				}
				if end > int64(len(left.Value)) {
					end = int64(len(left.Value))
				}
				return &object.String{Value: left.Value[start:end]}
			}

		case *object.Hashmap:
			key, ok := idx.(object.Hashable)
			if !ok {
				return newError(node.Index.Pos(), "key not hashable: %s", idx.Type())
			}

			pair, ok := left.Pairs[key.HashKey()]
			if !ok {
				return object.NULL
			}
			return pair.Value

		}
		return newError(node.Index.Pos(), "index operator not supported: %s", idx.Type())

	case *ast.Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}
		if builtin, ok := builtins[node.Value]; ok {
			return builtin
		}
		return newError(node.Pos(), "identifier not found: "+node.Value)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		result := evalPrefixExpression(node.Operator, right, env.IsYoloMode())
		if errObj, ok := result.(*object.Error); ok {
			errObj.Pos = node.Pos()
		}
		return result

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		result := evalInfixExpression(node.Operator, left, right, env.IsYoloMode())
		if errObj, ok := result.(*object.Error); ok {
			errObj.Pos = node.Pos()
		}
		return result

	case *ast.YoloExpression:
		extendedEnv := object.NewEnclosedEnvironment(env)
		extendedEnv.SetYoloMode()
		return Eval(node.Body, extendedEnv)

	// CONTROL FLOW

	case *ast.AndExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		if !isTruthy(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return right

	case *ast.OrExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		if isTruthy(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return right

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

	case *ast.YoyoExpression:
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
			return newError(node.Iterable.Pos(), "cannot iterate over %s, type of %s", iter, iter.Type())
		}

		return result

	// LITERALS

	case *ast.NullLiteral:
		return object.NULL

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

	case *ast.BooleanLiteral:
		return toYeetBool(node.Value)

	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}

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
			return newError(node.Start.Pos(), "only integers can be used to create a range (got %s..%s)", start.Type(), end.Type())
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
				return newError(k.Pos(), "key not hashable: %s", key.Type())
			}

			val := Eval(v, env)
			if isError(val) {
				return val
			}

			pair := object.HashPair{Key: key, Value: val}
			hashmap.Pairs[hashKey.HashKey()] = pair
		}
		return hashmap

	case nil:
		return newErrorWithoutPos("unexpected error: something went wrong somewhere (that's all we know).")

	default:
		return newError(node.Pos(), "ast object not supported %q %T", node, node)
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
			return newError(
				callExpr.Pos(),
				"wrong number of args for %s (got %d, want %d)",
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
		result := fn.Fn(args...)
		if errObj, ok := result.(*object.Error); ok {
			errObj.Pos = callExpr.Function.Pos()
		}
		return result

	default:
		return newError(callExpr.Pos(), "not a function: %s", fn.Type())
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

	return newErrorWithoutPos("unknown operator: %s%s", op, right.Type())
}

func evalInfixExpression(op string, left, right object.Object, yoloOK bool) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && op == "<<":
		left := left.(*object.Array)
		left.Elements = append(left.Elements, right)
		return left

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
		case "<=":
			return toYeetBool(left.Value <= rVal)
		case ">=":
			return toYeetBool(left.Value >= rVal)
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
			return &object.Number{Value: float64(left.Value % int64(right.Value))}
		case "<":
			return toYeetBool(lVal < right.Value)
		case ">":
			return toYeetBool(lVal > right.Value)
		case "<=":
			return toYeetBool(lVal <= right.Value)
		case ">=":
			return toYeetBool(lVal >= right.Value)
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
		return newErrorWithoutPos("type mismatch: %s %s %s", left.Type(), op, right.Type())

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
		case "<=":
			return toYeetBool(left.Value <= right.Value)
		case ">=":
			return toYeetBool(left.Value >= right.Value)
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
		case "<=":
			return toYeetBool(left.Value <= right.Value)
		case ">=":
			return toYeetBool(left.Value >= right.Value)
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

	return newErrorWithoutPos("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		return obj.Value
	case *object.String:
		return len(obj.Value) > 0
	case *object.Array:
		return len(obj.Elements) > 0
	case *object.Hashmap:
		return len(obj.Pairs) > 0
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

func newErrorWithoutPos(format string, args ...any) *object.Error {
	return newError(-1, format, args...)
}

func newError(position int, format string, args ...any) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, args...), Pos: position}
}
