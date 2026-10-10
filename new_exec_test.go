package expr

import (
	"fmt"
	"reflect"
	"testing"
)

func BenchmarkExecNew(b *testing.B) {

	v, err := DefaultEnv.ParseValue(`
 a= 1;
if a == 1 then b = 1 end;
a= 1;
if a == 1 then b = 1 end;
a= 1;
if a == 1 then b = 1 end;
a= 1;
if a == 1 then b = 1 end ;
str = 'hello';

str2 = str + 'world';
str2 = str + 'world';
str2 = str + 'world';
`)
	if err != nil {
		b.Fatal(err)
	}
	c := DefaultEnv.NewContext(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ExecValue(v)
	}
}

var (
	resE EV
)

func BenchmarkEString(b *testing.B) {
	a1 := EFloat(1)
	a2 := EFloat(1)
	for i := 0; i < b.N; i++ {
		resE = Add(a1, a2)
	}
}

func TestReturnInLambda(t *testing.T) {

	c := DefaultEnv.NewContext(nil)
	parseAndExec(DefaultEnv, `fun  = func() 
return(0,"err")
end;
call(fun);
a = 1;
`, c)
	assertEqual(t, c, "a", nil)
}

func BenchmarkFuncObj(b *testing.B) {
	type ssr struct{}

	SelfDefine0(DefaultEnv, "gg", func(ctx *Context, self *ssr) any {
		return self
	})
	e, err := DefaultEnv.ParseValue(`number("123")`)
	if err != nil {
		b.Fatal(err)
	}
	c := DefaultEnv.NewContext(nil)
	c.SetByString("ssr", new(ssr))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.ExecValue(e)
	}
}

func BenchmarkSetCtx(b *testing.B) {
	as := DefaultEnv.NewHashKey("aaaa").Hash
	as2 := DefaultEnv.NewHashKey("bbbb").Hash
	as3 := DefaultEnv.NewHashKey("cccc").Hash
	as4 := DefaultEnv.NewHashKey("dddd").Hash
	c := DefaultEnv.NewContext(nil)
	fmt.Println(as, as2, as3, as4)
	for i := 0; i < b.N; i++ {
		c.Set(as, c.Get(as))

	}
}

func TestForRange(t *testing.T) {
	e := mustParse(`
for span in  const [{code:0,kvs:[{k:'aa',v:1}]},{code:0},{code:0},{code:0},{code:0}] do 
	for e in span.kvs do 
		switch e.k 
		case 'aa':
			return (e.v)
		end
	end
end;
`)
	c := DefaultEnv.NewContext(nil)
	res := c.ExecValue(e)
	if !reflect.DeepEqual(res, &Return{
		Var: []any{1.0},
	}) {
		t.Errorf("unexpected return %v,%T", res, res)
	}
}

func TestForRange2(t *testing.T) {
	e := mustParse(`
for span in  const [{code:0,kvs:[{k:'aa',v:1}]},{code:0},{code:0},{code:0},{code:0}] do 
	for e in span.kvs do 
		switch e.v 
		case 0~2:
			return (e.v)
		end
	end
end;
`)
	c := DefaultEnv.NewContext(nil)
	res := c.ExecValue(e)
	if !reflect.DeepEqual(res, &Return{
		Var: []any{1.0},
	}) {
		t.Errorf("unexpected return %v,%T", res, res)
	}
}
func TestForRange3(t *testing.T) {
	e := mustParse(`
for v in arr do 
	if v > 2 then
		return(1)
	end;
end;
`)
	c := DefaultEnv.NewContext(map[string]any{
		"arr": []float64{1.0, 2.0, 3.0},
	})
	res := c.ExecValue(e)
	if !reflect.DeepEqual(res, &Return{
		Var: []any{1.0},
	}) {
		t.Errorf("unexpected return %v,%T", res, res)
	}
}

// call()
func BenchmarkForRange(b *testing.B) {
	e := mustParse(`
for span in  const [
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},
{code:0},

] do 
	if span.code > 0 then
		return(1)
	end
end;
`)
	c := DefaultEnv.NewContext(nil)

	b.ReportAllocs()

	fmt.Println(c.ExecValue(e))
	fmt.Println(c.GetTable())

	for i := 0; i < b.N; i++ {
		c.ExecValue(e)
	}
}
func filter(arr []any, fn func(any) bool) any {
	for _, a := range arr {
		if !fn(a) {
			return a
		}
	}
	return nil
}
func BenchmarkForRangeGo(b *testing.B) {
	mm := []any{
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
		map[string]any{"code": 0.0},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		filter(mm, func(a any) bool {
			if NumberOf(a.(map[string]any)["code"]) > 0.0 {
				return false
			}
			return true
		})
	}
}

//15
