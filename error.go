package expr

import "fmt"

type RuntimeError struct {
	Err string
}

func (r *RuntimeError) Error() string {
	//TODO implement me
	return r.Err
}

// newErrorf 创建错误（不带 Context，使用默认行为：panic）
func newErrorf(format string, args ...interface{}) *Error {
	err := &Error{Err: &RuntimeError{Err: fmt.Sprintf(format, args...)}}
	return err
}

// newErrorfWithCtx 创建错误（带 Context，根据 Context.PanicWhenError 决定是否 panic）
func newErrorfWithCtx(ctx *Context, format string, args ...interface{}) *Error {
	err := &Error{Err: &RuntimeError{Err: fmt.Sprintf(format, args...)}}
	if ctx != nil && ctx.PanicWhenError {
		panic(err)
	}
	return err
}
