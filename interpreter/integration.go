package interpreter

// ExprEnvAdapter 将 expr.Env 适配为 interpreter.Env
// 使用示例:
//
//	exprEnv := expr.DefaultEnv
//	interpEnv := interpreter.NewExprEnvAdapter(exprEnv)
//	interp := interpreter.NewInterpreter(interpEnv)
type ExprEnvAdapter struct {
	env interface {
		CalcHash(s string) uint64
		// 其他需要的方法可以在这里声明
	}
}

// NewExprEnvAdapter 创建适配器
// 参数应为 *expr.Env，但为了避免循环依赖，这里使用 interface{}
func NewExprEnvAdapter(env any) *ExprEnvAdapter {
	return &ExprEnvAdapter{env: env.(interface {
		CalcHash(s string) uint64
	})}
}

// CalcHash 实现 Env 接口
func (e *ExprEnvAdapter) CalcHash(s string) uint64 {
	return e.env.CalcHash(s)
}

// GetFunc 实现 Env 接口（查找全局函数）
// 注意：此方法需要访问 expr.Env 的内部函数表，当前为占位实现
func (e *ExprEnvAdapter) GetFunc(hash uint64, name string) (ScriptFunc, bool) {
	// TODO: 需要通过反射或导出方法访问 expr.Env.funtables
	// 当前返回 nil 以避免 panic
	return nil, false
}

// GetLibFunc 实现 Env 接口（查找库函数）
func (e *ExprEnvAdapter) GetLibFunc(objHash, funcHash uint64) ScriptFunc {
	// TODO: 需要访问 expr.Env.libs
	return nil
}

// ExprContextAdapter 将 expr.Context 适配为 interpreter.Context
// expr.Context 已经实现了所需的 Get/Set 方法，可以直接使用
type ExprContextAdapter struct {
	ctx interface {
		Get(hash uint64) any
		Set(hash uint64, val any)
		GetByString(key string) any
		SetByString(key string, val any)
	}
}

// NewExprContextAdapter 创建上下文适配器
func NewExprContextAdapter(ctx any) *ExprContextAdapter {
	return &ExprContextAdapter{ctx: ctx.(interface {
		Get(hash uint64) any
		Set(hash uint64, val any)
		GetByString(key string) any
		SetByString(key string, val any)
	})}
}

// Get 实现 Context 接口
func (c *ExprContextAdapter) Get(hash uint64) any {
	return c.ctx.Get(hash)
}

// Set 实现 Context 接口
func (c *ExprContextAdapter) Set(hash uint64, val any) {
	c.ctx.Set(hash, val)
}

// GetByString 实现 Context 接口
func (c *ExprContextAdapter) GetByString(key string) any {
	return c.ctx.GetByString(key)
}

// SetByString 实现 Context 接口
func (c *ExprContextAdapter) SetByString(key string, val any) {
	c.ctx.SetByString(key, val)
}

// 注意事项
//
// 1. 函数调用集成
//
// 当前版本的 GetFunc 和 GetLibFunc 是占位实现。要完全集成函数调用，
// 需要 expr 包暴露或通过反射访问其内部函数表。有以下几种方案：
//
// 方案 A: 在 expr.Env 添加导出方法
//   func (e *Env) GetFuncByHash(hash uint64, name string) (ScriptFunc, bool)
//
// 方案 B: 使用反射访问私有字段（不推荐，脆弱）
//   reflect.ValueOf(env).Elem().FieldByName("funtables")
//
// 方案 C: 将 interpreter 直接放入 expr 包内（需要调整目录结构）
//
// 2. 对象方法调用
//
// objFuncMap 是包级全局变量，VM 中的 callMethod 需要访问它来调用
// 如 str.split()、arr.map() 等方法。需要：
//   - expr 包导出 GetObjFunc(typeHash, funcHash) 方法
//   - 或将 objFuncMap 改为 Env 的字段
//
// 3. 性能对比
//
// 建议添加基准测试，对比字节码 VM 与原 Val 闭包树的性能：
//   - 编译时间
//   - 首次执行时间
//   - 重复执行时间（字节码优势所在）
//   - 内存占用
//
// 4. 错误处理
//
// expr 的 Val.Val(ctx) 返回 any，错误通过 *Error 类型表示。
// VM 需要识别并正确传播这类错误。
