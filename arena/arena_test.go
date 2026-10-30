package arena

import (
	"fmt"
	"testing"

	"github.com/ortuman/nuke"
)

func BenchmarkNuke(b *testing.B) {
	a := nuke.NewMonotonicArena(1024*4, 1)
	for i := 0; i < b.N; i++ {

		nuke.New[string](a)
		a.Reset(false)
	}
}

func BenchmarkArena(b *testing.B) {
	a := &Arena{
		buf: make([]byte, 1024*4),
	}
	fmt.Println(a.StringToAny("asdfsf"))
	for i := 0; i < b.N; i++ {

		a.Float64ToAny(1.2)
		a.free()
	}
}
func TestArena(t *testing.T) {
	a := &Arena{
		buf: make([]byte, 1024),
	}
	fmt.Println(a.alloc(1))
	fmt.Println(a.alloc(1))
	fmt.Println(a.alloc(1))
	fmt.Println(a.alloc(1))
	fmt.Println(a.StringToAny("asdfsf"))

	b := a.MakeBytes(24)
	pf := New[float64](a)
	*pf = 4
	fmt.Println(cap(b))
	fmt.Println(len(b))

	fmt.Println(*pf)
	b = append(b, "helloworld"...)
	fmt.Println(*pf)

	c := a.MakeBytes(15)
	a.Append(c, []byte("324")...)
}
