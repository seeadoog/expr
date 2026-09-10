package interpreter

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

// VM 虚拟机，执行字节码
type VM struct {
	bc     *ByteCode // 待执行的字节码
	env    Env       // 宿主环境（函数查找）
	ctx    Context   // 执行上下文（变量存取）
	stack  []any     // 操作数栈
	sp     int       // 栈指针
	pc     int       // 程序计数器
	locals []any     // 局部变量槽
}

// NewVM 创建虚拟机实例
func NewVM(bc *ByteCode, env Env, ctx Context) *VM {
	return &VM{
		bc:     bc,
		env:    env,
		ctx:    ctx,
		stack:  make([]any, 0, 128),
		sp:     0,
		pc:     0,
		locals: make([]any, bc.MaxLocals),
	}
}

// Reset 重置虚拟机状态，用于复用同一 VM 实例多次执行
func (vm *VM) Reset() {
	vm.stack = vm.stack[:0]
	vm.sp = 0
	vm.pc = 0
	// 清空 locals
	for i := range vm.locals {
		vm.locals[i] = nil
	}
}

// Run 执行字节码，返回最终栈顶值
func (vm *VM) Run() (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("VM panic: %v", r)
		}
	}()

	for vm.pc < len(vm.bc.Instructions) {
		inst := vm.bc.Instructions[vm.pc]
		vm.pc++

		switch inst.Op {
		case OpPush:
			vm.push(vm.bc.Constants[inst.A])

		case OpPop:
			vm.pop()

		case OpLoadVar:
			nameRef := vm.bc.Names[inst.A]
			val := vm.ctx.Get(nameRef.Hash)
			vm.push(val)

		case OpStoreVar:
			nameRef := vm.bc.Names[inst.A]
			val := vm.peek()
			vm.ctx.Set(nameRef.Hash, val)

		case OpLoadFast:
			vm.push(vm.locals[inst.A])

		case OpStoreFast:
			vm.locals[inst.A] = vm.peek()

		case OpAdd:
			vm.binaryOp(vm.add)
		case OpSub:
			vm.binaryOp(func(a, b any) any { return vm.numberOf(a) - vm.numberOf(b) })
		case OpMul:
			vm.binaryOp(func(a, b any) any { return vm.numberOf(a) * vm.numberOf(b) })
		case OpDiv:
			vm.binaryOp(func(a, b any) any { return vm.numberOf(a) / vm.numberOf(b) })
		case OpMod:
			vm.binaryOp(func(a, b any) any { return math.Mod(vm.numberOf(a), vm.numberOf(b)) })
		case OpPow:
			vm.binaryOp(func(a, b any) any { return math.Pow(vm.numberOf(a), vm.numberOf(b)) })
		case OpNeg:
			vm.unaryOp(func(a any) any { return -vm.numberOf(a) })

		case OpEq:
			vm.binaryOp(func(a, b any) any { return a == b })
		case OpNeq:
			vm.binaryOp(func(a, b any) any { return a != b })
		case OpLt:
			vm.binaryOp(func(a, b any) any { return vm.numberOf(a) < vm.numberOf(b) })
		case OpLte:
			vm.binaryOp(func(a, b any) any { return vm.numberOf(a) <= vm.numberOf(b) })
		case OpGt:
			vm.binaryOp(func(a, b any) any { return vm.numberOf(a) > vm.numberOf(b) })
		case OpGte:
			vm.binaryOp(func(a, b any) any { return vm.numberOf(a) >= vm.numberOf(b) })
		case OpEqT:
			vm.binaryOp(func(a, b any) any { return vm.eqt(a, b) })
		case OpNeqT:
			vm.binaryOp(func(a, b any) any { return !vm.eqt(a, b) })

		case OpNot:
			vm.unaryOp(func(a any) any { return !vm.boolOf(a) })

		case OpBitAnd:
			vm.binaryOp(func(a, b any) any {
				return float64(int64(vm.numberOf(a)) & int64(vm.numberOf(b)))
			})
		case OpBitOr:
			vm.binaryOp(func(a, b any) any {
				return float64(int64(vm.numberOf(a)) | int64(vm.numberOf(b)))
			})
		case OpBitXor:
			vm.binaryOp(func(a, b any) any {
				return float64(int64(vm.numberOf(a)) ^ int64(vm.numberOf(b)))
			})
		case OpBitNot:
			vm.unaryOp(func(a any) any {
				return float64(^int64(vm.numberOf(a)))
			})
		case OpShl:
			vm.binaryOp(func(a, b any) any {
				return float64(int64(vm.numberOf(a)) << uint64(vm.numberOf(b)))
			})
		case OpShr:
			vm.binaryOp(func(a, b any) any {
				return float64(int64(vm.numberOf(a)) >> uint64(vm.numberOf(b)))
			})

		case OpGetAttr:
			nameRef := vm.bc.Names[inst.A]
			obj := vm.pop()
			val := vm.getAttr(obj, nameRef.Name)
			vm.push(val)

		case OpSetAttr:
			nameRef := vm.bc.Names[inst.A]
			val := vm.pop()
			obj := vm.pop()
			vm.setAttr(obj, nameRef.Name, val)
			vm.push(val)

		case OpGetIndex:
			idx := vm.pop()
			obj := vm.pop()
			val := vm.getIndex(obj, idx)
			vm.push(val)

		case OpSetIndex:
			val := vm.pop()
			idx := vm.pop()
			obj := vm.pop()
			vm.setIndex(obj, idx, val)
			vm.push(val)

		case OpSlice:
			ed := vm.pop()
			st := vm.pop()
			obj := vm.pop()
			val := vm.slice(obj, st, ed)
			vm.push(val)

		case OpCall:
			callInfo := vm.bc.Calls[inst.A]
			args := make([]any, callInfo.Argc)
			for i := callInfo.Argc - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			result := vm.callFunc(callInfo.Hash, callInfo.Name, args)
			vm.push(result)

		case OpCallMethod:
			callInfo := vm.bc.Calls[inst.A]
			args := make([]any, callInfo.Argc)
			for i := callInfo.Argc - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			obj := vm.pop()
			result := vm.callMethod(obj, callInfo.Hash, callInfo.Name, args)
			vm.push(result)

		case OpJump:
			vm.pc = int(inst.A)

		case OpJumpIfTrue:
			cond := vm.pop()
			if vm.boolOf(cond) {
				vm.pc = int(inst.A)
			}

		case OpJumpIfFalse:
			cond := vm.pop()
			if !vm.boolOf(cond) {
				vm.pc = int(inst.A)
			}

		case OpMakeArray:
			count := int(inst.A)
			arr := make([]any, count)
			for i := count - 1; i >= 0; i-- {
				arr[i] = vm.pop()
			}
			vm.push(arr)

		case OpMakeMap:
			count := int(inst.A)
			m := make(map[string]any, count)
			for i := 0; i < count; i++ {
				val := vm.pop()
				key := vm.pop().(string)
				m[key] = val
			}
			vm.push(m)

		case OpNotNil:
			val := vm.pop()
			if val == nil {
				return nil, fmt.Errorf("NotNil assertion failed: value is nil")
			}
			vm.push(val)

		case OpHalt:
			if vm.sp > 0 {
				return vm.pop(), nil
			}
			return nil, nil

		default:
			return nil, fmt.Errorf("unknown opcode: %v", inst.Op)
		}
	}

	if vm.sp > 0 {
		return vm.pop(), nil
	}
	return nil, nil
}

// 栈操作
func (vm *VM) push(val any) {
	vm.stack = append(vm.stack, val)
	vm.sp++
}

func (vm *VM) pop() any {
	if vm.sp == 0 {
		panic("stack underflow")
	}
	vm.sp--
	val := vm.stack[vm.sp]
	vm.stack = vm.stack[:vm.sp]
	return val
}

func (vm *VM) peek() any {
	if vm.sp == 0 {
		panic("stack underflow")
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) binaryOp(op func(a, b any) any) {
	b := vm.pop()
	a := vm.pop()
	vm.push(op(a, b))
}

func (vm *VM) unaryOp(op func(a any) any) {
	a := vm.pop()
	vm.push(op(a))
}

// 类型转换（与 expr.NumberOf/StringOf/BoolOf 对齐）
func (vm *VM) numberOf(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case uint64:
		return float64(val)
	case uint:
		return float64(val)
	case int32:
		return float64(val)
	case uint32:
		return float64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		if val == "true" {
			return 1
		}
		return f
	}
	return 0
}

func (vm *VM) boolOf(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case int:
		return val != 0
	case string:
		return val != ""
	case nil:
		return false
	}
	return true
}

func (vm *VM) stringOf(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case *string:
		return *val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case nil:
		return ""
	case []byte:
		return unsafe.String(unsafe.SliceData(val), len(val))
	}
	return fmt.Sprintf("%v", v)
}

func (vm *VM) add(a, b any) any {
	// 字符串拼接
	if _, ok := a.(string); ok {
		return vm.stringOf(a) + vm.stringOf(b)
	}
	if _, ok := b.(string); ok {
		return vm.stringOf(a) + vm.stringOf(b)
	}
	return vm.numberOf(a) + vm.numberOf(b)
}

// eqt 类型转换相等（与 expr/yacc.go 的 eqt 对齐）
func (vm *VM) eqt(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		// nil 与 "" / 0 / false 相等
		switch bv := b.(type) {
		case string:
			return bv == ""
		case float64:
			return bv == 0
		case int:
			return bv == 0
		case bool:
			return !bv
		}
		switch av := a.(type) {
		case string:
			return av == ""
		case float64:
			return av == 0
		case int:
			return av == 0
		case bool:
			return !av
		}
		return false
	}

	// 数字类型互相转换比较
	switch a.(type) {
	case float64, int, int64, uint64:
		switch b.(type) {
		case float64, int, int64, uint64:
			return vm.numberOf(a) == vm.numberOf(b)
		}
	}
	return a == b
}

// 对象属性访问
func (vm *VM) getAttr(obj any, name string) any {
	if obj == nil {
		return nil
	}
	switch o := obj.(type) {
	case map[string]any:
		return o[name]
	default:
		rv := reflect.ValueOf(obj)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			fv := rv.FieldByName(name)
			if fv.IsValid() && fv.CanInterface() {
				return fv.Interface()
			}
		}
	}
	return nil
}

func (vm *VM) setAttr(obj any, name string, val any) {
	if obj == nil {
		return
	}
	switch o := obj.(type) {
	case map[string]any:
		o[name] = val
	default:
		rv := reflect.ValueOf(obj)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			fv := rv.FieldByName(name)
			if fv.IsValid() && fv.CanSet() {
				fv.Set(reflect.ValueOf(val))
			}
		}
	}
}

func (vm *VM) getIndex(obj, idx any) any {
	if obj == nil {
		return nil
	}
	switch o := obj.(type) {
	case []any:
		i := int(vm.numberOf(idx))
		if i >= 0 && i < len(o) {
			return o[i]
		}
	case map[string]any:
		return o[vm.stringOf(idx)]
	case string:
		i := int(vm.numberOf(idx))
		if i >= 0 && i < len(o) {
			return string(o[i])
		}
	default:
		rv := reflect.ValueOf(obj)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			i := int(vm.numberOf(idx))
			if i >= 0 && i < rv.Len() {
				return rv.Index(i).Interface()
			}
		}
	}
	return nil
}

func (vm *VM) setIndex(obj, idx, val any) {
	if obj == nil {
		return
	}
	switch o := obj.(type) {
	case []any:
		i := int(vm.numberOf(idx))
		if i >= 0 && i < len(o) {
			o[i] = val
		}
	case map[string]any:
		o[vm.stringOf(idx)] = val
	default:
		rv := reflect.ValueOf(obj)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			i := int(vm.numberOf(idx))
			if i >= 0 && i < rv.Len() {
				rv.Index(i).Set(reflect.ValueOf(val))
			}
		}
	}
}

func (vm *VM) slice(obj, st, ed any) any {
	if obj == nil {
		return nil
	}
	start := 0
	if st != nil {
		start = int(vm.numberOf(st))
	}

	switch o := obj.(type) {
	case []any:
		end := len(o)
		if ed != nil {
			end = int(vm.numberOf(ed))
		}
		if start < 0 {
			start = 0
		}
		if end > len(o) {
			end = len(o)
		}
		if start > end {
			return []any{}
		}
		return o[start:end]
	case string:
		end := len(o)
		if ed != nil {
			end = int(vm.numberOf(ed))
		}
		if start < 0 {
			start = 0
		}
		if end > len(o) {
			end = len(o)
		}
		if start > end {
			return ""
		}
		return o[start:end]
	default:
		rv := reflect.ValueOf(obj)
		if rv.Kind() == reflect.Slice {
			end := rv.Len()
			if ed != nil {
				end = int(vm.numberOf(ed))
			}
			if start < 0 {
				start = 0
			}
			if end > rv.Len() {
				end = rv.Len()
			}
			if start > end {
				return reflect.MakeSlice(rv.Type(), 0, 0).Interface()
			}
			return rv.Slice(start, end).Interface()
		}
	}
	return nil
}

// callFunc 调用全局函数（通过 Env）
func (vm *VM) callFunc(hash uint64, name string, args []any) any {
	// 这里需要 Env 提供函数查找接口
	// 由于我们定义的 Env 接口还未实现，先返回 nil 占位
	// 实际集成时需要调用 env.GetFunc(hash, name) 并执行
	// 目前只是演示架构
	_ = vm.env
	return fmt.Errorf("callFunc not implemented: %s", name)
}

// callMethod 调用对象方法（通过 objFuncMap 或反射）
func (vm *VM) callMethod(obj any, hash uint64, name string, args []any) any {
	// 同样需要与 expr 的 objFuncMap 集成
	// 目前仅占位
	_ = vm.env
	return fmt.Errorf("callMethod not implemented: %s on %T", name, obj)
}
