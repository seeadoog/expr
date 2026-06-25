# Benchmark 项目说明

## 项目结构

```
benchmark/
├── go.mod                  # 独立的 Go module
├── benchmark_test.go       # 性能测试代码
├── README.md              # 详细的性能对比报告
├── Makefile               # 便捷的测试命令
├── run_benchmark.sh       # 自动化测试脚本
└── .gitignore            # Git 忽略文件
```

## 快速开始

### 1. 安装依赖

```bash
cd benchmark
go mod tidy
```

### 2. 运行基准测试

```bash
# 使用 Makefile（推荐）
make bench

# 或直接使用 go test
go test -bench=. -benchmem -benchtime=2s
```

## 可用命令

### 使用 Makefile

```bash
make bench          # 运行完整基准测试 (2s)
make bench-short    # 快速测试 (100ms)
make bench-all      # 扩展测试 (5s)
make bench-save     # 运行并保存结果
make bench-compare  # 与上次结果对比
make bench-cpu      # CPU 性能分析
make bench-mem      # 内存性能分析
make clean          # 清理生成文件
make help           # 显示帮助信息
```

### 单项测试

```bash
make bench-simple   # 只测试简单算术
make bench-bool     # 只测试布尔逻辑
make bench-string   # 只测试字符串操作
make bench-compile  # 只测试编译性能
```

## 测试场景

1. **SimpleArithmetic** - 简单算术运算
2. **ComplexBoolean** - 复杂布尔逻辑
3. **StringOperations** - 字符串操作
4. **NestedConditions** - 嵌套条件
5. **FunctionCalls** - 函数调用
6. **Compilation** - 编译性能
7. **LargeExpression** - 大型表达式
8. **TernaryOperator** - 三元运算符

## 性能分析

### CPU 性能分析

```bash
make bench-cpu
go tool pprof cpu.prof
```

在 pprof 交互模式中：
```
top10        # 显示 CPU 占用前 10
list funcName # 查看函数详细信息
web          # 生成调用图（需要 graphviz）
```

### 内存性能分析

```bash
make bench-mem
go tool pprof mem.prof
```

## 测试结果解读

### 输出示例

```
BenchmarkSimpleArithmetic/seeadoog/expr-10    60558144    39.58 ns/op    32 B/op    4 allocs/op
```

- `60558144`: 测试迭代次数
- `39.58 ns/op`: 每次操作耗时（纳秒）
- `32 B/op`: 每次操作内存分配（字节）
- `4 allocs/op`: 每次操作分配次数

### 性能对比

查看 [README.md](README.md) 获取详细的性能对比报告和分析。

## 添加新的测试场景

1. 在 `benchmark_test.go` 中添加新的测试函数：

```go
func BenchmarkYourTest(b *testing.B) {
    expression := "your expression here"
    params := map[string]interface{}{
        "var1": value1,
    }

    b.Run("seeadoog/expr", func(b *testing.B) {
        env := expr.NewEnv()
        compiled, _ := env.ParseValue(expression)
        ctx := env.NewContext(params)

        b.ResetTimer()
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = compiled.Val(ctx)
        }
    })

    // 添加其他引擎的测试...
}
```

2. 运行测试：

```bash
go test -bench=BenchmarkYourTest -benchmem
```

## 环境要求

- Go 1.21+
- 支持的操作系统：Linux, macOS, Windows

## 依赖版本

- `github.com/seeadoog/expr` - 本地版本
- `github.com/antonmedv/expr` v1.15.5
- `github.com/Knetic/govaluate` v3.0.0

## 注意事项

1. **预热**: 所有测试都会自动预热，前几次迭代不计入结果
2. **时间**: 使用 `-benchtime` 控制测试时长，默认 2s
3. **CPU**: 结果会因 CPU 架构和频率而异
4. **GC**: Go 的 GC 可能影响结果，多次运行取平均值
5. **负载**: 测试时关闭其他高负载程序

## 持续集成

可以将基准测试集成到 CI/CD 流程：

```yaml
# GitHub Actions 示例
- name: Run Benchmarks
  run: |
    cd benchmark
    go test -bench=. -benchmem -benchtime=2s
```

## 性能回归检测

使用 `benchstat` 工具对比两次测试结果：

```bash
# 安装 benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# 对比结果
make bench-save  # 第一次
# ... 修改代码 ...
make bench > new.txt
benchstat benchmark_results.txt new.txt
```

## 贡献

欢迎提交新的测试场景或改进建议！

## 许可证

与主项目保持一致
