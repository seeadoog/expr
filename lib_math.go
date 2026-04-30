package expr

import "math"

func NewMathLib(name string) *Lib {
	if name == "" {
		name = "math"
	}
	lib := NewLib(name)
	initMathLib(lib)
	return lib
}

func initMathLib(lib *Lib) {

	RegisterOptFuncDefine1(lib, "log", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Log(a)
	})
	// ===== 1个参数 =====
	RegisterOptFuncDefine1(lib, "abs", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Abs(a)
	})

	RegisterOptFuncDefine1(lib, "sqrt", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Sqrt(a)
	})

	RegisterOptFuncDefine1(lib, "log10", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Log10(a)
	})

	RegisterOptFuncDefine1(lib, "exp", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Exp(a)
	})

	RegisterOptFuncDefine1(lib, "sin", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Sin(a)
	})

	RegisterOptFuncDefine1(lib, "cos", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Cos(a)
	})

	RegisterOptFuncDefine1(lib, "tan", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Tan(a)
	})

	RegisterOptFuncDefine1(lib, "asin", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Asin(a)
	})

	RegisterOptFuncDefine1(lib, "acos", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Acos(a)
	})

	RegisterOptFuncDefine1(lib, "atan", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Atan(a)
	})

	RegisterOptFuncDefine1(lib, "ceil", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Ceil(a)
	})

	RegisterOptFuncDefine1(lib, "floor", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Floor(a)
	})

	RegisterOptFuncDefine1(lib, "round", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Round(a)
	})

	RegisterOptFuncDefine1(lib, "trunc", func(ctx *Context, a float64, opt *Options) float64 {
		return math.Trunc(a)
	})

	RegisterOptFuncDefine1(lib, "sign", func(ctx *Context, a float64, opt *Options) float64 {
		if a > 0 {
			return 1
		}
		if a < 0 {
			return -1
		}
		return 0
	})

	// ===== 2个参数 =====
	RegisterOptFuncDefine2(lib, "pow", func(ctx *Context, a, b float64, opt *Options) float64 {
		return math.Pow(a, b)
	})

	RegisterOptFuncDefine2(lib, "max", func(ctx *Context, a, b float64, opt *Options) float64 {
		return math.Max(a, b)
	})

	RegisterOptFuncDefine2(lib, "min", func(ctx *Context, a, b float64, opt *Options) float64 {
		return math.Min(a, b)
	})

	RegisterOptFuncDefine2(lib, "mod", func(ctx *Context, a, b float64, opt *Options) float64 {
		return math.Mod(a, b)
	})

	RegisterOptFuncDefine2(lib, "atan2", func(ctx *Context, y, x float64, opt *Options) float64 {
		return math.Atan2(y, x)
	})

	RegisterOptFuncDefine2(lib, "hypot", func(ctx *Context, a, b float64, opt *Options) float64 {
		return math.Hypot(a, b)
	})

	// ===== 3个参数 =====
	RegisterOptFuncDefine3(lib, "clamp", func(ctx *Context, x, min, max float64, opt *Options) float64 {
		if x < min {
			return min
		}
		if x > max {
			return max
		}
		return x
	})
}
