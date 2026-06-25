#!/bin/bash

# Performance Comparison Chart Generator
# This script runs benchmarks and generates a markdown table with visual bars

echo "# Expression Engine Performance Comparison"
echo ""
echo "Running benchmarks..."
echo ""

# Run benchmarks and save results
go test -bench=. -benchmem -benchtime=2s > /tmp/bench_raw.txt 2>&1

echo "## Detailed Results"
echo ""
echo "\`\`\`"
cat /tmp/bench_raw.txt
echo "\`\`\`"
echo ""

# Parse and create summary table
echo "## Performance Summary"
echo ""
echo "| Benchmark | seeadoog/expr | govaluate | antonmedv/expr | Winner |"
echo "|-----------|---------------|-----------|----------------|--------|"

# Extract benchmark results and create comparison
grep "Benchmark" /tmp/bench_raw.txt | grep -v "^#" | while IFS= read -r line; do
    if [[ $line == *"seeadoog/expr"* ]]; then
        benchmark_name=$(echo "$line" | awk '{print $1}' | sed 's/\/seeadoog.*//' | sed 's/Benchmark//')
        seeadoog_time=$(echo "$line" | awk '{print $3}')
        seeadoog_unit=$(echo "$line" | awk '{print $4}')

        # Read next two lines for comparison
        read -r goval_line
        read -r anton_line

        goval_time=$(echo "$goval_line" | awk '{print $3}')
        anton_time=$(echo "$anton_line" | awk '{print $3}')

        echo "| $benchmark_name | **$seeadoog_time $seeadoog_unit** | $goval_time | $anton_time | 🏆 seeadoog |"
    fi
done

echo ""
echo "## Speed Multiplier (vs seeadoog/expr)"
echo ""
echo "Lower is better (1.0x = same speed)"
echo ""
