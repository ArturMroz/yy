package eval_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"yy/eval"
	"yy/lexer"
	"yy/object"
	"yy/parser"
)

type evalTestCase struct {
	input    string
	expected any
}

func TestEvalIntegerExpression(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"5", 5},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"5 % 5", 0},
		{"5 % 3", 2},
		{"5 % 5 + 7", 7},
		{"7 + 5 % 5", 7},
	})
}

func TestEvalFloatExpression(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"5.0", 5.0},
		{"-15.0", -15.0},
		{"2.0 - 2", 0.0},
		{"2 - 2.0", 0.0},
		{"5 + 5.0 + 5 + 5 - 10", 10.0},
		{"2 * 2 * 2.0 * 2 * 2", 32.0},
		{"(5 + 10.0 * 2 + 15 / 3) * 2 + -10", 50.0},
		{"5.0 % 5", 0.0},
		{"5.0 % 3", 2.0},
	})
}

func TestStringLiteral(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`""`, ""},
		{`"piece of yarn"`, "piece of yarn"},
		{`"Żółć ∈ 陽子, ようこ ヨウコ"`, "Żółć ∈ 陽子, ようこ ヨウコ"},
	})
}

func TestTemplateStringLiteral2(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"age := 69; `i'm {age} yr old`", "i'm 69 yr old"},
		{"age := 69; `i'm { age + 2 } yr old`", "i'm 71 yr old"},
		{"`i'm { 8 + 2 * 3 } yr old`", "i'm 14 yr old"},
		{
			"age := 69; `i'm { age + 2 } yr old and have { 2 * 3 } dogs`",
			"i'm 71 yr old and have 6 dogs",
		},
		// {"`apples := 1; kiwis := 2; mangos := 3; `{apples}{kiwis}{mangos}`", "123"},
	})
}

func TestTemplateStringLiteral(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`age := 69; "i'm $age yr old"`, "i'm 69 yr old"},
		{`age := 69; "i'm $age! yr old"`, "i'm 69! yr old"},
		{`age := 69; "i'm $age"`, "i'm 69"},
		{
			`age := 69; "i'm Żółć ∈ 陽子 ようこ 陽 $age 陽 yr old"`,
			"i'm Żółć ∈ 陽子 ようこ 陽 69 陽 yr old",
		},
		{`name := "Yolanda"; "Hello, $name!"`, "Hello, Yolanda!"},
		{
			`age := 69; name := "Yolanda"; "i'm $name and i'm $age yr old"`,
			"i'm Yolanda and i'm 69 yr old",
		},
		{
			`apples := 69; pears := 8; "i've got $apples apples and $pears pears"`,
			"i've got 69 apples and 8 pears",
		},
		{
			`n1 := 69; n2 := 8; "i've got $n1 apples and $n2 pears"`,
			"i've got 69 apples and 8 pears",
		},
		{
			`n1 := 69; n2 := 8; n3 := 7; "i've got $n1 apples and $n2, $n3 other things"`,
			"i've got 69 apples and 8, 7 other things",
		},
		{`apples := 1; kiwis := 2; mangos := 3; "$apples$kiwis$mangos"`, "123"},
		{`n1 := 69; n2 := 8; n3 := 420; "$n1$n2$n3"`, "698420"},
		{`"i'm $$age yr old"`, "i'm $age yr old"},
		{`age := 69; "$$age = $age"`, "$age = 69"},
		{`"i'm $"`, "i'm $"},
		{`"$"`, "$"},
		{`"$$"`, "$"},
		{`"$$$"`, "$$"},
		{`"t $$"`, "t $"},
		{`"t $$$"`, "t $$"},
		{`"this will be $15"`, "this will be $15"},
		{`cost := 9; "this will be $$ $cost"`, "this will be $ 9"},
		{`cost := 9; "this will be $$$cost"`, "this will be $9"},
		{`"this will be 15$"`, "this will be 15$"},
		{`cost := 9; "this will be $cost$"`, "this will be 9$"},
	})
}

func TestStringConcatenation(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`"con" + "cat"`, "concat"},
		{`"" + "cat"`, "cat"},
		{`"" + ""`, ""},
		{`"con" + "cat" + "enation"`, "concatenation"},
	})
}

func TestEvalBooleanExpression(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"yoink" == "yoink"`, true},
		{`"yoink" == "yeet"`, false},
		{`"yoink" != "yoink"`, false},
		{`"yoink" != "yeet"`, true},
		{`"[1, 2, 3]" == "[1, 2, 3]"`, true},
		{`"[1, 2, 3]" == "[4, 5, 6]"`, false},
		{`"[1, 2, 3]" != "[1, 2, 3]"`, false},
		{`"[1, 2, 3]" != "[4, 5, 6]"`, true},
		{"null == null", true},
		{"null != null", false},

		// truthiness rules: false, nil and empty collections are falsy
		{"!null", true},
		{"!!null", false},
		{`!""`, true},
		{`!!""`, false},
		{"![]", true},
		{"!![]", false},
		{"!%{}", true},
		{"!!%{}", false},

		// mixed types
		{"true != null", true},
		{"true == null", false},
		{"2 != null", true},
		{"2 == null", false},
		{"[] != null", true},
		{"[] == null", false},
		{"[1, 2, 3] != null", true},
		{"[1, 2, 3] == null", false},
		{`"Testy McTestface" != null`, true},
		{`"Testy McTestface" == null`, false},
		{`"" != null`, true},
		{`"" == null`, false},
		{`"1" == 1`, false},
		{`"1" != 1`, true},
		{`"1" == [1]`, false},
		{`"1" != [1]`, true},
	})
}

func TestEvalAndExpression(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"true && false", false},
		{"true && null", nil},
		{"true && 8", 8},
		{"true && 8 && 4", 4},
		{"true && 8*2 && 4+5", 9},
		{"a := false && 8; a", false},
		{"a := true && 8; a", 8},
		{"a := 1; false && (a = 8) && 5", false},
		{"a := 1; false && (a = 8) && 5; a", 1},
	})
}

func TestEvalOrExpression(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"true || false", true},
		{"true || null", true},
		{"false || 8", 8},
		{"null || 8", 8},
		{"false || 8+2 || 4*3", 10},
		{"a := 1; false || (a = 8) || 5", 8},
		{"a := 1; false || (a = 8) || 5; a", 8},
		{"a := false || 8; a", 8},
		{"a := true || 8; a", true},
	})
}

func TestEvalOrAndExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"3 || 8 && 5", 3},
		{"false || 8 && 5", 5},
		{"false || 8 && 5 || 7", 5},
		{"(false || 8) && (9 || 7)", 9},
		{"(false || 8) && (false || 7)", 7},
	})
}

func TestBangOperator(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!0", false},
		{"!null", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!0", true},
		{"!!null", false},
	})
}

func TestIntegerArrayLiterals(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{
			"[1, 2, 3, 4, 5]",
			[]int64{1, 2, 3, 4, 5},
		},
		{
			"[1, 2 * 2, 3 + 3]",
			[]int64{1, 4, 6},
		},
		{
			"[4 / 2, 5 - 1, 8 * 4]",
			[]int64{2, 4, 32},
		},
		{
			"[1 + 1, 2 + 2, 3 + 3]",
			[]int64{2, 4, 6},
		},
	})
}

func TestArrayAppend(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{
			"a := []; a << 1",
			[]int64{1},
		},
		{
			"a := []; a << 1; a << 2",
			[]int64{1, 2},
		},
		{
			"a := []; a << 1; a << 2; a",
			[]int64{1, 2},
		},
		{
			"a := [9]; a << 1; a << 2; a",
			[]int64{9, 1, 2},
		},
	})
}

func TestHashLiterals(t *testing.T) {
	input := `
two := "two";
%{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
}`
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
	}

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Hashmap)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got %T (%+v)", evaluated, evaluated)
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got %d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key %q in Pairs", expectedKey)
		}
		if err := testIntegerObject(pair.Value, expectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"[1, 2, 3,][2]", 3},
		{"i := 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{
			"myArray := [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"myArray := [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"myArray := [1, 2, 3]; i := myArray[0]; myArray[i]",
			2,
		},

		{"a := [1, 2, 3]; a[0..len(a)]", []int64{1, 2, 3}},
		{"[1, 2, 3][0..999]", []int64{1, 2, 3}},
		{"[1, 2, 3][0..2]", []int64{1, 2}},
		{"[1, 2, 3][1..2]", []int64{2}},
		{"[1, 2, 3][-10..2]", []int64{1, 2}},
		{"a := [1, 2, 3]; b := a[0..len(a)]; a[2] == b[2]", true},
		{"a := [1, 2, 3]; b := a[0..len(a)]; b[2] = 9; a[2] == b[2]", false},

		// out of bounds access returns nil
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},

		{"arr[-1]", errmsg{"identifier not found: arr"}},
		{"arr := [2]; arr[idx]", errmsg{"identifier not found: idx"}},
		{"arr := [2]; arr[[]]", errmsg{"index operator not supported: ARRAY"}},
	})
}

func TestStringIndexExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`"Yolo McYoloface"[2]`, "l"},
		{`"Yarn"[1 + 1]`, "r"},
		{`y := "Yarn"; y[1 + 1]`, "r"},
		{`"Yolo McYoloface"[0..4]`, "Yolo"},
		{`"Yolo McYoloface"[5..11]`, "McYolo"},
		{`s := "Yolo McYoloface"; s[5..len(s)]`, "McYoloface"},
		{`s := "Yolo McYoloface"; s[-5..len(s)+5]`, "Yolo McYoloface"},
		{`"Yolo McYoloface"[69]`, nil},

		{`s1 := "Yolo McYoloface"; s2 := s1[0..len(s1)]; len(s1) == len(s2)`, true},
		{`s1 := "Yolo McYoloface"; s2 := s1[0..len(s1)]; s1[1] == s2[1]`, true},
		{`s1 := "Yolo McYoloface"; s2 := s1[0..len(s1)]; s2[1] = "X"; s1[1] == s2[1]`, false},
	})
}

func TestHashIndexExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`%{"foo": 5}["foo"]`, 5},
		{`%{"foo": 5}["bar"]`, nil},
		{`%{"ąźż": 5}["ąźż"]`, 5},
		{`%{"∈ 陽子": 5}["∈ 陽子"]`, 5},
		{`key := "foo"; %{"foo": 5}[key]`, 5},
		{`%{}["foo"]`, nil},
		{`%{5: 5}[5]`, 5},

		{
			`%{"name": "Vars McVariable"}[\x { x }];`,
			errmsg{"key not hashable: FUNCTION"},
		},
	})
}

func TestYifYelsExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"yif true { 10 }", 10},
		{"yif false { 10 }", nil},
		{"yif null { 10 }", nil},
		{"yif 1 { 10 }", 10},
		{"yif 1 < 2 { 10 }", 10},
		{"yif 1 > 2 { 10 }", nil},
		{"yif 1 > 2 { 10 } yels { 20 }", 20},
		{`yif 1 > 2 { "nope" } yels { yif 2 > 5 { "nope" } yels { 20 } }`, 20},
		{`yif 1 > 2 { "nope" } yels yif 2 > 5 { "nope" } yels { 20 }`, 20},
		{"yif 1 < 2 { 10 } yels { 20 }", 10},
		{"yif null { 10 } yels { 20 }", 20},
		{"result := yif null { 10 } yels { 20 }; result", 20},
		{"5 + yif null { 10 } yels { 20 }", 25},
		{"yif null { 10 } yels { 20 } * 2", 40},
		{"5 + yif null { 10 } yels { 20 } * 2", 45},
		{"5 + yif null { a := 10; a } yels { 20 } * 2", 45},
		{"result := 3 * yif null { 10 } yels { 20 } + 9; result", 69},
		{"yif true { a := 10; a } yels { 20 }", 10},
		{"yif true { a := 10; a } yels { 20 }; a", errmsg{"identifier not found: a"}},
	})
}

func TestYoyoExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"i := 0; yoyo i < 5 { i = i + 1 }", 5},
		{"sum := 0; i := 1; yoyo i < 5 { sum = sum + i; i = i + 1 }; sum", 10},
		{"sum := 0; i := 1; yoyo i < 5 { sum += i; i += 1 }; sum", 10},
		{"i := 1; yoyo false { i = 69 }; i", 1},
		{"i := 0; yoyo i < 5 { i += 1; yif i == 2 { yeet 69 }; -1 }", 69},
	})
}

func TestYallExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"yall [1, 2, 3] { yt }", 3},
		{"arr := [1, 2, 3]; yall arr { yt }", 3},
		{"sum := 0; yall [1, 2, 3] { sum = sum + yt }; sum", 6},
		{"sum := 0; yall [1, 2, 3] { sum += yt }; sum", 6},

		{`s := ""; yall "test" { s = s + yt + "-" }; s`, "t-e-s-t-"},
		{`s := ""; yall "test" { s += yt + "-" }; s`, "t-e-s-t-"},
		{`yall "testme" { yt }`, "e"},
		{`my_str := "swag"; yall my_str { yt }`, "g"},

		{`yall 0..5 { yt }`, 5},
		{`arr := []; yall 0..5 { arr << yt }; arr`, []int64{0, 1, 2, 3, 4, 5}},
		{`yall 4..4 { yt }`, 4},
		{`sum := 0; yall 1..4 { sum += yt }; sum`, 10},
		{`yall i: 0..5 { i }`, 5},
		{`sum := 0; yall j: 1..4 { sum += j }; sum`, 10},

		{`sum := 0; yall 5 { sum += 1 }; sum`, 6},
		{`sum := 0; yall 5 { sum += yt }; sum`, 15},
		{`arr := []; yall 5 { arr << yt }; arr`, []int64{0, 1, 2, 3, 4, 5}},
		{`sum := 0; yall -5 { sum += 1 }; sum`, 6},
		{`sum := 0; yall -5 { sum += yt }; sum`, -15},
		{`arr := []; yall -5 { arr << yt }; arr`, []int64{-5, -4, -3, -2, -1, 0}},

		// early exit
		{`yall 1..4 { yif yt == 1 { yeet 69 }; -1 }`, 69},
		{`yall 4..1 { yif yt == 3 { yeet 69 }; -1 }`, 69},
		{`yall [1, 2, 3] { yif yt == 1 { yeet 69 }; -1 }`, 69},
		{`yall "testme" { yif yt == "t" { yeet 69 }; -1 }`, 69},

		// scope leaking
		{`yall 0..5 { x }`, errmsg{"identifier not found: x"}},
		{`yall i: 0..5 { yt }`, errmsg{"identifier not found: yt"}},
	})
}

func TestYeetStatements(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"yeet 10;", 10},
		{"yeet 10; 9;", 10},
		{"yeet 2 * 5; 9;", 10},
		{"9; yeet 2 * 5; 9;", 10},
		{
			`
yif 10 > 1 {
	yif 10 > 1 {
		yeet 10;
	}
	yeet 1;
}`,
			10,
		},
	})
}

func TestRangeLiterals(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"0..5", rng{0, 5}},
		{"5..0", rng{5, 0}},
		{"(2+2)..(5*2)", rng{4, 10}},
		{"3 + 2 * 2 .. 5 - 2 * 2", rng{7, 1}},
		{"-5..-2", rng{-5, -2}},
		{"a := 1; b := 8; a..b", rng{1, 8}},
		{"a := 1; b := 8; (a+3)..(b-1)", rng{4, 7}},
		{"a := 1; b := 8; a+3 .. b-1", rng{4, 7}},
		{"r := 1+3 .. 9-1; r", rng{4, 8}},
	})
}

func TestDeclareExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"a := 8", 8},
		{"a := 8; a", 8},
		{"a := 8 * 5 + 3 / 2 - 2 * (2 + 3) * 3", 11},
		{"a := 8; b := a", 8},
		{"a := 8; b := a + 2", 10},
		{"a := 8; b := 2; c := a + b", 10},
		{"a := 8 == 5", false},
		{"a := 8 != 5", true},
		{"a := 8 > 5", true},
		{"a := 8; a = 15", 15},
		{"a := 8; b := 2; a = b", 2},
		{"a := b := c := 8; a + b + c", 24},
		{"a := b", errmsg{"identifier not found: b"}},
	})
}

func TestAssignExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"x := 8; x += 2; x", 10},
		{"x := 8; x -= 2; x", 6},
		{"x := 8; x *= 2; x", 16},
		{"x := 8; x /= 2; x", 4},
		{"x := 8; x %= 5; x", 3},

		{"x = 8", errmsg{"identifier not found: x (to declare a variable use := operator)"}},
		{"x += 8", errmsg{"identifier not found: x"}},

		{"a := [1, 2, 3]; a[1] = 69; a", []int64{1, 69, 3}},
		{"a := [1, 2, 3]; a[8] = 69", errmsg{"attempted to assign out of bounds for array 'a'"}},

		{`h := %{ "a": 1 }; h["a"] = 2; h["a"]`, 2},
		{`h := %{ "a": 1 }; h["b"] = 2; h["b"]`, 2},
		{`h := %{ "a": 1 }; h[[]] = 2`, 2},

		{`s := "yeet"; s[1] = "z"; s`, "yzet"},
		{`s := "yeet"; s[69] = "z"`, errmsg{"attempted to assign out of bounds for string 's'"}},
	})
}

func TestBlockExpressions(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{"{ 5 }", 5},
		{"{ a := 6; b := 9; a + b }", 15},
		{"x := { 5 }; x", 5},
		{"x := { a := 6; b := 9; a + b }; x", 15},

		{"x := { a := 6; a }; a", errmsg{"identifier not found: a"}},
	})
}

func TestLambdaApplication(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{`nope := \ { 69 }; nope();`, 69},
		{`identity := \x { x; }; identity(5);`, 5},
		{`identity := \x { yeet x; }; identity(5);`, 5},
		{`double := \x { x * 2; }; double(5);`, 10},
		{`add := \x, y { x + y; }; add(5, 5);`, 10},
		{`add := \x, y { x + y; }; add(5 + 5, add(5, 5));`, 20},
		{`add := \x, y, z { x + y + z; }; add(1, 2, 3);`, 6},

		{`add := \x y { x + y; }; add(5, 5);`, 10},
		{`add := \x y { x + y; }; add(5+5, add(5, 5));`, 20},
		{`add := \x y z { x + y + z }; add(1, 2, 3);`, 6},

		{`\x { x }(5)`, 5},
		{`\x, y { x * y }(3, 5)`, 15},
		{`\x y { x * y }(3, 5)`, 15},
		{`\x y { x * y }(3+2, 5+1)`, 30},
		{`\x y { x * y }(3 + 2, 5 + 1)`, 30},
	})
}

func TestRecursion(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{
			`
fib := \n {
    yif n < 2 { n } yels { fib(n-1) + fib(n-2) }
};
[fib(0), fib(1), fib(2), fib(3), fib(4), fib(8)]
`,
			[]int64{0, 1, 1, 2, 3, 21},
		},
		{
			`
factorial := \n { 
    yif n == 0 { 1 } yels { n * factorial(n-1) }
};
[factorial(0), factorial(1), factorial(5)]
`,
			[]int64{1, 1, 120},
		},
	})
}

func TestClosures(t *testing.T) {
	runEvalTests(t, []evalTestCase{
		{
			`
newAdder := \x { 
    \n { x + n } 
}
addTwo := newAdder(2)
addTwo(2)`,
			4,
		},
		{
			`
newGenerator := \ {
    n := 0
    \ { n = n + 2 }
}
gen := newGenerator()
gen() + gen() + gen()`,
			12, // 2 + 4 + 6
		},
	})
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"yif (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
yif (10 > 1) {
	yif (10 > 1) {
		yeet true + false;
	}
	yeet 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			"foobar()",
			"identifier not found: foobar",
		},
		{
			"foobar2(x, y)",
			"identifier not found: foobar2",
		},
		{
			`f := \a b { a }; f()`,
			"wrong number of args for f (got 0, want 2)",
		},
		{
			`fn := \a b { a }; fn(5)`,
			"wrong number of args for fn (got 1, want 2)",
		},
		{
			`f := \a b { a }; f(5, 6, 7)`,
			"wrong number of args for f (got 3, want 2)",
		},
		{
			`fun := \{ 8 }; fun(5)`,
			"wrong number of args for fun (got 1, want 0)",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		if err := testErrorObject(evaluated, tt.expectedMessage); err != nil {
			t.Error(err)
		}
	}
}

//
// PROGRAMS FROM EXAMPLES/
//

const examplesDir = "../examples"

func TestExampleFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing files in short mode")
	}

	testFiles, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("couldn't read example files dir: %s", err)
	}

	for _, f := range testFiles {
		t.Run(f.Name(), func(t *testing.T) {
			f := f
			t.Parallel()

			filename := filepath.Join(examplesDir, f.Name())
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("couldn't read test file: %s", err)
			}

			result := testEval(t, string(src))
			if evalError, ok := result.(*object.Error); ok {
				t.Errorf("runtime error: %q", evalError.Msg)
			}
		})
	}
}

//
// BENCHMARKS
//

func BenchmarkEval(b *testing.B) {
	testFiles, err := os.ReadDir(examplesDir)
	if err != nil {
		b.Fatalf("couldn't read example files dir: %s", err)
	}

	for _, f := range testFiles {
		b.Run(f.Name(), func(b *testing.B) {
			b.StopTimer()
			filename := filepath.Join(examplesDir, f.Name())
			src, err := os.ReadFile(filename)
			if err != nil {
				b.Fatalf("couldn't read test file: %s", err)
			}

			sSrc := string(src)
			l := lexer.New(sSrc)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				_ = eval.Eval(program, env)
			}
		})
	}
}

//
// HELPERS
//

type errmsg struct {
	msg string
}

type rng struct {
	start int64
	end   int64
}

func runEvalTests(t *testing.T, tests []evalTestCase) {
	t.Helper()

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch expected := tt.expected.(type) {
		case int:
			if err := testIntegerObject(evaluated, int64(expected)); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case []int64:
			if err := testIntegerArray(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case float64:
			if err := testNumberObject(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case []float64:
			if err := testNumberArray(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case bool:
			if err := testBooleanObject(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case string:
			if err := testStringObject(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case errmsg:
			if err := testErrorObject(evaluated, expected.msg); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case rng:
			if err := testRangeObject(evaluated, expected); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		case nil:
			if err := testNullObject(evaluated); err != nil {
				t.Errorf("%s (%s)", err, tt.input)
			}

		default:
			t.Errorf("unexpected type, got %T", expected)
		}
	}
}

func testEval(t *testing.T, input string) object.Object {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Error(err)
		}
		t.Fatalf("parsing errors, bailing")
	}

	env := object.NewEnvironment()

	return eval.Eval(program, env)
}

func testIntegerObject(obj object.Object, expected int64) error {
	result, ok := obj.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got %T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("Integer object has wrong value. got %d, want %d", result.Value, expected)
	}
	return nil
}

func testIntegerArray(obj object.Object, expected []int64) error {
	result, ok := obj.(*object.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got %T (%+v)", obj, obj)
	}
	if len(result.Elements) != len(expected) {
		return fmt.Errorf("Array has wrong number of elements. got %d, want %d",
			len(result.Elements), len(expected))
	}
	for i := range expected {
		if err := testIntegerObject(result.Elements[i], expected[i]); err != nil {
			return err
		}
	}
	return nil
}

func testNumberObject(obj object.Object, expected float64) error {
	result, ok := obj.(*object.Number)
	if !ok {
		return fmt.Errorf("object is not Number. got %T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("Number object has wrong value. got %g, want %g", result.Value, expected)
	}
	return nil
}

func testNumberArray(obj object.Object, expected []float64) error {
	result, ok := obj.(*object.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got %T (%+v)", obj, obj)
	}
	for i := range expected {
		if err := testNumberObject(result.Elements[i], expected[i]); err != nil {
			return err
		}
	}
	return nil
}

func testStringObject(obj object.Object, expected string) error {
	result, ok := obj.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got %T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("String has wrong value. want %q, got %q", expected, result.Value)
	}
	return nil
}

func testRangeObject(obj object.Object, expectedRng rng) error {
	result, ok := obj.(*object.Range)
	if !ok {
		return fmt.Errorf("object is not Range. got %T (%+v)", obj, obj)
	}
	if result.Start != expectedRng.start {
		return fmt.Errorf("wrong range start. want %d, got %d", expectedRng.start, result.Start)
	}
	if result.End != expectedRng.end {
		return fmt.Errorf("wrong range end. want %d, got %d", expectedRng.end, result.End)
	}
	return nil
}

func testBooleanObject(obj object.Object, expected bool) error {
	result, ok := obj.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got %T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. want %t, got %t", expected, result.Value)
	}
	return nil
}

func testNullObject(obj object.Object) error {
	if obj != object.NULL {
		return fmt.Errorf("object is not NULL. got %T (%+v)", obj, obj)
	}
	return nil
}

func testErrorObject(obj object.Object, expectedMsg string) error {
	result, ok := obj.(*object.Error)
	if !ok {
		return fmt.Errorf("object is not Error. got %T (%+v)", obj, obj)
	}
	if result.Msg != expectedMsg {
		return fmt.Errorf("wrong error message. want %q, got %q", expectedMsg, result.Msg)
	}
	return nil
}
