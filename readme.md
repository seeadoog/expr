# Expr - 高性能 Go 表达式引擎

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Performance](https://img.shields.io/badge/Performance-4.77x_faster-brightgreen)](benchmark/)

一个高性能、功能丰富的 Go 表达式引擎和规则引擎，专为高频求值和低延迟场景设计。

## ✨ 核心特性

- 🚀 **极致性能** - 比主流表达式引擎快 3-6 倍，部分场景零内存分配
- 📝 **语法丰富** - 支持算术、逻辑、字符串、数组、对象、Lambda、条件语句等
- 🔧 **易于扩展** - 支持自定义函数、变量、类型
- 💾 **内存高效** - 对象池设计，GC 友好
- 🎯 **类型系统** - 完善的类型系统和方法调用
- 📊 **生产就绪** - 经过大规模生产环境验证

## 🎯 性能对比

| 引擎 | 执行速度 | 内存分配 | 编译速度 |
|------|---------|---------|---------|
| **seeadoog/expr** | **38 ns/op** | **32 B/op** | **2.5 μs** |
| antonmedv/expr | 227 ns/op | 208 B/op | 7.5 μs |
| govaluate | 116 ns/op | 48 B/op | 4.3 μs |

*基于简单算术表达式测试，完整性能报告见 [benchmark/](benchmark/)*

**性能优势：**
- 🏆 比 antonmedv/expr 快 **4.77x**
- 🏆 比 govaluate 快 **3.22x**
- 🏆 布尔逻辑场景：**零内存分配**

## 📦 安装

```bash
go get github.com/seeadoog/expr
```

## 🚀 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/seeadoog/expr"
)

func main() {
    // 使用默认环境
    compiled, err := expr.DefaultEnv.ParseValue("a + b * c")
    if err != nil {
        panic(err)
    }
    
    // 创建上下文并设置变量
    ctx := expr.DefaultEnv.NewContext(map[string]any{
        "a": 10.0,
        "b": 20.0,
        "c": 30.0,
    })
    
    // 求值
    result := compiled.Val(ctx)
    fmt.Println(result) // 输出: 610
}
```

### 字符串模板

```go
// 支持变量嵌入和函数调用
expr := `'Hello ${name}, current time is ${get_cur_time()::format("2006-01-02 15:04:05")}'`
compiled, _ := expr2.DefaultEnv.ParseValue(expr)

ctx := expr2.DefaultEnv.NewContext(map[string]any{
    "name": "world",
})

result := compiled.Val(ctx)
fmt.Println(result)
```

## 📖 语法特性详解

### 1. 基本类型

```go
// 数字
123
45.67
-10

// 字符串（支持三种引号）
'hello world'
"hello world"
`hello world`

// 布尔值
true
false

// Nil
nil

// 数组
[1, 2, 3, 4, 5]
['a', 'b', 'c']

// Map/对象
{name: 'John', age: 30}
{key1: value1, key2: value2}
```

### 2. 算术运算

```go
a + b        // 加法
a - b        // 减法
a * b        // 乘法
a / b        // 除法
a % b        // 取模
a ^ b        // 幂运算

// 复合赋值运算符（新增）
a += 5       // a = a + 5
a -= 3       // a = a - 3
a *= 2       // a = a * 2
a /= 4       // a = a / 4

// 自增自减
a++          // 后缀自增
a--          // 后缀自减
```

### 3. 比较运算

```go
a == b       // 相等
a != b       // 不等
a > b        // 大于
a >= b       // 大于等于
a < b        // 小于
a <= b       // 小于等于
a === b      // 严格相等（类型 + 值）
a !== b      // 严格不等
```

### 4. 逻辑运算

```go
a && b       // 逻辑与
a || b       // 逻辑或
!a           // 逻辑非
a or b       // 或运算（a 为 nil/false 返回 b）
a orr b      // 同 or
```

### 5. 位运算

```go
a & b        // 按位与
a | b        // 按位或
a ^ b        // 按位异或
```

### 6. 字符串操作

```go
// 字符串拼接
'Hello' + ' ' + 'World'

// 字符串模板
'Hello ${name}'
'Result: ${a + b}'

// 字符串方法
name.len()                    // 长度
name.has_prefix('Mr.')        // 前缀判断
name.has_suffix('.com')       // 后缀判断
name.contains('test')         // 包含判断
name.split(' ', 2)            // 分割
name.trim_space()             // 去空格
name.slice(0, 5)              // 切片
name.base64()                 // Base64 编码
name.hex()                    // 十六进制编码
name.md5()                    // MD5 哈希
```

### 7. 数组操作

```go
// 数组定义
arr = [1, 2, 3, 4, 5]

// 索引访问
arr[0]           // 第一个元素
arr[3]           // 第四个元素

// 切片
arr[0:3]         // 索引 0 到 3（不含 3）
arr[:3]          // 从开始到索引 3
arr[2:]          // 从索引 2 到结束

// 数组赋值
arr[3] = 100     // 修改元素

// 数组方法
arr.len()                     // 长度
arr.get(2)                    // 获取元素
arr.slice(1, 3)               // 切片
arr.sort(compare_func)        // 排序
arr.all(condition)            // 过滤
arr.foreach({k, v} => print(k, v)) // 遍历
```

### 8. 对象/Map 操作

```go
// 对象定义
obj = {name: 'John', age: 30}

// 属性访问（多种方式）
obj.name         // 点操作符
obj['name']      // 索引访问
object->name     // 箭头操作符

// 属性赋值
obj.name = 'Jane'
obj['age'] = 31

// 多级属性（自动创建中间节点）
a.b.c.d.e.f.g = 1

// Map 方法
obj.get('name')              // 获取值
obj.set('key', 'value')      // 设置值
obj.delete('key')            // 删除键
obj.len()                    // 长度
obj.foreach({k, v} => print(k, v)) // 遍历
```

### 9. 三元运算符

```go
condition ? trueValue : falseValue

// 示例
age >= 18 ? 'Adult' : 'Minor'
score >= 60 ? 'Pass' : 'Fail'

// 嵌套
score >= 90 ? 'A' : score >= 80 ? 'B' : score >= 70 ? 'C' : 'D'

// 与 or 结合
name or 'Anonymous'
user ? user->name : 'Guest'
```

### 10. Lambda 表达式

```go
// 单参数 Lambda
x => x * 2

// 多参数 Lambda
{x, y} => x + y

// 自定义函数（以 $ 开头）
$func_def = {a, b} => (a + b)
$func_def(1, 2)  // 调用自定义函数

// Lambda 用于数组操作
arr.foreach({k, v} => print(k, v))
arr.all({item} => item > 10)
```

### 11. 条件语句

```go
// if-elseif-else 链式调用
if a == b then
   c= d 
elseif c == d then 
   
end 
```

### 12. 变量赋值

```go
// 简单赋值
a = 5

// 对象属性赋值
object.name = 'John'

// 数组元素赋值
arr[3] = 100

// 多个表达式（返回最后一个）
a = 4; b = 5; c = 6
```

### 13. 注释

```go
a = 5  # 这是注释
# 支持单行注释
```

### 14. 非空断言

```go
a.b.c!!        // 取值并要求不为 nil，否则退出执行
value@         // 另一种非空断言语法
```

## 🔧 内置函数

### 字符串函数

```go
str_to_upper(s)              // 转大写
str_to_lower(s)              // 转小写
str_trim(s)                  // 去除空格
str_split(s, sep, n)         // 分割字符串
str_join(arr, sep)           // 连接字符串
str_fields(s)                // 按空格分割
has_prefix(s, prefix)        // 前缀判断
has_suffix(s, suffix)        // 后缀判断
sprintf(format, args...)     // 格式化字符串
```

### 数学函数

```go
add(a, b)                    // 加法
sub(a, b)                    // 减法
mul(a, b)                    // 乘法
div(a, b)                    // 除法
mod(a, b)                    // 取模
pow(a, b)                    // 幂运算
```

### 类型转换

```go
int(x)                       // 转整数
number(x)                    // 转数字
string(x)                    // 转字符串
bool(x) / boolean(x)         // 转布尔值
bytes(x)                     // 转字节数组
```

### 数组/集合函数

```go
len(arr)                     // 长度
slice_new(items...)          // 创建数组
slice_init(size)             // 初始化数组
append(arr, items...)        // 追加元素
get(arr, index)              // 获取元素
set(arr, index, value)       // 设置元素
set_index(arr, index, value) // 同 set
slice_cut(arr, start, end)   // 切片
foreach(arr, lambda)             // 遍历
all(arr, condition)          // 过滤
range(n)                     // 生成 0 到 n-1 的数组
```

### JSON 操作

```go
json_to(obj)                 // 对象转 JSON 字符串
json_from(str)               // JSON 字符串转对象
to_json_str(obj)             // 同 json_to
to_json_obj(str)             // 同 json_from
```

### 时间函数

```go
time_now()                   // 当前时间
time_now_mill()              // 当前毫秒时间戳
time_from_unix(timestamp)    // Unix 时间戳转时间
time_parse(layout, value)    // 解析时间字符串
time_format(time, layout)    // 格式化时间

// 时间对象方法
time.year()                  // 年
time.month()                 // 月
time.day()                   // 日
time.hour()                  // 小时
time.minute()                // 分钟
time.second()                // 秒
time.unix()                  // Unix 时间戳
time.unix_mill()             // 毫秒时间戳
time.format(layout)          // 格式化
time.add_mill(ms)            // 添加毫秒
time.sub(other)              // 时间差
```

### 编码函数

```go
base64_encode(data)          // Base64 编码
base64_decode(data)          // Base64 解码
hex_encode(data)             // 十六进制编码
hex_decode(data)             // 十六进制解码
md5(data)                    // MD5 哈希
sha256(data)                 // SHA256 哈希
hmac_sha256(data, key)       // HMAC-SHA256
```

### 正则表达式

```go
regexp_new(pattern)          // 创建正则表达式
regex.match(str)             // 匹配字符串
```

### 控制流函数

```go
if(cond, expr...)            // 条件判断
ternary(cond, true, false)   // 三元运算
loop(expr...)                // 循环
return(value...)             // 返回
catch(expr)                  // 捕获异常
recover(expr)                // 恢复执行
```

### 其他函数

```go
print(args...)               // 打印
type(value)                  // 获取类型
new()                        // 创建新 Map
sleep(ms)                    // 休眠（毫秒）
go(expr)                     // 异步执行
exec(expr...)                // 执行表达式
unwrap(value)                // 解包装
```

## 💡 实际应用场景

### 1. 规则引擎

```go
// 定义业务规则
rules := []struct {
    name string
    expr string
    action string
}{
    {"VIP 折扣", "totalSpent > 10000 && orderCount > 50", "discount = 0.2"},
    {"新用户优惠", "registerDays < 30", "discount = 0.1"},
    {"满减活动", "totalAmount > 99", "shipping = 0"},
}

env := expr.NewEnv()

// 编译规则
for _, rule := range rules {
    condition, _ := env.ParseValue(rule.expr)
    action, _ := env.ParseValue(rule.action)
    
    ctx := env.NewContext(userData)
    
    if condition.Val(ctx).(bool) {
        action.Val(ctx)
        fmt.Printf("应用规则: %s\n", rule.name)
    }
}
```

### 2. 动态表单验证

```go
// 表单验证规则
validations := map[string]string{
    "email":    "len(email) > 0 && email.contains('@')",
    "age":      "age >= 18 && age <= 100",
    "password": "len(password) >= 8",
    "username": "len(username) >= 3 && len(username) <= 20",
}

env := expr.NewEnv()
formData := map[string]any{
    "email":    "user@example.com",
    "age":      25.0,
    "password": "secret123",
    "username": "john",
}

ctx := env.NewContext(formData)

for field, rule := range validations {
    compiled, _ := env.ParseValue(rule)
    if !compiled.Val(ctx).(bool) {
        fmt.Printf("%s 验证失败\n", field)
    }
}
```

### 3. 配置驱动的业务逻辑

```go
// 从配置文件加载规则
config := map[string]string{
    "canCheckout": "cartTotal > 0 && user.isVerified && !user.isBanned",
    "shippingFee": "cartTotal >= 99 ? 0 : cartTotal >= 50 ? 5 : 10",
    "discount":    "user.vipLevel * 0.05",
}

env := expr.NewEnv()
compiledRules := make(map[string]expr.Val)

for key, exprStr := range config {
    compiled, _ := env.ParseValue(exprStr)
    compiledRules[key] = compiled
}

// 使用规则
orderData := map[string]any{
    "cartTotal": 120.0,
    "user": map[string]any{
        "isVerified": true,
        "isBanned":   false,
        "vipLevel":   3.0,
    },
}

ctx := env.NewContext(orderData)
canCheckout := compiledRules["canCheckout"].Val(ctx).(bool)
shippingFee := compiledRules["shippingFee"].Val(ctx).(float64)
discount := compiledRules["discount"].Val(ctx).(float64)

fmt.Printf("Can Checkout: %v, Shipping: %.2f, Discount: %.0f%%\n", 
    canCheckout, shippingFee, discount*100)
```

### 4. 权限系统

```go
permissions := map[string]string{
    "read":   "true",  // 所有人可读
    "write":  "role == 'admin' || role == 'editor'",
    "delete": "role == 'admin'",
    "admin":  "role == 'admin' && department == 'IT'",
}

env := expr.NewEnv()
user := map[string]any{
    "role":       "editor",
    "department": "Marketing",
}

ctx := env.NewContext(user)

for action, rule := range permissions {
    compiled, _ := env.ParseValue(rule)
    allowed := compiled.Val(ctx).(bool)
    fmt.Printf("%s: %v\n", action, allowed)
}
```

## 📊 性能优化建议

### 1. 预编译表达式

```go
// ❌ 不好：每次都编译
for _, data := range dataset {
    compiled, _ := env.ParseValue("a + b")
    result := compiled.Val(ctx)
}

// ✅ 好：只编译一次
compiled, _ := env.ParseValue("a + b")
for _, data := range dataset {
    ctx.SetByString("a", data.A)
    ctx.SetByString("b", data.B)
    result := compiled.Val(ctx)
}
```

### 2. 使用对象池

```go
env := expr.NewEnv()
compiled, _ := env.ParseValue("complex expression")

for _, data := range dataset {
    // 从对象池获取
    ctx := env.GetContextFromPool()
    
    ctx.SetByString("data", data)
    result := compiled.Val(ctx)
    
    // 归还到对象池
    env.PutContext2Pool(ctx)
}
```

### 3. 使用 DefaultEnv

```go
// 直接使用全局默认环境
compiled, _ := expr.DefaultEnv.ParseValue("expression")
ctx := expr.DefaultEnv.NewContext(data)
result := compiled.Val(ctx)
```

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行性能测试
cd benchmark
make bench

# 查看测试覆盖率
go test -cover ./...
```

## 📈 性能基准测试

完整的性能测试报告请查看 [benchmark/](benchmark/) 目录。

快速运行性能测试：

```bash
cd benchmark

# 快速测试
make bench-short

# 完整测试（推荐）
make bench

# 扩展测试
make bench-all

# 查看帮助
make help
```

**测试结果摘要：**

| 场景 | seeadoog/expr | antonmedv/expr | 性能提升 |
|------|---------------|----------------|---------|
| 简单算术 | 38 ns/op | 227 ns/op | 5.96x |
| 布尔逻辑 | 17 ns/op | 92 ns/op | 5.56x |
| 字符串操作 | 53 ns/op | 135 ns/op | 2.54x |
| 三元运算符 | 16 ns/op | 102 ns/op | 6.20x |

详细报告：[benchmark/README.md](benchmark/README.md)

## 📚 文档

- [完整使用文档](docs/USAGE.md)
- [性能测试报告](benchmark/README.md)
- [性能优化指南](benchmark/BENCHMARK_GUIDE.md)
- [API 参考](docs/API.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证

## 🙏 致谢

- 感谢所有贡献者的支持
- 性能测试对比了 [antonmedv/expr](https://github.com/antonmedv/expr) 和 [govaluate](https://github.com/Knetic/govaluate)

## ⭐ Star History

如果这个项目对你有帮助，请给我们一个 Star ⭐！
