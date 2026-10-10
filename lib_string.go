package expr

import "unicode/utf8"

func init() {

	RegisterOptFuncDefine1(DefaultEnv, "str_is_utf8", func(ctx *Context, a string, opt *Options) bool {
		return utf8.ValidString(a)
	})

}
