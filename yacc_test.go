package expr

import (
	"strings"
	"testing"
)

func BenchmarkPost3(b *testing.B) {
	v := strings.Repeat("1==1 && ", 0)
	e, err := DefaultEnv.ParseValue(v + "1==1 || 1==1 || 1==1 || 1 == 1 ")
	if err != nil {
		b.Fatal(err)
	}
	//ss := stack[any]{}
	ctx := DefaultEnv.NewContext(nil)
	ctx.SetByString("a", 1.0)
	ctx.SetByString("b", 3.0)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e.Val(ctx)
	}
}

func sw(i int) int {
	switch i {
	case 0:
		return 1
	case 2:
		return 3
	case 1:
		return 2
	case 3:
		return 5
	case 4:
		return 3
	case 5:
		return 0
	case 6:
		return 7
	case 7:
		return 8
	case 8:
		return 7

	case 9:
		return 1
	default:
		panic("xx")
	}
}

// kl,k
func BenchmarkSw(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := DefaultEnv.GetContextFromPool()

		DefaultEnv.PutContext2Pool(ctx)
	}
}
