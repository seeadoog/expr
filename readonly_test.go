package expr

import "testing"

func TestNewReadOnlyVal(t *testing.T) {

	c := DefaultEnv.NewContext(map[string]any{
		"readmap": ReadOnlyMap{
			"name": "hello",
		},
		"readarr": ReadOnlyArray{
			"a1",
		},
	})
	DefaultEnv.GetContextFromPool()
	parseAndExec(DefaultEnv, `a1 = readarr[0]; m1 = readmap.name  `, c)

	assertEqual(t, c, "a1", "a1")
	assertEqual(t, c, "m1", "hello")

	parseAndExec(DefaultEnv, `readmap.for({k,v}=>_)`, c)
	assertEqual(t, c, "k", "name")
	assertEqual(t, c, "v", "hello")

	parseAndExec(DefaultEnv, `readarr.for({k,v}=>_)`, c)
	assertEqual(t, c, "k==0", true)
	assertEqual(t, c, "v", "a1")

	parseAndExec(DefaultEnv, `readarr[0] = 1;readmap.name  = 'chg'`, c)
	assertEqual(t, c, "readarr[0]", "a1")
	assertEqual(t, c, "readmap.name", "hello")

	parseAndExec(DefaultEnv, `constmap = const {name:'ns'};constmap.name = 'changed'`, c)
	parseAndExec(DefaultEnv, `constarr = const [1,2,3];constarr[0]=100'`, c)

	assertEqual(t, c, "constmap.name", "ns")
	assertEqual(t, c, "constarr[0]==1", true)

	//parseAndExec(DefaultEnv, ``, c)

	type session struct {
	}

	SelfDefine1(DefaultEnv, "get", func(ctx *Context, self *session, opt Option) any {
		o := NewOptions(opt)
		return o.Get("name")
	})

	c.SetByString("sess", &session{})
	parseAndExec(DefaultEnv, `sname = sess.get({name:"oss"})`, c)
	assertEqual(t, c, "sname", "oss")

}

func parseAndExec(env *Env, valExpr string, ctx *Context) {
	v, err := env.ParseValue(valExpr)
	if err != nil {
		panic(err)
	}
	v.Val(ctx)
}
