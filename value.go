package expr

import "reflect"

type exprValue struct {
	data any
}

func (e exprValue) String() string {
	return StringOf(e.data)
}
func (e exprValue) Set(ctx *Context, k string, val any) {
	switch m := e.data.(type) {
	case map[string]any:
		m[k] = val
	case map[string]string:
		m[k] = StringOf(val)
	case ReadOnlyMap:
	case Setter:
		m.SetField(ctx, k, val)
	default:
		setFieldOfStruct(ctx, reflect.ValueOf(e.data), k, val)
	}
}

func (e exprValue) Get(ctx *Context, key string) any {
	switch m := e.data.(type) {
	case map[string]any:
		return m[key]
	case map[string]string:
		return m[key]
	case ReadOnlyMap:
		return m[key]
	case Getter:
		return m.GetField(ctx, key)
	default:
		return getFieldOfStruct(ctx.ForceType, reflect.ValueOf(e.data), key)
	}
}

func (e exprValue) Contains(ctx *Context, v any) bool {
	switch m := e.data.(type) {
	case []any:
		for _, ev := range m {
			if ev == v {
				return true
			}
		}
	case ReadOnlyArray:
		for _, ev := range m {
			if ev == v {
				return true
			}
		}
	case ReadOnlyMap:
		sarg := StringOf(v)
		_, ok := m[sarg]
		if ok {
			return true
		}
	default:
		return e.data == v
	}
	return false
}

func (e exprValue) RangeMap(f func(k string, v any) bool) {
	//v := f.target.Val(c)
	switch m := e.data.(type) {
	case map[string]any:
		for k, v := range m {
			if !f(k, v) {
				return
			}
		}
	case ReadOnlyMap:
		for k, v := range m {
			if !f(k, v) {
				return
			}
		}
	case map[string]string:
		for k, v := range m {
			if !f(k, v) {
				return
			}
		}
	}
}

func (e exprValue) RangeArr(f func(k int, v any) bool) {
	switch m := e.data.(type) {
	case []any:
		for i, ev := range m {
			if !f(i, ev) {
				return
			}
		}
	case ReadOnlyArray:
		for i, ev := range m {
			if !f(i, ev) {
				return
			}
		}
	}
}
