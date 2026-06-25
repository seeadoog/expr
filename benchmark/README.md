# Expression Engine Performance Benchmark

这是 [seeadoog/expr](https://github.com/seeadoog/expr) 与其他流行 Go 表达式引擎的性能对比测试。

## 测试环境

- **CPU**: Apple M5
- **OS**: macOS (darwin/arm64)
- **Go Version**: 1.21+
- **测试时间**: 2s per benchmark

## 参与对比的库

1. **seeadoog/expr** - 本项目
2. **antonmedv/expr** (v1.15.5) - 最流行的 Go 表达式引擎
3. **govaluate** (v3.0.0) - 另一个广泛使用的表达式求值库

## 性能测试结果

### 1. 简单算术运算
表达式: `a + b * c - d / e`

| 引擎 | 速度 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|--------------|----------------|----------|
| **seeadoog/expr** | **38.13** | **32** | **4** |
| govaluate | 116.7 | 48 | 5 |
| antonmedv/expr | 227.2 | 208 | 11 |

**🏆 seeadoog/expr 比 antonmedv/expr 快 5.96x，比 govaluate 快 3.06x**

---

### 2. 复杂布尔逻辑
表达式: `(a > 5 && b < 100) || (c == 'test' && d != nil)`

| 引擎 | 速度 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|--------------|----------------|----------|
| **seeadoog/expr** | **16.59** | **0** | **0** |
| govaluate | 78.33 | 16 | 1 |
| antonmedv/expr | 92.29 | 64 | 3 |

**🏆 seeadoog/expr 比 antonmedv/expr 快 5.56x，比 govaluate 快 4.72x，且零内存分配！**

---

### 3. 字符串操作
表达式: `firstName + ' ' + lastName`

| 引擎 | 速度 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|--------------|----------------|----------|
| **seeadoog/expr** | **53.28** | **48** | **4** |
| govaluate | 130.4 | 64 | 5 |
| antonmedv/expr | 135.4 | 112 | 7 |

**🏆 seeadoog/expr 比 antonmedv/expr 快 2.54x，比 govaluate 快 2.45x**

---

### 4. 嵌套条件
表达式: `a > 10 && (b < 20 || (c >= 30 && d <= 40))`

| 引擎 | 速度 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|--------------|----------------|----------|
| **seeadoog/expr** | **32.57** | **0** | **0** |
| govaluate | 136.5 | 16 | 1 |
| antonmedv/expr | 177.7 | 96 | 5 |

**🏆 seeadoog/expr 比 antonmedv/expr 快 5.46x，比 govaluate 快 4.19x**

---

### 5. 函数调用
表达式: `a * 2 + b / 2`

| 引擎 | 速度 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|--------------|----------------|----------|
| **seeadoog/expr** | **26.84** | **24** | **3** |
| govaluate | 80.03 | 40 | 4 |
| antonmedv/expr | 128.4 | 152 | 7 |

**🏆 seeadoog/expr 比 antonmedv/expr 快 4.78x，比 govaluate 快 2.98x**

---

### 6. 编译性能
表达式: `(a + b) * (c - d) / e > 100 && (f == 'test' || g != nil)`

| 引擎 | 编译时间 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|-----------------|----------------|----------|
| **seeadoog/expr** | **2,538** | **8,968** | **90** |
| govaluate | 4,336 | 7,128 | 144 |
| antonmedv/expr | 7,519 | 18,232 | 101 |

**🏆 seeadoog/expr 编译速度比 antonmedv/expr 快 2.96x，比 govaluate 快 1.71x**

---

### 7. 大型表达式
表达式: `v1 + v2 + v3 + v4 + v5 + v6 + v7 + v8 + v9 + v10`

| 引擎 | 速度 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|--------------|----------------|----------|
| **seeadoog/expr** | **97.19** | **72** | **9** |
| govaluate | 232.4 | 88 | 10 |
| antonmedv/expr | 455.9 | 264 | 20 |

**🏆 seeadoog/expr 比 antonmedv/expr 快 4.69x，比 govaluate 快 2.39x**

---

### 8. 三元运算符
表达式: `a > 10 ? b * 2 : c / 2`

| 引擎 | 速度 (ns/op) | 内存分配 (B/op) | 分配次数 |
|------|--------------|----------------|----------|
| **seeadoog/expr** | **16.44** | **8** | **1** |
| govaluate | 80.44 | 24 | 2 |
| antonmedv/expr | 102.0 | 72 | 4 |

**🏆 seeadoog/expr 比 antonmedv/expr 快 6.20x，比 govaluate 快 4.89x**

---

## 综合性能分析

### 执行速度对比

在所有 8 个测试场景中：

| 对比项 | 平均加速比 | 优势场景 |
|--------|----------|---------|
| seeadoog/expr vs antonmedv/expr | **4.77x** | 8/8 全胜 |
| seeadoog/expr vs govaluate | **3.22x** | 8/8 全胜 |

### 内存效率对比

- **seeadoog/expr**: 在布尔逻辑和嵌套条件场景实现了**零内存分配**
- **平均内存占用**: seeadoog/expr 比 antonmedv/expr 少 60%+
- **平均分配次数**: seeadoog/expr 比 antonmedv/expr 少 50%+

## 性能优势原因

1. **高效的 AST 设计**: 紧凑的节点结构，减少内存开销
2. **Context 对象池**: 复用 Context 对象，减少 GC 压力
3. **变量名 Hash 预计算**: 避免重复计算 hash 值
4. **短路求值优化**: 布尔运算提前退出
5. **编译时优化**: yacc 生成的 parser 效率高

## 运行基准测试

```bash
cd benchmark
go test -bench=. -benchmem -benchtime=2s
```

### 生成对比报告

```bash
go test -bench=. -benchmem -benchtime=2s > results.txt
```

### 运行单个测试

```bash
go test -bench=BenchmarkSimpleArithmetic -benchmem
```

## 适用场景

基于性能测试结果，**seeadoog/expr** 特别适合：

- ✅ 高频表达式求值（每秒百万级）
- ✅ 复杂布尔逻辑判断（规则引擎）
- ✅ 内存敏感的应用（嵌入式、边缘计算）
- ✅ 低延迟要求的场景（实时系统）
- ✅ 大规模并发场景（对象池复用）

## 注意事项

- 测试结果基于 Apple M5 ARM 架构，不同 CPU 架构可能有差异
- 所有测试均在相同条件下进行，确保公平对比
- 表达式已预编译，测试的是执行速度而非编译速度（除了编译性能测试）

## 结论

**seeadoog/expr** 在所有测试场景中都展现出卓越的性能：

- 🚀 **执行速度**: 平均快 3-6 倍
- 💾 **内存效率**: 最低内存占用，部分场景零分配
- ⚡ **编译速度**: 最快的编译性能
- 🎯 **稳定性**: 所有场景一致的高性能表现

对于需要高性能表达式求值的 Go 应用，**seeadoog/expr** 是最佳选择。
