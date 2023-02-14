package eval

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"yy/ast"
	"yy/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	ABYSS = &object.String{Value: "Stare at the abyss long enough, and it starts to stare back at you."}
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
		if node.Function.TokenLiteral() == "quote" { // TODO this is ugly
			return quote(node.Arguments[0], env) // quote only supprots 1 arg
			// arg := node.Arguments[0] // quote only takes 1 arg
			// return &object.Quote{Node: arg}
		}

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
		if env.YoloMode() {
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

		return evalInfixExpression(node.Operator, left, right, env.YoloMode())

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
			return NULL
		}

	case *ast.YoloExpression:
		extendedEnv := object.NewEnclosedEnvironment(env)
		env.Set(object.YoloKey, TRUE)

		result := Eval(node.Body, extendedEnv)
		return result

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
			return newError("cannot iterate over %s, type of %s", iter, iter.Type())
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
		elts := []object.Object{}
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

func evalInfixExpression(op string, left, right object.Object, yoloOK bool) object.Object {
	switch {
	case left.Type() != right.Type():
		switch op {
		case "==":
			return toYeetBool(false)
		case "!=":
			return toYeetBool(true)
		}

		if yoloOK {
			return yoloMode(op, left, right)
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
			// TODO DeepEqual's performance isn't great
			return toYeetBool(reflect.DeepEqual(left.Elements, right.Elements))
		case "!=":
			return toYeetBool(!reflect.DeepEqual(left.Elements, right.Elements))
		}
	}

	return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func yoloMode(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.INTEGER_OBJ:
		return yoloMode(op, right, left) // handle below

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.ARRAY_OBJ:
		left := left.(*object.Integer)
		right := right.(*object.Array)

		switch op {
		case "+", "-", "*", "/":
			result := &object.Array{
				Elements: make([]object.Object, len(right.Elements)),
			}
			for i, v := range right.Elements {
				result.Elements[i] = evalInfixExpression(op, left, v, true)
			}
			return result

		case "<":
			return toYeetBool(true)
		case ">":
			return toYeetBool(false)
		}

	case left.Type() == object.STRING_OBJ && right.Type() == object.INTEGER_OBJ:
		left := left.(*object.String)

		if v, err := strconv.Atoi(left.Value); err == nil {
			return evalInfixExpression(op, &object.Integer{Value: int64(v)}, right, true)
		}

		return yoloMode(op, right, left) // handle below

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.STRING_OBJ:
		left := left.(*object.Integer)
		right := right.(*object.String)

		if v, err := strconv.Atoi(right.Value); err == nil {
			return evalInfixExpression(op, left, &object.Integer{Value: int64(v)}, true)
		}

		switch op {
		case "*":
			if left.Value < 0 {
				return ABYSS
			}

			if collective, ok := collectiveNouns[strings.TrimSpace(right.Value)]; ok {
				return &object.String{Value: collective}
			}

			result := strings.Repeat(right.Value, int(left.Value))
			return &object.String{Value: result}

		case "/":
			if left.Value <= 0 {
				return ABYSS
			}

			ss := strings.Split(right.Value, "")
			elems := make([]object.Object, len(ss))
			for i, s := range ss {
				elems[i] = &object.String{Value: s}
			}
			result := &object.Array{Elements: elems}
			return result

		case "<":
			return toYeetBool(true)
		case ">":
			return toYeetBool(false)
		}

		// TODO handle other type combinations
	}

	// catch all: just convert to string and concatenate
	return &object.String{Value: left.Inspect() + right.Inspect()}
}

func isTruthy(obj object.Object) bool {
	// Ruby's truthiness rule: nil & false are falsy, everything else is truthy
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}

func toYeetBool(b bool) object.Object {
	if b {
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

// TODO find a better place for this
var collectiveNouns = map[string]string{
	"actor":        "cast",
	"alligator":    "congregation",
	"angel":        "choir",
	"ant":          "army",
	"architect":    "argument",
	"arsonist":     "conflagration",
	"artillery":    "battery",
	"asteroid":     "belt",
	"baboon":       "congress",
	"bacteria":     "culture",
	"badger":       "cete",
	"banana":       "bunch",
	"barracuda":    "battery",
	"bat":          "colony",
	"batterie":     "bank",
	"beaver":       "colony",
	"bee":          "commonwealth",
	"beer":         "brew",
	"bishop":       "bench",
	"bobolink":     "chain",
	"bomb":         "cluster",
	"bread":        "batch",
	"buck":         "brace",
	"budgerigar":   "chatter",
	"bullfinche":   "bellowing",
	"camel":        "caravan",
	"cat":          "destruction",
	"cheetah":      "coalition",
	"chick":        "chattering",
	"chicken":      "cluck",
	"chimpanzee":   "cartload",
	"circuit":      "bank",
	"clam":         "bed",
	"classicist":   "codex",
	"clergy":       "assembly",
	"computer":     "cluster",
	"conie":        "bury",
	"coyote":       "band",
	"crocodile":    "bask",
	"crow":         "murder",
	"cur":          "cowardice",
	"cutlery":      "canteen",
	"deer":         "bevy",
	"director":     "board",
	"diver":        "bubble",
	"doctor":       "confab",
	"donkey":       "drove",
	"dove":         "bevy",
	"drawer":       "chest",
	"duck":         "badelynge",
	"eagle":        "aerie",
	"economist":    "clashing",
	"eel":          "bind",
	"egg":          "clutch",
	"event":        "chain",
	"fairie":       "charm",
	"ferret":       "business",
	"finche":       "charm",
	"flie":         "business",
	"flour":        "boll",
	"flower":       "bouquet",
	"game":         "bag",
	"giraffe":      "corps",
	"goat":         "drove",
	"gorilla":      "band",
	"grape":        "bunch",
	"grasshopper":  "cloud",
	"grouse":       "brood",
	"guillemot":    "bazaar",
	"gun":          "arsenal",
	"hawk":         "aerie",
	"hedgehog":     "array",
	"hen":          "brood",
	"herring":      "army",
	"hide":         "dicker",
	"hippopotami":  "bloat",
	"hippopotamus": "crash",
	"historian":    "argumentation",
	"horsemen":     "cavalcade",
	"hound":        "cry",
	"hummingbird":  "charm",
	"hyena":        "clan",
	"information":  "bits",
	"island":       "archipelago",
	"jay":          "band",
	"jewel":        "cache",
	"judge":        "bench",
	"knight":       "banner",
	"lark":         "ascension",
	"leper":        "colony",
	"magistrate":   "bench",
	"man":          "band",
	"manager":      "cost",
	"matche":       "chain",
	"monitor":      "bank",
	"monkey":       "cartload",
	"mormon":       "branch",
	"motorcyclist": "clutch",
	"mourner":      "cortege",
	"mule":         "barren",
	"musician":     "band",
	"onlooker":     "crowd",
	"otter":        "bevy",
	"oyster":       "bed",
	"paper":        "budget",
	"partridge":    "bew",
	"people":       "community",
	"pheasant":     "brace",
	"pigeon":       "bunch",
	"plum":         "basket",
	"polar bear":   "aurora",
	"polecat":      "chine",
	"prairie dog":  "coterie",
	"preacher":     "converting",
	"ptarmigan":    "covey",
	"puffin":       "circus",
	"quail":        "bevy",
	"rabbit":       "bury",
	"raven":        "conspiracy",
	"reed":         "clump",
	"relative":     "descent",
	"rhinoceros":   "crash",
	"sailor":       "crew",
	"saint":        "communion",
	"salmon":       "bind",
	"seal":         "crash",
	"ship":         "armada",
	"slug":         "cornucopia",
	"smoker":       "confraternity",
	"snake":        "bed",
	"soldier":      "brigade",
	"spider":       "cluster",
	"star":         "constellation",
	"starling":     "clutter",
	"student":      "class",
	"swan":         "bevy",
	"teal":         "diving",
	"thief":        "den",
	"tiger":        "ambush",
	"toucan":       "durante",
	"tree":         "forest",
	"truck":        "convoy",
	"turkey":       "brood",
	"turtle":       "bale",
	"unicorn":      "blessing",
	"weapon":       "cache",
	"widow":        "ambush",
	"wigeon":       "coil",
	"wizard":       "argument",
	"woodcock":     "covey",
	"worm":         "clew",
}
