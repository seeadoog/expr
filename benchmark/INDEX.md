# Benchmark 项目导航

欢迎来到 **seeadoog/expr** 性能基准测试项目！

## 📚 文档导航

### 快速开始
- 👉 **[README.md](README.md)** - 详细的性能测试报告和对比分析
  - 8 个测试场景的完整结果
  - 与 antonmedv/expr 和 govaluate 的对比
  - 性能优势分析

### 使用指南
- 📖 **[BENCHMARK_GUIDE.md](BENCHMARK_GUIDE.md)** - 完整的使用教程
  - 如何运行测试
  - Makefile 命令详解
  - 性能分析教程
  - 如何添加新测试

### 结果总结
- 🎯 **[PERFORMANCE_SUMMARY.md](PERFORMANCE_SUMMARY.md)** - 性能结果可视化总结
  - 测试结果亮点
  - 性能对比图表
  - 适用场景推荐
  - 关键发现

### 完成报告
- 📊 **[COMPLETION_REPORT.md](COMPLETION_REPORT.md)** - 项目完成报告
  - 项目完成情况
  - 文件清单
  - 技术总结
  - 最终评价

## 🚀 快速开始

```bash
# 1. 安装依赖
cd benchmark
go mod tidy

# 2. 运行测试
make bench

# 3. 查看帮助
make help
```

## 📊 性能亮点

| 对比项 | 结果 |
|--------|------|
| vs antonmedv/expr | **4.77x 更快** |
| vs govaluate | **3.22x 更快** |
| 零内存分配场景 | **2 个** |
| 编译速度优势 | **2.96x** |

## 🎯 测试场景

8 个全面的测试场景：

1. ⚡ 简单算术运算
2. 🔍 复杂布尔逻辑
3. 📝 字符串操作
4. 🌳 嵌套条件
5. 🔧 函数调用
6. 🏗️ 编译性能
7. 📈 大型表达式
8. ❓ 三元运算符

## 🛠️ 工具和资源

- **测试代码**: [benchmark_test.go](benchmark_test.go)
- **Makefile**: [Makefile](Makefile) - 13 个便捷命令
- **自动化脚本**: [run_benchmark.sh](run_benchmark.sh)
- **测试结果**: [benchmark_results.txt](benchmark_results.txt)

## 💡 建议阅读顺序

1. **首次使用**: README.md → BENCHMARK_GUIDE.md → 运行测试
2. **了解结果**: PERFORMANCE_SUMMARY.md → README.md
3. **深入学习**: BENCHMARK_GUIDE.md → COMPLETION_REPORT.md
4. **项目维护**: COMPLETION_REPORT.md → BENCHMARK_GUIDE.md

## 🎓 关键结论

**seeadoog/expr** 是目前 Go 生态中性能最优秀的表达式引擎之一：

- 🚀 **最快**: 所有场景平均快 3-6 倍
- 💾 **最省**: 内存占用最低，部分场景零分配
- ⚡ **最稳**: 编译快，执行稳定

**适合场景**:
- ✅ 高频规则引擎
- ✅ 实时决策系统
- ✅ 内存敏感应用
- ✅ 低延迟场景

---

**需要帮助？** 查看 [BENCHMARK_GUIDE.md](BENCHMARK_GUIDE.md) 获取详细使用说明。

**想了解更多？** 阅读 [COMPLETION_REPORT.md](COMPLETION_REPORT.md) 查看完整的技术分析。
