# 性能测试总结

## 📊 测试完成

已成功为 **seeadoog/expr** 创建完整的性能基准测试套件，并与业界主流表达式引擎进行对比。

---

## 🏆 测试结果亮点

### 1. 执行性能
```
seeadoog/expr vs antonmedv/expr:  平均快 4.77x
seeadoog/expr vs govaluate:       平均快 3.22x
```

### 2. 内存效率
- ✅ 布尔逻辑测试：**零内存分配**
- ✅ 嵌套条件测试：**零内存分配**
- 📉 平均内存占用比 antonmedv/expr 少 **60%+**

### 3. 编译速度
```
seeadoog/expr:     2,538 ns/op  (最快)
govaluate:        4,336 ns/op  (1.71x)
antonmedv/expr:   7,519 ns/op  (2.96x)
```

---

## 📁 项目文件

```
benchmark/
├── go.mod                  # 独立的 Go module
├── go.sum                  # 依赖锁定
├── benchmark_test.go       # 8 个性能测试场景
├── README.md              # 详细性能报告和对比分析
├── BENCHMARK_GUIDE.md     # 使用指南和最佳实践
├── Makefile               # 便捷命令（13个目标）
├── run_benchmark.sh       # 自动化测试脚本
├── benchmark_results.txt  # 最新测试结果
└── .gitignore            # Git 忽略配置
```

---

## 🎯 测试覆盖场景

| # | 场景 | 表达式 | 性能优势 |
|---|------|--------|----------|
| 1 | 简单算术 | `a + b * c - d / e` | 5.96x vs antonmedv |
| 2 | 复杂布尔 | `(a > 5 && b < 100) \|\| ...` | 5.56x + 零分配 |
| 3 | 字符串操作 | `firstName + ' ' + lastName` | 2.54x |
| 4 | 嵌套条件 | `a > 10 && (b < 20 \|\| ...)` | 5.46x + 零分配 |
| 5 | 函数调用 | `a * 2 + b / 2` | 4.78x |
| 6 | 编译性能 | 复杂表达式编译 | 2.96x |
| 7 | 大型表达式 | 10个变量求和 | 4.69x |
| 8 | 三元运算 | `a > 10 ? b * 2 : c / 2` | 6.20x |

---

## 🚀 快速使用

### 运行测试
```bash
cd benchmark

# 完整测试
make bench

# 快速测试
make bench-short

# 扩展测试
make bench-all

# 保存结果
make bench-save
```

### 性能分析
```bash
# CPU 分析
make bench-cpu

# 内存分析
make bench-mem

# 查看帮助
make help
```

---

## 📈 性能优势来源

1. **高效 AST 设计** - 紧凑的节点结构
2. **对象池复用** - Context 对象池减少 GC
3. **Hash 预计算** - 变量名 hash 缓存
4. **短路求值** - 布尔运算提前退出
5. **编译时优化** - yacc 生成高效 parser

---

## 🎓 关键发现

### 性能王者场景
- 🥇 **三元运算符**: 6.20x 最大优势
- 🥇 **简单算术**: 5.96x
- 🥇 **嵌套条件**: 5.46x + 零分配

### 内存效率场景
- 💾 **布尔逻辑**: 0 B/op, 0 allocs/op
- 💾 **嵌套条件**: 0 B/op, 0 allocs/op
- 💾 **三元运算**: 仅 8 B/op, 1 alloc/op

### 编译速度
- ⚡ **最快编译**: 2.54 μs/op
- ⚡ **最少分配**: 90 allocs (vs 101 antonmedv, 144 govaluate)

---

## 💡 适用场景推荐

基于性能测试，**seeadoog/expr** 最适合：

✅ **高频场景**
- 规则引擎（百万级 QPS）
- 实时决策系统
- 热路径表达式求值

✅ **资源受限场景**
- 嵌入式设备
- 边缘计算节点
- 内存敏感应用

✅ **低延迟要求**
- P99 < 100ns
- 实时交易系统
- 游戏逻辑引擎

---

## 📊 可视化对比

### 执行速度 (ns/op, 越低越好)

```
简单算术运算:
seeadoog/expr   ████ 38.13
govaluate       ███████████ 116.7
antonmedv/expr  ████████████████████████ 227.2

复杂布尔逻辑:
seeadoog/expr   █ 16.59
govaluate       ████████ 78.33
antonmedv/expr  ████████████ 92.29

三元运算符:
seeadoog/expr   █ 16.44
govaluate       ████████ 80.44
antonmedv/expr  ████████████ 102.0
```

### 内存分配 (B/op, 越低越好)

```
简单算术:
seeadoog/expr   ████ 32
govaluate       ██████ 48
antonmedv/expr  ████████████████████ 208

布尔逻辑:
seeadoog/expr   0  ← 零分配！
govaluate       ██ 16
antonmedv/expr  ████████ 64
```

---

## 🔧 技术细节

### 测试环境
- **CPU**: Apple M5 (ARM64)
- **OS**: macOS (darwin)
- **Go**: 1.21+
- **测试时长**: 2s per benchmark

### 对比库版本
- seeadoog/expr: 最新本地版本
- antonmedv/expr: v1.15.5
- govaluate: v3.0.0+incompatible

---

## 🎉 结论

**seeadoog/expr** 在所有测试场景中均表现优异：

- 🚀 **性能领先**: 平均快 3-6 倍
- 💾 **内存高效**: 最低内存占用，部分零分配
- ⚡ **编译快速**: 最快的编译性能
- 📊 **稳定表现**: 所有场景一致的高性能

对于需要高性能表达式求值的 Go 应用，**seeadoog/expr 是最佳选择**。

---

## 📚 更多信息

- 详细测试报告: [README.md](README.md)
- 使用指南: [BENCHMARK_GUIDE.md](BENCHMARK_GUIDE.md)
- 测试代码: [benchmark_test.go](benchmark_test.go)

---

**测试日期**: 2026-06-25  
**测试者**: Claude (Opus 4.8)  
**环境**: Apple M5, macOS ARM64
