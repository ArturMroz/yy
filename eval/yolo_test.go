package eval

import (
	"testing"

	"yy/object"
)

func TestYoloExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		// normal operations in yolo mode are still ok
		{"yolo { 2 + 2 }", 4},
		{"yolo { 2 + 2; 8 }", 8},
		{"yolo { a := 1; a }", 1},
		{"yolo { a := 1; b := 2; a + b }", 3},
		{"result := yolo { a := 1; b := 2; a + b }; result", 3},

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

		// functions
		{`5 + \x { x + 2 } `, errmsg{"type mismatch: INTEGER + FUNCTION"}},
		{
			`yolo { 
				fn := \a { x := 2; yeet a + x };
				fn2 := 69 + fn;
				[fn(1), fn2(1)]
			}`,
			[]int64{3, 72},
		},
		{
			`yolo { 
				double 			  := \a { yeet a * 2 };
				triple_the_double := double * 3;
				[double(2), triple_the_double(2)]
			}`,
			[]int64{4, 12},
		},
		{
			`yolo {
				max := \a b {
					yif a > b {
						yeet a
					}
					yeet b
				};
				max_plus := max + 5;
				[max(2, 10), max_plus(2, 10)]
			}`,
			[]int64{10, 15},
		},
		{
			`yolo {
				fn := \a { a + "hello" };
				fn2 := fn + "bake"; // baking in args 
				fn("test") // prints "testhello"
				fn2()	   // prints "bakehello", notice empty arg list
			}`,
			"bakehello",
		},

		// auto declaring variables if they don't exsist
		{"a = 1; a", errmsg{"identifier not found: a"}},
		{"yolo { a = 1; a };", 1},
		{"a := 5; yolo { a = 69 }; a", 69},

		// what happens in yolo, stays in yolo
		{"yolo { a := 1; a }; a", errmsg{"identifier not found: a"}},

		// prefix
		{`yolo { -"Gurer'f Lrrg va rirel Lbvax."}`, "There's Yeet in every Yoink."},
		{`yolo { -[1, 2, 3]}`, []int64{-1, -2, -3}},
		{`yolo { -null }`, object.ABYSS.Value},
	})
}
