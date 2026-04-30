package expr

import (
	"fmt"
	"sync"
)

type Env struct {
	funtables map[string]*innerFunc

	allTypeFuncs map[string]bool

	isDefault bool

	pool sync.Pool

	libs *envMap
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
		libs:         newEnvMap(1024 * 4),
	}
)

func NewEnv() *Env {
	env := &Env{
		funtables:    make(map[string]*innerFunc),
		allTypeFuncs: make(map[string]bool),
		isDefault:    false,
		libs:         newEnvMap(1024 * 4),
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

func (e *Env) GetLibFunc(libHash uint64, funNameHash uint64) ScriptFunc {
	l, _ := e.libs.getHash(libHash).(*Lib)
	if l != nil {
		fun, _ := l.funs.getHash(funNameHash).(ScriptFunc)
		return fun
	}
	return nil
}

func (e *Env) AddLib(l *Lib) {
	e.libs.putHash(calcHash(l.libName), l.libName, l)
}

type Lib struct {
	funs    *envMap
	libName string
}

func (l *Lib) RegisterFuncWithOpt(funName string, f OptFunc, argsNum int, argDesc string, opts ...funcOpt) {
	l.RegisterFunc(funName, func(ctx *Context, args ...Val) any {

		if len(args) < argsNum {
			panic(fmt.Sprintf("call func %s.%s() failed, want args num %d but %d", l.libName, funName, argsNum, len(args)))
		}
		var opt *Options
		if len(args) == argsNum+1 {
			opt = args[argsNum].Val(ctx).(*Options)
		}
		return f(ctx, args, opt)
	}, argsNum)
}

func NewLib(name string) *Lib {
	return &Lib{
		funs:    newEnvMap(1024 * 4),
		libName: name,
	}
}

func (l *Lib) RegisterFunc(funName string, f ScriptFunc, argsNum int, opts ...funcOpt) {
	fhash := calcHash(funName)
	l.funs.putHash(fhash, funName, f)
}

type Register interface {
	RegisterFunc(funName string, f ScriptFunc, argsNum int, opts ...funcOpt)
	RegisterFuncWithOpt(funName string, f OptFunc, argsNum int, argDesc string, opts ...funcOpt)
}
