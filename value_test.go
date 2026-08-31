package expr

import "testing"

func BenchmarkValueOf(b *testing.B) {

	b.ReportAllocs()
	v := ValueOf([]any{1, 2, 3, 4, 5, 6})
	//arr := v.AnyArr()
	for i := 0; i < b.N; i++ {
		v.AnyArr()
	}
}
