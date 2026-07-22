package expr

import (
	"fmt"
	"reflect"

	"github.com/seeadoog/expr/ast"
)

// parser_handler.go 拆分了 Env.ParseValueFromNode 中各 AST 节点类型的解析逻辑。
// 每个节点类型对应一个独立的 handler 函数，在 init 中通过 RegisterParseFunc 注册。

func init() {
	RegisterParseFunc(parseNodeString)
	RegisterParseFunc(parseNodeNumber)
	RegisterParseFunc(parseNodeBool)
	RegisterParseFunc(parseNodeNil)
	RegisterParseFunc(parseNodeVariable)
	RegisterParseFunc(parseNodeAccess)
	RegisterParseFunc(parseNodeCall)
	RegisterParseFunc(parseNodeUnary)
	RegisterParseFunc(parseNodeBinary)
	RegisterParseFunc(parseNodeSet)
	RegisterParseFunc(parseNodeMapSet)
	RegisterParseFunc(parseNodeArrDef)
	RegisterParseFunc(parseNodeArrAccess)
	RegisterParseFunc(parseNodeSliceCut)
	RegisterParseFunc(parseNodeLambda)
	RegisterParseFunc(parseNodeLambda2)
	RegisterParseFunc(parseNodeTernary)
	RegisterParseFunc(parseNodeConst)
	RegisterParseFunc(parseNodeNotNil)
}

func parseNodeString(e *Env, n *ast.String, isAccess bool, pc *ParserContext) (Val, error) {
	sp := &strparser{
		str: []rune(n.Val),
	}

	err := sp.parser()
	if err != nil {
		return nil, fmt.Errorf("parse string: %w %s", err, n.Val)
	}
	if len(sp.vals) == 1 && sp.vals[0].kind == 0 {
		return &constraint{
			value: sp.vals[0].val,
		}, nil
	}
	if len(sp.vals) == 0 {
		return &constraint{""}, nil
	}
	v, err := e.parseStrVals(sp.vals)
	if err != nil {
		return nil, fmt.Errorf("parse string: %w %s", err, n.Val)
	}
	return v, nil
}

func parseNodeNumber(e *Env, n *ast.Number, isAccess bool, pc *ParserContext) (Val, error) {
	return &constraint{
		value: n.Val,
	}, nil
}

func parseNodeBool(e *Env, n *ast.Bool, isAccess bool, pc *ParserContext) (Val, error) {
	return &constraint{
		value: n.Val,
	}, nil
}

func parseNodeNil(e *Env, n *ast.Nil, isAccess bool, pc *ParserContext) (Val, error) {
	return &constraint{}, nil
}

func parseNodeVariable(e *Env, n *ast.Variable, isAccess bool, pc *ParserContext) (Val, error) {
	switch n.Name {
	case "break":
		return &breakVar{}, nil
	case "_":
		return &emptyVal{}, nil
	}

	return &variable{
		varName: n.Name,
		hash:    calcHash(n.Name),
	}, nil
}

func parseNodeAccess(e *Env, n *ast.Access, isAccess bool, pc *ParserContext) (Val, error) {
	lv, err := e.ParseValueFromNode(n.L, false, pc)
	if err != nil {
		return nil, fmt.Errorf("binary parse val L error:%w %s", err, lv)
	}
	rv, err := e.ParseValueFromNode(n.R, true, pc)
	if err != nil {
		return nil, fmt.Errorf("binary parse val R error:%w %s", err, lv)
	}
	rf, ok := rv.(*objFuncVal)
	if ok {
		if e.allTypeFuncs[rf.funcName] {
			fun := e.funtables[rf.funcName]
			if fun == nil {
				return nil, fmt.Errorf("binary parse val function is all type funcs but not defined '%s'", rf.funcName)
			}
			if fun.hasOpt {
				if fun.argsNum >= 0 && (fun.argsNum != len(rf.args)+1 && fun.argsNum != len(rf.args)) {
					return nil, fmt.Errorf("binary parse val function '%s' args num should be %d  but  %d", rf.funcName, fun.argsNum-1, len(rf.args))
				}
				if fun.argsNum == len(rf.args) { // has opt
					optArg, ok := rf.args[len(rf.args)-1].(*mapDefineVal)
					if !ok {
						return nil, fmt.Errorf("binary parse val function '%s' last extra arg type should be object but  %s", rf.funcName, reflect.TypeOf(optArg).String())
					}
					conv, ok := tryConvertToConst(optArg).(*constraint)
					if !ok {
						return nil, fmt.Errorf("binary parse val function '%s' last extra arg type should be const object", rf.funcName)
					}
					conv.value = newOption(conv.value)
					rf.args[len(rf.args)-1] = conv

				}
			} else {
				if fun.argsNum >= 0 && fun.argsNum != len(rf.args)+1 {
					return nil, fmt.Errorf("binary parse val function '%s' args num should be %d  but  %d", rf.funcName, fun.argsNum-1, len(rf.args))
				}
			}

			return &funcVariable{
				funcNameHash: calcHash(rf.funcName),
				funcName:     rf.funcName,
				fun:          fun.fun,
				args:         append([]Val{lv}, rf.args...),
			}, nil
		}
	}
	av := &accessVal{
		left:  lv,
		right: rv,
	}
	if shouldCompileIf(av) {
		return tryCompileVal(av), nil
	}

	return av, nil
}

func parseNodeCall(e *Env, n *ast.Call, isAccess bool, pc *ParserContext) (Val, error) {
	if isAccess {
		args := make([]Val, 0, len(n.Args))
		for _, arg := range n.Args {
			argv, err := e.ParseValueFromNode(arg, false, pc)
			if err != nil {
				return nil, err
			}
			args = append(args, argv)
		}
		return &objFuncVal{
			funNameHash: calcHash(n.Name),
			args:        args,
			funcName:    n.Name,
		}, nil
	}
	fun := e.funtables[n.Name]

	hasOpt := false
	if fun != nil {
		if fun.hasOpt {
			switch {
			case len(n.Args) == fun.argsNum:

			case len(n.Args) == fun.argsNum+1:
				switch an := n.Args[fun.argsNum].(type) {
				case *ast.MapSet:
				case *ast.Const:
					_, ok := an.L.(*ast.MapSet)
					if !ok {
						return nil, fmt.Errorf("func '%s' option arg should be  define as object", n.Name)
					}
					n.Args[fun.argsNum] = an.L
				default:
					return nil, fmt.Errorf("func '%s' option arg should be  define as object", n.Name)
				}

				hasOpt = true

			case fun.argsNum == -1:
				return nil, fmt.Errorf("func '%s' has option, args num cannot be -1", n.Name)
			default:
				return nil, fmt.Errorf("func '%s' args num should be '%d' but '%d'", n.Name, fun.argsNum, len(n.Args))
			}
		} else {
			if fun.argsNum != -1 && len(n.Args) != fun.argsNum {
				return nil, fmt.Errorf("func '%s' args num should be '%d' but '%d'", n.Name, fun.argsNum, len(n.Args))
			}
		}

	}

	args := make([]Val, 0, len(n.Args))
	for _, arg := range n.Args {
		argv, err := e.ParseValueFromNode(arg, false, pc)
		if err != nil {
			return nil, err
		}
		args = append(args, argv)
	}
	if hasOpt {
		ov := args[len(args)-1].(*mapDefineVal)
		ov.isOpt = true
		ovconst, ok := tryConvertToConst(ov).(*constraint)
		if ok {
			ovconst.value = newOption(ovconst.value)
			args[len(args)-1] = ovconst
		}

	}
	var f ScriptFunc
	if fun != nil {
		f = fun.fun
	}

	if fun != nil {
		if fun.compiledArgs > 0 {
			var err error
			args, err = newCompileFunc(args, fun.compiledArgs, fun.compileFunc)
			if err != nil {
				return nil, fmt.Errorf("compile func '%s' error: %w", fun.name, err)
			}
		}
	}

	return &funcVariable{
		funcNameHash: calcHash(n.Name),
		funcName:     n.Name,
		fun:          f,
		args:         args,
	}, nil
}

func parseNodeUnary(e *Env, n *ast.Unary, isAccess bool, pc *ParserContext) (Val, error) {
	val, err := e.ParseValueFromNode(n.X, false, pc)
	if err != nil {
		return nil, fmt.Errorf("unary parse val error:%w", err)
	}
	switch n.Op {
	case "!":
		return newUnaryValue("!", val, func(ctx *Context, a Val) any {
			return !BoolCond(val.Val(ctx))
		}), nil
	case "-":
		return newUnaryValue("-", val, func(ctx *Context, a Val) any {
			v, _ := val.Val(ctx).(float64)
			return -v
		}), nil
	case "...":
		return &VariadicVal{val}, nil
	case "++":
		return &addAddVal{val}, nil
	case "--":
		return &subSubVal{val}, nil
	}
	return nil, fmt.Errorf("unknown unary operator:%s", n.Op)
}

func parseNodeBinary(e *Env, n *ast.Binary, isAccess bool, pc *ParserContext) (Val, error) {
	lv, err := e.ParseValueFromNode(n.L, false, pc)
	if err != nil {
		return nil, fmt.Errorf("binary parse val L error:%w %s", err, lv)
	}
	rv, err := e.ParseValueFromNode(n.R, false, pc)
	if err != nil {
		return nil, fmt.Errorf("binary parse val R error:%w %s", err, lv)
	}
	var fun ScriptFunc
	switch n.Op {
	case "+":
		fun = add2Func
	case "-":
		fun = subFunc
	case "*":
		fun = mulFunc
	case "/":
		fun = divFunc
	case "^":
		fun = powFunc
	case "&&":
		fun = andFunc
		return newBinaryValue("&&", lv, rv, func(ctx *Context, a, b Val) any {
			if !BoolCond(a.Val(ctx)) {
				return false
			}
			return BoolCond(b.Val(ctx))
		}), nil
	case "||":
		return newBinaryValue("||", lv, rv, func(ctx *Context, a, b Val) any {
			if BoolCond(a.Val(ctx)) {
				return true
			}
			return BoolCond(b.Val(ctx))
		}), nil
	case "==":
		return &eqVal{
			L: lv, R: rv,
		}, nil
	case "<":
		fun = lessFunc
	case "<=":
		fun = lessOrEqual
	case ">":
		fun = largeFunc
	case ">=":
		fun = largeOrEqual
	case "!=":
		return newBinaryValue("!=", lv, rv, func(ctx *Context, a, b Val) any {
			return a.Val(ctx) != b.Val(ctx)
		}), nil
	case "%":
		fun = modFunc

	case "orr":
		return newBinaryValue("orr", lv, rv, func(ctx *Context, a, b Val) any {
			v := a.Val(ctx)
			switch v.(type) {
			case nil:
				return b.Val(ctx)
			}
			return v
		}), nil
	case ";":
		return parAsList(lv, rv), nil
	case "in":
		fun = inFunc

		switch rv.(type) {
		case *arrDefVal, *mapDefineVal:
			rv = tryConvertToConst(rv)
		}

	case "|":
		return newBinaryValue("|", lv, rv, func(ctx *Context, a, b Val) any {
			return float64(int(NumberOf(a.Val(ctx))) | int(NumberOf(b.Val(ctx))))
		}), nil
	case "&":
		return newBinaryValue("&", lv, rv, func(ctx *Context, a, b Val) any {
			return float64(int(NumberOf(a.Val(ctx))) & int(NumberOf(b.Val(ctx))))
		}), nil
	case "+=":
		return &setValue{
			key: lv,
			val: newBinaryValue("+=", lv, rv, func(ctx *Context, a, b Val) any {
				return add2(a.Val(ctx), b.Val(ctx))
			}),
		}, nil
	case "-=":
		return &setValue{
			key: lv,
			val: newBinaryValue("-=", lv, rv, func(ctx *Context, a, b Val) any {
				return NumberOf(a.Val(ctx)) - NumberOf(b.Val(ctx))
			}),
		}, nil
	case "*=":
		return &setValue{
			key: lv,
			val: newBinaryValue("*=", lv, rv, func(ctx *Context, a, b Val) any {
				return NumberOf(a.Val(ctx)) * NumberOf(b.Val(ctx))
			}),
		}, nil
	case "/=":
		return &setValue{
			key: lv,
			val: newBinaryValue("/=", lv, rv, func(ctx *Context, a, b Val) any {
				return NumberOf(a.Val(ctx)) / NumberOf(b.Val(ctx))
			}),
		}, nil
	case "as":
		varv, ok := rv.(*variable)
		if !ok {
			return nil, fmt.Errorf("as right is not variable:%v", n)
		}
		return &asVal{
			key:     varv.varName,
			keyHash: varv.hash,
			val:     lv,
		}, nil
	case "===":
		return &eqValT{
			L: lv,
			R: rv,
		}, nil
	case "!==":
		return &notEqValT{
			L: lv,
			R: rv,
		}, nil
	default:
		return nil, fmt.Errorf("unknown operator of binary :%v %v", n.Op, n)
	}
	return &funcVariable{
		funcNameHash: calcHash(n.Op),
		funcName:     n.Op,
		fun:          fun,
		args:         []Val{lv, rv},
	}, nil
}

func parseNodeSet(e *Env, n *ast.Set, isAccess bool, pc *ParserContext) (Val, error) {
	key, err := e.ParseValueFromNode(n.L, false, pc)
	if err != nil {
		return nil, fmt.Errorf("set parse key error:%w %s", err, key)
	}
	val, err := e.ParseValueFromNode(n.R, false, pc)
	if err != nil {
		return nil, fmt.Errorf("set parse val error:%w", err)
	}
	if n.Const {
		val = tryConvertToConst(val)
		_, ok := val.(*constraint)
		if !ok {
			return nil, fmt.Errorf("set parse val error,val cannot parse as const %T", n.R)
		}
	}
	return &setValue{
		key: key,
		val: val,
	}, nil
}

func parseNodeMapSet(e *Env, n *ast.MapSet, isAccess bool, pc *ParserContext) (Val, error) {
	mapkvs := make([]mapKv, 0, len(n.Kvs))
	for _, kv := range n.Kvs {
		kk, err := e.ParseValueFromNode(kv.K, false, pc)
		if err != nil {
			return nil, fmt.Errorf("map parse key error:%w", err)
		}
		vv, err := e.ParseValueFromNode(kv.V, false, pc)
		if err != nil {
			return nil, fmt.Errorf("map parse value error:%w", err)
		}
		mapkvs = append(mapkvs, mapKv{kk, vv})
	}
	mv := &mapDefineVal{
		kvs: mapkvs,
	}
	return mv, nil
}

func parseNodeArrDef(e *Env, n *ast.ArrDef, isAccess bool, pc *ParserContext) (Val, error) {
	arrV := &arrDefVal{}
	for i, n2 := range n.V {
		v, err := e.ParseValueFromNode(n2, false, pc)
		if err != nil {
			return nil, fmt.Errorf("array parse error:%w %v", err, i)
		}
		arrV.vs = append(arrV.vs, v)
	}
	return arrV, nil
}

func parseNodeArrAccess(e *Env, n *ast.ArrAccess, isAccess bool, pc *ParserContext) (Val, error) {
	arrV := &arrAccessVal{}
	lv, err := e.ParseValueFromNode(n.L, false, pc)
	if err != nil {
		return nil, fmt.Errorf("array access parse left error:%w %v", err, n.L)
	}
	rv, err := e.ParseValueFromNode(n.R, false, pc)
	if err != nil {
		return nil, fmt.Errorf("array access parse right error:%w %v", err, n.R)
	}
	arrV.left = lv
	arrV.right = rv
	return arrV, nil
}

func parseNodeSliceCut(e *Env, n *ast.SliceCut, isAccess bool, pc *ParserContext) (Val, error) {
	v, err := e.ParseValueFromNode(n.V, false, pc)
	if err != nil {
		return nil, fmt.Errorf("slice cut parse value error:%w %v", err, n.V)
	}
	var st, ed Val
	if n.St != nil {
		st, err = e.ParseValueFromNode(n.St, false, pc)
		if err != nil {
			return nil, fmt.Errorf("slice cut parse st  error:%w %v", err, n.V)
		}
	}

	if n.Ed != nil {
		ed, err = e.ParseValueFromNode(n.Ed, false, pc)
		if err != nil {
			return nil, fmt.Errorf("slice cut parse ed error:%w %v", err, n.V)
		}
	}

	return &sliceCutVal{
		st:  st,
		ed:  ed,
		val: v,
	}, nil
}

func parseNodeLambda(e *Env, n *ast.Lambda, isAccess bool, pc *ParserContext) (Val, error) {
	r, err := e.ParseValueFromNode(n.R, false, pc)
	if err != nil {
		return nil, fmt.Errorf("lambda parse right error:%w %v", err, n.R)
	}
	lm := &lambda{
		Lefts:     n.L,
		Right:     r,
		leftsHash: hashOfStrings(n.L),
	}
	return lm, nil
}

func parseNodeLambda2(e *Env, n *ast.Lambda2, isAccess bool, pc *ParserContext) (Val, error) {
	r, err := e.ParseValueFromNode(n.R, false, pc)
	if err != nil {
		return nil, fmt.Errorf("lambda parse right error:%w %v", err, n.R)
	}
	lm := &lambda{
		Right: r,
	}

	switch l := n.L.(type) {
	case *ast.Variable:
		lm.Lefts = []string{l.Name}
	case *ast.ArrDef:
		for _, le := range l.V {

			ln, ok := le.(*ast.Variable)
			if ok {
				return nil, fmt.Errorf("lambda parse left array elem type is not variable but :%v", reflect.TypeOf(l))
			}
			lm.Lefts = append(lm.Lefts, ln.Name)
		}
	default:
		return nil, fmt.Errorf("lambda parse right array elem type is invaid:%s", reflect.TypeOf(n.R))
	}
	lm.leftsHash = hashOfStrings(lm.Lefts)

	return lm, nil
}

func parseNodeTernary(e *Env, n *ast.Ternary, isAccess bool, pc *ParserContext) (Val, error) {
	c, err := e.ParseValueFromNode(n.C, false, pc)
	if err != nil {
		return nil, fmt.Errorf("ternary parse cond error:%w %v", err, n.R)
	}
	l, err := e.ParseValueFromNode(n.L, false, pc)
	if err != nil {
		return nil, fmt.Errorf("ternary parse left error:%w %v", err, n.R)
	}
	var r Val
	if n.R != nil {
		r, err = e.ParseValueFromNode(n.R, false, pc)
		if err != nil {
			return nil, fmt.Errorf("ternary parse right error:%w %v", err, n.R)
		}
	}

	return &ternaryVal{
		c: c,
		l: l,
		r: r,
	}, nil
}

func parseNodeConst(e *Env, n *ast.Const, isAccess bool, pc *ParserContext) (Val, error) {
	lv, err := e.ParseValueFromNode(n.L, false, pc)
	if err != nil {
		return nil, fmt.Errorf("const val parse left error:%w %v", err, n.L)
	}

	cv := tryConvertToConst(lv)
	return cv, nil
}

func parseNodeNotNil(e *Env, n *ast.NotNil, isAccess bool, pc *ParserContext) (Val, error) {
	lv, err := e.ParseValueFromNode(n.N, isAccess, pc)
	if err != nil {
		return nil, fmt.Errorf("not nil val parse  error:%w %v", err, n.N)
	}

	return &notNil{
		val: lv,
	}, nil
}
