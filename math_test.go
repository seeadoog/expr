package expr

import (
	"testing"
)

// TestAddEq tests += operator
func TestAddEq(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		initial  map[string]any
		expected any
	}{
		{
			name:     "number addition",
			expr:     "a += 5",
			initial:  map[string]any{"a": 10.0},
			expected: 15.0,
		},
		{
			name:     "string concatenation",
			expr:     "s += ' world'",
			initial:  map[string]any{"s": "hello"},
			expected: "hello world",
		},
		{
			name:     "zero addition",
			expr:     "x += 0",
			initial:  map[string]any{"x": 100.0},
			expected: 100.0,
		},
		{
			name:     "negative addition",
			expr:     "n += -10",
			initial:  map[string]any{"n": 5.0},
			expected: -5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnv()
			compiled, err := env.ParseValue(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			ctx := env.NewContext(tt.initial)
			result := compiled.Val(ctx)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestSubEq tests -= operator
func TestSubEq(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		initial  map[string]any
		expected any
	}{
		{
			name:     "basic subtraction",
			expr:     "a -= 5",
			initial:  map[string]any{"a": 10.0},
			expected: 5.0,
		},
		{
			name:     "subtract from zero",
			expr:     "x -= 10",
			initial:  map[string]any{"x": 0.0},
			expected: -10.0,
		},
		{
			name:     "subtract negative",
			expr:     "n -= -5",
			initial:  map[string]any{"n": 10.0},
			expected: 15.0,
		},
		{
			name:     "subtract decimal",
			expr:     "d -= 0.5",
			initial:  map[string]any{"d": 10.5},
			expected: 10.0,
		},
		{
			name:     "large subtraction",
			expr:     "b -= 100",
			initial:  map[string]any{"b": 50.0},
			expected: -50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnv()
			compiled, err := env.ParseValue(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			ctx := env.NewContext(tt.initial)
			result := compiled.Val(ctx)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestMulEq tests *= operator
func TestMulEq(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		initial  map[string]any
		expected any
	}{
		{
			name:     "basic multiplication",
			expr:     "a *= 5",
			initial:  map[string]any{"a": 10.0},
			expected: 50.0,
		},
		{
			name:     "multiply by zero",
			expr:     "x *= 0",
			initial:  map[string]any{"x": 100.0},
			expected: 0.0,
		},
		{
			name:     "multiply by one",
			expr:     "y *= 1",
			initial:  map[string]any{"y": 42.0},
			expected: 42.0,
		},
		{
			name:     "multiply by negative",
			expr:     "n *= -2",
			initial:  map[string]any{"n": 5.0},
			expected: -10.0,
		},
		{
			name:     "multiply decimal",
			expr:     "d *= 0.5",
			initial:  map[string]any{"d": 10.0},
			expected: 5.0,
		},
		{
			name:     "multiply large numbers",
			expr:     "b *= 1000",
			initial:  map[string]any{"b": 3.0},
			expected: 3000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnv()
			compiled, err := env.ParseValue(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			ctx := env.NewContext(tt.initial)
			result := compiled.Val(ctx)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestDivEq tests /= operator
func TestDivEq(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		initial  map[string]any
		expected any
	}{
		{
			name:     "basic division",
			expr:     "a /= 5",
			initial:  map[string]any{"a": 10.0},
			expected: 2.0,
		},
		{
			name:     "divide by one",
			expr:     "x /= 1",
			initial:  map[string]any{"x": 42.0},
			expected: 42.0,
		},
		{
			name:     "divide by negative",
			expr:     "n /= -2",
			initial:  map[string]any{"n": 10.0},
			expected: -5.0,
		},
		{
			name:     "divide decimal",
			expr:     "d /= 0.5",
			initial:  map[string]any{"d": 10.0},
			expected: 20.0,
		},
		{
			name:     "divide into fraction",
			expr:     "b /= 3",
			initial:  map[string]any{"b": 10.0},
			expected: 10.0 / 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnv()
			compiled, err := env.ParseValue(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			ctx := env.NewContext(tt.initial)
			result := compiled.Val(ctx)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestMathOpsCombined tests combined math operations
func TestMathOpsCombined(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		initial  map[string]any
		expected any
	}{
		{
			name:     "multiple operations in sequence",
			expr:     "a += 5; a *= 2; a -= 4",
			initial:  map[string]any{"a": 10.0},
			expected: 26.0,
		},
		{
			name:     "division then addition",
			expr:     "x /= 2; x += 10",
			initial:  map[string]any{"x": 20.0},
			expected: 20.0,
		},
		{
			name:     "subtract then multiply",
			expr:     "y -= 5; y *= 3",
			initial:  map[string]any{"y": 15.0},
			expected: 30.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnv()
			compiled, err := env.ParseValue(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			ctx := env.NewContext(tt.initial)
			result := compiled.Val(ctx)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestMathOpsWithVarAccess tests math operations with variable access
func TestMathOpsWithVarAccess(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		initial  map[string]any
		expected any
	}{
		{
			name:     "access field and add",
			expr:     "obj->a += 10",
			initial:  map[string]any{"obj": map[string]any{"a": 5.0}},
			expected: 15.0,
		},
		{
			name:     "access field and subtract",
			expr:     "data->x -= 3",
			initial:  map[string]any{"data": map[string]any{"x": 10.0}},
			expected: 7.0,
		},
		{
			name:     "access field and multiply",
			expr:     "num->val *= 4",
			initial:  map[string]any{"num": map[string]any{"val": 5.0}},
			expected: 20.0,
		},
		{
			name:     "access field and divide",
			expr:     "calc->result /= 2",
			initial:  map[string]any{"calc": map[string]any{"result": 100.0}},
			expected: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnv()
			compiled, err := env.ParseValue(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			ctx := env.NewContext(tt.initial)
			result := compiled.Val(ctx)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestMathOpsWithArrayAccess tests math operations with array access
func TestMathOpsWithArrayAccess(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		initial  map[string]any
		expected any
	}{
		{
			name:     "array element add",
			expr:     "arr[0] += 5",
			initial:  map[string]any{"arr": []any{10.0, 20.0, 30.0}},
			expected: 15.0,
		},
		{
			name:     "array element subtract",
			expr:     "nums[1] -= 3",
			initial:  map[string]any{"nums": []any{5.0, 10.0, 15.0}},
			expected: 7.0,
		},
		{
			name:     "array element multiply",
			expr:     "vals[2] *= 2",
			initial:  map[string]any{"vals": []any{1.0, 2.0, 3.0}},
			expected: 6.0,
		},
		{
			name:     "array element divide",
			expr:     "data[0] /= 4",
			initial:  map[string]any{"data": []any{20.0, 40.0}},
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnv()
			compiled, err := env.ParseValue(tt.expr)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			ctx := env.NewContext(tt.initial)
			result := compiled.Val(ctx)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	c := DefaultEnv.NewContext(nil)
	parseAndExec(DefaultEnv, `
a = [1,2,3] .equals([1,2,3]);
af = [1,2,3,4] .equals([1,2,3]);
b = {name:2}.equals({name:2});
bf = {name:2}.equals({name:3});
c =  d.equals(g);
e = (1).equals(1);

ef = 345;
gf = ef ;
`, c)

	assertEqual(t, c, "a", true)
	assertEqual(t, c, "af", false)
	assertEqual(t, c, "b", true)
	assertEqual(t, c, "bf", false)
	assertEqual(t, c, "c", true)
	assertEqual(t, c, "e", true)
	assertEqual(t, c, "ef", 345.0)
	assertEqual(t, c, "gf", 345.0)

	c.SetByJp("name.age", 6)
	c.SetByJp("name.ase[0]", 6)

	assertEqual(t, c, "name.age", 6)
	assertEqual(t, c, "name.ase[0]", 6)
	// name:"router.zh_cn_en","566"
}
