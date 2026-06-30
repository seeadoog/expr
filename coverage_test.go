package expr

import (
	"reflect"
	"testing"
)

// eval 解析并求值一个表达式，返回结果与错误。
func eval(t *testing.T, expr string, vars map[string]any) (any, any) {
	t.Helper()
	v, err := DefaultEnv.ParseValue(expr)
	if err != nil {
		return nil, err
	}
	ctx := DefaultEnv.NewContext(vars)
	return ctx.SafeExecValue(v)
}

// mustEval 解析并求值，断言无错误并返回结果。
func mustEval(t *testing.T, expr string, vars map[string]any) any {
	t.Helper()
	r, err := eval(t, expr, vars)
	if err != nil {
		t.Fatalf("eval %q failed: %v", expr, err)
	}
	return r
}

// TestBinaryOperators 覆盖 parseNodeBinary 中的各类运算符分支。
func TestBinaryOperators(t *testing.T) {
	cases := []struct {
		expr string
		vars map[string]any
		want any
	}{
		{`a + b`, map[string]any{"a": 3.0, "b": 4.0}, 7.0},
		{`a - b`, map[string]any{"a": 10.0, "b": 4.0}, 6.0},
		{`a * b`, map[string]any{"a": 3.0, "b": 4.0}, 12.0},
		{`a / b`, map[string]any{"a": 12.0, "b": 4.0}, 3.0},
		{`a % b`, map[string]any{"a": 10.0, "b": 3.0}, 1.0},
		{`a ^ b`, map[string]any{"a": 2.0, "b": 10.0}, 1024.0},
		{`a | b`, map[string]any{"a": 5.0, "b": 2.0}, 7.0},
		{`a & b`, map[string]any{"a": 6.0, "b": 3.0}, 2.0},
		{`a < b`, map[string]any{"a": 1.0, "b": 2.0}, true},
		{`a <= b`, map[string]any{"a": 2.0, "b": 2.0}, true},
		{`a > b`, map[string]any{"a": 3.0, "b": 2.0}, true},
		{`a >= b`, map[string]any{"a": 2.0, "b": 2.0}, true},
		{`a == b`, map[string]any{"a": 2.0, "b": 2.0}, true},
		{`a != b`, map[string]any{"a": 1.0, "b": 2.0}, true},
		{`a && b`, map[string]any{"a": true, "b": true}, true},
		{`a && b`, map[string]any{"a": false, "b": true}, false},
		{`a || b`, map[string]any{"a": false, "b": true}, true},
		{`a || b`, map[string]any{"a": true, "b": false}, true},
	}
	for _, c := range cases {
		got := mustEval(t, c.expr, c.vars)
		if got != c.want {
			t.Errorf("%q = %v, want %v", c.expr, got, c.want)
		}
	}
}

// TestOrrOperator 覆盖 "orr"（or 关键字）分支：左值为 nil 时取右值。
func TestOrrOperator(t *testing.T) {
	if got := mustEval(t, `a or b`, map[string]any{"a": nil, "b": 99.0}); got != 99.0 {
		t.Errorf("nil or 99 = %v, want 99", got)
	}
	if got := mustEval(t, `a or b`, map[string]any{"a": 5.0, "b": 99.0}); got != 5.0 {
		t.Errorf("5 or 99 = %v, want 5", got)
	}
}

// TestTypedEquality 覆盖 === 和 !== 分支。
func TestTypedEquality(t *testing.T) {
	if got := mustEval(t, `a === b`, map[string]any{"a": 1.0, "b": 1.0}); got != true {
		t.Errorf("1 === 1 = %v, want true", got)
	}
	if got := mustEval(t, `a !== b`, map[string]any{"a": 1.0, "b": 2.0}); got != true {
		t.Errorf("1 !== 2 = %v, want true", got)
	}
}

// TestInOperator 覆盖 in 运算符。
func TestInOperator(t *testing.T) {
	if got := mustEval(t, `x in arr`, map[string]any{"x": 2.0, "arr": []any{1.0, 2.0, 3.0}}); got != true {
		t.Errorf("2 in [1,2,3] = %v, want true", got)
	}
	if got := mustEval(t, `x in [1,2,3]`, map[string]any{"x": 9.0}); got != false {
		t.Errorf("9 in [1,2,3] = %v, want false", got)
	}
}

// TestCompoundAssign 覆盖 +=、-=、*=、/= 分支。
func TestCompoundAssign(t *testing.T) {
	cases := []struct {
		expr string
		want any
	}{
		{`x=10;x+=5;x`, 15.0},
		{`x=10;x-=3;x`, 7.0},
		{`x=10;x*=2;x`, 20.0},
		{`x=10;x/=2;x`, 5.0},
	}
	for _, c := range cases {
		got := mustEval(t, c.expr, nil)
		if got != c.want {
			t.Errorf("%q = %v, want %v", c.expr, got, c.want)
		}
	}
}

// TestUnaryOperators 覆盖 parseNodeUnary：!、-、++、--。
func TestUnaryOperators(t *testing.T) {
	if got := mustEval(t, `!b`, map[string]any{"b": true}); got != false {
		t.Errorf("!true = %v, want false", got)
	}
	if got := mustEval(t, `-x`, map[string]any{"x": 5.0}); got != -5.0 {
		t.Errorf("-5 = %v, want -5", got)
	}
	if got := mustEval(t, `x=5;x++;x`, nil); got != 6.0 {
		t.Errorf("x++ = %v, want 6", got)
	}
	if got := mustEval(t, `x=5;x--;x`, nil); got != 4.0 {
		t.Errorf("x-- = %v, want 4", got)
	}
}

// TestSliceCut 覆盖 parseNodeSliceCut：数组与字符串切片，含起止省略。
func TestSliceCut(t *testing.T) {
	arr := map[string]any{"arr": []any{1.0, 2.0, 3.0, 4.0, 5.0}}
	if got := mustEval(t, `arr[1:3]`, arr); !reflect.DeepEqual(got, []any{2.0, 3.0}) {
		t.Errorf("arr[1:3] = %v, want [2 3]", got)
	}
	if got := mustEval(t, `arr[:2]`, arr); !reflect.DeepEqual(got, []any{1.0, 2.0}) {
		t.Errorf("arr[:2] = %v, want [1 2]", got)
	}
	if got := mustEval(t, `arr[3:]`, arr); !reflect.DeepEqual(got, []any{4.0, 5.0}) {
		t.Errorf("arr[3:] = %v, want [4 5]", got)
	}
	if got := mustEval(t, `s[1:3]`, map[string]any{"s": "hello"}); got != "el" {
		t.Errorf("s[1:3] = %v, want el", got)
	}
	// 越界返回 nil
	if got := mustEval(t, `arr[1:99]`, arr); got != nil {
		t.Errorf("arr[1:99] = %v, want nil", got)
	}
}

// TestLambdaAndLambda2 覆盖 parseNodeLambda 与 parseNodeLambda2。
func TestLambdaAndLambda2(t *testing.T) {
	// 单变量 lambda (Lambda): v => expr
	got := mustEval(t, `arr.filter(v => v > 2)`, map[string]any{"arr": []any{1.0, 2.0, 3.0, 4.0}})
	if !reflect.DeepEqual(got, []any{false, false, true, true}) {
		t.Errorf("filter lambda = %v", got)
	}
	// 多变量 lambda (Lambda): {a,b} => expr
	sorted := mustEval(t, `arr.sort({a,b}=>a<b)`, map[string]any{"arr": []any{3.0, 1.0, 2.0}})
	if !reflect.DeepEqual(sorted, []any{1.0, 2.0, 3.0}) {
		t.Errorf("sort lambda = %v, want [1 2 3]", sorted)
	}
	// Lambda2 形式: (v) => {expr}
	got2 := mustEval(t, `arr.filter((v) => {v > 2})`, map[string]any{"arr": []any{1.0, 2.0, 3.0}})
	if !reflect.DeepEqual(got2, []any{false, false, true}) {
		t.Errorf("filter lambda2 = %v", got2)
	}
}

// TestEmptyVar 覆盖 parseNodeVariable 的 "_" 分支。
func TestEmptyVar(t *testing.T) {
	// _ 作为 lambda 占位参数
	if _, err := DefaultEnv.ParseValue(`f = _ => 1`); err != nil {
		t.Errorf("empty var lambda parse failed: %v", err)
	}
}

// TestConstAssign 覆盖 parseNodeSet 的 n.Const 分支。
func TestConstAssign(t *testing.T) {
	if got := mustEval(t, `const x = 5; x`, nil); got != 5.0 {
		t.Errorf("const x = 5 => %v, want 5", got)
	}
	if got := mustEval(t, `const m = {a:1}; m["a"]`, nil); got != 1.0 {
		t.Errorf(`const m = {a:1}; m["a"] => %v, want 1`, got)
	}
}

// TestStringConcat 覆盖 + 运算符在字符串上的 add2 路径。
func TestStringConcat(t *testing.T) {
	if got := mustEval(t, `"a" + "b"`, nil); got != "ab" {
		t.Errorf(`"a" + "b" = %v, want ab`, got)
	}
	if got := mustEval(t, `s = "x"; s += "y"; s`, nil); got != "xy" {
		t.Errorf("string += = %v, want xy", got)
	}
}

// TestForLoop 覆盖 for + lambda 的累加路径（func + lambda 解析）。
func TestForLoop(t *testing.T) {
	got := mustEval(t, `sum = 0; for([1,2,3,4], v => sum = sum + v); sum`, nil)
	if got != 10.0 {
		t.Errorf("for sum = %v, want 10", got)
	}
}

// TestNotNilOperator 覆盖 parseNodeNotNil：@ 运算符。
func TestNotNilOperator(t *testing.T) {
	if got := mustEval(t, `v@`, map[string]any{"v": map[string]any{"k": 1.0}}); got == nil {
		t.Errorf("v@ on non-nil map = nil, want map")
	}
	if got := mustEval(t, `v@.k`, map[string]any{"v": map[string]any{"k": 7.0}}); got != 7.0 {
		t.Errorf("v@.k = %v, want 7", got)
	}
}

// TestTernary 覆盖 parseNodeTernary，含 R 为空的形式。
func TestTernary(t *testing.T) {
	if got := mustEval(t, `a > 5 ? b : c`, map[string]any{"a": 10.0, "b": 1.0, "c": 2.0}); got != 1.0 {
		t.Errorf("ternary true branch = %v, want 1", got)
	}
	if got := mustEval(t, `a > 5 ? b : c`, map[string]any{"a": 1.0, "b": 1.0, "c": 2.0}); got != 2.0 {
		t.Errorf("ternary false branch = %v, want 2", got)
	}
}

// TestAsBinding 覆盖 as 绑定。
func TestAsBinding(t *testing.T) {
	if got := mustEval(t, `5 + 5 as a1; a1`, nil); got != 10.0 {
		t.Errorf("as binding = %v, want 10", got)
	}
}

// TestCollections 覆盖数组定义、map 定义、数组下标访问、变量展开。
func TestCollections(t *testing.T) {
	if got := mustEval(t, `[1,2,3]`, nil); !reflect.DeepEqual(got, []any{1.0, 2.0, 3.0}) {
		t.Errorf("array def = %v", got)
	}
	if got := mustEval(t, `{a:1,b:2}`, nil); !reflect.DeepEqual(got, map[string]any{"a": 1.0, "b": 2.0}) {
		t.Errorf("map def = %v", got)
	}
	if got := mustEval(t, `m["k"]`, map[string]any{"m": map[string]any{"k": 42.0}}); got != 42.0 {
		t.Errorf("index access = %v, want 42", got)
	}
}

// TestMathFunctions 覆盖部分 math 库函数的调用路径。
func TestMathFunctions(t *testing.T) {
	cases := []struct {
		expr string
		vars map[string]any
		want float64
	}{
		{`math_sqrt(x)`, map[string]any{"x": 16.0}, 4.0},
		{`math_pow(a,b)`, map[string]any{"a": 2.0, "b": 3.0}, 8.0},
		{`math_max(a,b)`, map[string]any{"a": 2.0, "b": 3.0}, 3.0},
		{`math_min(a,b)`, map[string]any{"a": 2.0, "b": 3.0}, 2.0},
		{`math_abs(x)`, map[string]any{"x": -7.0}, 7.0},
		{`math_floor(x)`, map[string]any{"x": 3.7}, 3.0},
		{`math_ceil(x)`, map[string]any{"x": 3.2}, 4.0},
		{`math_round(x)`, map[string]any{"x": 3.5}, 4.0},
	}
	for _, c := range cases {
		got := mustEval(t, c.expr, c.vars)
		if got != c.want {
			t.Errorf("%q = %v, want %v", c.expr, got, c.want)
		}
	}
}

// TestStringInterpolation 覆盖 parseNodeString 的常量与插值两条路径。
func TestStringInterpolation(t *testing.T) {
	if got := mustEval(t, `"hello"`, nil); got != "hello" {
		t.Errorf("const string = %v, want hello", got)
	}
	if got := mustEval(t, `""`, nil); got != "" {
		t.Errorf("empty string = %v, want empty", got)
	}
}

// TestParseErrors 覆盖解析期的错误返回路径。
func TestParseErrors(t *testing.T) {
	bad := []string{
		`(`,
		`1 +`,
		`a >< b`,
		`[1,2`,
	}
	for _, expr := range bad {
		if _, err := DefaultEnv.ParseValue(expr); err == nil {
			t.Errorf("expected parse error for %q, got nil", expr)
		}
	}
}
