package eval

import (
	"ylang/ast"
	"ylang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		var result object.Object
		for _, stmt := range node.Statements {
			result = Eval(stmt)
		}
		return result

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return toYeetBool(node.Value)

	default:
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		switch right {
		case TRUE:
			return FALSE
		case FALSE, NULL:
			return TRUE
		default:
			return FALSE
		}

	case "-":
		right, ok := right.(*object.Integer)
		if !ok {
			return NULL // TODO handle errors better than just returning null
		}
		return &object.Integer{Value: -right.Value}

	default:
		return NULL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
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
		default:
			return NULL
		}

	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		right := right.(*object.Boolean)
		left := left.(*object.Boolean)
		switch operator {
		case "==":
			return toYeetBool(left.Value == right.Value)
		case "!=":
			return toYeetBool(left.Value != right.Value)
		default:
			return NULL
		}

	default:
		return NULL
	}
}

func toYeetBool(v bool) object.Object {
	if v {
		return TRUE
	}
	return FALSE
}

func evalInfixExpression_old(operator string, left, right object.Object) object.Object {
	switch operator {
	case "+":
		right, ok := right.(*object.Integer)
		if !ok {
			return NULL // TODO handle errors better than just returning null
		}
		left, ok := left.(*object.Integer)
		if !ok {
			return NULL // TODO handle errors better than just returning null
		}
		return &object.Integer{Value: left.Value + right.Value}

	default:
		return NULL
	}
}
