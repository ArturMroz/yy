package eval

import (
	"testing"

	"yy/object"
)

func TestYoloNormal(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		// normal operations in yolo mode are still ok and yolo block returns a value
		{"yolo { 2 + 2 }", 4},
		{"yolo { 2 + 2; 8 }", 8},
		{"yolo { a := 1; a }", 1},
		{"yolo { a := 1; b := 2; a + b }", 3},
		{"result := yolo { a := 1; b := 2; a + b }; result", 3},
	})
}

func TestYoloDeclarations(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		// legit operations in yolo mode are still ok
		{"yolo { 2 + 2 }", 4},
		{"yolo { 2 + 2; 8 }", 8},
		{"yolo { a := 1; a }", 1},
		{"yolo { a := 1; b := 2; a + b }", 3},
		{"result := yolo { a := 1; b := 2; a + b }; result", 3},

		// auto declaring variables if they don't exsist
		{"a = 1; a", errmsg{"identifier not found: a"}},
		{"yolo { a = 1; a };", 1},
		{"a := 5; yolo { a = 69 }; a", 69},

		// what happens in yolo, stays in yolo
		{"yolo { a := 1 }; a", errmsg{"identifier not found: a"}},
	})
}

func TestYoloPrefixExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`yolo { -"Gurer'f Lrrg va rirel Lbvax."}`, "There's Yeet in every Yoink."},
		{`yolo { -[1, 2, 3]}`, []int64{-1, -2, -3}},
		{`yolo { -null }`, object.ABYSS.Value},
		// {`yolo { -null }`, int64(^uint(0) >> 1)}, // max int value
		{`yolo { -(0..5) }`, rng{5, 0}},
		{`yolo { -(5..0) }`, rng{0, 5}},
	})
}

func TestYoloInfixExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		// arrays
		{"yolo { 3 + [1, 2, 3] }", []int64{4, 5, 6}},
		{"yolo { [1, 2, 3] + 4 }", []int64{5, 6, 7}},
		{"yolo { 3 * [1, 2, 3] }", []int64{3, 6, 9}},
		{"yolo { [1, 2, 3] * 3 }", []int64{3, 6, 9}},

		// strings
		{`2 + "troll"`, errmsg{"type mismatch: INTEGER + STRING"}},
		{`yolo { 3 * "22" }`, 66},
		{`yolo { "22" * 3 }`, 66},
		{`yolo { 3 * "troll" }`, "trolltrolltroll"},
		{`yolo { "troll" * 3 }`, "trolltrolltroll"},
		{`yolo { "tree" * 3 }`, "forest"},
		{`yolo { "tree   " * 3 }`, "forest"},
		{`yolo { "   tree" * 3 }`, "forest"},
		{`yolo { "   tree   " * 3 }`, "forest"},
		{`yolo { "crow" * 3 }`, "murder"},
		{`yolo { "test" * 0 }`, ""},
		{`yolo { "test" * -5 }`, object.ABYSS.Value},
		{`yolo { "test" / 0 }`, object.ABYSS.Value},
		{`yolo { 2 + "troll" }`, "2troll"},

		// bools
		{`yolo { 5 + true }`, 6},
		{`yolo { 5 + false }`, 5},
		{`yolo { 5 - true }`, 4},
		{`yolo { 5 - false }`, 5},

		// ranges
		{`5 + (0..5) `, errmsg{"type mismatch: INTEGER + RANGE"}},
		{`yolo { 5 + (0..5) }`, rng{5, 10}},
		{`yolo { (0..5) + 5 }`, rng{5, 10}},
		{`yolo { 5 - (0..5) }`, rng{5, 0}},
		{`yolo { (0..5) - 5 }`, rng{-5, 0}},
	})
}

func TestYoloFunctionObjects(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`5 + \x { x + 2 } `, errmsg{"type mismatch: INTEGER + FUNCTION"}},

		// ints
		{
			`yolo {
				add   := \a, b { a + b }
				add11 := add + 11
				// above line is equivalent to:
				// add11 := \b { 11 + b }  
				add11(6)
			}`,
			17,
		},
		{
			`yolo {
				add   := \a, b { a + b }
				add11 := 11 + add
				add11(6)
			}`,
			17,
		},

		// strings
		{
			`yolo {
				fn := \a { a + "hello" };
				fn2 := fn + "bake"; // baking in args 
				fn("test") // prints "testhello"
				fn2()      // prints "bakehello", notice empty arg list
			}`,
			"bakehello",
		},
		{
			`yolo {
				greet     := \name { "Hello, {name}!" }
				greet_yan := greet + "Yan"; // baking arg "Yan" into function 
				greet_yan()
			}`,
			"Hello, Yan!",
		},
		{
			`yolo {
				greet     := \name { "Hello, {name}!" }
				greet_yan := "Yan" + greet 
				greet_yan()
			}`,
			"Hello, Yan!",
		},

		// hashmaps
		{
			`yolo {
				add   := \a, b { a + b }
				add11 := add + %{ "a": 11 }
				// above line is equivalent to:
				// add11 := \b { 11 + b }  
				add11(6)
			}`,
			17,
		},
		{
			`yolo {
				add   := \a, b { a + b }
				add11 := %{ "a": 11 } + add
				add11(6)
			}`,
			17,
		},
		{
			`yolo {
				add      := \a, b { a + b }
				add11to5 := add + %{ "a": 11, "b": 5 }
				// above line is equivalent to:
				// add11to5 := \ { 11 + 5 }  
				add11to5()
			}`,
			16,
		},
		{
			`yolo {
				add      := \a, b { a + b }
				add11to5 := add + %{ "a": 11, "b": 5 }
				// above line is equivalent to:
				// add11to5 := \ { 11 + 5 }  
				add11to5()
			}`,
			16,
		},

		// arrays
		{
			`yolo {
				add     := \a b c { a + b + c }
				add_all := add + [1, 2, 3]
				add_all()
			}`,
			6,
		},
		{
			`yolo {
				add     := \a b c { a + b + c }
				add_all := add + [1, 2]
				add_all(3)
			}`,
			6,
		},
		{
			`yolo {
				add     := \a b c { a + b + c }
				add_all := [1, 2, 3] + add
				add_all()
			}`,
			6,
		},
		{
			`yolo {
				add     := \a b c { a + b + c }
				add_all := [1, 2, 3, 4, 5] + add // surplus is ignored
				add_all()
			}`,
			6,
		},
		{
			`yolo {
				add     := \a b c { a + b + c }
				add_all := [1, 2] + add
				add_all(3)
			}`,
			6,
		},

		// null
		{
			`yolo {
				add      := \a, b { a + b }
				add_null := add + null
				add_null(5, 6)
			}`,
			11,
		},
	})
}
