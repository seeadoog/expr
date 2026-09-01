package expr

import (
	"testing"
)

//func TestNewLib(t *testing.T) {
//	lib := NewMathLib("")
//
//	RegisterOptFuncDefine2(lib, "add", func(ctx *Context, a float64, b float64, opt *Options) any {
//
//		return a + b
//	})
//
//	e := NewEnv()
//	e.AddLib(lib)
//
//	val, err := e.ParseValue(`math.log10(10)===4`)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ctx := e.NewContext(nil)
//
//	fmt.Println(ctx.ExecValue(val))
//}

func TestEqT(t *testing.T) {
	exp, err := DefaultEnv.ParseValue(`
a = 1; b = '1';  

c1 = a === b ; 
c2 = a !== b;
c3 = a === d ;

b1 = b === a ; 
b2 = b !== a;

d1 = a == b ;

# ---------- nil ----------
n1 = nil === nil;
n2 = nil === '';
n3 = nil === 0;
n4 = nil === false;
n5 = nil !== 1;

# ---------- bool ----------
t1 = true === 1;
t2 = false === 0;
t3 = true === 'true';
t4 = false === 'false';
t5 = true !== 0;

# ---------- string 无法解析 ----------
s1 = 'abc' === 1;
s2 = 'abc' !== 1;

# ---------- float ----------
f1 = 1.0 === 1;
f2 = 1.5 === '1.5';
f3 = 1.5 !== '1.6';

# ---------- 混合优先级（容易出坑） ----------
m1 = '1' === true;
m2 = '0' === false;

# ---------- 对称性 ----------
sym1 = 1 === '1';
sym2 = '1' === 1;
sym3 = 1 !== '2';
sym4 = '2' !== 1;
`)
	if err != nil {
		panic(err)
	}

	ctx := DefaultEnv.NewContext(nil)
	ctx.SetByString("d", 1)

	ctx.ExecValue(exp)

	// ---------- 原有 ----------
	assertEqual(t, ctx, "c1", true)
	assertEqual(t, ctx, "c2", false)
	assertEqual(t, ctx, "c3", true)

	assertEqual(t, ctx, "b1", true)
	assertEqual(t, ctx, "b2", false)
	assertEqual(t, ctx, "d1", false)

	// ---------- nil ----------
	assertEqual(t, ctx, "n1", true)
	assertEqual(t, ctx, "n2", true)
	assertEqual(t, ctx, "n3", true)
	assertEqual(t, ctx, "n4", true)
	assertEqual(t, ctx, "n5", true)

	// ---------- bool ----------
	assertEqual(t, ctx, "t1", true)
	assertEqual(t, ctx, "t2", true)
	assertEqual(t, ctx, "t3", true)
	assertEqual(t, ctx, "t4", true)
	assertEqual(t, ctx, "t5", true)

	// ---------- string 解析失败 ----------
	assertEqual(t, ctx, "s1", false)
	assertEqual(t, ctx, "s2", true)

	// ---------- float ----------
	assertEqual(t, ctx, "f1", true)
	assertEqual(t, ctx, "f2", true)
	assertEqual(t, ctx, "f3", true)

	// ---------- 混合 ----------
	assertEqual(t, ctx, "m1", true)
	assertEqual(t, ctx, "m2", true)

	// ---------- 对称性 ----------
	assertEqual(t, ctx, "sym1", true)
	assertEqual(t, ctx, "sym2", true)
	assertEqual(t, ctx, "sym3", true)
	assertEqual(t, ctx, "sym4", true)
}

/*

db.query('')


*/
