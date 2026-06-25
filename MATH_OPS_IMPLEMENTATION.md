# 数学赋值操作符实现说明

## 概述
本次更新为 expr 规则引擎新增了三个数学赋值操作符：`-=`、`*=`、`/=`，与现有的 `+=` 操作符保持一致的实现风格。

## 修改的文件

### 1. `ast/expr.y` (语法定义文件)
- 添加了三个新的 token：`SUBEQ`、`MULEQ`、`DIVEQ`
- 在优先级声明中添加了这些操作符
- 为每个操作符添加了对应的语法规则（支持 Ident、Var 和 ArrIndex）

### 2. `lex.go` (词法分析器)
- 在 `-` 字符处理中添加了 `-=` 的识别
- 在 `*` 和 `/` 字符处理中添加了 `*=` 和 `/=` 的识别

### 3. `yacc.go` (语义分析器)
- 为 `-=`、`*=`、`/=` 添加了对应的 case 处理
- 使用 `setValue` 结构配合 `newBinaryValue` 实现操作
- `-=` 使用 `NumberOf(a.Val(ctx)) - NumberOf(b.Val(ctx))`
- `*=` 使用 `NumberOf(a.Val(ctx)) * NumberOf(b.Val(ctx))`
- `/=` 使用 `NumberOf(a.Val(ctx)) / NumberOf(b.Val(ctx))`

### 4. `ast/parser.go` (自动生成)
- 通过 `goyacc -o parser.go expr.y` 命令重新生成

### 5. `math_test.go` (新增测试文件)
- 创建了全面的单元测试，包括：
  - `TestAddEq`: 测试 `+=` 操作符（数字加法、字符串拼接）
  - `TestSubEq`: 测试 `-=` 操作符（基本减法、负数、小数）
  - `TestMulEq`: 测试 `*=` 操作符（乘0、乘1、负数、小数）
  - `TestDivEq`: 测试 `/=` 操作符（基本除法、负数、小数）
  - `TestMathOpsCombined`: 测试组合操作
  - `TestMathOpsWithVarAccess`: 测试对象属性访问的操作
  - `TestMathOpsWithArrayAccess`: 测试数组元素访问的操作

## 使用示例

### 基本用法
```go
env := expr.NewEnv()

// 减法赋值
ctx := env.NewContext(map[string]any{"x": 100.0})
val, _ := env.ParseValue("x -= 30")
result := val.Val(ctx) // 70

// 乘法赋值
ctx := env.NewContext(map[string]any{"y": 5.0})
val, _ := env.ParseValue("y *= 4")
result := val.Val(ctx) // 20

// 除法赋值
ctx := env.NewContext(map[string]any{"z": 50.0})
val, _ := env.ParseValue("z /= 10")
result := val.Val(ctx) // 5
```

### 组合操作
```go
ctx := env.NewContext(map[string]any{"a": 10.0})
val, _ := env.ParseValue("a += 5; a *= 2; a -= 4; a /= 2")
result := val.Val(ctx) // 13
```

### 对象属性访问
```go
ctx := env.NewContext(map[string]any{
    "obj": map[string]any{"count": 100.0}
})
val, _ := env.ParseValue("obj->count -= 25")
result := val.Val(ctx) // 75
```

### 数组元素访问
```go
ctx := env.NewContext(map[string]any{
    "arr": []any{10.0, 20.0, 30.0}
})
val, _ := env.ParseValue("arr[0] *= 3")
result := val.Val(ctx) // 30
```

## 测试结果

所有测试均通过：
- ✅ `TestAddEq`: 4个子测试全部通过
- ✅ `TestSubEq`: 5个子测试全部通过
- ✅ `TestMulEq`: 6个子测试全部通过
- ✅ `TestDivEq`: 5个子测试全部通过
- ✅ `TestMathOpsCombined`: 3个子测试全部通过
- ✅ `TestMathOpsWithVarAccess`: 4个子测试全部通过
- ✅ `TestMathOpsWithArrayAccess`: 4个子测试全部通过

完整测试套件也全部通过，确保没有破坏现有功能。

## 技术细节

### 操作符优先级
所有赋值操作符（`+=`、`-=`、`*=`、`/=`）具有相同的优先级，处于表达式优先级的较低层级。

### 类型转换
- `-=`、`*=`、`/=` 使用 `NumberOf()` 函数将操作数转换为数字类型
- `+=` 保持原有行为，支持数字加法和字符串拼接

### 语法支持
所有操作符都支持三种左值形式：
1. 标识符：`x -= 5`
2. 对象属性访问：`obj->field *= 2`
3. 数组元素访问：`arr[0] /= 3`

## 兼容性
本次更新完全向后兼容，不影响任何现有功能和 API。
