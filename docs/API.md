# Expr API 参考文档

本文档提供 seeadoog/expr 的完整 API 参考。

## 目录

- [核心类型](#核心类型)
- [环境 (Env)](#环境-env)
- [上下文 (Context)](#上下文-context)
- [值接口 (Val)](#值接口-val)
- [内置函数完整列表](#内置函数完整列表)
- [工具函数](#工具函数)

---

## 核心类型

### Env

表达式环境，管理函数注册和表达式编译。

```go
type Env struct {
    // 内部字段
}
```

### Context

执行上下文，存储变量和执行状态。

```go
type Context struct {
    // 内部字段
}
```

### Val

编译后的表达式接口。

```go
type Val interface {
    Val(ctx *Context) any
}
```

---

## 环境 (Env)

### NewEnv

创建新的表达式环境。

```go
func NewEnv() *Env
```

**返回：**
- `*Env` - 新的环境实例

**示例：**
```go
env := expr.NewEnv()
```

### DefaultEnv

全局默认环境实例。

```go
var DefaultEnv = NewEnv()
```

**示例：**
```go
compiled, _ := expr.DefaultEnv.ParseValue("a + b")
```

### ParseValue

编译表达式字符串为可执行对象。

```go
func (e *Env) ParseValue(s string) (Val, error)
```

**参数：**
- `s string` - 表达式字符串

**返回：**
- `Val` - 编译后的表达式对象
- `error` - 编译错误（如果有）

**示例：**
```go
compiled, err := env.ParseValue("a + b * c")
if err != nil {
    log.Fatal(err)
}
```

### ParseValueToAstNode

编译表达式为 AST 节点。

```go
func (e *Env) ParseValueToAstNode(val string) (ast.Node, error)
```

**参数：**
- `val string` - 表达式字符串

**返回：**
- `ast.Node` - AST 节点
- `error` - 编译错误（如果有）

**示例：**
```go
node, err := env.ParseValueToAstNode("a + b")
```

### RegFunc

注册自定义函数。

```go
func (e *Env) RegFunc(name string, fn any)
```

**参数：**
- `name string` - 函数名
- `fn any` - 函数实现（必须是函数类型）

**支持的函数签名：**
```go
func()                          // 无参数
func() R                        // 无参数，有返回值
func(T1) R                      // 单参数
func(T1, T2) R                  // 多参数
func(T1, T2, ...Tn) R           // 可变参数
func(T1) (R, error)             // 带错误返回
```

**示例：**
```go
// 简单函数
env.RegFunc("double", func(x float64) float64 {
    return x * 2
})

// 多参数函数
env.RegFunc("add", func(a, b float64) float64 {
    return a + b
})

// 可变参数
env.RegFunc("sum", func(nums ...float64) float64 {
    total := 0.0
    for _, n := range nums {
        total += n
    }
    return total
})

// 带错误返回
env.RegFunc("divide", func(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
})
```

### NewContext

创建新的执行上下文。

```go
func (e *Env) NewContext(data map[string]any) *Context
```

**参数：**
- `data map[string]any` - 初始变量（可以为 nil）

**返回：**
- `*Context` - 新的上下文实例

**示例：**
```go
ctx := env.NewContext(map[string]any{
    "name": "John",
    "age": 30.0,
})
```

### GetContextFromPool

从对象池获取上下文（性能优化）。

```go
func (e *Env) GetContextFromPool() *Context
```

**返回：**
- `*Context` - 来自对象池的上下文实例

**示例：**
```go
ctx := env.GetContextFromPool()
defer env.PutContext2Pool(ctx)

ctx.SetByString("key", value)
```

### PutContext2Pool

将上下文归还到对象池。

```go
func (e *Env) PutContext2Pool(ctx *Context)
```

**参数：**
- `ctx *Context` - 要归还的上下文

**示例：**
```go
ctx := env.GetContextFromPool()
defer env.PutContext2Pool(ctx)

// 使用 ctx
```

---

## 上下文 (Context)

### SetByString

通过字符串 key 设置变量。

```go
func (c *Context) SetByString(key string, val any)
```

**参数：**
- `key string` - 变量名
- `val any` - 变量值

**示例：**
```go
ctx.SetByString("name", "John")
ctx.SetByString("age", 30.0)
ctx.SetByString("scores", []any{90.0, 85.0, 95.0})
```

### Set

通过 hash 值设置变量（更快）。

```go
func (c *Context) Set(hash uint64, val any)
```

**参数：**
- `hash uint64` - 变量名的 hash 值
- `val any` - 变量值

**示例：**
```go
nameHash := expr.CalcHash("name")
ctx.Set(nameHash, "John")
```

### Get

通过 hash 值获取变量。

```go
func (c *Context) Get(hash uint64) any
```

**参数：**
- `hash uint64` - 变量名的 hash 值

**返回：**
- `any` - 变量值

**示例：**
```go
nameHash := expr.CalcHash("name")
value := ctx.Get(nameHash)
```

---

## 值接口 (Val)

### Val

执行表达式并返回结果。

```go
func (v Val) Val(ctx *Context) any
```

**参数：**
- `ctx *Context` - 执行上下文

**返回：**
- `any` - 表达式结果

**返回类型：**
- `float64` - 数字
- `string` - 字符串
- `bool` - 布尔值
- `[]any` - 数组
- `map[string]any` - 对象
- `nil` - 空值
- `error` - 运行时错误

**示例：**
```go
compiled, _ := env.ParseValue("a + b")
ctx := env.NewContext(map[string]any{
    "a": 10.0,
    "b": 20.0,
})

result := compiled.Val(ctx)
fmt.Println(result) // 30
```

---

## 内置函数完整列表

### 数学运算

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `add(a, b)` | 2 | number | 加法 |
| `sub(a, b)` | 2 | number | 减法 |
| `mul(a, b)` | 2 | number | 乘法 |
| `div(a, b)` | 2 | number | 除法 |
| `mod(a, b)` | 2 | number | 取模 |
| `pow(a, b)` | 2 | number | 幂运算 |
| `neg(x)` | 1 | number | 取负 |

### 类型转换

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `int(x)` | 1 | number | 转整数 |
| `number(x)` | 1 | number | 转数字 |
| `string(x)` | 1 | string | 转字符串 |
| `bool(x)` / `boolean(x)` | 1 | bool | 转布尔 |
| `bytes(x)` | 1 | []byte | 转字节数组 |
| `type(x)` | 1 | string | 获取类型名 |

### 字符串函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `str_to_upper(s)` | 1 | string | 转大写 |
| `str_to_lower(s)` | 1 | string | 转小写 |
| `str_trim(s)` | 1 | string | 去除首尾空格 |
| `str_fields(s)` | 1 | []string | 按空格分割 |
| `str_split(s, sep, n)` | 3 | []string | 分割字符串 |
| `str_join(parts...)` | -1 | string | 连接字符串 |
| `has_prefix(s, prefix)` | 2 | bool | 前缀判断 |
| `has_suffix(s, suffix)` | 2 | bool | 后缀判断 |
| `sprintf(format, args...)` | -1 | string | 格式化 |
| `len(s)` | 1 | number | 长度 |

### 字符串方法

| 方法 | 说明 |
|------|------|
| `str.len()` | 字符串长度 |
| `str.has_prefix(prefix)` | 前缀判断 |
| `str.has_suffix(suffix)` | 后缀判断 |
| `str.contains(substr)` | 包含判断 |
| `str.has(substr)` | 同 contains |
| `str.split(sep, n)` | 分割字符串 |
| `str.fields()` | 按空格分割 |
| `str.trim(cutset)` | 去除指定字符 |
| `str.trim_left(cutset)` | 去除左侧字符 |
| `str.trim_right(cutset)` | 去除右侧字符 |
| `str.trim_space()` | 去除空格 |
| `str.slice(start, end)` | 子串 |
| `str.bytes()` | 转字节数组 |
| `str.base64()` | Base64 编码 |
| `str.base64d()` | Base64 解码 |
| `str.hex()` | 十六进制编码 |
| `str.md5()` | MD5 哈希 |
| `str.json_str()` | JSON 序列化 |

### 数组函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `slice_new(items...)` | -1 | []any | 创建数组 |
| `slice_init(size)` | -1 | []any | 初始化数组 |
| `append(arr, items...)` | -1 | []any | 追加元素 |
| `len(arr)` | 1 | number | 数组长度 |
| `get(arr, index)` | 2 | any | 获取元素 |
| `set(arr, index, value)` | 3 | any | 设置元素 |
| `set_index(arr, index, value)` | 3 | any | 同 set |
| `slice_cut(arr, start, end)` | 3 | []any | 切片 |
| `for(arr, lambda)` | 2 | void | 遍历 |
| `all(arr, condition)` | 2 | []any | 过滤 |
| `range(n)` | 1 | []any | 生成序列 |

### 数组方法

| 方法 | 说明 |
|------|------|
| `arr.len()` | 数组长度 |
| `arr.get(index)` | 获取元素 |
| `arr.slice(start, end)` | 切片 |
| `arr.for(lambda)` | 遍历 |
| `arr.all(condition)` | 过滤 |
| `arr.sort(compare)` | 排序 |
| `arr.json_str()` | JSON 序列化 |

### Map/对象函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `new()` | 0 | map | 创建空 Map |
| `get(obj, key)` | 2 | any | 获取值 |
| `set(obj, key, value)` | 3 | map | 设置值 |
| `delete(obj, key)` | 2 | map | 删除键 |
| `len(obj)` | 1 | number | 键数量 |

### Map/对象方法

| 方法 | 说明 |
|------|------|
| `obj.get(key)` | 获取值 |
| `obj.set(key, value)` | 设置值 |
| `obj.delete(key)` | 删除键 |
| `obj.len()` | 键数量 |
| `obj.for(lambda)` | 遍历 |
| `obj.json_str()` | JSON 序列化 |

### JSON 函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `json_to(obj)` | 1 | string | 对象转 JSON |
| `json_from(str)` | 1 | any | JSON 转对象 |
| `to_json_str(obj)` | 1 | string | 同 json_to |
| `to_json_obj(str)` | 1 | any | 同 json_from |

### 时间函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `time_now()` | 0 | Time | 当前时间 |
| `time_now_mill()` | 0 | number | 当前毫秒时间戳 |
| `time_from_unix(ts)` | 1 | Time | Unix 时间戳转时间 |
| `time_parse(layout, value)` | 2 | Time | 解析时间 |
| `time_format(time, layout)` | 2 | string | 格式化时间 |

### 时间方法

| 方法 | 说明 |
|------|------|
| `time.year()` | 年份 |
| `time.month()` | 月份 |
| `time.day()` | 日期 |
| `time.hour()` | 小时 |
| `time.minute()` | 分钟 |
| `time.second()` | 秒 |
| `time.unix()` | Unix 时间戳（秒） |
| `time.unix_micro()` | 微秒时间戳 |
| `time.unix_mill()` | 毫秒时间戳 |
| `time.format(layout)` | 格式化 |
| `time.add_mill(ms)` | 添加毫秒 |
| `time.sub(other)` | 时间差（毫秒） |
| `time.local()` | 转本地时区 |
| `time.utc()` | 转 UTC 时区 |

### 编码函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `base64_encode(data)` | 1 | string | Base64 编码 |
| `base64_decode(data)` | 1 | []byte | Base64 解码 |
| `hex_encode(data)` | 1 | string | 十六进制编码 |
| `hex_decode(data)` | 1 | []byte | 十六进制解码 |
| `md5(data)` | 1 | []byte | MD5 哈希 |
| `sha256(data)` | 1 | []byte | SHA256 哈希 |
| `hmac_sha256(data, key)` | 2 | []byte | HMAC-SHA256 |

### 字节数组方法

| 方法 | 说明 |
|------|------|
| `bytes.base64()` | Base64 编码 |
| `bytes.base64d()` | Base64 解码 |
| `bytes.hex()` | 十六进制编码 |
| `bytes.slice(start, end)` | 切片 |
| `bytes.copy()` | 复制 |
| `bytes.bytes()` | 转字节数组 |

### 逻辑函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `eq(a, b)` | 2 | bool | 相等 |
| `neq(a, b)` | 2 | bool | 不等 |
| `eqs(a, b)` | 2 | bool | 严格相等 |
| `neqs(a, b)` | 2 | bool | 严格不等 |
| `gt(a, b)` | 2 | bool | 大于 |
| `gte(a, b)` | 2 | bool | 大于等于 |
| `lt(a, b)` | 2 | bool | 小于 |
| `lte(a, b)` | 2 | bool | 小于等于 |
| `and(args...)` | -1 | bool | 逻辑与 |
| `or(args...)` | -1 | bool | 逻辑或 |
| `not(x)` | 1 | bool | 逻辑非 |
| `orr(a, b)` | 2 | any | a 为 nil 返回 b |
| `in(value, arr)` | -1 | bool | 包含判断 |

### 控制流函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `if(cond, expr...)` | -1 | any | 条件判断 |
| `ternary(cond, true, false)` | 3 | any | 三元运算 |
| `loop(expr...)` | -1 | any | 循环 |
| `return(value...)` | -1 | any | 返回 |
| `catch(expr)` | 1 | any | 捕获异常 |
| `recover(expr)` | 1 | any | 恢复执行 |

### 正则表达式

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `regexp_new(pattern)` | 1 | *Regexp | 创建正则 |

| 方法 | 说明 |
|------|------|
| `regex.match(str)` | 匹配字符串 |

### 工具函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `print(args...)` | -1 | void | 打印（调试用） |
| `sleep(ms)` | 1 | void | 休眠（毫秒） |
| `go(expr)` | 1 | void | 异步执行 |
| `exec(expr...)` | -1 | any | 执行表达式 |
| `unwrap(value)` | 1 | any | 解包装 |

### URL 函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `url_new_values()` | 0 | url.Values | 创建 URL 参数 |

| 方法 | 说明 |
|------|------|
| `values.get(key)` | 获取参数 |
| `values.set(key, value)` | 设置参数 |
| `values.encode()` | 编码为查询字符串 |

### HTTP 函数

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `http_request(method, url, headers, body, timeout)` | 5 | any | HTTP 请求 |
| `response_write(data)` | 1 | void | 写响应 |

### StringBuilder

| 函数 | 参数 | 返回 | 说明 |
|------|------|------|------|
| `str_builder()` | 0 | *Builder | 创建字符串构建器 |

| 方法 | 说明 |
|------|------|
| `builder.write(str)` | 追加字符串 |

---

## 工具函数

### CalcHash

计算字符串的 hash 值。

```go
func CalcHash(s string) uint64
```

**参数：**
- `s string` - 要计算 hash 的字符串

**返回：**
- `uint64` - hash 值

**示例：**
```go
nameHash := expr.CalcHash("name")
ctx.Set(nameHash, "John")
```

### NumberOf

将值转换为数字。

```go
func NumberOf(v any) float64
```

**参数：**
- `v any` - 任意值

**返回：**
- `float64` - 数字值

**示例：**
```go
num := expr.NumberOf("123")  // 123.0
num = expr.NumberOf(true)    // 1.0
num = expr.NumberOf(false)   // 0.0
```

### StringOf

将值转换为字符串。

```go
func StringOf(v any) string
```

**参数：**
- `v any` - 任意值

**返回：**
- `string` - 字符串值

**示例：**
```go
str := expr.StringOf(123)     // "123"
str = expr.StringOf(true)     // "true"
```

---

## 类型定义

### 函数签名类型

#### FuncDefine1
单参数函数。

```go
type FuncDefine1[T1, R any] func(T1) R
```

#### FuncDefine2
双参数函数。

```go
type FuncDefine2[T1, T2, R any] func(T1, T2) R
```

#### FuncDefine3
三参数函数。

```go
type FuncDefine3[T1, T2, T3, R any] func(T1, T2, T3) R
```

---

## 使用示例

### 完整示例：规则引擎

```go
package main

import (
    "fmt"
    "github.com/seeadoog/expr"
)

func main() {
    // 1. 创建环境并注册自定义函数
    env := expr.NewEnv()
    
    env.RegFunc("is_vip", func(spent, orders float64) bool {
        return spent > 10000 && orders > 50
    })
    
    // 2. 编译规则
    rules := map[string]string{
        "vip_discount": "is_vip(totalSpent, orderCount) ? 0.2 : 0",
        "free_shipping": "totalAmount >= 99 || vipLevel >= 3",
    }
    
    compiled := make(map[string]expr.Val)
    for name, exprStr := range rules {
        val, err := env.ParseValue(exprStr)
        if err != nil {
            panic(fmt.Sprintf("规则 %s 编译失败: %v", name, err))
        }
        compiled[name] = val
    }
    
    // 3. 应用规则
    userData := map[string]any{
        "totalSpent":  15000.0,
        "orderCount":  60.0,
        "totalAmount": 120.0,
        "vipLevel":    2.0,
    }
    
    ctx := env.NewContext(userData)
    
    discount := compiled["vip_discount"].Val(ctx).(float64)
    freeShipping := compiled["free_shipping"].Val(ctx).(bool)
    
    fmt.Printf("折扣: %.0f%%, 免运费: %v\n", discount*100, freeShipping)
}
```

---

## 相关文档

- [使用文档](USAGE.md)
- [性能测试](../benchmark/README.md)
- [主 README](../README.md)

---

**文档版本**: v1.0  
**更新时间**: 2026-06-25
