package eval

import (
	"fmt"

	"yy/ast"
	"yy/object"
	"yy/token"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	found := []int{}

	// 1. find macros and note down their indexes
	for i, stmt := range program.Statements {
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			continue
		}
		assExpr, ok := exprStmt.Expression.(*ast.AssignExpression)
		if !ok {
			continue
		}
		macroLiteral, ok := assExpr.Value.(*ast.MacroLiteral)
		if !ok {
			continue
		}

		macro := &object.Macro{
			Parameters: macroLiteral.Parameters,
			Env:        env,
			Body:       macroLiteral.Body,
		}

		env.Set(assExpr.Name.Value, macro)
		found = append(found, i)
	}

	// 2. remove identified macros from ast
	for i := len(found) - 1; i >= 0; i = i - 1 {
		idx := found[i]
		program.Statements = append(
			program.Statements[:idx],
			program.Statements[idx+1:]...,
		)
	}
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpr, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		identifier, ok := callExpr.Function.(*ast.Identifier)
		if !ok {
			return node
		}
		obj, ok := env.Get(identifier.Value)
		if !ok {
			return node
		}
		macro, ok := obj.(*object.Macro)
		if !ok {
			return node
		}

		// quote args
		args := []*object.Quote{}
		for _, a := range callExpr.Arguments {
			args = append(args, &object.Quote{Node: a})
		}

		extendedEnv := object.NewEnclosedEnvironment(macro.Env)
		for i, param := range macro.Parameters {
			extendedEnv.Set(param.Value, args[i])
		}

		evaluated := Eval(macro.Body, extendedEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			// TODO handle in less panicky way
			panic("runtime error: only quoted objects can be returned from macros")
		}

		return quote.Node
	})
}

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
