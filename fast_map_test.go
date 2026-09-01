package expr

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkFS(b *testing.B) {
	f := newFuncMap(4)
	a := "hsd"
	aa := calcHash(a)
	f.put(aa, a, nil)
	for i := 0; i < b.N; i++ {
		//f.getS(aa, a)
		f.get(aa)
	}
}

func TestFF(t *testing.T) {
	return
	for _, datum := range objFuncMap.data {
		if len(datum) > 0 {
			fmt.Println("prt:", len(datum), len(datum[0].val.data), datum[0].val.size)
		}
	}
	objFuncMap.foreach(func(f *funcMap) bool {
		for _, datum := range f.data {
			if len(datum) > 0 {
				fmt.Println(len(datum))
			}
		}
		return true
	})
}

func TestFuncMap(t *testing.T) {
	f := newFuncMap(4)
	f.puts("a", nil)
	f.puts("b", nil)
	f.puts("c", nil)
	f.puts("d", nil)
	f.puts("e", nil)
	f.puts("f", nil)
	f.puts("g", nil)
	f.puts("h", nil)
	f.puts("i", nil)

	assertEqual2(t, f.mod, uint64(127))
	assertEqual2(t, f.size, (9))
}

func TestEnvMap(t *testing.T) {
	//m := newEnvMap(8)
	//
	//for i := 0; i < 10000; i++ {
	//	ss := strconv.Itoa(i) + "xxxadsf"
	//	ha := calcHash(ss)
	//	m.putHash(ha, ss, i)
	//}
	//confilct := make(map[int][]int)
	//for i, datum := range m.data {
	//	confilct[len(datum)] = append(confilct[len(datum)], i)
	//}
	//for i, i2 := range confilct {
	//	if i == 3 {
	//		fmt.Println(m.data[i2[2]][2].key)
	//	}
	//}

}

func BenchmarkEnvMap(b *testing.B) {
	//5477xxxadsf
	m := newEnvMap(8)

	for i := 0; i < 10000; i++ {
		ss := strconv.Itoa(i) + "xxxadsf"
		ha := calcHash(ss)
		m.putHash(ha, i)
	}
	b.ReportAllocs()

	ha := calcHash("7706xxxadsfsdf")
	m.putHashOnly(ha, nil)
	for i := 0; i < b.N; i++ {
		m.putHashOnly(ha, nil)
		//m.getHash(ha)

	}
}

func BenchmarkEnvMap2(b *testing.B) {
	//5477xxxadsf
	var m1 = make(map[string]any, 0)
	for i := 0; i < 10000; i++ {
		ss := strconv.Itoa(i) + "xxxadsf"
		m1[ss] = i
	}
	b.ReportAllocs()

	//ha := calcHash("7706xxxadsf")

	for i := 0; i < b.N; i++ {

		_ = m1["5477xxxadsf"]
	}
	fmt.Println(m1["x"])
}

func TestExpr(t *testing.T) {
	env := NewEnv()

	exp, err := env.ParseValue(" $.a = 1")
	if err != nil {
		panic(err)
	}

	c := env.NewContext(nil)

	fmt.Println(exp.Val(c))
	fmt.Println(c.GetByString("$"))
}

/*
route:[
{
	"define":{
		"domain_route":{
			"bm3.5":""
		}
	},
},
"domain_route[domain] or 'ent_domain'"
]

*/
// [[user adddsd] split: -1]
// service： channel
func BenchmarkExpr33(b *testing.B) {

	exp, err := DefaultEnv.ParseFromJSONStr(`
[
    {
      "if": "model_id == 'x57904567'",
      "then": [
        "rand_n(100) < (ase_x597_rate)*100? model_id_2 = 's322'  :_"
      ]
    }
  ]
`)
	if err != nil {
		panic(err)
	}
	c := DefaultEnv.NewContext(map[string]interface{}{
		"model_id":      "x57904567",
		"ase_x597_rate": 1,
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Exec(exp)
	}
	fmt.Println(c.GetByString("model_id_2"))
}

func TestRand100(t *testing.T) {
	e, err := DefaultEnv.ParseValue(`
rand_n(100) *rate * 100
`)
	if err != nil {
		panic(err)
	}

	fmt.Println(DefaultEnv.NewContext(map[string]interface{}{
		"rate": 0.001,
	}).ExecValue(e))

}

func BenchmarkRand33(b *testing.B) {
	e, err := DefaultEnv.ParseValue(`
rand_n(100) < (rate)*100? model_id = 's322'  :_
`)

	if err != nil {
		panic(err)
	}
	evn := DefaultEnv.NewContext(map[string]interface{}{
		"rate": 0.001,
	})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		evn.ExecValue(e)
	}
}

var (
	_ss = ""
	ss1 = "helo"
)

func BenchmarkCutter(b *testing.B) {
	m := map[string]interface{}{
		"ss": _ss,
	}
	for i := 0; i < b.N; i++ {
		getterOf(m)("ss")
	}
}
func BenchmarkCutter2(b *testing.B) {
	arr := []any{1, 2, 3}

	for i := 0; i < b.N; i++ {
		f, _ := cutterOf(arr)
		f(1, 2)
	}
}

func getVal(m any, k string) any {
	switch m := m.(type) {
	case map[string]interface{}:
		return m[k]
	default:
		return nil
	}
}

func getterOf(m any) func(k string) any {
	switch m := m.(type) {
	case map[string]interface{}:
		return func(k string) any {
			return m[k]
		}
	}
	return nil
}

func BenchmarkArrSet(b *testing.B) {
	c := DefaultEnv.NewContext(nil)
	c.ForceType = true
	c.SetByString("$", map[string]any{
		"req": map[string]any{
			"ss":  2.0,
			"app": "xxx",
		},
		"pa": map[string]any{
			"b1": map[string]any{
				"a": 1.0,
				"b": 1.0,
				"c": 1.0,
			},
		},
	})
	v, err := DefaultEnv.ParseValue(`

 a = 1;b = 2 

`)
	if err != nil {
		panic(err)
	}
	c.ExecValue(v)

	fmt.Println(c.GetTable())

	fmt.Println(DefaultEnv.NewHashKey("xxxxx"))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.ExecValue(v)
	}
}

var (
	arr = make([]any, 100)
)

func set(i int, v any) {
	arr[i] = v
}
func BenchmarkSet(b *testing.B) {
	var v any = 3
	for i := 0; i < b.N; i++ {
		set(i%100, v)
	}
}

func TestStringBuild(t *testing.T) {

	c := DefaultEnv.NewContext(nil)
	parseAndExec(DefaultEnv, `
name = 'lix';
age = 18890287457893;
class = 6;
desc = '${name}(${age}):${class}';
`, c)
	assertEqual(t, c, "desc", "lix(18890287457893):6")
}
