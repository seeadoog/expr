package expr

import (
	"reflect"
	"strconv"
)

type ExprValue struct {
	data any
}

func ValueOf(data any) ExprValue {
	return ExprValue{data: data}
}

func (e ExprValue) String() string {
	return StringOf(e.data)
}

func (e ExprValue) Number() float64 {
	return NumberOf(e.data)
}
func (e ExprValue) Bool() bool {
	return BoolOf(e.data)
}
func (e ExprValue) Set(ctx *Context, k string, val any) {
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

func (e ExprValue) Get(ctx *Context, key string) any {
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

func (e ExprValue) Contains(ctx *Context, v any) bool {
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

func (e ExprValue) RangeMap(f func(k string, v any) bool) {
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

func (e ExprValue) Len() int {
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

func (e ExprValue) IndexGet(i int) any {
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

func (e ExprValue) RangeArr(f func(k int, v any) bool) {
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

func (e ExprValue) AnyArr() []any {
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
		if v.Kind() == reflect.Slice {
			dst := make([]any, v.Len())

			for i := 0; i < v.Len(); i++ {
				dst[i] = v.Index(i).Interface()
			}
			return dst
		}
		return nil
	}
}

func (e ExprValue) F64Arr() []float64 {
	switch m := e.data.(type) {
	case []float64:
		return m
	case []any:
		d := make([]float64, len(m))
		for i, v := range m {
			d[i] = NumberOf(v)
		}
		return d
	case []int:
		d := make([]float64, len(m))
		for i, v := range m {
			d[i] = float64(v)
		}
		return d
	default:
		v := reflect.ValueOf(e.data)
		if v.Kind() == reflect.Slice {
			dst := make([]float64, v.Len())

			for i := 0; i < v.Len(); i++ {
				dst[i] = NumberOf(v.Index(i).Interface())
			}
			return dst
		}
		return nil
	}
}

func (e ExprValue) IntArr() []int {
	switch m := e.data.(type) {
	case []int:
		return m
	case []float64:
		d := make([]int, len(m))
		for i, v := range m {
			d[i] = int(v)
		}
		return d
	case []int64:
		d := make([]int, len(m))
		for i, v := range m {
			d[i] = int(v)
		}
		return d
	default:
		v := reflect.ValueOf(e.data)
		if v.Kind() == reflect.Slice {
			dst := make([]int, v.Len())
			for i := 0; i < v.Len(); i++ {
				dst[i] = IntOf(v.Index(i).Interface())
			}
			return dst
		}
		return nil
	}
}

func IntOf(v interface{}) int {
	switch m := v.(type) {
	case int:
		return m
	case *int:
		return *m
	case int64:
		return int(m)
	case int32:
		return int(m)
	case int16:
		return int(m)
	case int8:
		return int(m)
	case uint:
		return int(m)
	case uint64:
		return int(m)
	case uint32:
		return int(m)
	case uint16:
		return int(m)
	case uint8:
		return int(m)
	case float64:
		return int(m)
	case float32:
		return int(m)
	case bool:
		if m {
			return 1
		}
		return 0
	case string:
		i, err := strconv.Atoi(m)
		if err != nil {
			return 0
		}
		return i
	}
	return 0
}
