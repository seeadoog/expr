package benchmark

import (
	"testing"

	"github.com/Knetic/govaluate"
	antonmedv "github.com/antonmedv/expr"
	"github.com/seeadoog/expr"
)

// Simple arithmetic expression
func BenchmarkSimpleArithmetic(b *testing.B) {
	expression := "a + b * c - d / e"
	params := map[string]interface{}{
		"a": 10.0,
		"b": 20.0,
		"c": 30.0,
		"d": 40.0,
		"e": 5.0,
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue(expression)
		if err != nil {
			b.Fatal(err)
		}
		ctx := env.NewContext(params)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = compiled.Val(ctx)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		program, err := antonmedv.Compile(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, params)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		compiled, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Evaluate(params)
		}
	})
}

// Complex boolean logic
func BenchmarkComplexBoolean(b *testing.B) {
	expression := "(a > 5 && b < 100) || (c == 'test' && d != nil)"
	params := map[string]interface{}{
		"a": 10.0,
		"b": 50.0,
		"c": "test",
		"d": "value",
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue(expression)
		if err != nil {
			b.Fatal(err)
		}
		ctx := env.NewContext(params)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = compiled.Val(ctx)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		program, err := antonmedv.Compile(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, params)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		compiled, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Evaluate(params)
		}
	})
}

// String operations
func BenchmarkStringOperations(b *testing.B) {
	expression := "firstName + ' ' + lastName"
	params := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue(expression)
		if err != nil {
			b.Fatal(err)
		}
		ctx := env.NewContext(params)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = compiled.Val(ctx)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		program, err := antonmedv.Compile(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, params)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		compiled, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Evaluate(params)
		}
	})
}

// Nested conditions
func BenchmarkNestedConditions(b *testing.B) {
	expression := "a > 10 && (b < 20 || (c >= 30 && d <= 40))"
	params := map[string]interface{}{
		"a": 15.0,
		"b": 25.0,
		"c": 35.0,
		"d": 35.0,
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue(expression)
		if err != nil {
			b.Fatal(err)
		}
		ctx := env.NewContext(params)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = compiled.Val(ctx)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		program, err := antonmedv.Compile(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, params)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		compiled, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Evaluate(params)
		}
	})
}

// Function calls
func BenchmarkFunctionCalls(b *testing.B) {
	// Simple math functions
	expression := "a * 2 + b / 2"
	params := map[string]interface{}{
		"a": 10.0,
		"b": 20.0,
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue(expression)
		if err != nil {
			b.Fatal(err)
		}
		ctx := env.NewContext(params)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = compiled.Val(ctx)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		program, err := antonmedv.Compile(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, params)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		compiled, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Evaluate(params)
		}
	})
}

// Compilation benchmark
func BenchmarkCompilation(b *testing.B) {
	expression := "(a + b) * (c - d) / e > 100 && (f == 'test' || g != nil)"

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = env.ParseValue(expression)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Compile(expression)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = govaluate.NewEvaluableExpression(expression)
		}
	})
}

// Large expression
func BenchmarkLargeExpression(b *testing.B) {
	expression := "v1 + v2 + v3 + v4 + v5 + v6 + v7 + v8 + v9 + v10"
	params := map[string]interface{}{
		"v1": 1.0, "v2": 2.0, "v3": 3.0, "v4": 4.0, "v5": 5.0,
		"v6": 6.0, "v7": 7.0, "v8": 8.0, "v9": 9.0, "v10": 10.0,
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue(expression)
		if err != nil {
			b.Fatal(err)
		}
		ctx := env.NewContext(params)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = compiled.Val(ctx)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		program, err := antonmedv.Compile(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, params)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		compiled, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Evaluate(params)
		}
	})
}

// Ternary operator
func BenchmarkTernaryOperator(b *testing.B) {
	expression := "a > 10 ? b * 2 : c / 2"
	params := map[string]interface{}{
		"a": 15.0,
		"b": 20.0,
		"c": 30.0,
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue(expression)
		if err != nil {
			b.Fatal(err)
		}
		ctx := env.NewContext(params)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = compiled.Val(ctx)
		}
	})

	b.Run("antonmedv/expr", func(b *testing.B) {
		program, err := antonmedv.Compile(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, params)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		compiled, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = compiled.Evaluate(params)
		}
	})
}
