package expr

import (
	"fmt"
	"hash/crc64"
	"strconv"
	"sync"
	"testing"
	"unicode/utf8"
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

func BenchmarkHash(b *testing.B) {
	bs := []byte("hello world ")

	fmt.Println(xxhash.Sum64(bs))
	fmt.Println(xxhash.Sum64(bs))
	for i := 0; i < b.N; i++ {
		xxhash.Sum64(bs)
	}
}

func BenchmarkHash2(b *testing.B) {
	bs := []byte("hello world ")
	for i := 0; i < b.N; i++ {
		crc64.Checksum(bs, table)
	}
}

func BenchmarkHash3(b *testing.B) {

	b.ReportAllocs()

	m := make(map[string]any, 64)
	for i := 0; i < 32; i++ {
		m[strconv.Itoa(i)] = i
	}
	c := DefaultEnv.NewContext(m)
	for i := 0; i < b.N; i++ {
		c.GetByString("32")
	}
}

func BenchmarkFas(b *testing.B) {
	fm := newFuncMap(64)
	for i := 0; i < 32; i++ {
		fm.puts(strconv.Itoa(i), nil)
	}
	k := calcHash("20")
	for i := 0; i < b.N; i++ {
		fm.get(k)
	}
}
func BenchmarkMap3(b *testing.B) {
	fm := newMap2(256)
	for i := 0; i < 32; i++ {
		fm.put(strconv.Itoa(i), i)
	}

	for i := 0; i < b.N; i++ {
		fm.get("32")
	}
}

type map2Elem struct {
	key     string
	val     any
	keyHash uint64
}

type map2 struct {
	data [][]*map2Elem
	mod  uint64
}

func newMap2(cap uint64) *map2 {
	return &map2{
		data: make([][]*map2Elem, cap),
		mod:  cap - 1,
	}
}

func (m *map2) put(key string, val any) {
	idx := xxhash.Sum64(ToBytes(key)) & m.mod

	for _, v := range m.data[idx] {
		if v.key == key {
			v.val = val
			return
		}
	}
	m.data[idx] = append(m.data[idx], &map2Elem{
		key:     key,
		val:     val,
		keyHash: calcHash(key),
	})
}

func (m *map2) get(key string) any {
	idx := xxhash.Sum64(ToBytes(key)) & m.mod

	for _, v := range m.data[idx] {
		if v.key == key {
			return v.val
		}
	}
	return nil
}

func Test222(t *testing.T) {
}

func BenchmarkReset(b *testing.B) {
	m := DefaultEnv.NewContext(nil)

	fmt.Println(m.stack.mod)
	for i := 0; i < b.N; i++ {
		//for j := 0; j <= int(m.table.mod); j++ {
		//	m.Set(uint64(j), "", nil)
		//}
		m.Reset()
	}
	fmt.Println(m.stack.mod)
}

func BenchmarkReuseVm(b *testing.B) {
	p := sync.Pool{
		New: func() interface{} {
			return DefaultEnv.NewContext(nil)
		},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		d := p.Get().(*Context)

		d.SetByString("sdfsdf", 1)
		d.SetByString("dsfsdf", 1)

		d.Reset()
		p.Put(d)
	}
}

func BenchmarkParallel(b *testing.B) {
	p := sync.Pool{
		New: func() interface{} {
			return DefaultEnv.NewContext(nil)
		},
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			d := p.Get().(*Context)
			d.Reset()
			p.Put(d)
		}
	})
}

func TestSyncMap(t *testing.T) {

	m := Map[string, any]{}
	m.Store("name", 1)
	m.Store("age", 2)
	assertEqual2(t, m.Get("name"), 1)
	assertEqual2(t, m.Get("age"), 2)
	m.Store("age", 3)
	assertEqual2(t, m.Get("age"), 3)
	assertEqual2(t, m.Get("age"), 3)
	assertEqual2(t, m.Get("age"), 3)
	m.Delete("name")
	assertEqual2(t, m.Get("name"), nil)

}

// 假设你有这个构造函数
// func newMap() mapp

func TestMapBVT(t *testing.T) {
	m := Map[string, any]{}

	t.Run("Get on empty map", func(t *testing.T) {
		if v := m.Get("not-exist"); v != nil {
			t.Fatalf("expect nil, got %v", v)
		}
	})

	t.Run("Delete on empty map", func(t *testing.T) {
		m.Delete("not-exist")
	})

	t.Run("Basic Store and Get", func(t *testing.T) {
		m.Store("k1", "v1")

		if v := m.Get("k1"); v != "v1" {
			t.Fatalf("expect v1, got %v", v)
		}
	})

	t.Run("Store override existing value", func(t *testing.T) {
		m.Store("k2", 1)
		m.Store("k2", 2)

		if v := m.Get("k2"); v != 2 {
			t.Fatalf("expect 2, got %v", v)
		}
	})

	t.Run("Delete returns old value", func(t *testing.T) {
		m.Store("k3", "old")

		m.Delete("k3")

		if v := m.Get("k3"); v != nil {
			t.Fatalf("expect nil after delete, got %v", v)
		}
	})

	t.Run("Delete twice is safe", func(t *testing.T) {
		m.Store("k4", "v4")

		m.Delete("k4")
	})

	t.Run("Store after Delete", func(t *testing.T) {
		m.Store("k5", "v5")
		m.Delete("k5")
		m.Store("k5", "v5-new")

		if v := m.Get("k5"); v != "v5-new" {
			t.Fatalf("expect v5-new, got %v", v)
		}
	})

	t.Run("Value types", func(t *testing.T) {
		type S struct {
			A int
		}

		s := &S{A: 10}

		m.Store("int", 123)
		m.Store("struct", S{A: 1})
		m.Store("ptr", s)
		m.Store("nil", nil)

		if v := m.Get("int"); v != 123 {
			t.Fatalf("expect 123, got %v", v)
		}

		if v := m.Get("struct"); v.(S).A != 1 {
			t.Fatalf("unexpected struct value: %v", v)
		}

		if v := m.Get("ptr"); v.(*S).A != 10 {
			t.Fatalf("unexpected ptr value: %v", v)
		}

		if v := m.Get("nil"); v != nil {
			t.Fatalf("expect nil value, got %v", v)
		}
	})

	t.Run("Empty string key", func(t *testing.T) {
		m.Store("", "empty-key")

		if v := m.Get(""); v != "empty-key" {
			t.Fatalf("expect empty-key, got %v", v)
		}

		m.Delete("")
	})

	t.Run("Special characters key", func(t *testing.T) {
		key := "中文-key-!@#$%^&*()"
		m.Store(key, "ok")

		if v := m.Get(key); v != "ok" {
			t.Fatalf("expect ok, got %v", v)
		}
	})
}

var (
	_map = map[string]int{}
)

func BenchmarkMap4(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_map = map[string]int{
			"a": 1,
			"b": 1,
			"c": 1,
			"d": 1,
			"e": 1,
		}
	}
}

func BenchmarkHash6(b *testing.B) {
	for i := 0; i < b.N; i++ {

		xxhash.Sum64String("hello world")

	}
}

func BenchmarkHash67(b *testing.B) {
	data := []byte("hello world")
	for i := 0; i < b.N; i++ {
		crc64.Checksum(data, table)
	}
}

func mapLogic(m map[string]any) {
	chanVal, _ := m["chan"].(string)
	funcVal, _ := m["func"].(string)

	if chanVal == "iat" && funcVal == "cbm" {
		appid, _ := m["appid"].(string)
		useOld, _ := m["use_old"].(bool)

		switch appid {
		case "super":
			m["pass"] = true
		case "forbid":
			m["pass"] = false
		case "root":
			m["pass"] = useOld || false
		default:
			m["pass"] = false
		}
	}
}

func Benchmark_MapLogic(b *testing.B) {
	m := map[string]any{
		"chan":    "iat",
		"func":    "cbm",
		"appid":   "root",
		"use_old": true,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mapLogic(m)
	}

}

func TestUTF(t *testing.T) {

	s := "地方"
	e, i := utf8.DecodeRuneInString(s)

	fmt.Println(e)
	fmt.Println(s[i:])

	for i := 0; i < len(s); i++ {
		b := s[i]

		if b < 0x80 {

		} else {

		}
	}

}

func TestFastMapClone(t *testing.T) {

}

func BenchmarkExprGet(b *testing.B) {
	e, err := DefaultEnv.ParseValue(`handler = func(e) e.ret != 0 end`)
	if err != nil {
		b.Fatal(err)
	}
	spans := map[string]any{
		"sp1": map[string]any{
			"ret": 0.0,
		},
		"sp2": map[string]any{
			"ret": 0.0,
		},
		"sp3": map[string]any{
			"ret": 0.0,
		},
	}
	c := DefaultEnv.NewContext(nil)
	c.SetByString("b", 1.0)
	c.SetByString("spans", spans)

	c.ExecValue(e)
	fmt.Println(DefaultEnv.NewHashKey("___"))

	handler := c.GetByString("handler").(*LambdaVal)

	RunLambda(c, handler)
	for i := 0; i < b.N; i++ {
		for _, sp := range spans {
			RunLambda(c, handler, sp)
		}
	}
}

func BenchmarkExprGet2(b *testing.B) {
	e, err := DefaultEnv.ParseValue(`
do_set = func() 
	a = 1 ;
    b = 2; 
    c = 3 ;
end;

call(do_set);
`)
	if err != nil {
		b.Fatal(err)
	}
	spans := map[string]any{
		"sp1": map[string]any{
			"ret": 0.0,
		},
		"sp2": map[string]any{
			"ret": 0.0,
		},
		"sp3": map[string]any{
			"ret": 0.0,
		},
	}
	spans2 := []any{
		map[string]any{
			"ret": 0.0,
		},
		map[string]any{
			"ret": 0.0,
		},
		map[string]any{
			"ret": 0.0,
		},
		map[string]any{
			"ret": 0.0,
		},
		map[string]any{
			"ret": 0.0,
		},
		map[string]any{
			"ret": 0.0,
		},
	}
	c := DefaultEnv.NewContext(nil)

	c.SetByString("b", 1.0)
	c.SetByString("spans", spans)
	c.SetByString("spans2", spans2)

	c.ExecValue(e)
	fmt.Println(DefaultEnv.NewHashKey("___"))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.ExecValue(e)

	}
}

type EV struct {
	T unsafe.Pointer
	V unsafe.Pointer
}

func (e *EV) Number() float64 {
	return float64(uintptr(e.V))
}

func (e *EV) Interface() interface{} {

	switch e.T {
	case intType:
		return int(uintptr(e.V))
	case floatType:
		return float64(uintptr(e.V))
	default:
		return *(*any)(unsafe.Pointer(e))
	}
}

func EValueOf(v interface{}) EV {
	switch v := v.(type) {
	case int:
		return EIntOf(v)
	case float64:
		return EFloat(v)
	}
	return *(*EV)(unsafe.Pointer(&v))
}

func evalueOf(v interface{}) EV {
	return *(*EV)(unsafe.Pointer(&v))
}

var (
	intType   = evalueOf(int(1)).T
	floatType = evalueOf(float64(1)).T
)

func EIntOf(v int) EV {
	ptr := EV{
		T: intType,
		V: unsafe.Pointer(uintptr(v)),
	}
	return ptr
}
func EFloat(v float64) EV {
	ptr := EV{
		T: floatType,
		V: unsafe.Pointer(uintptr(v)),
	}
	return ptr
}

func Add(a, b EV) EV {
	return EFloat(a.Number() + b.Number())
}
