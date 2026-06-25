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

}
