package expr

import (
	"context"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	ic := 0
	c := NewInstanceCache[string, any, any](func(ctx context.Context, k string, c any) (v any, err error) {
		ic++
		return ic, nil
	})

	var v any
	var err error
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err = c.Get(context.Background(), "1", nil)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	wg.Wait()

	assertEqual2(t, v, 1)
	assertEqual2(t, err, nil)

}

func execVal(k int, v Val, c *Context) any {
	switch k {
	case 0:
		RunLambda(c, v)
	case 1:
		return c.Get(v.(*variable).hash)
	case 2:
		return v.(*constraint).value
	}
	return v
}

func BenchmarkExecRaw(b *testing.B) {
	e, err := DefaultEnv.parseValueV(`1`)
	if err != nil {
		b.Fatal(err)
	}

	c := DefaultEnv.NewContext(nil)
	for i := 0; i < b.N; i++ {
		execVal(2, e, c)
	}
}

func BenchmarkExecRaw2(b *testing.B) {
	e, err := DefaultEnv.parseValueV(`a=1`)
	if err != nil {
		b.Fatal(err)
	}

	c := DefaultEnv.NewContext(nil)
	for i := 0; i < b.N; i++ {
		c.ExecValue(e)
	}
}
