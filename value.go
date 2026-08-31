package expr

import "reflect"

type exprValue struct {
	data any
}

func ValueOf(data any) exprValue {
	return exprValue{data: data}
}

func (e exprValue) String() string {
	return StringOf(e.data)
}

func (e exprValue) Number() float64 {
	return NumberOf(e.data)
}
func (e exprValue) Bool() bool {
	return BoolOf(e.data)
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

func (e exprValue) Len() int {
	switch m := e.data.(type) {
	case []any:
		return len(m)
	case ReadOnlyArray:
		return len(m)
	case ReadOnlyMap:
		return len(m)
	case []float64:
		return len(m)
	default:
		v := reflect.ValueOf(e.data)
		if v.Kind() == reflect.Slice {
			return v.Len()
		}
		return 0
	}
}

func (e exprValue) IndexGet(i int) any {
	switch m := e.data.(type) {
	case []any:
		return m[i]
	case ReadOnlyArray:
		return m[i]
	case []float64:
		return m[i]
	case []int:
		return m[i]
	default:
		v := reflect.ValueOf(e.data)
		if v.Kind() == reflect.Slice {
			return v.Index(i).Interface()
		}
		return 0
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
	case []int:
		for i, ev := range m {
			if !f(i, ev) {
				return
			}
		}
	case []string:
		for i, ev := range m {
			if !f(i, ev) {
			}
		}
	case []float64:
		for i, ev := range m {
			if !f(i, ev) {
				return
			}
		}
	}
}

func (e exprValue) AnyArr() []any {
	switch m := e.data.(type) {
	case []any:
		return m
	case []float64:
		d := make([]any, len(m))
		for i, v := range m {
			d[i] = v
		}
		return d
	case []int:
		d := make([]any, len(m))
		for i, v := range m {
			d[i] = v
		}
		return d
	case []string:
		d := make([]any, len(m))
		for i, v := range m {
			d[i] = v
		}
		return d
	default:
		v := reflect.ValueOf(e.data)
		dst := make([]any, v.Len())
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				dst[i] = v.Index(i).Interface()
			}
			return dst
		}
		return nil
	}
}
