package eval

import (
	"fmt"

	"yy/ast"
	"yy/object"
	"yy/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if call.Function.TokenLiteral() != "unquote" { // TODO ugly
			return node
		}

		if len(call.Arguments) != 1 { // only 1 arg is supported atm
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return objectToASTNode(unquoted)
	})
}

func objectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}

	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}

	case *object.Quote:
		return obj.Node

	default:
		return nil
	}
}
