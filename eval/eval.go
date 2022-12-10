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
		return evalStatements(node.Statements)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		condition := Eval(node.Condition)
		if isTruthy(condition) {
			return Eval(node.Consequence)
		} else if node.Alternative != nil {
			return Eval(node.Alternative)
		} else {
			return NULL
		}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return toYeetBool(node.Value)

	default:
		return NULL
	}
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range statements {
		result = Eval(stmt)
	}
	return result
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

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return toYeetBool(!isTruthy(right))

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

// func evalInfixExpression_old(operator string, left, right object.Object) object.Object {
// 	switch operator {
// 	case "+":
// 		right, ok := right.(*object.Integer)
// 		if !ok {
// 			return NULL // TODO handle errors better than just returning null
// 		}
// 		left, ok := left.(*object.Integer)
// 		if !ok {
// 			return NULL // TODO handle errors better than just returning null
// 		}
// 		return &object.Integer{Value: left.Value + right.Value}

// 	default:
// 		return NULL
// 	}
// }
