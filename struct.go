package expr

import (
	"fmt"
	"reflect"
	"strings"
)

// anyToInt64 尽量无精度损失地把 v 转成 int64（整型直接取值，避免走 float64）。
func anyToInt64(v any) int64 {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(rv.Float())
	case reflect.Bool:
		if rv.Bool() {
			return 1
		}
		return 0
	default:
		return int64(NumberOf(v))
	}
}

// anyToUint64 尽量无精度损失地把 v 转成 uint64。
func anyToUint64(v any) uint64 {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint()
	case reflect.Float32, reflect.Float64:
		return uint64(rv.Float())
	case reflect.Bool:
		if rv.Bool() {
			return 1
		}
		return 0
	default:
		return uint64(NumberOf(v))
	}
}

// asStringMap 把 v 规整成 map[string]any（支持任意 string 为 key 的 map）。
func asStringMap(v any) (map[string]any, bool) {
	if m, ok := v.(map[string]any); ok {
		return m, true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
		m := make(map[string]any, rv.Len())
		for _, k := range rv.MapKeys() {
			m[k.String()] = rv.MapIndex(k).Interface()
		}
		return m, true
	}
	return nil, false
}

// asAnySlice 把 v 规整成 []any（支持任意 slice/array）。
func asAnySlice(v any) ([]any, bool) {
	if s, ok := v.([]any); ok {
		return s, true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		s := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			s[i] = rv.Index(i).Interface()
		}
		return s, true
	}
	return nil, false
}

func getFieldOfStruct(forceType bool, rv reflect.Value, name string) any {
	switch rv.Kind() {
	case reflect.Struct:
		fv := rv.FieldByName(name)
		if !fv.IsValid() {
			mv := rv.MethodByName(name)
			if mv.IsValid() {
				return mv.Interface()
			}
			return nil
		}

		if fv.Kind() == reflect.Struct {
			return fv.Addr().Interface()
		}
		return structValueToVm(forceType, fv.Interface())
	case reflect.Ptr:
		if !rv.IsNil() {

			v := getFieldOfStruct(forceType, rv.Elem(), name)
			if v != nil {
				return v
			}
			mv := rv.MethodByName(name)
			if mv.IsValid() {
				return mv.Interface()
			}
		}
		return nil
	case reflect.Map:

		kv, ok := structValConvert(nil, rv.Type().Key(), name)
		if !ok {
			return nil
		}

		data := rv.MapIndex(kv)
		if !data.IsValid() {
			return nil
		}
		return structValueToVm(forceType, data.Interface())
	default:
		return nil
	}
}

func setValStruct(ctx *Context, fv reflect.Value, val any) {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(StringOf(val))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fv.SetInt(int64(NumberOf(val)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fv.SetUint(uint64(NumberOf(val)))
	case reflect.Float32, reflect.Float64:
		fv.SetFloat(NumberOf(val))
	case reflect.Bool:
		fv.SetBool(BoolOf(val))
	default:
		v, ok := structValConvert(ctx, fv.Type(), val)
		if ok {
			fv.Set(v)
		} else {
			panic(fmt.Sprintf("cannot set val '%v':%s to type %s", val, reflect.TypeOf(val).String(), fv.Type().String()))
		}
	}
}

// 将 any 类型的 v 根据 t 转成 t 对应的 reflect.Value。
// 支持：精确类型、可赋值/可转换类型、指针自动取址与解引用、接口、
// 基础类型互转、struct(map 填充)、slice/array、map、func(lambda)。
func structValConvert(ctx *Context, t reflect.Type, v any) (vv reflect.Value, ok bool) {
	// nil 源：返回目标类型的零值（指针/接口/map/slice 为 nil）。
	if v == nil {
		return reflect.Zero(t), true
	}

	rv := reflect.ValueOf(v)
	vt := rv.Type()

	// 1. 类型完全一致。
	if vt == t {
		return rv, true
	}

	// 2. 直接可赋值（命名类型、空接口等）。
	if vt.AssignableTo(t) {
		return rv, true
	}

	// 3. 目标是接口：v 实现了该接口即可。
	if t.Kind() == reflect.Interface {
		if vt.Implements(t) {
			return rv, true
		}
		// 非空接口未实现则失败，避免后续 Set 时 panic。
		if t.NumMethod() > 0 {
			return vv, false
		}
		// 空接口（any）兜底。
		return rv, true
	}

	// 4. 源是指针：解引用后再尝试匹配目标。
	if vt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Zero(t), true
		}
		if vt.Elem() == t || vt.Elem().AssignableTo(t) {
			return rv.Elem(), true
		}
	}

	var tv reflect.Value
	isPtr := false
	if t.Kind() == reflect.Ptr {
		// 源恰好是目标指向的元素类型：新建一个指针并写入。
		if t.Elem() == vt || vt.AssignableTo(t.Elem()) {
			p := reflect.New(t.Elem())
			p.Elem().Set(rv)
			return p, true
		}
		tv = reflect.New(t.Elem()).Elem()
		t = t.Elem()
		isPtr = true
	} else {
		tv = reflect.New(t).Elem()
	}

	// ret 根据是否为指针目标，返回 tv 本身或其地址。
	ret := func() (reflect.Value, bool) {
		if isPtr {
			return tv.Addr(), true
		}
		return tv, true
	}

	// 5. 可直接转换（如 type MyInt int、不同数值别名、string<->[]byte 等）。
	//    注意：避免数值/字符串互转走这里导致语义偏差，这类由下方分支显式处理。
	if vt.ConvertibleTo(t) && !isLossyConvert(vt, t) {
		tv.Set(rv.Convert(t))
		return ret()
	}

	switch t.Kind() {
	case reflect.String:
		tv.SetString(StringOf(v))
		return ret()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tv.SetInt(anyToInt64(v))
		return ret()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		tv.SetUint(anyToUint64(v))
		return ret()
	case reflect.Float32, reflect.Float64:
		tv.SetFloat(NumberOf(v))
		return ret()
	case reflect.Bool:
		tv.SetBool(BoolOf(v))
		return ret()
	case reflect.Struct:
		obj, ok := asStringMap(v)
		if !ok {
			return vv, false
		}
		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)
			if !fi.IsExported() {
				continue
			}
			fieldV, exist := lookupStructField(obj, fi)
			if !exist || fieldV == nil {
				continue
			}
			ftv, ok := structValConvert(ctx, fi.Type, fieldV)
			if !ok {
				return vv, false
			}
			tv.Field(i).Set(ftv)
		}
		return ret()
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			if s, isStr := v.(string); isStr {
				tv.SetBytes(ToBytes(s))
				return ret()
			}
		}
		obj, ok := asAnySlice(v)
		if !ok {
			return vv, false
		}
		tvp := reflect.MakeSlice(t, 0, len(obj))
		for _, a := range obj {
			ftv, ok := structValConvert(ctx, t.Elem(), a)
			if !ok {
				return vv, false
			}
			tvp = reflect.Append(tvp, ftv)
		}
		tv.Set(tvp)
		return ret()
	case reflect.Array:
		obj, ok := asAnySlice(v)
		if !ok {
			return vv, false
		}
		n := tv.Len()
		for i, a := range obj {
			if i >= n {
				break
			}
			ftv, ok := structValConvert(ctx, t.Elem(), a)
			if !ok {
				return vv, false
			}
			tv.Index(i).Set(ftv)
		}
		return ret()
	case reflect.Map:
		obj, ok := asStringMap(v)
		if !ok {
			return vv, false
		}
		tvp := reflect.MakeMapWithSize(t, len(obj))
		for key, val := range obj {
			mk, ok := structValConvert(ctx, t.Key(), key)
			if !ok {
				return vv, false
			}
			nv, ok := structValConvert(ctx, t.Elem(), val)
			if !ok {
				return vv, false
			}
			tvp.SetMapIndex(mk, nv)
		}
		tv.Set(tvp)
		return ret()

	case reflect.Interface:
		tv.Set(rv)
		return ret()

	case reflect.Func:
		val, ok := v.(Val)
		if !ok {
			return vv, false
		}
		fn := reflect.MakeFunc(t, func(args []reflect.Value) []reflect.Value {
			argvs := make([]any, len(args))
			for i, arg := range args {
				argvs[i] = arg.Interface()
			}
			res := RunLambda(ctx, val, argvs...)
			switch t.NumOut() {
			case 0:
				return nil
			case 1:
				data, ok := structValConvert(ctx, t.Out(0), res)
				if !ok {
					data = reflect.Zero(t.Out(0))
				}
				return []reflect.Value{data}
			default:
				rr := make([]reflect.Value, t.NumOut())
				arr, _ := res.([]any)
				for i := 0; i < t.NumOut(); i++ {
					var vvv any
					if i < len(arr) {
						vvv = arr[i]
					}
					data, ok := structValConvert(ctx, t.Out(i), vvv)
					if !ok {
						data = reflect.Zero(t.Out(i))
					}
					rr[i] = data
				}
				return rr
			}
		})
		tv.Set(fn)
		return ret()
	default:
		return vv, false
	}
}

// isLossyConvert 判断从 vt 到 t 的 reflect 转换是否会丢失语义，
// 这类转换不应走通用 Convert 分支，而由显式基础类型分支处理。
// 典型：数值<->字符串（Go 会把 int 当 rune 转成字符）。
func isLossyConvert(vt, t reflect.Type) bool {
	vk, tk := vt.Kind(), t.Kind()
	isNum := func(k reflect.Kind) bool {
		switch k {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64:
			return true
		}
		return false
	}
	// 数值 <-> 字符串 之间不要用 Convert。
	if (isNum(vk) && tk == reflect.String) || (vk == reflect.String && isNum(tk)) {
		return true
	}
	return false
}

// lookupStructField 从 map 中按字段名或 json tag 查找字段值。
// 优先精确字段名，其次 json tag 名。
func lookupStructField(obj map[string]any, fi reflect.StructField) (any, bool) {
	if val, exist := obj[fi.Name]; exist {
		return val, true
	}
	tag := fi.Tag.Get("json")
	if tag != "" && tag != "-" {
		name := tag
		if idx := strings.IndexByte(tag, ','); idx >= 0 {
			name = tag[:idx]
		}
		if name != "" {
			if val, exist := obj[name]; exist {
				return val, true
			}
		}
	}
	return nil, false
}

func structValueToVm(force bool, vv any) any {
	if !force {
		return vv
	}
	switch v := vv.(type) {
	case string, float64, bool, []any, map[string]interface{}, []byte:
		return vv
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case *string:
		return *v
	case *int:
		return float64(*v)
	case *int8:
		return float64(*v)
	default:
		return v
	}
}

func setFieldOfStruct(ctx *Context, rv reflect.Value, name string, val any) {
	switch rv.Kind() {
	case reflect.Struct:
		fv := rv.FieldByName(name)
		if !fv.IsValid() {
			return
		}
		setValStruct(ctx, fv, val)
	case reflect.Ptr:
		if !rv.IsNil() {
			setFieldOfStruct(ctx, rv.Elem(), name, val)
			return
		}
	case reflect.Map:
		v, ok := structValConvert(ctx, rv.Type().Elem(), val)
		if !ok {
			return
		}
		if rv.Type().Key().Kind() != reflect.String {
			return
		}
		rv.SetMapIndex(reflect.ValueOf(name), v)
	default:
	}
}

func getIndexOfSlice(ctx *Context, forceType bool, rv reflect.Value, v any) any {
	switch rv.Kind() {
	case reflect.Ptr:
		if !rv.IsNil() {
			return getIndexOfSlice(ctx, forceType, rv.Elem(), v)
		}
		return nil
	case reflect.Slice:
		idx := int(NumberOf(v))
		if idx >= rv.Len() {
			return nil
		}
		v := rv.Index(idx)
		if v.Kind() == reflect.Struct {
			return v.Addr().Interface()
		}
		return structValueToVm(forceType, v.Interface())
	case reflect.Map:
		keyType := rv.Type().Key()

		kk, ok := structValConvert(ctx, keyType, v)
		if !ok {
			return nil
		}
		ssv := rv.MapIndex(kk)
		if ssv.IsValid() {
			return ssv.Interface()
		}
		return nil
	default:
		return nil
	}
}

func setIndexOfStruct(ctx *Context, rv reflect.Value, idx int, val any) {
	switch rv.Kind() {
	case reflect.Ptr:
		if !rv.IsNil() {
			setIndexOfStruct(ctx, rv.Elem(), idx, val)
		}
	case reflect.Slice:
		if idx >= rv.Len() {
			//rvc := rv
			//var e reflect.Value
			//for i := 0; i <= idx-rv.Len(); i++ {
			//	e = reflect.New(rv.Type().Elem()).Elem()
			//	rvc = reflect.Append(rvc, e)
			//}
			//v, ok := structValConvert(rv.Type().Elem(), val)
			//if ok {
			//	e.Set(v)
			//	return
			//}
			//
			//rv.Set(rvc)

			return
		}
		v, ok := structValConvert(ctx, rv.Type().Elem(), val)
		if !ok {
			return
		}
		rv.Index(idx).Set(v)

	default:

	}
}

func lenOfStruct(rv reflect.Value) int64 {
	switch rv.Kind() {
	case reflect.Ptr:
		if !rv.IsNil() {
			return lenOfStruct(rv.Elem())
		}
	case reflect.Slice:
		return int64(rv.Len())
	}
	return 0
}

var (
	contextType = reflect.TypeOf(new(Context))
	optionType  = reflect.TypeOf(new(Options))
)

func callFunc(ctx *Context, fv reflect.Value, args []Val) (ress any) {
	ft := fv.Type()
	fvls := make([]reflect.Value, 0, ft.NumIn())

	offset := 0
	isVariadic := ft.IsVariadic()
	for i := 0; i < ft.NumIn(); i++ {

		argi := ft.In(i)

		if i == 0 && argi == contextType {
			fvls = append(fvls, reflect.ValueOf(ctx))
			offset = 1
			continue
		}

		var argv any

		if isVariadic && i-offset == ft.NumIn()-1 {
			argvs := []any{}
			for j := i - offset; j < len(args); j++ {

				vadv, ok := args[j].(*VariadicVal)
				if ok {
					argvs = append(argvs, vadv.ArrVal(ctx)...)
				} else {
					e := args[j].Val(ctx)
					argvs = append(argvs, e)
				}

			}
			argv = argvs
		} else {
			if i-offset < len(args) {
				argv = args[i-offset].Val(ctx)
			}
		}

		if i == ft.NumIn()-1 && argi == optionType {
			o, _ := argv.(map[string]any)
			fvls = append(fvls, reflect.ValueOf(&Options{data: o}))
			continue
		}

		v, ok := structValConvert(ctx, argi, argv)
		if !ok {
			return newErrorfWithCtx(ctx, "faile to call '%s' arg  type is not support: %v", argi.Name(), argi.String())
		}
		fvls = append(fvls, v)
	}
	var res []reflect.Value
	if isVariadic {
		res = fv.CallSlice(fvls)
	} else {
		res = fv.Call(fvls)
	}
	if len(res) == 0 {
		return nil
	}
	if len(res) == 1 {
		return res[0].Interface()
	}
	rrss := make([]any, len(res))
	for i, re := range res {
		rrss[i] = re.Interface()
	}
	return rrss
}

func callFuncByReflect(ctx *Context, f *objFuncVal, v any, args []Val) (ress any, ok bool) {

	if v == nil {
		return nil, false
	}
	vv := reflect.ValueOf(v)
	fv := vv.MethodByName(f.funcName)
	if !fv.IsValid() {
		return nil, false
	}
	return callFunc(ctx, fv, args), true
}
