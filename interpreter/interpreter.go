package interpreter

import (
	"fmt"

	"github.com/seeadoog/expr/ast"
)

// Interpreter 字节码解释器入口
type Interpreter struct {
	env Env
}

// Env 环境接口（与 expr.Env 对接）
type Env interface {
	CalcHash(s string) uint64
	GetFunc(hash uint64, name string) (ScriptFunc, bool)
	GetLibFunc(objHash, funcHash uint64) ScriptFunc
}

// Context 执行上下文接口（与 expr.Context 对接）
type Context interface {
	Get(hash uint64) any
	Set(hash uint64, val any)
	GetByString(key string) any
	SetByString(key string, val any)
}

// ScriptFunc 脚本函数类型（与 expr.ScriptFunc 兼容）
type ScriptFunc func(ctx Context, args ...any) any

// NewInterpreter 创建解释器
func NewInterpreter(env Env) *Interpreter {
	return &Interpreter{env: env}
}

// Compile 编译 AST 到字节码
func (interp *Interpreter) Compile(node ast.Node) (*ByteCode, error) {
	compiler := NewCompiler(interp.env.CalcHash)
	return compiler.Compile(node)
}

// Execute 执行字节码
func (interp *Interpreter) Execute(bc *ByteCode, ctx Context) (any, error) {
	vm := NewVM(bc, interp.env, ctx)
	return vm.Run()
}

// CompileAndRun 编译并执行 AST
func (interp *Interpreter) CompileAndRun(node ast.Node, ctx Context) (any, error) {
	bc, err := interp.Compile(node)
	if err != nil {
		return nil, err
	}
	return interp.Execute(bc, ctx)
}

// Disassemble 反汇编字节码（调试用）
func (interp *Interpreter) Disassemble(bc *ByteCode) string {
	return bc.Disassemble()
}

// MustCompile 编译 AST，失败则 panic
func (interp *Interpreter) MustCompile(node ast.Node) *ByteCode {
	bc, err := interp.Compile(node)
	if err != nil {
		panic(fmt.Errorf("compile failed: %w", err))
	}
	return bc
}
