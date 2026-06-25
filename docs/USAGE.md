# Expr 使用文档

本文档详细介绍 seeadoog/expr 表达式引擎的使用方法、API 和最佳实践。

## 目录

- [快速开始](#快速开始)
- [环境管理](#环境管理)
- [表达式编译](#表达式编译)
- [上下文管理](#上下文管理)
- [语法详解](#语法详解)
- [自定义函数](#自定义函数)
- [类型系统](#类型系统)
- [错误处理](#错误处理)
- [性能优化](#性能优化)
- [最佳实践](#最佳实践)

---

## 快速开始

### 安装

```bash
go get github.com/seeadoog/expr
```

### 第一个示例

```go
package main

import (
    "fmt"
    "github.com/seeadoog/expr"
)

func main() {
    // 1. 创建环境（或使用 DefaultEnv）
    env := expr.NewEnv()
    
    // 2. 编译表达式
    compiled, err := env.ParseValue("a + b * c")
    if err != nil {
        panic(err)
    }
    
    // 3. 创建上下文并设置变量
    ctx := env.NewContext(map[string]any{
        "a": 10.0,
        "b": 20.0,
        "c": 30.0,
    })
    
    // 4. 求值
    result := compiled.Val(ctx)
    fmt.Println(result) // 输出: 610
}
```

---

## 环境管理

### 创建环境

#### 使用默认环境

```go
// 最简单的方式，使用全局默认环境
compiled, err := expr.DefaultEnv.ParseValue("expression")
ctx := expr.DefaultEnv.NewContext(data)
```

#### 创建新环境

```go
// 创建独立的环境实例
env := expr.NewEnv()
```

**使用场景：**
- 使用 `DefaultEnv`：适合简单应用，全局共享
- 创建新环境：适合多租户场景，需要隔离函数和配置

### 环境配置

```go
env := expr.NewEnv()

// 注册自定义函数
env.RegFunc("myFunc", func(x float64) float64 {
    return x * 2
})

// 设置全局变量（所有表达式共享）
// 注意：全局变量在所有 Context 中可见
```

---

## 表达式编译

### 编译方法

#### ParseValue - 编译表达式

```go
compiled, err := env.ParseValue("a + b * c")
if err != nil {
    // 处理语法错误
    fmt.Println("编译错误:", err)
    return
}
```

#### ParseValueToAstNode - 编译为 AST

```go
astNode, err := env.ParseValueToAstNode("a + b")
if err != nil {
    panic(err)
}
// astNode 是抽象语法树节点
```

### 预编译优化

```go
// ❌ 错误：每次循环都编译
for _, data := range dataset {
    compiled, _ := env.ParseValue("a + b")
    result := compiled.Val(ctx)
}

// ✅ 正确：预编译一次，多次使用
compiled, err := env.ParseValue("a + b")
if err != nil {
    panic(err)
}

for _, data := range dataset {
    ctx.SetByString("a", data.A)
    ctx.SetByString("b", data.B)
    result := compiled.Val(ctx)
}
```

### 表达式缓存

```go
type ExpressionCache struct {
    cache map[string]expr.Val
    env   *expr.Env
    mu    sync.RWMutex
}

func (c *ExpressionCache) Get(expression string) (expr.Val, error) {
    // 先尝试从缓存读取
    c.mu.RLock()
    if compiled, ok := c.cache[expression]; ok {
        c.mu.RUnlock()
        return compiled, nil
    }
    c.mu.RUnlock()
    
    // 编译新表达式
    compiled, err := c.env.ParseValue(expression)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    c.mu.Lock()
    c.cache[expression] = compiled
    c.mu.Unlock()
    
    return compiled, nil
}
```

---

## 上下文管理

### 创建上下文

#### 方式 1: NewContext

```go
ctx := env.NewContext(map[string]any{
    "name":  "John",
    "age":   30.0,
    "score": 95.5,
})
```

#### 方式 2: 空上下文后设置

```go
ctx := env.NewContext(nil)
ctx.SetByString("name", "John")
ctx.SetByString("age", 30.0)
```

### 使用对象池

**对象池可以显著提升性能，减少 GC 压力**

```go
env := expr.NewEnv()
compiled, _ := env.ParseValue("expression")

// 从对象池获取
ctx := env.GetContextFromPool()
defer env.PutContext2Pool(ctx) // 使用完毕归还

// 设置变量
ctx.SetByString("a", 10)
ctx.SetByString("b", 20)

// 求值
result := compiled.Val(ctx)
```

### 上下文操作

#### 设置变量

```go
// 按字符串 key 设置
ctx.SetByString("key", value)

// 按 hash 设置（更快，但需要预先计算 hash）
hash := expr.CalcHash("key")
ctx.Set(hash, value)
```

#### 获取变量

```go
// 在表达式中，直接使用变量名
compiled, _ := env.ParseValue("name + ' is ' + string(age)")
```

#### 嵌套对象

```go
ctx := env.NewContext(map[string]any{
    "user": map[string]any{
        "name":  "John",
        "email": "john@example.com",
        "profile": map[string]any{
            "city": "Beijing",
        },
    },
})

// 在表达式中访问
compiled, _ := env.ParseValue("user->name + ' from ' + user->profile->city")
result := compiled.Val(ctx)
```

---

## 语法详解

### 数据类型

#### 数字

```go
// 整数
env.ParseValue("42")

// 浮点数
env.ParseValue("3.14159")

// 负数
env.ParseValue("-10")

// 科学记数法
env.ParseValue("1.23e-4")
```

#### 字符串

```go
// 单引号
env.ParseValue("'hello world'")

// 双引号
env.ParseValue("\"hello world\"")

// 反引号
env.ParseValue("`hello world`")

// 字符串模板
env.ParseValue("'Hello ${name}, you are ${age} years old'")
```

#### 布尔值

```go
env.ParseValue("true")
env.ParseValue("false")
```

#### Nil

```go
env.ParseValue("nil")
env.ParseValue("value or nil")
```

#### 数组

```go
// 字面量定义
env.ParseValue("[1, 2, 3, 4, 5]")
env.ParseValue("['a', 'b', 'c']")

// 混合类型
env.ParseValue("[1, 'hello', true, nil]")

// 嵌套数组
env.ParseValue("[[1, 2], [3, 4]]")
```

#### 对象/Map

```go
// 字面量定义
env.ParseValue("{name: 'John', age: 30}")

// 嵌套对象
env.ParseValue("{user: {name: 'John', profile: {city: 'Beijing'}}}")

// 动态键
env.ParseValue("{[key]: value}")
```

### 运算符详解

#### 算术运算符

```go
// 基本运算
"a + b"      // 加法
"a - b"      // 减法
"a * b"      // 乘法
"a / b"      // 除法
"a % b"      // 取模
"a ^ b"      // 幂运算

// 复合赋值
"a += 5"     // a = a + 5
"a -= 3"     // a = a - 3
"a *= 2"     // a = a * 2
"a /= 4"     // a = a / 4

// 自增自减
"a++"        // 后缀自增
"a--"        // 后缀自减
```

**优先级（从高到低）：**
1. `^` 幂运算
2. `*`, `/`, `%` 乘、除、取模
3. `+`, `-` 加、减

#### 比较运算符

```go
"a == b"     // 相等（值相等）
"a != b"     // 不等
"a === b"    // 严格相等（类型和值都相等）
"a !== b"    // 严格不等
"a > b"      // 大于
"a >= b"     // 大于等于
"a < b"      // 小于
"a <= b"     // 小于等于
```

#### 逻辑运算符

```go
"a && b"     // 逻辑与（短路）
"a || b"     // 逻辑或（短路）
"!a"         // 逻辑非
"a or b"     // 或运算（a 为 nil/false 返回 b）
"a orr b"    // 同 or
```

**短路求值：**
```go
// && 短路：a 为 false 时不计算 b
"false && expensive_function()"

// || 短路：a 为 true 时不计算 b
"true || expensive_function()"
```

#### 位运算符

```go
"a & b"      // 按位与
"a | b"      // 按位或
"a ^ b"      // 按位异或
```

### 控制流

#### 三元运算符

```go
// 基本形式
"condition ? trueValue : falseValue"

// 示例
"age >= 18 ? 'Adult' : 'Minor'"
"score >= 60 ? 'Pass' : 'Fail'"

// 嵌套三元
"score >= 90 ? 'A' : score >= 80 ? 'B' : score >= 70 ? 'C' : 'D'"
```

#### if-elseif-else

```go
// 链式调用
"if(a == 5, c = 5).elseif(a == 6, c = 6).else(c = 9).end()"

// 多条语句
"if(score >= 90, grade = 'A'; message = 'Excellent').else(grade = 'F'; message = 'Failed').end()"
```

#### switch-case

```go
// 基本用法
"switch(value).case(1, result = 'one').case(2, result = 'two').default(result = 'other').end()"

// 多个 case
"switch(day).case(1, name = 'Monday').case(2, name = 'Tuesday').case(3, name = 'Wednesday').default(name = 'Unknown').end()"
```

### 数组操作详解

#### 访问元素

```go
"arr[0]"           // 第一个元素
"arr[len(arr)-1]"  // 最后一个元素
```

#### 切片操作

```go
"arr[1:4]"         // 索引 1 到 4（不含 4）
"arr[:3]"          // 从开始到索引 3
"arr[2:]"          // 从索引 2 到结束
"arr[:]"           // 整个数组（复制）
```

#### 数组方法

```go
"arr.len()"                          // 获取长度
"arr.get(2)"                         // 获取索引 2 的元素
"arr.slice(1, 3)"                    // 切片
"arr.for({k, v} => print(k, v))"     // 遍历
"arr.all({item} => item > 10)"       // 过滤（返回满足条件的元素）
"arr.sort({a, b} => a - b)"          // 排序
```

### 对象操作详解

#### 属性访问

```go
// 点操作符
"user.name"

// 箭头操作符（推荐）
"user->name"

// 索引访问
"user['name']"

// 多级访问
"user->profile->address->city"
```

#### 属性赋值

```go
// 简单赋值
"user.name = 'John'"
"user->name = 'John'"

// 多级赋值（自动创建中间节点）
"a.b.c.d.e = 100"

// 数组索引赋值
"arr[0] = 100"
"obj['key'] = 'value'"
```

#### 对象方法

```go
"obj.get('key')"                     // 获取值
"obj.set('key', 'value')"            // 设置值
"obj.delete('key')"                  // 删除键
"obj.len()"                          // 获取键数量
"obj.for({k, v} => print(k, v))"     // 遍历
```

### Lambda 表达式

#### 单参数 Lambda

```go
// 基本形式
"x => x * 2"

// 使用示例
"arr.for(item => print(item))"
```

#### 多参数 Lambda

```go
// 基本形式
"{x, y} => x + y"

// 使用示例
"arr.for({index, value} => print(index, value))"
```

#### 自定义函数

```go
// 定义函数（必须以 $ 开头）
"$add = {a, b} => a + b"

// 调用函数
"$add(10, 20)"

// 复杂示例
"$factorial = {n} => n <= 1 ? 1 : n * $factorial(n - 1); $factorial(5)"
```

### 字符串操作

#### 拼接

```go
"'Hello' + ' ' + 'World'"
"firstName + ' ' + lastName"
```

#### 模板字符串

```go
"'Hello ${name}'"
"'Your score is ${score} points'"
"'Time: ${time_now()::format(\"2006-01-02 15:04:05\")}'"
```

#### 字符串方法

```go
"str.len()"                          // 长度
"str.has_prefix('Mr.')"              // 前缀判断
"str.has_suffix('.com')"             // 后缀判断
"str.contains('test')"               // 包含判断
"str.split(' ', 2)"                  // 分割
"str.trim_space()"                   // 去除空格
"str.slice(0, 5)"                    // 子串
"str.base64()"                       // Base64 编码
"str.hex()"                          // 十六进制编码
"str.md5()"                          // MD5 哈希
```

---

## 自定义函数

### 注册函数

#### 简单函数

```go
env := expr.NewEnv()

// 注册单参数函数
env.RegFunc("double", func(x float64) float64 {
    return x * 2
})

// 使用
compiled, _ := env.ParseValue("double(21)")
ctx := env.NewContext(nil)
result := compiled.Val(ctx) // 42
```

#### 多参数函数

```go
env.RegFunc("add", func(a, b float64) float64 {
    return a + b
})

env.RegFunc("concat", func(a, b, c string) string {
    return a + b + c
})
```

#### 可变参数函数

```go
env.RegFunc("sum", func(nums ...float64) float64 {
    total := 0.0
    for _, n := range nums {
        total += n
    }
    return total
})

// 使用
compiled, _ := env.ParseValue("sum(1, 2, 3, 4, 5)")
result := compiled.Val(ctx) // 15
```

#### 返回错误的函数

```go
env.RegFunc("divide", func(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
})
```

### 函数命名规范

- 常规函数：使用小写字母和下划线，如 `my_function`
- Lambda 函数：必须以 `$` 开头，如 `$my_func`
- 内置函数：查看完整列表避免冲突

### 函数最佳实践

```go
// ✅ 好的实践
env.RegFunc("calculate_discount", func(price, rate float64) float64 {
    if price < 0 || rate < 0 || rate > 1 {
        return 0 // 返回安全的默认值
    }
    return price * rate
})

// ❌ 避免的做法
env.RegFunc("risky_func", func(x any) any {
    // 不要使用 any 类型，缺乏类型安全
    return x
})
```

---

## 类型系统

### 基本类型

| 类型 | Go 类型 | 示例 |
|------|---------|------|
| 数字 | `float64` | `42`, `3.14` |
| 字符串 | `string` | `'hello'`, `"world"` |
| 布尔 | `bool` | `true`, `false` |
| Nil | `nil` | `nil` |
| 数组 | `[]any` | `[1, 2, 3]` |
| 对象 | `map[string]any` | `{name: 'John'}` |

### 类型转换

```go
// 显式转换
"int(3.14)"          // 转整数 -> 3
"float(42)"          // 转浮点数 -> 42.0
"string(123)"        // 转字符串 -> "123"
"bool(1)"            // 转布尔 -> true
"bytes('hello')"     // 转字节数组

// 隐式转换
"10 + 5"             // 整数运算 -> 15
"'Score: ' + string(95)"  // 字符串拼接
```

### 类型判断

```go
"type(value)"        // 返回类型名称字符串

// 判断类型
"type(value) == 'number'"
"type(value) == 'string'"
"type(value) == 'array'"
"type(value) == 'map'"
```

### 严格相等

```go
// == 只比较值
"1 == 1.0"           // true
"'1' == 1"           // 可能 true（取决于实现）

// === 比较类型和值
"1 === 1"            // true
"1 === 1.0"          // true（都是 number）
"'1' === 1"          // false（类型不同）
```

---

## 错误处理

### 编译错误

```go
compiled, err := env.ParseValue("invalid expression +++")
if err != nil {
    fmt.Println("语法错误:", err)
    return
}
```

### 运行时错误

```go
// 某些运行时错误会返回 error 类型的值
result := compiled.Val(ctx)

// 检查是否为错误
if err, ok := result.(error); ok {
    fmt.Println("运行时错误:", err)
    return
}
```

### 错误处理最佳实践

```go
func EvaluateExpression(env *expr.Env, exprStr string, data map[string]any) (any, error) {
    // 1. 编译错误处理
    compiled, err := env.ParseValue(exprStr)
    if err != nil {
        return nil, fmt.Errorf("编译错误: %w", err)
    }
    
    // 2. 创建上下文
    ctx := env.NewContext(data)
    
    // 3. 求值
    result := compiled.Val(ctx)
    
    // 4. 检查运行时错误
    if err, ok := result.(error); ok {
        return nil, fmt.Errorf("运行时错误: %w", err)
    }
    
    return result, nil
}
```

---

## 性能优化

### 1. 表达式预编译

**影响：巨大**

```go
// ❌ 每次编译：慢
for i := 0; i < 1000000; i++ {
    compiled, _ := env.ParseValue("a + b")
    result := compiled.Val(ctx)
}

// ✅ 预编译：快
compiled, _ := env.ParseValue("a + b")
for i := 0; i < 1000000; i++ {
    result := compiled.Val(ctx)
}
```

**性能提升：100-1000x**

### 2. 使用对象池

**影响：中等**

```go
// ✅ 使用对象池
compiled, _ := env.ParseValue("expression")

for _, data := range bigDataset {
    ctx := env.GetContextFromPool()
    ctx.SetByString("data", data)
    result := compiled.Val(ctx)
    env.PutContext2Pool(ctx)
}
```

**性能提升：减少 GC 压力，提升 20-50%**

### 3. 避免复杂嵌套

```go
// ❌ 过度嵌套
"a.b.c.d.e.f.g.h.i.j + x.y.z.w.v"

// ✅ 简化访问路径
"user->name + user->email"
```

### 4. 使用 Hash 优化

```go
// 如果同一个变量名会被频繁访问
nameHash := expr.CalcHash("name")

// 使用 Hash 设置（更快）
ctx.Set(nameHash, value)
```

### 5. 批量操作

```go
// ❌ 逐个处理
for _, item := range items {
    compiled, _ := env.ParseValue("process(item)")
    result := compiled.Val(ctx)
}

// ✅ 批量处理
compiled, _ := env.ParseValue("items.for({item} => process(item))")
result := compiled.Val(ctx)
```

---

## 最佳实践

### 1. 环境管理

```go
// 全局单例
var GlobalEnv = expr.NewEnv()

func init() {
    // 注册所有自定义函数
    GlobalEnv.RegFunc("myFunc", myFuncImpl)
}

// 使用
func Process(exprStr string, data map[string]any) (any, error) {
    compiled, err := GlobalEnv.ParseValue(exprStr)
    if err != nil {
        return nil, err
    }
    
    ctx := GlobalEnv.NewContext(data)
    return compiled.Val(ctx), nil
}
```

### 2. 表达式缓存

```go
var (
    exprCache = make(map[string]expr.Val)
    cacheMu   sync.RWMutex
)

func GetCompiled(env *expr.Env, exprStr string) (expr.Val, error) {
    cacheMu.RLock()
    if compiled, ok := exprCache[exprStr]; ok {
        cacheMu.RUnlock()
        return compiled, nil
    }
    cacheMu.RUnlock()
    
    compiled, err := env.ParseValue(exprStr)
    if err != nil {
        return nil, err
    }
    
    cacheMu.Lock()
    exprCache[exprStr] = compiled
    cacheMu.Unlock()
    
    return compiled, nil
}
```

### 3. 安全的表达式执行

```go
func SafeEval(env *expr.Env, exprStr string, data map[string]any) (result any, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    compiled, err := env.ParseValue(exprStr)
    if err != nil {
        return nil, fmt.Errorf("compile error: %w", err)
    }
    
    ctx := env.NewContext(data)
    result = compiled.Val(ctx)
    
    if errVal, ok := result.(error); ok {
        return nil, fmt.Errorf("runtime error: %w", errVal)
    }
    
    return result, nil
}
```

### 4. 类型断言

```go
result := compiled.Val(ctx)

// 安全的类型断言
switch v := result.(type) {
case float64:
    fmt.Printf("数字: %.2f\n", v)
case string:
    fmt.Printf("字符串: %s\n", v)
case bool:
    fmt.Printf("布尔: %v\n", v)
case []any:
    fmt.Printf("数组: %v\n", v)
case map[string]any:
    fmt.Printf("对象: %v\n", v)
case error:
    fmt.Printf("错误: %v\n", v)
default:
    fmt.Printf("未知类型: %T\n", v)
}
```

### 5. 配置驱动

```go
// config.json
{
  "rules": {
    "vip_check": "totalSpent > 10000 && orderCount > 50",
    "discount": "vipLevel * 0.05",
    "free_shipping": "totalAmount >= 99"
  }
}

// 加载配置
type Config struct {
    Rules map[string]string `json:"rules"`
}

var compiledRules = make(map[string]expr.Val)

func LoadConfig(configPath string) error {
    data, _ := os.ReadFile(configPath)
    var cfg Config
    json.Unmarshal(data, &cfg)
    
    for name, exprStr := range cfg.Rules {
        compiled, err := expr.DefaultEnv.ParseValue(exprStr)
        if err != nil {
            return fmt.Errorf("rule %s: %w", name, err)
        }
        compiledRules[name] = compiled
    }
    
    return nil
}

// 使用规则
func ApplyRule(ruleName string, data map[string]any) (any, error) {
    compiled, ok := compiledRules[ruleName]
    if !ok {
        return nil, fmt.Errorf("rule not found: %s", ruleName)
    }
    
    ctx := expr.DefaultEnv.NewContext(data)
    return compiled.Val(ctx), nil
}
```

---

## 常见问题

### Q: 如何处理除零错误？

```go
// 在自定义函数中处理
env.RegFunc("safe_divide", func(a, b float64) float64 {
    if b == 0 {
        return 0 // 或返回特殊值
    }
    return a / b
})

// 或在表达式中检查
"b != 0 ? a / b : 0"
```

### Q: 如何访问嵌套很深的对象？

```go
// 使用箭头操作符链式访问
"user->profile->address->city->name"

// 使用非空断言
"user->profile->address!!->city"
```

### Q: 表达式性能不够怎么办？

1. 检查是否预编译了表达式
2. 使用对象池
3. 简化表达式逻辑
4. 考虑缓存中间结果

### Q: 如何调试表达式？

```go
// 使用 print 函数
"print('value:', value); value * 2"

// 分步执行
"step1 = a + b; print('step1:', step1); step2 = step1 * c; print('step2:', step2); step2"
```

---

## 下一步

- 查看 [API 参考](API.md)
- 阅读 [性能优化指南](../benchmark/BENCHMARK_GUIDE.md)
- 查看 [示例代码](../examples/)

---

**文档版本**: v1.0  
**更新时间**: 2026-06-25
