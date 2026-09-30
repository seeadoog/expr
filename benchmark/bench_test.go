package benchmark

import (
	"fmt"
	"testing"

	"github.com/seeadoog/expr"
	"github.com/seeadoog/jsonschema/v2/jsonpath"
)

var (
	DefaultEnv = expr.DefaultEnv
)
var (
	exp = "a.b.c.d.e.f.g"
)

func BenchmarkGETFromMap(b *testing.B) {

	exp2 := "$." + exp
	v, err := DefaultEnv.ParseValue("a=1")
	if err != nil {
		b.Fatal(err)
	}

	c := DefaultEnv.NewContext(nil)
	c.SetByJp(exp2, 1.0)
	b.ReportMetric(0, "MB/s")
	b.ReportMetric(1, "MB/s")
	fmt.Println(c.ExecValue(v))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ExecValue(v)
	}
}

func BenchmarkGETFromMap2(b *testing.B) {

	ced, err := jsonpath.Compile(exp)
	if err != nil {
		b.Fatal(err)
	}
	c := DefaultEnv.NewContext(nil)
	c.SetByJp("$."+exp, 1.0)
	b.ResetTimer()
	o := c.GetByString("$")
	b.ReportAllocs()
	fmt.Println(ced.Get(o))
	for i := 0; i < b.N; i++ {
		ced.Get(o)
	}
}

func BenchmarkGETFromMapRW(b *testing.B) {

	c := DefaultEnv.NewContext(nil)
	c.SetByJp("$."+exp, 1.0)
	b.ResetTimer()
	o := c.GetByString("$")
	b.ReportAllocs()
	rv := RowV{o}

	fmt.Println(rv.Get("a").Get("b").Get("c").Get("d").V())
	for i := 0; i < b.N; i++ {
		rv.Get("a").Get("b").Get("c").Get("d").Get("e").Get("f").Get("g").V()
	}
}

type RowV struct {
	v any
}

func (r RowV) Get(key string) RowV {
	switch v := r.v.(type) {
	case map[string]interface{}:
		return RowV{v[key]}
	}
	return RowV{nil}
}

func (r RowV) V() any {
	return r.v
}
