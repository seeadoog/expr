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

// Function calls - invokes a real user-defined function in each library.
// Expression: addone(a) + addone(b) * addone(c)  => 11 + 21*31 = 662
func BenchmarkFunctionCalls(b *testing.B) {
	expression := "addone(a) + addone(b) * addone(c)"
	params := map[string]interface{}{
		"a": 10.0,
		"b": 20.0,
		"c": 30.0,
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		env.RegisterFunc("addone", expr.FuncDefine1(func(ctx *expr.Context, a float64) float64 {
			return a + 1
		}), 1)
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
		env := map[string]interface{}{
			"a":      10.0,
			"b":      20.0,
			"c":      30.0,
			"addone": func(x float64) float64 { return x + 1 },
		}
		program, err := antonmedv.Compile(expression, antonmedv.Env(env))
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, env)
		}
	})

	b.Run("govaluate", func(b *testing.B) {
		functions := map[string]govaluate.ExpressionFunction{
			"addone": func(args ...interface{}) (interface{}, error) {
				return args[0].(float64) + 1, nil
			},
		}
		compiled, err := govaluate.NewEvaluableExpressionWithFunctions(expression, functions)
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

// Lambda / higher-order function - filters a slice with an inline predicate.
// govaluate has no lambda/closure support, so it is omitted here.
func BenchmarkLambda(b *testing.B) {
	params := map[string]interface{}{
		"nums": []interface{}{1.0, 5.0, 8.0, 3.0, 10.0, 2.0, 7.0, 4.0},
	}

	b.Run("seeadoog/expr", func(b *testing.B) {
		env := expr.NewEnv()
		compiled, err := env.ParseValue("nums.filter(x => x > 5)")
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
		env := map[string]interface{}{
			"nums": []interface{}{1.0, 5.0, 8.0, 3.0, 10.0, 2.0, 7.0, 4.0},
		}
		program, err := antonmedv.Compile("filter(nums, # > 5)", antonmedv.Env(env))
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = antonmedv.Run(program, env)
		}
	})

	// govaluate: no lambda support
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
