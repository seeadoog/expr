package expr

import (
	"fmt"
	"reflect"
	"testing"
)

func BenchmarkValueOf(b *testing.B) {

	b.ReportAllocs()
	v := ValueOf(make([]int, 128))
	//arr := v.AnyArr()
	for i := 0; i < b.N; i++ {
		v.RangeArr(func(k int, v any) bool {
			return true
		})
	}
}

// v0.2.0-rc11
func TestSwitchRangeReturn(t *testing.T) {
	e := mustParse(`
const {
	cc:{
		a: 'xxxx.xx',
		b: 'xx.xx.xx' 
	}
}

`)
	c := DefaultEnv.NewContext(nil)

	res := c.ExecValue(e)
	fmt.Println(res)

	v := ValueOf(res)
	vv := v.GetV(c, "cc").GetV(c, "a").String()
	fmt.Println(reflect.TypeOf(vv), vv)

}

func BenchmarkArgSet(b *testing.B) {
	var (
		ka = DefaultEnv.NewHashKey("a").Hash
		kb = DefaultEnv.NewHashKey("b").Hash
		kc = DefaultEnv.NewHashKey("c").Hash
		kd = DefaultEnv.NewHashKey("d").Hash
		ke = DefaultEnv.NewHashKey("e").Hash
		kf = DefaultEnv.NewHashKey("f").Hash
	)

	c := DefaultEnv.NewContext(nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Set(ka, i)
		c.Set(kb, i)
		c.Set(kc, i)
		c.Set(kd, i)
		c.Set(ke, i)
		c.Set(kf, i)
		c.Get(ka)
		c.Get(kb)
		c.Get(kc)
		c.Get(kd)
		c.Get(ke)
		c.Get(ke)
		c.Get(ke)
		c.Get(ke)
		c.Get(ke)
	}
	fmt.Println(kf)
}
