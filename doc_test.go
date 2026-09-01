package expr

import (
	"fmt"
	"testing"
)

func TestDoc2(t *testing.T) {
	showDocOf("ctx", &Usr{})
	//fmt.Println(showDocOf("ctx.", &Usr{}))
}

func TestDOc3(t *testing.T) {

	e := NewEnv()
	n, err := e.ParseValueToAstNode("ctx.check_flow('aa','bb') && check_cnt('aa',bb+22)")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(n)
	//fmt.Println(showDocOf("", addV))
}

type V struct {
	typ     Type
	integer int
	str     string

	Name string
}

func (v V) Int() int {
	return v.integer
}

func (v V) String() string {
	return v.str
}

func addV(a V, b V) (r V) {
	r.integer = a.integer + b.integer
	return r
}

var (
	c V
)

func BenchmarkAddV2(bb *testing.B) {
	a := V{}
	b := V{}

	for i := 0; i < bb.N; i++ {

		c = addV(a, b)
	}
}

func TestRule(t *testing.T) {
	c := DefaultEnv.NewContext(nil)

	SelfDefine1(DefaultEnv, "adds", func(ctx *Context, self float64, a float64) float64 {
		return self + a
	})
	parseAndExec(DefaultEnv, `
map =  const {a:1, b:2,c:3,d:4,e:6};
for k,v in map do
	map2[k] = v;
end;
cc = 2;
dd = cc.adds(3);

arr = const [1,2,3,4];
arr[1]=6; #readonly not set
for i,v in arr do
	arr2[i] = v;
end;

`, c)
	assertEqual(t, c, "map.equals(map2)", true)
	assertEqual(t, c, "map2.equals(map)", true)
	assertEqual(t, c, "arr.equals(arr2)", true)
	assertEqual(t, c, "arr2.equals(arr)", true)
	assertEqual(t, c, "dd", 5.0)
	fmt.Println(DefaultEnv.NewHashKey("xx"))
}
