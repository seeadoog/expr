package main

import (
	"fmt"
	"github.com/seeadoog/expr"
	"math"
	"sync"
)

type HandlerContext struct {
	Class int
	Args  map[string]interface{}
}

var (
	env = expr.NewEnv()
)

func init() {

	expr.SelfDefine0(env, "get_class", func(ctx *expr.Context, self *HandlerContext) float64 {
		return float64(self.Class)
	})

	expr.SelfDefine0(env, "args", func(ctx *expr.Context, self *HandlerContext) any {
		return self.Args
	})

	expr.SelfDefine0(env, "lock", func(ctx *expr.Context, self *sync.Mutex) any {
		self.Lock()
		return nil
	})

	expr.SelfDefine0(env, "unlock", func(ctx *expr.Context, self *sync.Mutex) any {
		self.Unlock()
		return nil
	})

	expr.RegisterOptFuncDefine1(env, "math_log", func(ctx *expr.Context, f float64, opt *expr.Options) float64 {

		return math.Log(f)
	})
}

func main() {

	ctx := env.GetContextFromPool()
	defer env.PutContext2Pool(ctx)
	hc := &HandlerContext{Class: 1, Args: make(map[string]interface{})}
	ctx.SetByString("ctx", hc)
	ctx.SetByString("mu", &sync.Mutex{})

	v, err := env.ParseValue(`

ctx.args().name = 5;
args  = ctx.args();
args.age = '135';
ctx.get_class();
mu.lock();
mu.unlock();
math_log(4) / math_log(2)`)
	if err != nil {
		panic(err)
	}

	res, e := ctx.SafeValue(v)
	if e != nil {
		panic(e)
	}
	fmt.Println(res)
	fmt.Println(hc.Args)

}
