package eval_test

import "testing"

func TestBuiltinLenFunction(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len([1, 2, 3])`, 3},
		{`len(0..3)`, 4},
		{`len(3..0)`, 4},
		{`len(%{ "a": 1, "b": 2})`, 2},
		{`len(1)`, errmsg{"argument to `len` not supported, got INTEGER"}},
		{`len("one", "two")`, errmsg{"wrong number of args for len (got 2, want 1)"}},
	})
}

func TestBuiltinYassertFunction(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`yassert(1 == 1)`, nil},
		{`yassert(1 == 2)`, errmsg{"yassert failed"}},
		{`yassert(false)`, errmsg{"yassert failed"}},
		{`a := 5; b := 6; yassert(a == b)`, errmsg{"yassert failed"}},
		{`yassert(true)`, nil},
		{`yassert(1 == 2, "one isn't two")`, errmsg{"yassert failed: one isn't two"}},
	})
}

func TestBuiltinCastingFunctions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`yarn(5)`, "5"},
		{`yarn(true)`, "true"},
		{`yarn([1, 2, 3])`, "[1, 2, 3]"},
		{`yarn(0..2)`, "0..2"},
		{`yarn("test")`, "test"},

		{`chr(5)`, "\x05"},
		{`chr(15)`, "\x0f"},
		{`chr(60)`, "<"},
		{`chr(65)`, "A"},
		{`chr(90)`, "Z"},
		{`chr(97)`, "a"},
		{`chr(122)`, "z"},
		{`chr(322)`, "ł"},
		{`chr("5")`, "5"},
		{`chr("24")`, "2"},

		{`int(5)`, 5},
		{`int(5.0)`, 5},
		{`int(true)`, 1},
		{`int(false)`, 0},
		{`int("5")`, 5},
		{`int(0..2)`, errmsg{"unsupported argument type for int, got RANGE"}},
		{`int([1, 2, 3])`, errmsg{"unsupported argument type for int, got ARRAY"}},

		{`float(5)`, 5.0},
		{`float(5.0)`, 5.0},
		{`float(true)`, 1.0},
		{`float(false)`, 0.0},
		{`float("5")`, 5.0},
		{`float(0..2)`, errmsg{"unsupported argument type for float, got RANGE"}},
		{`float([1, 2, 3])`, errmsg{"unsupported argument type for float, got ARRAY"}},
	})
}

func TestBuiltinYoinkFunction(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`arr := [1, 2, 3]; x := yoink(arr); x`, 3},
		{`arr := [1, 2, 3]; x := yoink(arr); arr`, []int64{1, 2}},
		{`arr := [1, 2, 3]; x := yoink(arr, 1); x`, 2},
		{`arr := [1, 2, 3]; x := yoink(arr, 1); arr`, []int64{1, 3}},
		{`arr := [1, 2, 3]; x := yoink(arr, 100); x`, nil},

		{`str := "howdy"; x := yoink(str); x`, "y"},
		{`str := "howdy"; x := yoink(str); str`, "howd"},
		{`str := "howdy"; x := yoink(str, 1); x`, "o"},
		{`str := "howdy"; x := yoink(str, 1); str`, "hwdy"},
		{`str := "howdy"; x := yoink(str, 100); x`, nil},

		{`a := 69; b := yoink(a); [a, b]`, []int64{0, 69}},
		{`a := -69; b := yoink(a); [a, b]`, []int64{0, -69}},
		{`a := 69.0; b := yoink(a); [a, b]`, []float64{0.0, 69.0}},

		{`a := null; b := yoink(a); b`, nil},

		{`yoink(0..69)`, errmsg{"cannot yoink from RANGE"}},
	})
}
