package eval_test

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
		{"yolo { a = 1; a };", 1},
		{"a := 5; yolo { a = 69 }; a", 69},
		// vars need to be declared before use outside of yolo
		{"a = 1; a", errmsg{"identifier not found: a (to declare a variable use := operator)"}},
		// yolo rules don't leak outside of yolo block
		{"yolo {}; a = 1", errmsg{"identifier not found: a (to declare a variable use := operator)"}},
		// what happens in yolo, stays in yolo
		{"yolo { a := 1 }; a", errmsg{"identifier not found: a"}},
	})
}

func TestYoloPrefixExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`yolo { -"Gurer'f Lrrg va rirel Lbvax."}`, "There's Yeet in every Yoink."},
		{`yolo { -[1, 2, 3]}`, []int64{-1, -2, -3}},
		{`yolo { -null }`, object.ABYSS.Value},
		{`yolo { -true }`, false},
		{`yolo { -false }`, true},
		{`yolo { -(0..5) }`, rng{5, 0}},
		{`yolo { -(5..0) }`, rng{0, 5}},
		{`yolo { hash := %{ "a": "z" }; (-hash)["z"] }`, "a"},
		{`yolo { hash := %{ "a": 6, "b": 9 }; (-hash)[6] }`, "a"},
		{`yolo { hash := %{ "a": 5, "c": [2, 3] }; len(-hash) }`, 2},
		{`yolo { hash := %{ "a": 5, "c": [2, 3] }; (-hash)[[2, 3]] }`, "c"},

		{`yolo { fn := \a { 2 * a }; fn_neg := -fn; fn_neg(5) }`, -10},
		{
			`yolo { 
				fn     := \a b { yif a > b { a } yels { b } }
				fn_neg := -fn;
				[fn_neg(5, 8), fn_neg(15, 8)]
			}`,
			[]int64{-8, -15},
		},
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

func TestYoloInfixFunctionObject(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{
			`yolo { 
				fn  := \a b { yif a > b { a } yels { b } }
				fn2 := fn * 3;
				[fn2(5, 8), fn2(15, 8)]
			}`,
			[]int64{24, 45},
		},
		{
			`yolo { 
				fn  := \a b { yif a > b { a } yels { b } }
				fn2 := 3 * fn;
				[fn2(5, 8), fn2(15, 8)]
			}`,
			[]int64{24, 45},
		},
		{
			`yolo { 
				fn  := \a b { yif a > b { a } yels { b } }
				fn2 := fn / 3;
				[fn2(6, 9), fn2(15, 9)]
			}`,
			[]int64{3, 5},
		},
		{
			`yolo { 
				fn  := \a b { yif a > b { a } yels { b } }
				fn2 := 15 / fn;
				[fn2(5, 3), fn2(15, 9)]
			}`,
			[]int64{3, 1},
		},
		{
			`yolo { 
				fn  := \a b { yif a > b { a } yels { b } }
				fn2 := fn - 3;
				[fn2(6, 9), fn2(15, 9)]
			}`,
			[]int64{6, 12},
		},
		{
			`yolo { 
				fn  := \a b { yif a > b { a } yels { b } }
				fn2 := 3 - fn;
				[fn2(6, 9), fn2(15, 9)]
			}`,
			[]int64{-6, -12},
		},
	})
}

func TestYoloFunctionAdding(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{
			`
			add3 := \a { a + 3 }
			mul5 := \b { b * 5 }
			yolo {
				add_mul := add3 + mul5
				add_mul(4)
			}`,
			35,
		},
		{
			`
			add3 := \a { a + 3 }
			mul5 := \b { b * 5 }
			sub2 := \c { c - 2 }
			yolo {
				add_mul_sub := add3 + mul5 + sub2
				add_mul_sub(4)
			}`,
			33,
		},
		{
			`
			add3 := \a { a + 3 }
			mul5 := \a { a * 5 }
			sub2 := \a { a - 2 }
			yolo {
				add_mul_sub := add3 + mul5 + sub2
				add_mul_sub(4)
			}`,
			33,
		},
		{
			`
			add3 := \a { a + 3 }
			mul5 := \a { a * 5 }
			sub2 := \a { a - 2 }
			yolo {
				add_mul_sub := add3 + mul5 + sub2
				add_mul_sub(4) == (add3 + mul5 + sub2)(4) 
			}`,
			true,
		},
		{
			`
			add3 := \a { a + 3 }
			mul5 := \a { a * 5 }
			sub2 := \a { a - 2 }
			yolo {
				sub2(mul5(add3(4))) == (add3 + mul5 + sub2)(4) 
			}`,
			true,
		},
		{
			`
			add  := \a, b { a + b }
			mul5 := \x { x * 5 }
			yolo {
				add_mul := add + mul5
				add_mul(1, 3)
			}`,
			20,
		},
		{
			`
			add := \a, b { a + b }
			mul := \x, y  { x * y }
			yolo {
				add_mul := add + mul
				add_mul(1, 3)
			}`,
			errmsg{"identifier not found: y"},
		},
		{
			`
			add := \a, b { a + b }
			mul := \x, y  { x * y }
			yolo {
				mul5    := mul + 5
				add_mul := add + mul5
				add_mul(1, 3)
			}`,
			20,
		},
		{
			`
			add  := \a, b { a + b }
			mul5 := \x { x * 5 }
			yolo {
				add3    := add + 3
				add_mul := add3 + mul5
				add_mul(1)
			}`,
			20,
		},
	})
}

func TestYoloArgumentsBaking(t *testing.T) {
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
				greet     := \name { "Hello, $name!" }
				greet_yan := greet + "Yan"; // baking arg "Yan" into function 
				greet_yan()
			}`,
			"Hello, Yan!",
		},
		{
			`yolo {
				greet     := \name { "Hello, $name!" }
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
			`
			add     := \a b c { a + b + c }
			add_all := yolo { [1, 2, 3] + add }
			add_all()
			`,
			6,
		},
		{
			`
			add     := \a b c { a + b + c }
			add_all := yolo { [1, 2, 3, 4, 5] + add } // surplus is ignored
			add_all()
			`,
			6,
		},
		{
			`
			add     := \a b c { a + b + c }
			add_all := yolo { [1, 2] + add }
			add_all(3)
			`,
			6,
		},

		// null
		{
			`yolo {
				add      := \a, b { a + b }
				add_null := add + null
				add_null(6)
			}`,
			"null6",
		},
		{
			`
			add      := \a, b { a + b }
			add_null := yolo { add + null }
			add_null(6)
			`,
			errmsg{"type mismatch: NULL + INTEGER"},
		},
	})
}
