package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	// TODO handle errors
	switch node := node.(type) {
	case *Program:
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}

	case *BlockStatement:
		for i := range node.Statements {
			node.Statements[i], _ = Modify(node.Statements[i], modifier).(Statement)
		}

	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)

	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)

	case *YifExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}

	case *YetExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)

	case *YallExpression:
		node.Iterable, _ = Modify(node.Iterable, modifier).(Expression)
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)

	case *DeclareExpression:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	case *YeetStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)

	// literals

	case *ArrayLiteral:
		for i := range node.Elements {
			node.Elements[i] = Modify(node.Elements[i], modifier).(Expression)
		}

	case *HashmapLiteral:
		newPairs := map[Expression]Expression{}
		for k, v := range node.Pairs {
			k, _ = Modify(k, modifier).(Expression)
			v, _ = Modify(v, modifier).(Expression)
			newPairs[k] = v
		}
		node.Pairs = newPairs

	case *RangeLiteral:
		node.Start, _ = Modify(node.Start, modifier).(Expression)
		node.End, _ = Modify(node.End, modifier).(Expression)
	}

	return modifier(node)
}
