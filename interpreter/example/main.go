package main

import (
	"fmt"

	"github.com/seeadoog/expr"
	"github.com/seeadoog/expr/interpreter"
)

func main() {
	fmt.Println("=== Bytecode Interpreter 示例 ===\n")

	// 1. 创建 expr 环境
	env := expr.DefaultEnv
	ctx := env.NewContext(map[string]any{
		"x":    10.0,
		"y":    20.0,
		"name": "Alice",
	})

	// 2. 定义测试表达式
	expressions := []string{
		"x + y",
		"x * y + 5",
		"x > 5 ? 100 : 200",
		"name + \" says hello\"",
		"[1, 2, 3, x, y]",
		`{"name": name, "sum": x + y}`,
	}

	// 3. 对每个表达式进行编译和执行
	for i, exprText := range expressions {
		fmt.Printf("表达式 %d: %s\n", i+1, exprText)

		// 解析为 AST
		node, err := env.ParseValueToAstNode(exprText)
		if err != nil {
			fmt.Printf("  解析错误: %v\n\n", err)
			continue
		}

		// 编译为字节码
		compiler := interpreter.NewCompiler(func(s string) uint64 {
			return env.NewHashKey(s).Hash
		})
		bc, err := compiler.Compile(node)
		if err != nil {
			fmt.Printf("  编译错误: %v\n\n", err)
			continue
		}

		// 显示字节码信息
		fmt.Printf("  字节码指令数: %d\n", len(bc.Instructions))
		fmt.Printf("  常量池大小: %d\n", len(bc.Constants))
		fmt.Printf("  名称表大小: %d\n", len(bc.Names))

		// 执行字节码
		vm := interpreter.NewVM(bc, &envAdapter{env}, &ctxAdapter{ctx})
		result, err := vm.Run()
		if err != nil {
			fmt.Printf("  执行错误: %v\n\n", err)
			continue
		}

		fmt.Printf("  执行结果: %v (%T)\n", result, result)
		fmt.Println()
	}

	// 4. 演示反汇编
	fmt.Println("=== 反汇编示例 ===\n")
	exprText := "x + y * 2"
	fmt.Printf("表达式: %s\n\n", exprText)

	node, _ := env.ParseValueToAstNode(exprText)
	compiler := interpreter.NewCompiler(func(s string) uint64 {
		return env.NewHashKey(s).Hash
	})
	bc, _ := compiler.Compile(node)

	fmt.Println(bc.Disassemble())
}

// envAdapter 适配 expr.Env 到 interpreter.Env
type envAdapter struct {
	*expr.Env
}

func (e *envAdapter) CalcHash(s string) uint64 {
	return e.Env.NewHashKey(s).Hash
}

func (e *envAdapter) GetFunc(hash uint64, name string) (interpreter.ScriptFunc, bool) {
	// 占位实现
	return nil, false
}

func (e *envAdapter) GetLibFunc(objHash, funcHash uint64) interpreter.ScriptFunc {
	// 占位实现
	return nil
}

// ctxAdapter 适配 expr.Context 到 interpreter.Context
type ctxAdapter struct {
	*expr.Context
}

// expr.Context 已经实现了 Get/Set/GetByString/SetByString，无需额外实现
