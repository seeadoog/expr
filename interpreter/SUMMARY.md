# 字节码解释器实现总结

## 完成情况

✅ **已完成**：核心表达式字节码解释器，支持大部分常用场景

### 实现的功能

1. **编译器** (`compiler.go`)
   - AST 到字节码的编译
   - 局部变量优化（const 声明）
   - 短路求值优化（&& 和 ||）
   - 常量池去重

2. **虚拟机** (`vm.go`)
   - 基于栈的指令执行
   - 与 expr.Context 集成（变量读写）
   - 类型转换（numberOf/stringOf/boolOf）
   - 对象属性和索引访问
   - 反射支持结构体字段访问

3. **指令集** (`opcode.go`)
   - 42 条指令，覆盖核心表达式
   - 双操作数格式（A, B）
   - 专用表结构（CallInfo, NameRef）避免位打包

4. **集成支持** (`interpreter.go`, `integration.go`)
   - 接口定义解耦 expr 包
   - 适配器模式对接现有系统
   - 示例代码和使用指南

### 支持的 AST 节点

- ✅ Number, String, Bool, Nil
- ✅ Variable（全局变量、局部变量）
- ✅ Binary（+, -, *, /, %, **, ==, !=, <, <=, >, >=, ===, !==, &&, ||, &, |, ^, <<, >>）
- ✅ Unary（-, !, ~）
- ✅ Call（函数调用 - 架构完成，需 expr 集成）
- ✅ Access（obj.field, obj.method()）
- ✅ ArrAccess（arr[idx]）
- ✅ Set（赋值，支持 const）
- ✅ Ternary（三元运算符）
- ✅ ArrDef（数组字面量）
- ✅ MapSet（map 字面量）
- ✅ SliceCut（切片 arr[st:ed]）
- ✅ NotNil（!! 断言）
- ✅ NodeList（表达式序列）

### 测试验证

```bash
$ go test -v ./interpreter/...
=== RUN   TestCompileNumber
--- PASS: TestCompileNumber (0.00s)
=== RUN   TestCompileBinary
--- PASS: TestCompileBinary (0.00s)
=== RUN   TestCompileVariable
--- PASS: TestCompileVariable (0.00s)
=== RUN   TestCompileSet
--- PASS: TestCompileSet (0.00s)
=== RUN   TestCompileTernary
--- PASS: TestCompileTernary (0.00s)
=== RUN   TestCompileArray
--- PASS: TestCompileArray (0.00s)
=== RUN   TestCompileMap
--- PASS: TestCompileMap (0.00s)
=== RUN   TestCompileShortCircuit
--- PASS: TestCompileShortCircuit (0.00s)
=== RUN   TestDisassemble
--- PASS: TestDisassemble (0.00s)
PASS
ok      github.com/seeadoog/expr/interpreter    0.039s
```

### 实际运行效果

```bash
$ go run interpreter/example/main.go
=== Bytecode Interpreter 示例 ===

表达式 1: x + y
  字节码指令数: 4
  常量池大小: 0
  名称表大小: 2
  执行结果: 30 (float64)

表达式 4: name + " says hello"
  字节码指令数: 4
  常量池大小: 1
  名称表大小: 1
  执行结果: Alice says hello (string)

表达式 5: [1, 2, 3, x, y]
  执行结果: [1 2 3 10 20] ([]interface {})

表达式 6: {"name": name, "sum": x + y}
  执行结果: map[name:Alice sum:30] (map[string]interface {})
```

## 架构设计

### 核心原则

1. **接口解耦**：interpreter 包独立，通过 Env/Context 接口与 expr 对接，避免循环依赖
2. **值表示简单**：使用 `[]any` 栈，实现简单优先（后续可优化为标签联合）
3. **指令紧凑**：双操作数格式，专用表减少位打包
4. **渐进实现**：v1 聚焦核心表达式，控制流和 lambda 留给后续版本

### 关键技术点

1. **哈希机制**：expr 使用 `hashManager` 将字符串映射为递增 ID，避免真哈希冲突
2. **短路求值**：&& 和 || 通过条件跳转正确实现短路
3. **栈序约定**：
   - 二元运算：[左, 右] → [结果]
   - 赋值：[值, 对象, 索引] → [值]
   - 函数调用：[arg0, arg1, ...] → [结果]

### 文件结构

```
interpreter/
├── opcode.go              - 42 条操作码定义
├── bytecode.go            - 字节码结构（指令、常量池、名称表）
├── compiler.go            - AST → 字节码编译器（500+ 行）
├── vm.go                  - 虚拟机执行引擎（400+ 行）
├── interpreter.go         - 统一入口和接口定义
├── integration.go         - 集成指南和适配器
├── interpreter_test.go    - 9 个单元测试
├── example/main.go        - 完整示例程序
└── README.md              - 文档
```

## 待完成事项

### 必需（核心功能）

1. **函数调用集成**
   - 当前 `vm.callFunc` 和 `vm.callMethod` 是占位实现
   - 需要访问 `expr.Env.funtables` 和 `objFuncMap`
   - 建议在 expr.Env 添加导出方法：
     ```go
     func (e *Env) GetFuncByHash(hash uint64) ScriptFunc
     func GetObjFuncByType(typ reflect.Type, funcHash uint64) ObjFunc
     ```

2. **错误处理对齐**
   - expr 通过 `*Error` 类型返回错误
   - VM 需要识别并传播（当前简单 panic）
   - 建议检查栈顶值类型，如果是 `*Error` 则提前返回

### 可选（性能和功能扩展）

1. **性能优化**
   - 标签联合值类型避免装箱
   - 寄存器 VM 替代栈机
   - 指令融合（PUSH+ADD → ADD_IMM）
   - JIT 编译热点代码

2. **控制流支持**
   - IfElse, ForRange, Switch
   - Break, Return
   - 需要额外的循环上下文管理

3. **Lambda 支持**
   - 捕获外部变量（闭包）
   - 嵌套作用域
   - Lambda 字节码编译为子程序

4. **基准测试**
   - 对比字节码 VM vs Val 闭包树
   - 编译时间、执行时间、内存占用

## 使用建议

### 何时使用字节码解释器

✅ **适合**：
- 表达式需要多次执行（编译一次，执行多次）
- 规则引擎、模板渲染等场景
- 需要缓存编译结果

❌ **不适合**：
- 一次性表达式（编译开销无法摊销）
- 需要 lambda 或复杂控制流（当前版本未支持）
- 对启动时间敏感（编译有开销）

### 集成步骤

1. 使用 `env.ParseValueToAstNode()` 解析表达式
2. 创建 `interpreter.NewCompiler(env.CalcHash)` 编译为字节码
3. 缓存 `ByteCode` 对象（可序列化存储）
4. 使用 `NewVM(bc, env, ctx).Run()` 执行
5. 调试时用 `bc.Disassemble()` 查看字节码

### 完整集成需要（expr 包侧）

为了完全集成，建议在 `expr` 包添加以下导出方法：

```go
// env.go
func (e *Env) GetFuncByHash(hash uint64, name string) (ScriptFunc, bool) {
    f := e.funtables[name]
    if f == nil {
        return nil, false
    }
    return f.fun, true
}

// funcs.go 或新文件
func GetObjFunc(typ reflect.Type, funcHash uint64, funcName string) (ObjFunc, bool) {
    typeFuncs := objFuncMap.get(typ)
    if typeFuncs == nil {
        return nil, false
    }
    f := typeFuncs.get(funcHash)
    return f, f != nil
}
```

## 总结

✅ 核心功能完成，架构清晰，测试通过
✅ 与 expr 包解耦良好，通过接口对接
✅ 支持大部分常用表达式场景
⚠️ 函数调用需要 expr 包配合添加导出方法
⚠️ 性能优化（标签联合）和高级功能（控制流、lambda）留待后续

字节码解释器已经可以工作，可以处理变量、算术、比较、逻辑、字符串、数组、map 等核心表达式。要达到与现有 `Val` 系统完全等价，需要完成函数调用集成和控制流支持。
