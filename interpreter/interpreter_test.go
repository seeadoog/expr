package interpreter_test

import (
	"fmt"
	"testing"

	"github.com/seeadoog/expr"
	"github.com/seeadoog/expr/ast"
	"github.com/seeadoog/expr/interpreter"
)

// mockContext 实现 interpreter.Context 接口用于测试
type mockContext struct {
	vars map[uint64]any
}

func newMockContext() *mockContext {
	return &mockContext{
		vars: make(map[uint64]any),
	}
}

func (m *mockContext) Get(hash uint64) any {
	return m.vars[hash]
}

func (m *mockContext) Set(hash uint64, val any) {
	m.vars[hash] = val
}

func (m *mockContext) GetByString(key string) any {
	// 简单实现：不使用
	return nil
}

func (m *mockContext) SetByString(key string, val any) {
	// 简单实现：不使用
}

// mockEnv 实现 interpreter.Env 接口用于测试
type mockEnv struct {
	hashFunc func(string) uint64
}

func newMockEnv() *mockEnv {
	return &mockEnv{
		hashFunc: defaultHash,
	}
}

func (m *mockEnv) CalcHash(s string) uint64 {
	return m.hashFunc(s)
}

func (m *mockEnv) GetFunc(hash uint64, name string) (interpreter.ScriptFunc, bool) {
	return nil, false
}

func (m *mockEnv) GetLibFunc(objHash, funcHash uint64) interpreter.ScriptFunc {
	return nil
}

func defaultHash(s string) uint64 {
	h := uint64(1469598103934665603)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

func TestCompileNumber(t *testing.T) {
	node := &ast.Number{Val: 42.0}
	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if result != 42.0 {
		t.Errorf("expected 42.0, got %v", result)
	}
}

func TestCompileBinary(t *testing.T) {
	// 3 + 5
	node := &ast.Binary{
		Op: "+",
		L:  &ast.Number{Val: 3.0},
		R:  &ast.Number{Val: 5.0},
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if result != 8.0 {
		t.Errorf("expected 8.0, got %v", result)
	}
}

func TestCompileVariable(t *testing.T) {
	// 变量 x
	node := &ast.Variable{Name: "x"}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	// 设置变量值
	ctx.Set(defaultHash("x"), 100.0)

	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if result != 100.0 {
		t.Errorf("expected 100.0, got %v", result)
	}
}

func TestCompileSet(t *testing.T) {
	// x = 42
	node := &ast.Set{
		L: &ast.Variable{Name: "x"},
		R: &ast.Number{Val: 42.0},
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()

	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if result != 42.0 {
		t.Errorf("expected 42.0, got %v", result)
	}

	// 验证变量已设置
	if ctx.Get(defaultHash("x")) != 42.0 {
		t.Errorf("variable x should be 42.0, got %v", ctx.Get(defaultHash("x")))
	}
}

func TestCompileTernary(t *testing.T) {
	// true ? 10 : 20
	node := &ast.Ternary{
		C: &ast.Bool{Val: true},
		L: &ast.Number{Val: 10.0},
		R: &ast.Number{Val: 20.0},
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if result != 10.0 {
		t.Errorf("expected 10.0, got %v", result)
	}
}

func TestCompileArray(t *testing.T) {
	// [1, 2, 3]
	node := &ast.ArrDef{
		V: []ast.Node{
			&ast.Number{Val: 1.0},
			&ast.Number{Val: 2.0},
			&ast.Number{Val: 3.0},
		},
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}

	if len(arr) != 3 {
		t.Errorf("expected length 3, got %d", len(arr))
	}

	if arr[0] != 1.0 || arr[1] != 2.0 || arr[2] != 3.0 {
		t.Errorf("expected [1, 2, 3], got %v", arr)
	}
}

func TestCompileMap(t *testing.T) {
	// {"a": 1, "b": 2}
	node := &ast.MapSet{
		Kvs: []ast.KV{
			{K: "a", V: &ast.Number{Val: 1.0}},
			{K: "b", V: &ast.Number{Val: 2.0}},
		},
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}

	if m["a"] != 1.0 || m["b"] != 2.0 {
		t.Errorf("expected {a:1, b:2}, got %v", m)
	}
}

func TestCompileShortCircuit(t *testing.T) {
	// false && (会 panic 的表达式)，应短路不执行右侧
	node := &ast.Binary{
		Op: "&&",
		L:  &ast.Bool{Val: false},
		R:  &ast.Number{Val: 42.0}, // 不会执行
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	vm := interpreter.NewVM(bc, env, ctx)
	result, err := vm.Run()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if result != false {
		t.Errorf("expected false, got %v", result)
	}
}

func TestDisassemble(t *testing.T) {
	// 简单表达式: x + 5
	node := &ast.Binary{
		Op: "+",
		L:  &ast.Variable{Name: "x"},
		R:  &ast.Number{Val: 5.0},
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	// 输出反汇编结果
	t.Logf("\n%s", bc.Disassemble())
}

func BenchmarkFFF(b *testing.B) {
	// 测试嵌套三元运算符性能
	node, err := expr.DefaultEnv.ParseValueToAstNode(`a`)
	if err != nil {
		b.Fatal(err)
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		b.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	ctx.Set(defaultHash("a"), false) // 设置变量 a

	//fmt.Println(bc.Instructions)
	for _, instruction := range bc.Instructions {
		fmt.Println(instruction)
	}
	vm := interpreter.NewVM(bc, env, ctx)

	b.ResetTimer() // 重置计时器，排除准备时间
	for i := 0; i < b.N; i++ {
		vm.Reset()
		_, err := vm.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestVM(b *testing.T) {
	// 测试嵌套三元运算符性能
	node, err := expr.DefaultEnv.ParseValueToAstNode(`
a=1
`)
	if err != nil {
		b.Fatal(err)
	}

	compiler := interpreter.NewCompiler(defaultHash)
	bc, err := compiler.Compile(node)
	if err != nil {
		b.Fatalf("compile failed: %v", err)
	}

	env := newMockEnv()
	ctx := newMockContext()
	ctx.Set(defaultHash("a"), false) // 设置变量 a

	//fmt.Println(bc.Instructions)
	for _, instruction := range bc.Instructions {
		fmt.Println(instruction)
	}
	bc2, err := compiler.Compile(node)
	if err != nil {
		b.Fatalf("compile failed: %v", err)
	}
	for _, instruction := range bc2.Instructions {
		fmt.Println(instruction)
	}

	vm := interpreter.NewVM(bc, env, ctx)
	vm.Reset()
}

var funs = []func() int{
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	}, func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
	func() int {
		return 1
	},
	func() int {
		return 2
	},
}

func dofunc(i int) int {
	return funs[i]()
}
func getString(i int) string {
	switch i {
	case 0:
		return "a"
	case 1:
		return "b"
	case 2:
		return "c"
	case 3:
		return "d"
	case 4:
		return "e"
	case 5:
		return "f"
	case 6:
		return "g"
	case 7:
		return "h"
	case 8:
		return "i"
	case 9:
		return "j"
	case 10:
		return "k"
	case 11:
		return "l"
	case 12:
		return "m"
	case 13:
		return "n"
	case 14:
		return "o"
	case 15:
		return "p"
	case 16:
		return "q"
	case 17:
		return "r"
	case 18:
		return "s"
	case 19:
		return "t"
	case 20:
		return "u"
	case 21:
		return "v"
	case 22:
		return "w"
	case 23:
		return "x"
	case 24:
		return "y"
	case 25:
		return "z"
	case 26:
		return "0"
	case 27:
		return "1"
	case 28:
		return "2"
	case 29:
		return "3"
	case 30:
		return "4"
	}
	return ""
}

var (
	s string
)
var (
	is int
)

func BenchmarkGetSrr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		is = dofunc(i % 31)
	}
}
