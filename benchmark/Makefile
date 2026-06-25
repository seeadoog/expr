.PHONY: bench bench-short bench-all bench-compare clean

# Run all benchmarks with 2s duration
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem -benchtime=2s

# Quick benchmark run (100ms)
bench-short:
	@echo "Running quick benchmarks..."
	@go test -bench=. -benchmem -benchtime=100ms

# Extended benchmark run (5s)
bench-all:
	@echo "Running extended benchmarks..."
	@go test -bench=. -benchmem -benchtime=5s

# Run benchmark and save to file
bench-save:
	@echo "Running benchmarks and saving to benchmark_results.txt..."
	@go test -bench=. -benchmem -benchtime=2s | tee benchmark_results.txt
	@echo ""
	@echo "Results saved to benchmark_results.txt"

# Compare with previous results
bench-compare:
	@if [ -f benchmark_results.txt ]; then \
		echo "Running new benchmark and comparing..."; \
		go test -bench=. -benchmem -benchtime=2s > benchmark_new.txt; \
		echo "=== Previous Results ==="; \
		cat benchmark_results.txt; \
		echo ""; \
		echo "=== New Results ==="; \
		cat benchmark_new.txt; \
	else \
		echo "No previous results found. Run 'make bench-save' first."; \
	fi

# Run specific benchmark
bench-simple:
	@go test -bench=BenchmarkSimpleArithmetic -benchmem

bench-bool:
	@go test -bench=BenchmarkComplexBoolean -benchmem

bench-string:
	@go test -bench=BenchmarkStringOperations -benchmem

bench-compile:
	@go test -bench=BenchmarkCompilation -benchmem

# CPU profiling
bench-cpu:
	@echo "Running benchmarks with CPU profiling..."
	@go test -bench=. -benchmem -cpuprofile=cpu.prof
	@echo "Profile saved to cpu.prof"
	@echo "View with: go tool pprof cpu.prof"

# Memory profiling
bench-mem:
	@echo "Running benchmarks with memory profiling..."
	@go test -bench=. -benchmem -memprofile=mem.prof
	@echo "Profile saved to mem.prof"
	@echo "View with: go tool pprof mem.prof"

# Clean generated files
clean:
	@rm -f benchmark_results.txt benchmark_new.txt
	@rm -f cpu.prof mem.prof
	@rm -f *.test
	@echo "Cleaned benchmark artifacts"

# Show help
help:
	@echo "Available targets:"
	@echo "  bench          - Run benchmarks (2s each)"
	@echo "  bench-short    - Quick benchmarks (100ms each)"
	@echo "  bench-all      - Extended benchmarks (5s each)"
	@echo "  bench-save     - Run and save results to file"
	@echo "  bench-compare  - Compare with previous results"
	@echo "  bench-simple   - Run only simple arithmetic benchmark"
	@echo "  bench-bool     - Run only boolean logic benchmark"
	@echo "  bench-string   - Run only string operations benchmark"
	@echo "  bench-compile  - Run only compilation benchmark"
	@echo "  bench-cpu      - Run with CPU profiling"
	@echo "  bench-mem      - Run with memory profiling"
	@echo "  clean          - Remove generated files"
	@echo "  help           - Show this help message"
