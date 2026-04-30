package expr

import (
	"fmt"
	"testing"
)

func TestNewLib(t *testing.T) {
	lib := NewMathLib("")

	RegisterOptFuncDefine2(lib, "add", func(ctx *Context, a float64, b float64, opt *Options) any {

		return a + b
	})

	e := NewEnv()
	e.AddLib(lib)

	val, err := e.ParseValue(`math.log10(10)`)
	if err != nil {
		t.Fatal(err)
	}

	ctx := e.NewContext(nil)

	fmt.Println(ctx.ExecValue(val))
}

func BenchmarkLib(b *testing.B) {
	lib := NewLib("math")
	e := NewEnv()

	RegisterOptFuncDefine2(e, "add", func(ctx *Context, a float64, b float64, opt *Options) any {

		return nil
	})

	e.AddLib(lib)

	b.ReportAllocs()
	val, err := e.ParseValue(`add(1,2)`)
	if err != nil {
		b.Fatal(err)
	}
	ctx := e.NewContext(nil)

	fmt.Println(ctx.ExecValue(val))
	for i := 0; i < b.N; i++ {
		ctx.ExecValue(val)
	}

}

/*

 */
