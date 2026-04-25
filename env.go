package expr

import "sync"

type Env struct {
	funtables map[string]*innerFunc

	allTypeFuncs map[string]bool

	isDefault bool

	pool sync.Pool
}

func (e *Env) GetContextFromPool() *Context {
	ctx, ok := e.pool.Get().(*Context)
	if ok {
		return ctx
	}
	return e.NewContext(nil)
}

func (e *Env) PutContext2Pool(ctx *Context) {
	ctx.Reset()
	e.pool.Put(ctx)
}

var (
	DefaultEnv = &Env{
		funtables:    funtables_,
		allTypeFuncs: allTypeFuncs_,
		isDefault:    true,
	}
)

func NewEnv() *Env {
	env := &Env{
		funtables:    make(map[string]*innerFunc),
		allTypeFuncs: make(map[string]bool),
		isDefault:    false,
	}
	for key, val := range DefaultEnv.funtables {
		env.funtables[key] = val
	}
	for key, val := range DefaultEnv.allTypeFuncs {
		env.allTypeFuncs[key] = val
	}
	return env
}

func (e *Env) RemoveFunc(f string) {
	delete(e.funtables, f)
}
