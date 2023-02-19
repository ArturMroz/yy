package eval

import "testing"

func TestBuiltinFunctions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len([1, 2, 3])`, 3},
		{`len(0..3)`, 4},
		{`len(3..0)`, 4},
		{`len({ "a": 1, "b": 2})`, 2},
		{`len(1)`, errmsg{"argument to `len` not supported, got INTEGER"}},
		{`len("one", "two")`, errmsg{"wrong number of arguments. got=2, want=1"}},

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

		{`yarn(5)`, "5"},
		{`yarn(true)`, "true"},
		{`yarn([1, 2, 3])`, "[1, 2, 3]"},
		{`yarn(0..2)`, "0..2"},
		{`yarn("test")`, "test"},
	})
}
