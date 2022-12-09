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

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE

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
		obj, ok := right.(*object.Integer)
		if !ok {
			return NULL // TODO handle errors better than just returning null
		}
		return &object.Integer{Value: -obj.Value}

	default:
		return NULL
	}
}
