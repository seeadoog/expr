package expr

import (
	"reflect"
	"testing"
)

func conv(t *testing.T, typ reflect.Type, v any) any {
	t.Helper()
	rv, ok := structValConvert(nil, typ, v)
	if !ok {
		t.Fatalf("convert failed: %v -> %s", v, typ)
	}
	return rv.Interface()
}

func TestConvertBasic(t *testing.T) {
	if got := conv(t, reflect.TypeOf(""), 123); got != "123" {
		t.Fatalf("int->string got %v", got)
	}
	if got := conv(t, reflect.TypeOf(int(0)), "42"); got != 42 {
		t.Fatalf("string->int got %v", got)
	}
	if got := conv(t, reflect.TypeOf(float64(0)), "3.14"); got != 3.14 {
		t.Fatalf("string->float got %v", got)
	}
	if got := conv(t, reflect.TypeOf(true), "true"); got != true {
		t.Fatalf("string->bool got %v", got)
	}
}

func TestConvertIntPrecision(t *testing.T) {
	// 大整数不能走 float64 损失精度
	big := int64(9007199254780000)
	if got := conv(t, reflect.TypeOf(int64(0)), big); got != big {
		t.Fatalf("int64 precision lost: got %v want %v", got, big)
	}
}

func TestConvertPtr(t *testing.T) {
	got := conv(t, reflect.TypeOf((*int)(nil)), 7)
	p, ok := got.(*int)
	if !ok || p == nil || *p != 7 {
		t.Fatalf("to *int got %#v", got)
	}
	// 源为指针，目标为值
	x := 9
	got2 := conv(t, reflect.TypeOf(int(0)), &x)
	if got2 != 9 {
		t.Fatalf("*int->int got %v", got2)
	}
}

type inner struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type outer struct {
	ID    int `json:"id"`
	Inner inner
	Tags  []string
	Attrs map[string]int
}

func TestConvertStructNested(t *testing.T) {
	src := map[string]any{
		"id": 5,
		"Inner": map[string]any{
			"name": "bob",
			"age":  float64(30),
		},
		"Tags":  []any{"a", "b"},
		"Attrs": map[string]any{"x": 1, "y": 2},
	}
	got := conv(t, reflect.TypeOf(outer{}), src).(outer)
	if got.ID != 5 || got.Inner.Name != "bob" || got.Inner.Age != 30 {
		t.Fatalf("nested struct wrong: %#v", got)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "a" {
		t.Fatalf("slice wrong: %#v", got.Tags)
	}
	if got.Attrs["x"] != 1 || got.Attrs["y"] != 2 {
		t.Fatalf("map wrong: %#v", got.Attrs)
	}
}

func TestConvertStructToPtr(t *testing.T) {
	src := map[string]any{"id": 1}
	got := conv(t, reflect.TypeOf((*outer)(nil)), src)
	p, ok := got.(*outer)
	if !ok || p == nil || p.ID != 1 {
		t.Fatalf("to *struct got %#v", got)
	}
}

func TestConvertTypedCollections(t *testing.T) {
	// 真实 typed slice / map 而非 []any / map[string]any
	got := conv(t, reflect.TypeOf([]int{}), []string{"1", "2", "3"})
	s := got.([]int)
	if len(s) != 3 || s[2] != 3 {
		t.Fatalf("typed slice got %#v", s)
	}
	got2 := conv(t, reflect.TypeOf(map[string]int{}), map[string]string{"a": "10"})
	m := got2.(map[string]int)
	if m["a"] != 10 {
		t.Fatalf("typed map got %#v", m)
	}
}

func TestConvertNamedType(t *testing.T) {
	type MyInt int
	got := conv(t, reflect.TypeOf(MyInt(0)), 5)
	if got.(MyInt) != 5 {
		t.Fatalf("named int got %v", got)
	}
}

func TestConvertInterface(t *testing.T) {
	var anyT = reflect.TypeOf((*any)(nil)).Elem()
	got := conv(t, anyT, "hello")
	if got != "hello" {
		t.Fatalf("to any got %v", got)
	}
	// 非空接口未实现应失败
	type Stringer interface{ String() string }
	st := reflect.TypeOf((*Stringer)(nil)).Elem()
	if _, ok := structValConvert(nil, st, 123); ok {
		t.Fatalf("int should not satisfy Stringer")
	}
}

func TestConvertNil(t *testing.T) {
	rv, ok := structValConvert(nil, reflect.TypeOf((*int)(nil)), nil)
	if !ok || !rv.IsNil() {
		t.Fatalf("nil->*int should be nil, got %#v ok=%v", rv, ok)
	}
}

func TestConvertBytes(t *testing.T) {
	got := conv(t, reflect.TypeOf([]byte(nil)), "abc")
	if string(got.([]byte)) != "abc" {
		t.Fatalf("string->[]byte got %v", got)
	}
}

func TestConvertArray(t *testing.T) {
	got := conv(t, reflect.TypeOf([3]int{}), []any{1, 2, 3})
	a := got.([3]int)
	if a[0] != 1 || a[2] != 3 {
		t.Fatalf("array got %#v", a)
	}
}
