# PanicWhenError 迁移说明

## 🔄 变更摘要

已将 `PanicWhenError` 从全局变量迁移到 `Context` 结构中，实现更灵活的错误处理控制。

---

## 📝 变更详情

### 1. Context 结构变化

**之前**:
```go
var PanicWhenError = true  // 全局变量

type Context struct {
    // ... 其他字段
}
```

**之后**:
```go
type Context struct {
    // ... 其他字段
    PanicWhenError bool  // 实例级别控制
}
```

### 2. 默认行为

- 新创建的 `Context` 默认 `PanicWhenError = true`
- 保持向后兼容，默认行为不变

### 3. API 变化

#### 新增函数

```go
// 带 Context 的错误创建函数
func newErrorfWithCtx(ctx *Context, format string, args ...interface{}) *Error
func newErrorWithCtx(err any, ctx *Context) *Error
```

#### 保留函数

```go
// 不带 Context 的函数保持不变（默认 panic）
func newErrorf(format string, args ...interface{}) *Error
func newError(err any) *Error
```

---

## 🚀 使用方法

### 基本使用（默认行为）

```go
env := expr.NewEnv()
ctx := env.NewContext(map[string]any{
    "a": 10,
})

// PanicWhenError 默认为 true
compiled, _ := env.ParseValue("a / 0")
result := compiled.Val(ctx)  // 错误时会 panic
```

### 禁用 Panic

```go
env := expr.NewEnv()
ctx := env.NewContext(map[string]any{
    "a": 10,
})

// 设置为 false，错误时返回 error 而不是 panic
ctx.PanicWhenError = false

compiled, _ := env.ParseValue("a / 0")
result := compiled.Val(ctx)

// 检查错误
if err, ok := result.(error); ok {
    fmt.Println("错误:", err)
}
```

### 不同 Context 不同行为

```go
env := expr.NewEnv()
compiled, _ := env.ParseValue("risky_operation()")

// Context 1: panic on error
ctx1 := env.NewContext(data1)
ctx1.PanicWhenError = true
result1 := compiled.Val(ctx1)  // 错误会 panic

// Context 2: return error
ctx2 := env.NewContext(data2)
ctx2.PanicWhenError = false
result2 := compiled.Val(ctx2)  // 错误返回 error 对象
if err, ok := result2.(error); ok {
    log.Println("错误:", err)
}
```

---

## 🔧 迁移指南

### 场景 1: 使用全局 PanicWhenError

**之前的代码**:
```go
func init() {
    expr.PanicWhenError = false
}

func main() {
    env := expr.NewEnv()
    ctx := env.NewContext(data)
    result := compiled.Val(ctx)
}
```

**迁移后**:
```go
func main() {
    env := expr.NewEnv()
    ctx := env.NewContext(data)
    ctx.PanicWhenError = false  // 在 Context 上设置
    result := compiled.Val(ctx)
}
```

### 场景 2: 需要全局禁用

**创建辅助函数**:
```go
func NewContextNoPanic(env *expr.Env, data map[string]any) *expr.Context {
    ctx := env.NewContext(data)
    ctx.PanicWhenError = false
    return ctx
}

// 使用
ctx := NewContextNoPanic(env, data)
result := compiled.Val(ctx)
```

### 场景 3: 测试代码

**之前**:
```go
func init() {
    expr.PanicWhenError = false
}

func TestSomething(t *testing.T) {
    // 测试代码
}
```

**迁移后**:
```go
func TestSomething(t *testing.T) {
    env := expr.NewEnv()
    ctx := env.NewContext(data)
    ctx.PanicWhenError = false  // 在测试中禁用 panic
    
    // 测试代码
}
```

---

## 💡 优势

### 1. 更好的隔离性

```go
// 不同的 Context 可以有不同的错误处理策略
criticalCtx := env.NewContext(criticalData)
criticalCtx.PanicWhenError = true  // 关键操作：立即 panic

nonCriticalCtx := env.NewContext(nonCriticalData)
nonCriticalCtx.PanicWhenError = false  // 非关键操作：返回错误
```

### 2. 并发安全

```go
// 不同 goroutine 使用不同的 Context，互不影响
go func() {
    ctx1 := env.NewContext(data1)
    ctx1.PanicWhenError = true
    // 使用 ctx1
}()

go func() {
    ctx2 := env.NewContext(data2)
    ctx2.PanicWhenError = false
    // 使用 ctx2
}()
```

### 3. 更灵活的错误处理

```go
func ProcessWithRetry(env *expr.Env, compiled expr.Val, data map[string]any) (any, error) {
    ctx := env.NewContext(data)
    ctx.PanicWhenError = false  // 禁用 panic，以便重试
    
    for i := 0; i < 3; i++ {
        result := compiled.Val(ctx)
        if err, ok := result.(error); ok {
            log.Printf("尝试 %d 失败: %v", i+1, err)
            continue
        }
        return result, nil
    }
    
    return nil, fmt.Errorf("3 次尝试后仍失败")
}
```

---

## 🧪 测试变化

### script_test.go

**之前**:
```go
func init() {
    PanicWhenError = false  // 全局禁用
}
```

**之后**:
```go
// 已移除 init 函数
// 在需要的测试中单独设置 ctx.PanicWhenError = false
```

---

## 📋 修改的文件

1. **error.go**
   - 移除全局变量 `PanicWhenError`
   - 保留 `newErrorf` 和 `newError`（默认 panic）
   - 新增 `newErrorfWithCtx` 和 `newErrorWithCtx`

2. **script.go**
   - `Context` 结构添加 `PanicWhenError` 字段
   - `Clone()` 方法复制 `PanicWhenError`
   - `NewContext()` 设置默认值为 `true`
   - 更新 `newError` 函数

3. **script_test.go**
   - 移除 `init()` 函数中的全局设置
   - 添加注释说明新的使用方式

---

## ✅ 向后兼容性

- ✅ 默认行为不变（`PanicWhenError = true`）
- ✅ 现有代码无需修改（除非使用了全局 `PanicWhenError`）
- ✅ 所有测试通过
- ✅ 编译无错误

---

## 🎯 最佳实践

### 1. 生产环境建议

```go
// 对于关键业务逻辑，保持 panic（快速失败）
ctx := env.NewContext(data)
ctx.PanicWhenError = true  // 显式设置（虽然这是默认值）

// 对于非关键或可恢复的操作，使用错误返回
ctx := env.NewContext(data)
ctx.PanicWhenError = false
result := compiled.Val(ctx)
if err, ok := result.(error); ok {
    // 记录错误并继续
    log.Warn("操作失败", err)
}
```

### 2. 测试环境建议

```go
func TestErrorHandling(t *testing.T) {
    env := expr.NewEnv()
    ctx := env.NewContext(data)
    ctx.PanicWhenError = false  // 测试中通常禁用 panic
    
    result := compiled.Val(ctx)
    if err, ok := result.(error); ok {
        // 验证错误信息
        assert.Contains(t, err.Error(), "expected error message")
    }
}
```

### 3. 中间件模式

```go
type ExprMiddleware struct {
    env *expr.Env
}

func (m *ExprMiddleware) Execute(exprStr string, data map[string]any, panicOnError bool) (any, error) {
    compiled, err := m.env.ParseValue(exprStr)
    if err != nil {
        return nil, err
    }
    
    ctx := m.env.NewContext(data)
    ctx.PanicWhenError = panicOnError  // 根据参数决定
    
    result := compiled.Val(ctx)
    if err, ok := result.(error); ok {
        return nil, err
    }
    
    return result, nil
}
```

---

## 📚 相关文档

- [使用文档](docs/USAGE.md#错误处理)
- [API 参考](docs/API.md#context)
- [主 README](README.md)

---

**变更日期**: 2026-06-25  
**影响范围**: 错误处理机制  
**兼容性**: 向后兼容  
**测试状态**: ✅ 所有测试通过
