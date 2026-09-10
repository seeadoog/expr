package expr

import "testing"

func BenchmarkValueOf(b *testing.B) {

	b.ReportAllocs()
	v := ValueOf(make([]int, 128))
	//arr := v.AnyArr()
	for i := 0; i < b.N; i++ {
		v.AnyArr()
	}
}
