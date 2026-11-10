package expr

import (
	"fmt"
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

//v0.2.0-rc11
func TestSwitchRangeReturn(t *testing.T) {
	e := mustParse(`

switch 11
case 0 ~ 10:
	return(1)
case 10 ~ 20:
	return(2)
end;
return(10);
`)
	c := DefaultEnv.NewContext(nil)

	res := c.ExecValue(e)
	fmt.Println(res)

}
