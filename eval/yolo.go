package eval

import (
	"strconv"
	"strings"

	"yy/object"
)

func yoloPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "-":
		switch right := right.(type) {
		case *object.Null:
			return ABYSS

		case *object.String:
			result := strings.Map(rot13, right.Value)
			return &object.String{Value: result}

		case *object.Array:
			result := &object.Array{
				Elements: make([]object.Object, len(right.Elements)),
			}
			for i, v := range right.Elements {
				result.Elements[i] = evalPrefixExpression(op, v, true)
			}
			return result

		case *object.Hash:
			// TODO reverse keys with values

		case *object.Range:
			// TODO invert

		case *object.Function:
			// TODO negate return value?
		}
	}

	return newError("unknown operator: %s%s", op, right.Type())
}

func yoloInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.INTEGER_OBJ:
		return yoloInfixExpression(op, right, left) // handle below

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

		if op == "+" {
			return &object.String{Value: left.String() + right.String()}
		}

		return yoloInfixExpression(op, right, left) // handle below

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

	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.INTEGER_OBJ:
		return yoloInfixExpression(op, right, left) // handle below

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.BOOLEAN_OBJ:
		boolVal := right.(*object.Boolean).Value
		bitSet := int64(0)
		if boolVal {
			bitSet = 1
		}
		boolAsInt := &object.Integer{Value: bitSet}
		return evalInfixExpression(op, left, boolAsInt, true)

		// TODO handle other type combinations
	}

	// catch all: just convert to string and concatenate
	return &object.String{Value: left.String() + right.String()}
}

func rot13(ch rune) rune {
	switch {
	case 'A' <= ch && ch <= 'M', 'a' <= ch && ch <= 'm':
		ch += 13
	case 'N' <= ch && ch <= 'Z', 'n' <= ch && ch <= 'z':
		ch -= 13
	}
	return ch
}

var collectiveNouns = map[string]string{
	"actor":        "cast",
	"angel":        "choir",
	"ant":          "army",
	"asteroid":     "belt",
	"bacteria":     "culture",
	"badger":       "cete",
	"balloon":      "festival",
	"banana":       "bunch",
	"barracuda":    "battery",
	"bat":          "colony",
	"beaver":       "colony",
	"bee":          "commonwealth",
	"book":         "library",
	"camel":        "caravan",
	"cat":          "destruction",
	"cheetah":      "coalition",
	"chick":        "chattering",
	"chicken":      "cluck",
	"chimpanzee":   "cartload",
	"clam":         "bed",
	"coyote":       "pack",
	"crocodile":    "bask",
	"crow":         "murder",
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
	"fish":         "haul",
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
	"hippopotamus": "crash",
	"horsemen":     "cavalcade",
	"hound":        "cry",
	"hummingbird":  "charm",
	"hyena":        "clan",
	"insect":       "swarm",
	"island":       "archipelago",
	"judge":        "bench",
	"knight":       "banner",
	"lark":         "ascension",
	"leper":        "colony",
	"matche":       "chain",
	"meerkat":      "mob",
	"monkey":       "cartload",
	"mule":         "barren",
	"musician":     "band",
	"native":       "tribe",
	"onlooker":     "crowd",
	"otter":        "bevy",
	"owl":          "wisdom",
	"oyster":       "bed",
	"paper":        "budget",
	"partridge":    "bew",
	"peasant":      "toil",
	"performer":    "troupe",
	"pheasant":     "brace",
	"pigeon":       "bunch",
	"polar bear":   "aurora",
	"prairie dog":  "coterie",
	"ptarmigan":    "covey",
	"puffin":       "circus",
	"quail":        "bevy",
	"rabbit":       "wrack",
	"raven":        "conspiracy",
	"reed":         "clump",
	"rhinoceros":   "crash",
	"sailor":       "crew",
	"salmon":       "bind",
	"savage":       "horde",
	"seal":         "harem",
	"ship":         "armada",
	"slug":         "cornucopia",
	"soldier":      "brigade",
	"spider":       "cluster",
	"star":         "constellation",
	"starling":     "clutter",
	"student":      "class",
	"swan":         "bevy",
	"thief":        "den",
	"tiger":        "ambush",
	"toucan":       "durante",
	"tree":         "forest",
	"truck":        "convoy",
	"turkey":       "brood",
	"turtle":       "bale",
	"unicorn":      "blessing",
	"vulture":      "wake",
	"widow":        "ambush",
	"wigeon":       "coil",
	"woodcock":     "covey",
	"worm":         "clew",
	"zebra":        "zeal",
}
