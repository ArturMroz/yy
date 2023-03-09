package eval

import "testing"

func TestBuiltinFunctions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len([1, 2, 3])`, 3},
		{`len(0..3)`, 4},
		{`len(3..0)`, 4},
		{`len(%{ "a": 1, "b": 2})`, 2},
		{`len(1)`, errmsg{"argument to `len` not supported, got INTEGER"}},
		{`len("one", "two")`, errmsg{"wrong number of args for len (got 2, want 1)"}},

		{`yassert(1 == 1)`, nil},
		{`yassert(1 == 2)`, errmsg{"yassert failed"}},
		{`yassert(false)`, errmsg{"yassert failed"}},
		{`a := 5; b := 6; yassert(a == b)`, errmsg{"yassert failed"}},
		{`yassert(true)`, nil},
		{`yassert(1 == 2, "one isn't two")`, errmsg{"yassert failed: one isn't two"}},

		{`arr := [1, 2, 3]; x := yoink(arr); x`, 3},
		{`arr := [1, 2, 3]; x := yoink(arr); arr`, []int64{1, 2}},
		{`arr := [1, 2, 3]; x := yoink(arr, 1); x`, 2},
		{`arr := [1, 2, 3]; x := yoink(arr, 1); arr`, []int64{1, 3}},
		{`str := "howdy"; x := yoink(str); x`, "y"},
		{`str := "howdy"; x := yoink(str); str`, "howd"},
		{`str := "howdy"; x := yoink(str, 1); x`, "o"},
		{`str := "howdy"; x := yoink(str, 1); str`, "hwdy"},
		{`yoink(69)`, errmsg{"cannot yoink from INTEGER"}},

		{`swap([1, 2, 3, 4], 0, 2)`, []int64{3, 2, 1, 4}},
		{`a := swap([1, 2, 3, 4], 1, 3); a`, []int64{1, 4, 3, 2}},
		{`a := [1, 2, 3, 4]; swap(a, 1, 3)`, []int64{1, 4, 3, 2}},
		{`a := [1, 2, 3, 4]; swap(a, 1, 69)`, []int64{1, 2, 3, 4}},
		{`a := [1, 2, 3, 4]; swap(a, 1, -3)`, []int64{1, 2, 3, 4}},
	})
}

func TestBuiltinCastingFunctions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`yarn(5)`, "5"},
		{`yarn(true)`, "true"},
		{`yarn([1, 2, 3])`, "[1, 2, 3]"},
		{`yarn(0..2)`, "0..2"},
		{`yarn("test")`, "test"},

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
