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

func init() {
	DefaultEnv.RegisterFunc("math_abs", FuncDefine1(math.Abs), 1)
	DefaultEnv.RegisterFunc("math_ceil", FuncDefine1(math.Ceil), 1)
	DefaultEnv.RegisterFunc("math_floor", FuncDefine1(math.Floor), 1)

	//DefaultEnv.RegisterFunc("math_abs", FuncDefine1(math.Abs), 1)
	DefaultEnv.RegisterFunc("math_acos", FuncDefine1(math.Acos), 1)
	DefaultEnv.RegisterFunc("math_acosh", FuncDefine1(math.Acosh), 1)
	DefaultEnv.RegisterFunc("math_asin", FuncDefine1(math.Asin), 1)
	DefaultEnv.RegisterFunc("math_asinh", FuncDefine1(math.Asinh), 1)
	DefaultEnv.RegisterFunc("math_atan", FuncDefine1(math.Atan), 1)
	DefaultEnv.RegisterFunc("math_atanh", FuncDefine1(math.Atanh), 1)
	DefaultEnv.RegisterFunc("math_cbrt", FuncDefine1(math.Cbrt), 1)
	//DefaultEnv.RegisterFunc("math_ceil", FuncDefine1(math.Ceil), 1)
	DefaultEnv.RegisterFunc("math_cos", FuncDefine1(math.Cos), 1)
	DefaultEnv.RegisterFunc("math_cosh", FuncDefine1(math.Cosh), 1)
	DefaultEnv.RegisterFunc("math_erf", FuncDefine1(math.Erf), 1)
	DefaultEnv.RegisterFunc("math_erfc", FuncDefine1(math.Erfc), 1)
	DefaultEnv.RegisterFunc("math_exp", FuncDefine1(math.Exp), 1)
	DefaultEnv.RegisterFunc("math_exp2", FuncDefine1(math.Exp2), 1)
	DefaultEnv.RegisterFunc("math_expm1", FuncDefine1(math.Expm1), 1)
	//DefaultEnv.RegisterFunc("math_floor", FuncDefine1(math.Floor), 1)
	DefaultEnv.RegisterFunc("math_gamma", FuncDefine1(math.Gamma), 1)
	DefaultEnv.RegisterFunc("math_j0", FuncDefine1(math.J0), 1)
	DefaultEnv.RegisterFunc("math_j1", FuncDefine1(math.J1), 1)
	DefaultEnv.RegisterFunc("math_log", FuncDefine1(math.Log), 1)
	DefaultEnv.RegisterFunc("math_log10", FuncDefine1(math.Log10), 1)
	DefaultEnv.RegisterFunc("math_log1p", FuncDefine1(math.Log1p), 1)
	DefaultEnv.RegisterFunc("math_log2", FuncDefine1(math.Log2), 1)
	DefaultEnv.RegisterFunc("math_round", FuncDefine1(math.Round), 1)
	DefaultEnv.RegisterFunc("math_roundtoeven", FuncDefine1(math.RoundToEven), 1)
	DefaultEnv.RegisterFunc("math_signbit", FuncDefine1(math.Signbit), 1)
	DefaultEnv.RegisterFunc("math_sin", FuncDefine1(math.Sin), 1)
	DefaultEnv.RegisterFunc("math_sinh", FuncDefine1(math.Sinh), 1)
	DefaultEnv.RegisterFunc("math_sqrt", FuncDefine1(math.Sqrt), 1)
	DefaultEnv.RegisterFunc("math_tan", FuncDefine1(math.Tan), 1)
	DefaultEnv.RegisterFunc("math_tanh", FuncDefine1(math.Tanh), 1)
	DefaultEnv.RegisterFunc("math_trunc", FuncDefine1(math.Trunc), 1)
	DefaultEnv.RegisterFunc("math_y0", FuncDefine1(math.Y0), 1)
	DefaultEnv.RegisterFunc("math_y1", FuncDefine1(math.Y1), 1)

	DefaultEnv.RegisterFunc("math_atan2", FuncDefine2(math.Atan2), 2)
	DefaultEnv.RegisterFunc("math_copysign", FuncDefine2(math.Copysign), 2)
	DefaultEnv.RegisterFunc("math_dim", FuncDefine2(math.Dim), 2)
	DefaultEnv.RegisterFunc("math_hypot", FuncDefine2(math.Hypot), 2)
	DefaultEnv.RegisterFunc("math_max", FuncDefine2(math.Max), 2)
	DefaultEnv.RegisterFunc("math_min", FuncDefine2(math.Min), 2)
	DefaultEnv.RegisterFunc("math_mod", FuncDefine2(math.Mod), 2)
	DefaultEnv.RegisterFunc("math_nextafter", FuncDefine2(math.Nextafter), 2)
	DefaultEnv.RegisterFunc("math_pow", FuncDefine2(math.Pow), 2)
	DefaultEnv.RegisterFunc("math_remainder", FuncDefine2(math.Remainder), 2)

	DefaultEnv.RegisterFunc("math_jn", FuncDefine2(math.Jn), 2)
	DefaultEnv.RegisterFunc("math_yn", FuncDefine2(math.Yn), 2)
}
