package expr

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/seeadoog/expr/ast"
)

const (
	variables = 0
	constant  = -1
	number    = -2
)

type lexer struct {
	tokens []tokenV
	pos    int
	err    []string
	root   ast.Node
}

func (l *lexer) Lex(lval *ast.YySymType) int {
	//TODO implement me
	if l.pos >= len(l.tokens) {
		return 0
	}
	tt := l.tokens[l.pos]
	l.pos++
	switch tt.kind {
	case variables:
		switch tt.tkn {
		case "true":
			lval.SetBool(true)
			lval.SetPos(tt.x, tt.y)
			return ast.BOOL
		case "false":
			lval.SetBool(false)
			lval.SetPos(tt.x, tt.y)
			return ast.BOOL
		case "nil":
			lval.SetPos(tt.x, tt.y)
			return ast.NIL
		}

		if len(tt.tkn) > 0 {
			c := tt.tkn[0]
			if c >= '0' && c <= '9' {
				nn, err := parseNumber(tt.tkn)
				if err == nil {
					lval.SetNum(nn)
					lval.SetPos(tt.x, tt.y)
					return ast.NUMBER
				} else {
					l.Error("invalid number:" + tt.tkn)
				}
			}
		}

		lval.SetStr(tt.tkn)
		lval.SetPos(tt.x, tt.y)
		return ast.IDENT
	case constant:
		lval.SetStr(tt.tkn)
		lval.SetPos(tt.x, tt.y)
		return ast.STRING
	case number:
		lval.SetNum(tt.num)
		lval.SetPos(tt.x, tt.y)
		return ast.NUMBER
	default:
		lval.SetPos(tt.x, tt.y)
		return tt.kind
	}
}

func parseNumber(s string) (float64, error) {
	if strings.HasPrefix(s, "0x") {
		n, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return 0, fmt.Errorf("parser invalid number: %s", s)
		}
		return float64(n), nil
	}
	if strings.Contains(s, ".") || strings.Contains(s, "e") {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("parser invalid number: %s", s)
		}
		return n, nil
	}
	n, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("parser invalid number: %s", s)
	}

	return float64(n), nil
}

func (l *lexer) SetRoot(node ast.Node) {
	l.root = node
}

func (l *lexer) Error(s string) {
	l.err = append(l.err, fmt.Sprintf("'%s' %s near: '%v'  ", func() string {
		if len(l.tokens) == 0 {
			return "no token"
		}
		if l.pos == 0 {
			return "no token and invalid pos 0"
		}
		t := l.tokens[l.pos-1]
		return t.tkn + fmt.Sprintf(" at %d:%d", t.y, t.x)
	}(), s, l.near()))
}

func (l *lexer) near() string {
	next := l.pos + 5
	pre := l.pos - 5
	if pre < 0 {
		pre = 0
	}
	if next > len(l.tokens) {
		next = len(l.tokens)
	}
	ss := l.tokens[pre:next]
	arr := make([]string, 0, 6)
	for _, s := range ss {
		arr = append(arr, s.tkn)
	}
	return strings.Join(arr, " ")
}

type emptyVal struct {
}

func (e *emptyVal) Set(c *Context, val any) {

}

func (e *emptyVal) Val(c *Context) any {
	return nil
}

type ParserContext struct {
}

func CalcKeyHash(key string) uint64 {
	hash := calcHash(key)
	return hash
}

func NewKeyHash(key string) HashKey {
	hash := CalcKeyHash(key)
	return HashKey{
		Key:  key,
		Hash: hash,
	}
}

type HashKey struct {
	Key  string
	Hash uint64
}

func NewParserContext() *ParserContext {
	//pc := &ParserContext{
	//	tb: newEnvMap(8),
	//}
	//for _, key := range mapKeys {
	//	pc.tb.putHash(calcHash(key), key, nil)
	//}
	//for _, key := range arrKeys {
	//	pc.tb.putHash(calcHash(key), key, nil)
	//}
	pc := &ParserContext{}
	return pc
}

type parseFunc func(e *Env, node ast.Node, isAccess bool, ctx *ParserContext) (Val, error)

var parseFuncs = map[reflect.Type]parseFunc{}

func RegisterParseFunc[T ast.Node](f func(e *Env, node T, isAccess bool, ctx *ParserContext) (Val, error)) {
	parseFuncs[reflect.TypeOf(new(T)).Elem()] = func(e *Env, node ast.Node, isAccess bool, ctx *ParserContext) (Val, error) {
		return f(e, node.(T), isAccess, ctx)
	}
}

func (e *Env) parseValueFromNode(node ast.Node, isAccess bool, pc *ParserContext) (Val, error) {
	pf := parseFuncs[reflect.TypeOf(node)]
	if pf == nil {
		return nil, fmt.Errorf("parser function not found,type: '%s'", reflect.TypeOf(node))
	}
	return pf(e, node, isAccess, pc)
}

func (e *Env) ParseValueFromNode(node ast.Node, isAccess bool, pc *ParserContext) (Val, error) {
	return e.parseValueFromNode(node, isAccess, pc)
}

type sliceCutVal struct {
	val Val
	st  Val
	ed  Val
}

func (s *sliceCutVal) Set(c *Context, val any) {

}

func (s *sliceCutVal) Val(c *Context) any {
	//TODO implement me
	f, length := cutterOf(s.val.Val(c))
	if f == nil {
		return nil
	}
	st := 0
	if s.st != nil {
		st = int(NumberOf(s.st.Val(c)))
	}
	ed := length
	if s.ed != nil {
		ed = int(NumberOf(s.ed.Val(c)))
	}
	if st > ed || st < 0 || ed > length {
		return nil
	}
	return f(st, ed)
}

func cutterOf(v any) (func(st, ed int) any, int) {
	switch vs := v.(type) {
	case []any:
		return func(st, ed int) any {
			return vs[st:ed]
		}, len(vs)
	case []byte:
		return func(st, ed int) any {
			return vs[st:ed]
		}, len(vs)
	case string:
		return func(st, ed int) any {
			return vs[st:ed]
		}, len(vs)
	case ReadOnlyArray:
		return func(st, ed int) any {
			return vs[st:ed]
		}, len(vs)
	default:
		return nil, 0
	}
}

type mapKv struct {
	k, v Val
}
type mapDefineVal struct {
	kvs   []mapKv
	isOpt bool
}

func (m *mapDefineVal) Set(c *Context, v any) {

}

func (m *mapDefineVal) Val(c *Context) any {
	mm := make(map[string]any)
	for _, kv := range m.kvs {
		key := ""
		vk, ok := kv.k.(*variable)
		if ok {
			//vvv := kv.k.Val(c)
			//_, ok := vvv.(string)
			//if vvv != nil  && {
			//	key = StringOf(vvv)
			//} else {
			key = vk.varName
			//}
		} else {
			key = StringOf(kv.k.Val(c))
		}
		mm[key] = kv.v.Val(c)
	}
	if m.isOpt {
		return newOption(mm)
	}
	return mm
}

type arrDefVal struct {
	vs []Val
}

func (a *arrDefVal) Val(c *Context) any {
	arr := make([]any, len(a.vs))
	for i, vv := range a.vs {
		arr[i] = vv.Val(c)
	}
	return arr
}

func (a *arrDefVal) Set(c *Context, v any) {}

type strval struct {
	kind int
	val  string
}

type stringFmtVal struct {
	vals []Val
}

var arrPool = sync.Pool{
	New: func() interface{} {
		return make([]string, 0, 3)
	},
}

func (s *stringFmtVal) Val(c *Context) any {

	//sb := strings.Builder{}
	//for _, val := range s.vals {
	//	sb.WriteString(StringOf(val.Val(c)))
	//}
	//return sb.String()
	arr := arrPool.Get().([]string)
	//arr := make([]string, 0, len(s.vals))
	for _, val := range s.vals {
		arr = append(arr, StringOf(val.Val(c)))
	}
	l := 0
	for _, s2 := range arr {
		l += len(s2)
	}
	res := make([]byte, 0, l)
	for _, s2 := range arr {
		res = append(res, s2...)
	}
	arrPool.Put(arr[:0])
	return ToString(res)
}

func (s *stringFmtVal) Set(c *Context, v any) {

}

func (e *Env) parseStrVals(vs []*strval) (Val, error) {

	smt := &stringFmtVal{}
	for _, v := range vs {
		switch v.kind {
		case 0:
			smt.vals = append(smt.vals, &constraint{
				value: v.val,
			})
		case 1:
			vv, err := e.parseValueV(v.val)
			if err != nil {
				return nil, fmt.Errorf("parse fmt value error:%w %s", err, v.val)
			}
			smt.vals = append(smt.vals, vv)
		}
	}
	return smt, nil
}

type strparser struct {
	str   []rune
	pos   int
	vals  []*strval
	token []rune
}

func (s *strparser) next() (rune, bool) {
	if s.pos >= len(s.str) {
		return 0, false
	}
	r := s.str[s.pos]
	s.pos++
	return r, true
}

func (s *strparser) parseVars() error {
	for {
		c, ok := s.next()
		if !ok {
			return fmt.Errorf("unexpected end in string format var ,need '}' to end '${' ")
		}
		switch c {
		//case '\'':
		//	return fmt.Errorf("invalid char ' in string format variable")
		case '}':
			s.appendToken(1)
			return nil
		default:
			s.token = append(s.token, c)
		}
	}
}

func (s *strparser) appendToken(kind int) {
	if len(s.token) == 0 {
		return
	}
	s.vals = append(s.vals, &strval{kind: kind, val: string(s.token)})
	s.token = s.token[:0]

}

func (s *strparser) parser() error {
	for {
		c, ok := s.next()
		if !ok {
			s.appendToken(0)
			return nil
		}
		switch c {
		case '$':
			cc, ok := s.next()
			if !ok {
				s.token = append(s.token, c)
				continue
			}
			if cc == '{' {
				s.appendToken(0)
				err := s.parseVars()
				if err != nil {
					return err
				}
			} else {
				s.token = append(s.token, c)
				s.pos--
			}
		case '\\':
			cc, ok := s.next()
			if !ok {
				return nil
			}
			s.token = append(s.token, cc)

		default:
			s.token = append(s.token, c)
		}
	}
}

type accessVal struct {
	left  Val
	right Val
}

func setForObject(left Val, lv any, right string, c *Context, val any) {
	//rvar, ok := right.(*variable)
	//if !ok {
	//	return
	//}
	//lv := left.Val(c)
	switch parent := lv.(type) {
	case map[string]any:
		parent[right] = val
	case nil:
		pr := map[string]any{}
		left.Set(c, pr)
		//set, ok := left.(parentValueSetter)
		//if ok {
		//	set.Set(c, pr)
		//}
		pr[right] = val
	case *Result:
		switch right {
		case "err":
			parent.Err = val
		case "data":
			parent.Data = val
		}

	case Setter:
		parent.SetField(c, right, val)
	case map[string]string:
		parent[right] = StringOf(val)
	case ReadOnlyMap:

	case ReadOnlyArray:

	default:
		setFieldOfStruct(c, reflect.ValueOf(lv), right, val)

	}
}

func (a *accessVal) Set(c *Context, val any) {

	switch rv := a.right.(type) {
	case *variable:
		setForObject(a.left, a.left.Val(c), rv.varName, c, val)
		//case *compiledVar:
		//	c.stackSet(rv.index, val)
	}
}

// abc::b()::c()::d

var (
	nilType = TypeOf(nil)
)

func callSelf(ctx *Context, self any, f *objFuncVal) (any, bool) {
	s, ok := self.(map[string]any)
	if !ok {
		return nil, false
	}

	ff := s[f.funcName]
	if ff == nil {
		return nil, false
	}
	fun, ok := ff.(*LambdaVal)
	if !ok {
		fv := reflect.ValueOf(ff)
		if fv.Kind() == reflect.Func {
			return callFunc(ctx, fv, f.args), true
		}
		panic(fmt.Sprintf("cannot call func '%s' ,type is not func but :%v", f.funcName, reflect.TypeOf(ff).String()))
	}
	args := make([]any, len(f.args))
	for i, arg := range f.args {
		args[i] = arg.Val(ctx)
	}
	return RunLambda(ctx, fun, args...), true
}

func (a *accessVal) Val(ctx *Context) any {

	switch v := a.right.(type) {
	case *objFuncVal:
		self := a.left.Val(ctx)
		se, ok := self.(*Error)
		if ok {
			return se
		}
		t := TypeOf(self)
		//f := objFuncMap[t]
		f := objFuncMap.get(t)
		if f == nil {

			se, ok := a.left.(*variable)
			if ok {
				lf := ctx.Env.GetLibFunc(se.hash, v.funNameHash)
				if lf != nil {
					return lf(ctx, v.args...)
				}
			}

			fv := reflect.ValueOf(self)
			if fv.Kind() == reflect.Func {
				return callFunc(ctx, fv, v.args)
			}

			data, ok := callFuncByReflect(ctx, v, self, v.args)
			if ok {
				return data
			}
			if ctx.IgnoreFuncNotFoundError {
				return nil
			}
			return newErrorf("var '%s' type '%v' do not define func '%s' ", nameOf(a.left), reflect.TypeOf(self), v.funcName)
		}
		//ff := f[v.funcName]
		ff := f.get(v.funNameHash)
		if ff == nil {

			data, ok := callFuncByReflect(ctx, v, self, v.args)
			if ok {
				return data
			}

			data, ok = callSelf(ctx, self, v)
			if ok {
				return data
			}
			if ctx.IgnoreFuncNotFoundError {
				return nil
			}
			return newErrorf("var '%s' type '%v' do not define func '%s'", nameOf(a.left), reflect.TypeOf(self), v.funcName)
		}
		return ff.fun(ctx, self, v.args...)
	case *variable:
		lv := a.left.Val(ctx)
		//lvv, ok := lv.(map[string]any)
		//if ok {
		//	return lvv[v.varName]
		//}
		switch data := lv.(type) {
		case map[string]any:
			return data[v.varName]

		case ReadOnlyMap:
			return data[v.varName]
		case *Result:
			switch v.varName {
			case "data":
				return data.Data
			case "err":
				return data.Err
			}
			return nil
		case nil:
			return nil

		case map[string]string:
			return data[v.varName]

		case Getter:
			return data.GetField(ctx, v.varName)
		default:

			return getFieldOfStruct(ctx.ForceType, reflect.ValueOf(lv), v.varName)
		}
	//return getFieldOfStruct(reflect.ValueOf(lv), v.varName)

	default:
		return nil
	}
}

type Setter interface {
	SetField(ctx *Context, name string, val any)
}

type Getter interface {
	GetField(c *Context, key string) any
}

type IndexGet interface {
	IndexGet(c *Context, key float64) any
}

type IndexSet interface {
	GetIndSet(ctx *Context, key float64, val any)
}

type parentValueSetter interface {
	Set(c *Context, val any)
}

//// a.b.c
//func (a *accessVal) SetSelf(ctx *Context, v any) {
//	lv := a.left.Val(ctx)
//	if lv == nil {
//		switch lvr := a.left.(type) {
//		case *accessVal:
//			lvrv := lvr.left.Val(ctx)
//			if lvrv == nil {
//				lvr.left.(SetSelf).SetSelf(ctx, map[string]any{})
//			}else{
//				lvr.right.
//				lvrv.(map[string]any)[]
//			}
//		}
//	}
//}

type arrAccessVal struct {
	left  Val
	right Val
}

func (a *arrAccessVal) Set(c *Context, val any) {
	lv := a.left.Val(c)
	rv := a.right.Val(c)

	switch rvv := rv.(type) {
	case string:
		setForObject(a.left, lv, rvv, c, val)
		return

	case float64:
		idx := int(rvv)
		parent, ok := lv.([]any)
		if !ok {
			if lv != nil {
				if _, ok := lv.(ReadOnlyArray); ok {
					return
				}
				setIndexOfStruct(c, reflect.ValueOf(lv), idx, val)
				return
			}
			parent = make([]any, idx+1)
			a.left.Set(c, parent)
		} else {
			if len(parent) <= idx {
				old := parent
				parent = make([]any, idx+1)
				copy(parent, old)
				a.left.Set(c, parent)
			}
		}
		parent[idx] = val
		return
	case nil:
	}
	return
}

func (a *arrAccessVal) Val(ctx *Context) any {
	lv := a.left.Val(ctx)
	rv := a.right.Val(ctx)
	switch v := lv.(type) {
	case []any:
		idx := int(NumberOf(rv))

		if idx >= len(v) {
			return nil
		}
		return v[idx]
	case ReadOnlyMap:
		idx := StringOf(rv)
		return v[idx]
	case ReadOnlyArray:
		idx := int(NumberOf(rv))

		if idx >= len(v) {
			return nil
		}
		return v[idx]
	case []string:
		idx := int(NumberOf(rv))

		if idx >= len(v) {
			return nil
		}
		return v[idx]
	case map[string]any:
		idx := StringOf(rv)
		return v[idx]
	case nil:
		return nil
	default:
		return getIndexOfSlice(ctx, ctx.ForceType, reflect.ValueOf(lv), rv)

	}
	return nil
}

func tryConvertToConst(val Val) (res Val) {

	switch vv := val.(type) {
	case *arrDefVal:
		res = tryCovertArrToConst(vv)
		convertExprReadOnlyVal(res)
		return res
	case *mapDefineVal:
		res = tryCovertMapToConst(vv)
		convertExprReadOnlyVal(res)
		return res
	case *setValue:
		vv.val = tryConvertToConst(vv.val)
		//convertExprReadOnlyVal(vv.val)
	}
	return val
}

func tryCovertArrToConst(val *arrDefVal) Val {
	dst := []any{}
	for _, v := range val.vs {
		vv, ok := tryConvertToConst(v).(*constraint)
		if ok {
			dst = append(dst, vv.value)
		} else {
			return val
		}
	}
	return &constraint{
		value: dst,
	}
}
func tryCovertMapToConst(val *mapDefineVal) Val {
	dst := map[string]any{}
	for _, v := range val.kvs {
		//cst, ok := v.(*constraint)
		//if !ok {
		//	return val
		//}
		var ckk any
		ck, ok1 := v.k.(*constraint)
		if ok1 {
			ckk = ck.value
		}
		ck2, ok2 := v.k.(*variable)
		if ok2 {
			ckk = ck2.varName
		}
		if !ok1 && !ok2 {
			return val
		}

		vcv, ok := tryConvertToConst(v.v).(*constraint)
		if ok {
			dst[StringOf(ckk)] = vcv.value
		} else {
			return val
		}
	}
	return &constraint{
		value: dst,
	}
}

type ternaryVal struct {
	c Val
	l Val
	r Val
}

func (t *ternaryVal) Val(c *Context) any {
	if BoolCond(t.c.Val(c)) {
		return t.l.Val(c)
	}
	if t.r == nil {
		return nil
	}

	return t.r.Val(c)
}

func (t *ternaryVal) Set(c *Context, val any) {
}

type binaryValue struct {
	name string
	fun  func(ctx *Context, a, b Val) any
	l, r Val
}

func (b *binaryValue) Val(c *Context) any {
	return b.fun(c, b.l, b.r)
}

func (b *binaryValue) Set(c *Context, val any) {
}

func newBinaryValue(name string, l, r Val, f func(ctx *Context, a, b Val) any) *binaryValue {
	return &binaryValue{
		name: name,
		fun:  f,
		l:    l,
		r:    r,
	}
}

type unaryValue struct {
	name string
	fun  func(ctx *Context, a Val) any
	v    Val
}

func (u *unaryValue) Val(c *Context) any {
	return u.fun(c, u.v)
}

func (u *unaryValue) Set(c *Context, val any) {
}
func newUnaryValue(name string, v Val, f func(ctx *Context, a Val) any) *unaryValue {
	return &unaryValue{
		name: name,
		v:    v,
		fun:  f,
	}
}

type notNil struct {
	val Val
}

func nameOf(val Val) string {
	switch vv := val.(type) {
	case *variable:
		return vv.varName
	case *funcVariable:
		return vv.funcName + "()"
	case *objFuncVal:
		return vv.funcName + "()"
	case *accessVal:
		return nameOf(vv.left) + "." + nameOf(vv.right)
	case *constraint:
		return StringOf(vv.value)
	case *notNil:
		return nameOf(vv.val) + "!!"
	default:
		return reflect.TypeOf(val).String()
	}
}

func (n *notNil) Val(c *Context) any {
	v := n.val.Val(c)
	if v != nil {
		return v
	}
	return newErrorf("%v", nameOf(n.val)+" val is nil")
}

func (n *notNil) Set(c *Context, val any) {
	//TODO implement me
	panic("implement me")
}

type expList struct {
	Vals []Val
}

func (e *expList) Val(c *Context) any {
	//TODO implement me
	var v any
	for _, e := range e.Vals {
		v = e.Val(c)
		if convertToError(v) != nil {
			return v
		}
	}
	return v
}

func (e *expList) Set(c *Context, val any) {

}

type eqVal struct {
	L Val
	R Val
}

func (e *eqVal) Set(c *Context, val any) {
	//TODO implement me
	return
}

func (e *eqVal) Val(c *Context) any {
	return e.L.Val(c) == e.R.Val(c)
}

type eqValT struct {
	L Val
	R Val
}

func (e *eqValT) Set(c *Context, val any) {
	//TODO implement me
	return
}

func (e *eqValT) Val(c *Context) any {
	return eqt(e.L.Val(c), e.R.Val(c))
}
func eqt(a, b any) bool {

	switch av := a.(type) {
	case nil:
		switch bv := b.(type) {
		case nil:
			return true
		case string:
			return bv == ""
		case float64:
			return bv == 0
		case int:
			return bv == 0
		case bool:
			return bv == false
		default:
			return a == b
		}

	case string:
		switch bv := b.(type) {
		case nil:
			return av == ""
		case string:
			return av == bv
		case float64:
			if f, err := strconv.ParseFloat(av, 64); err == nil {
				return f == bv
			}
			return false
		case int:
			if f, err := strconv.ParseFloat(av, 64); err == nil {
				return f == float64(bv)
			}
			return false
		case bool:
			if v, err := strconv.ParseBool(av); err == nil {
				return v == bv
			}
			return false
		default:
			return a == b
		}

	case bool:
		switch bv := b.(type) {
		case nil:
			return av == false
		case bool:
			return av == bv
		case float64:
			if av {
				return bv == 1
			}
			return bv == 0
		case int:
			if av {
				return bv == 1
			}
			return bv == 0
		case string:
			if v, err := strconv.ParseBool(bv); err == nil {
				return av == v
			}
			return false
		default:
			return a == b
		}

	case float64:
		switch bv := b.(type) {
		case nil:
			return av == 0
		case float64:
			return av == bv
		case int:
			return av == float64(bv)
		case bool:
			if bv {
				return av == 1
			}
			return av == 0
		case string:
			if f, err := strconv.ParseFloat(bv, 64); err == nil {
				return av == f
			}
			return false
		default:
			return a == b
		}

	case int:
		switch bv := b.(type) {
		case nil:
			return av == 0
		case int:
			return av == bv
		case float64:
			return float64(av) == bv
		case bool:
			if bv {
				return av == 1
			}
			return av == 0
		case string:
			if f, err := strconv.ParseFloat(bv, 64); err == nil {
				return float64(av) == f
			}
			return false
		default:
			return a == b
		}

	default:
		// fallback：仅同类型直接比较
		return a == b
	}
}

type notEqValT struct {
	L Val
	R Val
}

func (e *notEqValT) Set(c *Context, val any) {
	//TODO implement me
	return
}

func (e *notEqValT) Val(c *Context) any {
	return !eqt(e.L.Val(c), e.R.Val(c))
}

// a==5 && b == 6 && call(a,b,c) //    and eq a 5 and b 6

type binaryCode struct {
	op  int
	val Val
}

func toPost(v Val) []binaryCode {
	switch vv := v.(type) {
	case *variable:
		return []binaryCode{{val: vv}}
	case *eqVal:
		ss := make([]binaryCode, 0)
		ss = append(ss, toPost(vv.L)...)
		ss = append(ss, toPost(vv.R)...)
		ss = append(ss, binaryCode{val: vv, op: '='})
		return ss
	case *constraint:
		return []binaryCode{{val: v}}
	case *binaryValue:
		ss := make([]binaryCode, 0)
		ss = append(ss, toPost(vv.l)...)
		ss = append(ss, toPost(vv.r)...)
		switch vv.name {
		case "&&":
			ss = append(ss, binaryCode{val: vv.l, op: '&'})
			return ss
		default:
			panic("invalid v")
		}
	}
	panic("invalid val")
}

func shouldCompileIf(v *accessVal) bool {
	switch rv := v.right.(type) {
	case *objFuncVal:
		switch rv.funcName {
		case "end":
			return true
		}
	}
	return false
}

func tryCompileVal(v *accessVal) Val {
	name := getTopFuncName(v)
	switch name {
	case "if":
		ifc := &ifctx{}
		if compileIF(ifc, v) {
			return ifc
		}
		return v
	case "switch":
		swc := &switchCtx{}
		if compileSwitch(swc, v) {
			return swc
		}
		return v
	default:
		return v
	}
}

func getTopFuncName(lv Val) string {
	switch lv := lv.(type) {
	case *accessVal:
		return getTopFuncName(lv.left)
	case *funcVariable:
		return lv.funcName
	default:
		return ""
	}
}

func compileSwitch(ctx *switchCtx, v *accessVal) bool {
	switch rv := v.right.(type) {
	case *objFuncVal:
		switch rv.funcName {
		case "case":
			switch len(rv.args) {
			case 1:
				ctx.cases = append([]elfs{{cond: rv.args[0]}}, ctx.cases...)
			case 2:
				ctx.cases = append([]elfs{{cond: rv.args[0], Do: rv.args[1]}}, ctx.cases...)
			}

		case "default":
			if len(rv.args) > 0 {
				ctx.def = rv.args[0]
			}
		case "end":

		default:
			return false
		}
	}
	switch lv := v.left.(type) {
	case *accessVal:
		if !compileSwitch(ctx, lv) {
			return false
		}
		return true
	case *funcVariable:
		if lv.funcName == "switch" {
			ctx.sw = lv.args[0]
			return true
		} else {
			return false
		}
	default:
		return false
	}

}

func compileIF(ctx *ifctx, v *accessVal) bool {
	switch rv := v.right.(type) {
	case *objFuncVal:
		switch rv.funcName {
		case "then":
			if len(rv.args) > 0 {
				ctx.then = rv.args[0]
			}

		case "else":
			if len(rv.args) > 0 {
				ctx.EL = rv.args[0]
			}

		case "elseif":
			switch len(rv.args) {
			case 1:
				ctx.Elfs = append([]elfs{{
					cond: rv.args[0],
					Do:   nil,
				}}, ctx.Elfs...)
			case 2:
				ctx.Elfs = append([]elfs{{
					cond: rv.args[0],
					Do:   rv.args[1],
				}}, ctx.Elfs...)
			}

		case "end":

		default:
			return false
		}
	}
	switch lv := v.left.(type) {
	case *accessVal:
		if !compileIF(ctx, lv) {
			return false
		}
		return true
	case *funcVariable:
		if lv.funcName == "if" {

			switch len(lv.args) {
			case 1:
				ctx.cond = lv.args[0]
			case 2:
				ctx.cond = lv.args[0]
				ctx.then = lv.args[1]
			}

			return true
		} else {
			return false
		}
	default:
		return false
	}
}
func parAsList(l Val, r Val) *expList {
	newx := &expList{}
	switch lv := l.(type) {
	case *expList:
		newx.Vals = append(newx.Vals, lv.Vals...)
	default:
		newx.Vals = append(newx.Vals, lv)
	}

	switch rv := r.(type) {
	case *expList:
		newx.Vals = append(newx.Vals, rv.Vals...)
	default:
		newx.Vals = append(newx.Vals, rv)
	}
	return newx
}

type VariadicVal struct {
	V Val
}

func (v *VariadicVal) ArrVal(ctx *Context) []any {
	e := v.V.Val(ctx)
	switch e := e.(type) {
	case []any:
		return e
	case ReadOnlyArray:
		return e
	}

	return []any{e}

}

func (v *VariadicVal) Val(c *Context) any {
	//TODO implement me
	//return v.v.Val(c)
	return newErrorf("variadicVar called unexpeced: %s", nameOf(v.V))
}

func (v *VariadicVal) Set(c *Context, val any) {
	//TODO implement me
	return
}

type addAddVal struct {
	v Val
}

func (v *addAddVal) Val(c *Context) any {

	n := NumberOf(v.v.Val(c))
	v.v.Set(c, n+1)
	return n + 1
}

func (v *addAddVal) Set(c *Context, val any) {

}

type subSubVal struct {
	v Val
}

func (v *subSubVal) Val(c *Context) any {

	n := NumberOf(v.v.Val(c))
	v.v.Set(c, n-1)
	return n - 1
}

func (v *subSubVal) Set(c *Context, val any) {

}

type asVal struct {
	key     string
	keyHash uint64
	val     Val
}

func (a *asVal) Val(c *Context) any {
	v := a.val.Val(c)
	c.Set(a.keyHash, a.key, v)
	return v
}

func (a *asVal) Set(c *Context, val any) {

}
