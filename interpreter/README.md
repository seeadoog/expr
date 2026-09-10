# Bytecode Interpreter

高性能字节码解释器，用于执行 expr AST。

## 架构设计

### 核心组件

1. **Compiler** - 将 AST 编译成字节码
2. **ByteCode** - 字节码表示（指令序列、常量池、名称表等）
3. **VM** - 虚拟机，执行字节码
4. **Interpreter** - 统一入口，协调编译和执行

### 指令集

- **栈操作**: PUSH, POP, LOAD_VAR, STORE_VAR, LOAD_FAST, STORE_FAST
- **算术**: ADD, SUB, MUL, DIV, MOD, POW, NEG
- **比较**: EQ, NEQ, LT, LTE, GT, GTE, EQT, NEQT
- **逻辑**: NOT (&&, || 通过短路跳转实现)
- **位运算**: BIT_AND, BIT_OR, BIT_XOR, BIT_NOT, SHL, SHR
- **对象访问**: GET_ATTR, SET_ATTR, GET_INDEX, SET_INDEX, SLICE
- **函数调用**: CALL, CALL_METHOD
- **控制流**: JUMP, JUMP_IF_TRUE, JUMP_IF_FALSE
- **字面量**: MAKE_ARRAY, MAKE_MAP
- **特殊**: NOT_NIL, HALT

### 值表示

当前版本使用 `[]any` 作为操作数栈，实现简单但存在接口装箱开销。后续可优化为标签联合类型（tagged union）以提升性能。

## 支持的 AST 节点（v1）

✅ 核心表达式（已实现）:
- Number, String, Bool, Nil - 字面量
- Variable - 变量
- Binary - 二元运算（+, -, *, /, %, **, ==, !=, <, <=, >, >=, ===, !==, &&, ||, &, |, ^, <<, >>）
- Unary - 一元运算（-, !, ~）
- Call - 函数调用
- Access - 对象属性访问和方法调用
- ArrAccess - 数组/映射索引
- Set - 赋值（支持 const 声明）
- Ternary - 三元运算符
- ArrDef - 数组字面量
- MapSet - 映射字面量
- SliceCut - 切片操作
- NotNil - !! 非空断言
- NodeList - 表达式列表

❌ 暂未实现（未来版本）:
- IfElse, ForRange, Switch - 控制流
- Lambda, Lambda2 - Lambda 函数
- Break, Return - 流程控制

## 使用示例

### 基本使用（独立测试）

```go
import (
    "github.com/seeadoog/expr/ast"
    "github.com/seeadoog/expr/interpreter"
)

// 创建简单的环境和上下文
env := newMockEnv()
ctx := newMockContext()

// 设置变量
ctx.Set(hashFunc("x"), 10.0)

// 创建 AST: x + 5
node := &ast.Binary{
    Op: "+",
    L:  &ast.Variable{Name: "x"},
    R:  &ast.Number{Val: 5.0},
}

// 编译
compiler := interpreter.NewCompiler(hashFunc)
bc, err := compiler.Compile(node)
if err != nil {
    panic(err)
}

// 执行
vm := interpreter.NewVM(bc, env, ctx)
result, err := vm.Run()
// result = 15.0
```

### 反汇编（调试）

```go
bc, _ := compiler.Compile(node)
fmt.Println(bc.Disassemble())

// 输出:
// === Bytecode Disassembly ===
// Constants:
//   [0] 5
// Names:
//   [0] x (hash=4953295298048672641)
// Instructions:
//   0000  LOAD_VAR   A=0    B=0     ; x
//   0001  PUSH       A=0    B=0     ; 5
//   0002  ADD        A=0    B=0
//   0003  HALT       A=0    B=0
```

## 与 expr 包集成

interpreter 包通过接口与 expr 包解耦，需要创建适配器：

### 接口定义

```go
// interpreter.Env - 函数查找
type Env interface {
    CalcHash(s string) uint64
    GetFunc(hash uint64, name string) (ScriptFunc, bool)
    GetLibFunc(objHash, funcHash uint64) ScriptFunc
}

// interpreter.Context - 变量存取
type Context interface {
    Get(hash uint64) any
    Set(hash uint64, val any)
    GetByString(key string) any
    SetByString(key string, val any)
}
```

### 适配器实现（示例）

```go
// expr.Env 实现了所需接口，可直接使用
// expr.Context 也实现了所需接口

// 集成示例
func CompileAndRunWithExpr(exprEnv *expr.Env, exprText string, ctx *expr.Context) (any, error) {
    // 1. 使用 expr 的 parser
    node, err := exprEnv.ParseValueToAstNode(exprText)
    if err != nil {
        return nil, err
    }

    // 2. 创建 interpreter（传入 expr.Env 作为适配器）
    interp := interpreter.NewInterpreter(envAdapter{exprEnv})

    // 3. 编译并执行
    bc, err := interp.Compile(node)
    if err != nil {
        return nil, err
    }

    return interp.Execute(bc, ctxAdapter{ctx})
}

// 适配器实现
type envAdapter struct { *expr.Env }
func (e envAdapter) CalcHash(s string) uint64 { 
    return e.Env.NewHashKey(s).Hash 
}
// ... 实现其他方法

type ctxAdapter struct { *expr.Context }
// Context 接口 expr.Context 已经实现，可直接使用
```

## 性能特性

### 优势
- **编译一次，多次执行**: 字节码可缓存，避免重复解析和编译
- **紧凑的指令格式**: 双操作数设计，减少指令数量
- **短路求值优化**: && 和 || 正确实现短路跳转
- **常量池去重**: 相同常量仅存储一次

### 当前限制
- **接口装箱**: `[]any` 栈导致数值类型装箱，有性能开销
- **函数调用**: 当前版本函数调用未完全集成（占位实现）
- **对象方法**: objFuncMap 集成待完成

### 后续优化方向
1. **值类型优化**: 使用标签联合（tagged union）避免装箱
   ```go
   type Value struct {
       tag  byte  // 0=nil, 1=bool, 2=num, 3=ptr
       num  float64
       ptr  any
   }
   ```
2. **寄存器VM**: 从栈机改为寄存器机，减少栈操作
3. **JIT编译**: 热点代码编译为原生机器码
4. **指令融合**: 识别常见模式合并指令（如 PUSH+ADD → ADD_IMM）

## 测试

运行测试:
```bash
go test -v ./interpreter/...
```

基准测试（待添加）:
```bash
go test -bench=. ./interpreter/...
```

## 文件结构

```
interpreter/
├── opcode.go           - 操作码定义
├── bytecode.go         - 字节码结构
├── compiler.go         - AST -> 字节码编译器
├── vm.go               - 虚拟机执行引擎
├── interpreter.go      - 统一入口和接口定义
├── interpreter_test.go - 单元测试
└── README.md           - 本文档
```

## 许可

与主项目保持一致。
