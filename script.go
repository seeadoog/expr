package expr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

/*
"if":"eq(name,5)"

*/

type Context struct {
	pctx                    context.Context
	table                   *envMap
	returnVal               []any
	IgnoreFuncNotFoundError bool
	ForceType               bool // if false will disable convert struct type to vm type and improve performance
	NewCallEnv              bool // if enabled , will use new env to call lambda which will cause extra performance cost
	Debug                   bool
	PanicWhenError          bool // if true, will panic when error occurs
	Env                     *Env
}

func (c *Context) Clone() *Context {
	ctx := &Context{
		//table:                   make(map[string]any),
		IgnoreFuncNotFoundError: c.IgnoreFuncNotFoundError,
		ForceType:               c.ForceType,
		NewCallEnv:              c.NewCallEnv,
		PanicWhenError:          c.PanicWhenError,
		pctx:                    c.pctx,
		Env:                     c.Env,
		//funcs:                   c.funcs,
	}
	ctx.table = c.table.clone()
	return ctx
}

func (e *Env) NewContext(table map[string]any) *Context {

	f := newEnvMap(8)
	if table != nil {
		for s, a := range table {
			f.putHash(calcHash(s), s, a)
		}
	}

	return &Context{
		Env:                     e,
		table:                   f,
		IgnoreFuncNotFoundError: false,
		ForceType:               false,
		NewCallEnv:              false,
		PanicWhenError:          true, // 默认为 true，保持向后兼容
	}
}

func (c *Context) SetHash(key HashKey, value any) {
	c.Set(key.Hash, key.Key, value)
}

func (c *Context) GetHash(key HashKey) any {
	return c.Get(key.Hash)
}

func (c *Context) Get(key uint64) interface{} {
	v := c.table.getHash(key)
	return v
}
func (c *Context) GetByString(key string) interface{} {
	v := c.table.getHash(calcHash(key))
	return v
}

func (c *Context) GetByJp(key string) any {
	v, err := c.Env.parseValueV(key)
	if err != nil {
		return nil
	}
	return v.Val(c)
}

func (c *Context) SetByJp(key string, val any) error {
	v, err := c.Env.parseValueV(key)
	if err != nil {
		return err
	}
	v.Set(c, val)
	return nil
}

func (c *Context) Set(key uint64, skey string, value interface{}) {
	c.table.putHashOnly(key, skey, value)
}
func (c *Context) SetByString(skey string, value interface{}) {
	c.table.putHash(calcHash(skey), skey, value)
}

func (c *Context) Delete(key string) {
	c.table.del(calcHash(key))
}

func (c *Context) Reset() {
	c.pctx = nil
	c.returnVal = nil
	c.table.reset()
}

func (c *Context) GetAsString(key HashKey) string {
	res, _ := c.GetHash(key).(string)
	return res
}

func (c *Context) GetAsArray(key HashKey) []any {
	switch val := c.GetHash(key).(type) {
	case []any:
		return val
	case ReadOnlyArray:
		return val
	default:
		return nil
	}
}

func (c *Context) GetAsNumber(key HashKey) float64 {
	return NumberOf(c.GetHash(key))
}

func (c *Context) GetAsBool(key HashKey) bool {
	return BoolOf(c.GetHash(key))
}

func (c *Context) GetAsObject(key HashKey) map[string]interface{} {
	switch val := c.GetHash(key).(type) {
	case map[string]interface{}:
		return val
	case ReadOnlyMap:
		return val
	}
	return nil
}

//func (c *Context) SetFunc(key string, fn ScriptFunc) {
//	if funtables[key] == nil {
//		if !strings.HasPrefix(key, "$") {
//			panic(fmt.Sprintf("func '%s' not registerd by RegisterDynamicFunc", key))
//		}
//	}
//	if c.funcs == nil {
//		c.funcs = make(map[string]ScriptFunc)
//	}
//	c.funcs[key] = fn
//}

func (c *Context) Debugf(f string, args ...any) {
	if !c.Debug {
		return
	}
	fmt.Printf(f+"\n", args...)
}

func (c *Context) SafeExec(e Expr) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	return c.Exec(e)
}

func (c *Context) Exec(e Expr) (err error) {
	err = e.Exec(c)
	if err != nil {
		if err == errBreak || err == errReturn {
			return nil
		}
		re, ok := err.(*Return)
		if ok {
			c.returnVal = re.Var
			return nil
		}
	}
	return err
}

func (c *Context) GetReturn() []any {
	return c.returnVal
}

func (c *Context) GetTable() map[string]any {
	dst := make(map[string]any)
	c.table.foreach(func(key uint64, hk string, val any) bool {
		dst[hk] = val
		return true
	})
	return dst
}

func (c *Context) Range(fn func(keyHash uint64, key string, val any) bool) {
	c.table.foreach(fn)
}

func (c *Context) Done() <-chan struct{} {
	if c.pctx == nil {
		return nil
	}
	return c.pctx.Done()
}

func (c *Context) Err() error {
	if c.pctx == nil {
		return nil
	}
	return c.pctx.Err()
}

func (c *Context) Value(key interface{}) interface{} {

	k, ok := key.(string)
	if ok {
		return c.GetByString(k)
	}
	if c.pctx == nil {
		return nil
	}
	return c.pctx.Value(key)
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c.pctx == nil {
		return deadline, false
	}
	return c.pctx.Deadline()
}

func (c *Context) SetContext(ctx context.Context) {
	c.pctx = ctx
}

func (c *Context) SafeExecValue(v Val) (res any, err any) {
	defer func() {
		if r := recover(); r != nil {
			err = r
		}
	}()
	res = v.Val(c)
	return
}

func (c *Context) ExecValue(v Val) (res any) {
	res = v.Val(c)
	return
}

func (c *Context) EncodeToStruct(structType reflect.Type, src any) (reflect.Value, bool) {
	return structValConvert(c, structType, src)
}

type setValue struct {
	key Val
	val Val
}

// a.b = 5   accs
func (s *setValue) Val(c *Context) any {
	//v := s.val.Val(c)
	//if s.jp == nil {
	//	c.Set(s.key, v)
	//} else {
	//	c.SetJP(s.jp, v)
	//}
	//return v
	v := s.val.Val(c)
	s.Set(c, v)
	return v
}

func setFor(c *Context, left Val, v any) {
	//switch vs := (left).(type) {
	//case *accessVal:
	//	vs.Set(c, v)
	//case *variable:
	//	vs.Set(c, v)
	//	return
	//case *arrAccessVal:
	//	vs.Set(c, v)
	//	//case *compiledVar:
	//	//	vs.Set(c, v)
	//	//c.stackSet(vs.index, v)
	//}
}

func (s *setValue) Set(c *Context, val any) {
	//setFor(c, s.key, val)
	s.key.Set(c, val)
}

type Expr = exp
type exp interface {
	Exec(c *Context) error
}

type Val interface {
	Val(c *Context) any
	parentValueSetter
}

type variable struct {
	varName string
	hash    uint64
	//varPath *jsonpath.Complied
}

func (v *variable) Val(c *Context) any {
	return c.Get(v.hash)
}

//	type stackVariable struct {
//		index int
//		name  string
//	}
//
//	func (s *stackVariable) Val(c *Context) any {
//		return c.stack[c.sp-s.index]
//	}
func (v *variable) Set(c *Context, val any) {
	//c.Set(v.varName, val)
	//c.table[v.varName] = val
	c.Set(v.hash, v.varName, val)
}

type constraint struct {
	value any
}

func (c *constraint) Val(ctx *Context) any {
	return c.value
}

func (c *constraint) Set(ctx *Context, val any) {}

type ScriptFunc func(ctx *Context, args ...Val) any
type funcVariable struct {
	funcName     string
	funcNameHash uint64
	fun          func(ctx *Context, args ...Val) any
	args         []Val
}

type FuncVariable = funcVariable

func (c *funcVariable) Val(ctx *Context) any {
	if c.fun == nil {
		//if ctx.funcs != nil {
		//	f := ctx.funcs[c.funcName]
		//	if f != nil {
		//		return f(ctx, c.args...)
		//	}
		//}
		lm, ok := ctx.Get(c.funcNameHash).(*lambda)
		if ok {
			return lambaCall(lm, ctx, c.args)
		}

		if ctx.IgnoreFuncNotFoundError {
			return nil
		}
		return newErrorf("function '%s' not found in table", c.funcName)
	}
	return c.fun(ctx, c.args...)
}

func (c *funcVariable) Set(ctx *Context, val any) {

}

type ifCond struct {
	cond     Val
	thenCond exp
	elseCond exp
}

func (i *ifCond) Exec(c *Context) error {
	if BoolCond(i.cond.Val(c)) {
		if i.thenCond != nil {
			return i.thenCond.Exec(c)
		}
	} else {
		if i.elseCond != nil {
			return i.elseCond.Exec(c)
		}
	}
	return nil
}

type valCond struct {
	val Val
}

func convertToError(o any) error {
	switch o := o.(type) {
	case *Return:
		return o
	case *Error:
		return o
	}
	return nil
}

func (v *valCond) Exec(c *Context) error {
	o := v.val.Val(c)
	return convertToError(o)
}

//func (s *setCond) Exec(c *Context) error {
//	if s.nameJPath == nil {
//		c.Set(s.,s.varName, s.val.Val(c))
//		return nil
//	}
//	return c.SetJP(s.nameJPath, s.val.Val(c))
//}

type callCond struct {
	fun *funcVariable
}

type Error struct {
	Err any
}

// newError 创建错误（不带 Context，使用默认行为：panic）
func newError(err any) *Error {
	e := &Error{Err: err}
	// 默认行为：panic（保持向后兼容）
	panic(e)
}

// newErrorWithCtx 创建错误（带 Context，根据 Context.PanicWhenError 决定是否 panic）
func newErrorWithCtx(err any, ctx *Context) *Error {
	e := &Error{Err: err}
	if ctx != nil && ctx.PanicWhenError {
		panic(e)
	}
	return e
}

func (e *Error) Error() string {
	return fmt.Sprint(e.Err)
}

type Break struct {
}

type breakVar struct {
}

var (
	_break = &Break{}
)

func (b *breakVar) Val(c *Context) any {
	//TODO implement me
	return _break
}
func (b *breakVar) Set(c *Context, val any) {
}

func (c *callCond) Exec(ctx *Context) error {
	o := c.fun.Val(ctx)
	return convertToError(o)
}

type forRange struct {
	target  Val
	keyName string
	keyHash uint64
	valName string
	valHash uint64
	do      exp
}

var (
	errBreak  = errors.New("break")
	errReturn = errors.New("return")
)

func (f *forRange) Exec(c *Context) error {
	v := f.target.Val(c)
	length := 0
	var valueOf func(i int) any
	switch v := v.(type) {
	case []any:
		length = len(v)
		valueOf = func(i int) any {
			return v[i]
		}
	case []string:
		length = len(v)
		valueOf = func(i int) any {
			return v[i]
		}
	case []float64:
		length = len(v)
		valueOf = func(i int) any {
			return v[i]
		}
	case map[string]interface{}:
		for i, a := range v {
			c.Set(f.keyHash, f.keyName, i)
			c.Set(f.valHash, f.valName, a)
			err := f.do.Exec(c)
			if err == errBreak {
				return nil
			}
			if err != nil {
				return err
			}
		}
		return nil
	case ReadOnlyMap:
		for i, a := range v {
			c.Set(f.keyHash, f.keyName, i)
			c.Set(f.valHash, f.valName, a)
			err := f.do.Exec(c)
			if err == errBreak {
				return nil
			}
			if err != nil {
				return err
			}
		}
	case ReadOnlyArray:
		length = len(v)
		valueOf = func(i int) any {
			return v[i]
		}
	}

	if valueOf != nil {
		for i := 0; i < length; i++ {
			c.Set(f.keyHash, f.keyName, i)
			c.Set(f.valHash, f.valName, valueOf(i))
			err := f.do.Exec(c)
			if err == errBreak {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type errExpr struct {
	err error
}

func (b *errExpr) Exec(c *Context) error {
	return b.err
}

type switchExpr struct {
	val   Val
	cases map[any]exp
	def   exp
}

func (s *switchExpr) Exec(c *Context) error {
	v := s.val.Val(c)
	expr := s.cases[v]
	if expr == nil {
		if s.def != nil {
			return s.def.Exec(c)
		}
		return nil
	}
	return expr.Exec(c)
}

var (
	setCondReg = regexp.MustCompile(`([$_\-0-9a-zA-Z.\[\]]+\s*)=([^=])\s*(.*)`)

	funcCallReg = regexp.MustCompile(`^[$._0-9a-zA-Z]+\s*\((.*)\)`)
)

func isSetCond(e string) bool {
	return setCondReg.MatchString(e)
}
func isFuncCall(e string) bool {
	return funcCallReg.MatchString(e)
}

type exps struct {
	exps []exp
}

func (e *exps) Exec(c *Context) error {
	for _, e2 := range e.exps {
		err := e2.Exec(c)
		if err != nil {
			return err
		}
	}
	return nil
}

type ExpParseFunc func(e *Env, o map[string]any, val any) (Expr, error)

var (
	expParserFactory = map[string]ExpParseFunc{}
)

func init() {
	expParserFactory["if"] = parseIf
	expParserFactory["for"] = parseForRange
	expParserFactory["switch"] = parseSwitch
	expParserFactory["define"] = parseDefine

}

func RegisterExp(name string, exp ExpParseFunc) {
	expParserFactory[name] = exp
}

var parseIf ExpParseFunc = func(e *Env, o map[string]any, val any) (exp, error) {
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("if val must be a string,but %v", val)
	}
	v, err := e.parseValueV(s)
	if err != nil {
		return nil, fmt.Errorf("parse if value error:%w %v", err, val)
	}

	iff := &ifCond{
		cond:     v,
		thenCond: nil,
		elseCond: nil,
	}
	thenRaw := o["then"]
	if thenRaw != nil {
		exp, err := e.ParseFromJSONObj(thenRaw)
		if err != nil {
			return nil, fmt.Errorf("parse if then error:%w %v", err, thenRaw)
		}
		iff.thenCond = exp
	}
	elseRaw := o["else"]
	if elseRaw != nil {
		exp, err := e.ParseFromJSONObj(elseRaw)
		if err != nil {
			return nil, fmt.Errorf("parse if else error:%w %v", err, elseRaw)
		}
		iff.elseCond = exp
	}

	return iff, nil
}

var (
	forRegexp = regexp.MustCompile(`^(\w+)\s*,\s*(\w+)\s*in\s*(.+)$`)
)

var parseForRange ExpParseFunc = func(ev *Env, o map[string]any, val any) (exp, error) {
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("for val must be a string,but %v", val)
	}

	if !forRegexp.MatchString(s) {
		return nil, fmt.Errorf("invalid for exp %v", s)
	}

	values := forRegexp.FindAllStringSubmatch(s, -1)
	if len(values) == 0 || len(values[0]) < 3 {
		return nil, fmt.Errorf("invalid for exp reg err %v", s)
	}
	vv, err := ev.parseValueV(values[0][3])
	if err != nil {
		return nil, fmt.Errorf("parse for value as val error:%w %v", err, s)
	}

	do, err := ev.ParseFromJSONObj(o["do"])
	if err != nil {
		return nil, fmt.Errorf("parse for::do exp error:%w %v", err, o)
	}
	e := &forRange{
		target:  vv,
		keyName: values[0][1],
		valName: values[0][2],
		keyHash: calcHash(values[0][1]),
		valHash: calcHash(values[0][2]),
		do:      do,
	}
	return e, nil
}

type casesExprs struct {
	val  Val
	expr exp
}

type switchCasesExpr struct {
	cases       []casesExprs
	defaultExpr Expr
}

func (s *switchCasesExpr) Exec(c *Context) error {
	for _, exprs := range s.cases {
		if BoolCond(exprs.val.Val(c)) {
			return exprs.expr.Exec(c)
		}
	}
	if s.defaultExpr != nil {
		return s.defaultExpr.Exec(c)
	}
	return nil
}

func parseSwitchExpr(e *Env, o map[string]any, val map[string]any) (exp, error) {
	switchCases := &switchCasesExpr{}
	for cases, exp := range val {
		val, err := e.parseValueV(cases)
		if err != nil {
			return nil, fmt.Errorf("parse switchcases cases error:%w %v", err, cases)
		}
		expr, err := e.ParseFromJSONObj(exp)
		if err != nil {
			return nil, fmt.Errorf("parse switchcases expr error:%w %v:%v", err, cases, exp)
		}
		switchCases.cases = append(switchCases.cases, casesExprs{
			val:  val,
			expr: expr,
		})
	}

	def := o["default"]
	if def != nil {
		defExpr, err := e.ParseFromJSONObj(def)
		if err != nil {
			return nil, fmt.Errorf("parse switchcases default expr error:%w %v", err, def)
		}
		switchCases.defaultExpr = defExpr
	}
	return switchCases, nil
}

var parseSwitch ExpParseFunc = func(e *Env, o map[string]any, val any) (exp, error) {
	str, ok := val.(string)
	if !ok {
		ov, ok := val.(map[string]any)
		if ok {
			return parseSwitchExpr(e, o, ov)
		}
		return nil, fmt.Errorf("switch val type must be a string or object,but %T", val)
	}

	vale, err := e.parseValueV(str)
	if err != nil {
		return nil, fmt.Errorf("parse switch value error:%w %v", err, str)
	}

	cases, ok := o["case"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("switch case type must be a map,but %v", o)
	}

	sw := &switchExpr{
		val:   vale,
		cases: map[any]exp{},
	}
	for val, exprs := range cases {
		vv, err := e.parseValueV(val)
		if err != nil {
			return nil, fmt.Errorf("parse switch case value error:%w %v", err, val)
		}
		//vvcst, ok := vv.(*constraint)
		//if !ok {
		//	return nil, fmt.Errorf("parse switch cases value is not constraint:%v", val)
		//}
		vexpr, err := e.ParseFromJSONObj(exprs)
		if err != nil {
			return nil, fmt.Errorf("parse switch expression error:%w %v: %v", err, val, exprs)
		}
		switch vv := vv.(type) {
		case *constraint:
			sw.cases[vv.value] = vexpr
		case *arrDefVal:
			for _, v := range vv.vs {
				vc, ok := v.(*constraint)
				if !ok {
					return nil, fmt.Errorf("switch case array type elem is not all const:%s", val)
				}
				sw.cases[vc.value] = vexpr
			}
		default:
			return nil, fmt.Errorf("parse switch cases value is not constraint:%v", val)
		}

	}

	def := o["default"]
	if def != nil {
		defExpr, err := e.ParseFromJSONObj(def)
		if err != nil {
			return nil, fmt.Errorf("parse switch default error:%w %v", err, def)
		}
		sw.def = defExpr
	}
	return sw, nil
}

func (e *Env) ParseFromJSONStr(str string) (Expr, error) {
	var o any
	if err := json.Unmarshal([]byte(str), &o); err != nil {
		return nil, fmt.Errorf("parse from json str error:%w", err)
	}
	return e.ParseFromJSONObj(o)
}

func (e *Env) ParseFromJSONObj(o any) (Expr, error) {
	switch o := o.(type) {
	case map[string]interface{}:
		for key, val := range o {
			fac := expParserFactory[key]
			if fac != nil {
				e, err := fac(e, o, val)
				if err != nil {
					return nil, err
				}
				return e, nil
			}
		}
		return nil, fmt.Errorf("not found exp field: %v", o)
	case []any:
		es := &exps{}
		for _, a := range o {
			e, err := e.ParseFromJSONObj(a)
			if err != nil {
				return nil, err
			}
			_, ok := e.(*noneExpr)
			if ok {
				continue
			}
			es.exps = append(es.exps, e)
		}
		return es, nil
	case string:
		e, err := e.parseExpr(strings.TrimSpace(o))
		if err != nil {
			return nil, err
		}
		return e, nil
	default:
		return nil, fmt.Errorf("invalid type %T", o)
	}
}

type noneExpr struct {
}

func (n *noneExpr) Exec(c *Context) error {
	return nil
}

func isJsonPath(s string) bool {
	return strings.Contains(s, ".") || strings.Contains(s, "[")
}
func (ev *Env) parseExpr(e string) (Expr, error) {
	if strings.HasPrefix(e, "#") {
		return &noneExpr{}, nil
	}
	switch {
	default:
		switch e {
		case "break":
			return &errExpr{err: errBreak}, nil
		case "return":
			return &errExpr{err: errReturn}, nil
		}
		v, err := ev.parseValueV(e)
		if err != nil {
			return nil, fmt.Errorf("parse stmt as value error:%w :: %s", err, e)
		}
		return &valCond{
			val: v,
		}, nil
		//return nil, fmt.Errorf("invalid exp:%s", e)
	}
}

type defineElem struct {
	hash HashKey
	val  any
}
type defineExpr struct {
	defs []defineElem
}

func (e *defineExpr) Exec(c *Context) error {
	for _, def := range e.defs {
		c.SetHash(def.hash, def.val)
	}
	return nil
}

var parseDefine = func(e *Env, o map[string]any, val any) (exp, error) {
	parentMap, ok := val.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid define val ,type should be map but %T", val)
	}
	defexp := &defineExpr{}
	for key, val := range parentMap {
		defexp.defs = append(defexp.defs, defineElem{
			hash: CalcHash(key),
			val:  val,
		})
	}
	return defexp, nil
}

type stack[T any] struct {
	data []T
}

func (s *stack[T]) push(data T) {
	s.data = append(s.data, data)
}
func (s *stack[T]) top() T {
	return s.data[len(s.data)-1]
}
func (s *stack[T]) pop() T {
	if len(s.data) == 0 {
		panic("stack is empty")
	}
	data := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return data
}
func (s *stack[T]) empty() bool {
	return len(s.data) == 0
}

type exprStack struct {
	ss  *stack[any]
	tkn []any
}

type iterator interface {
	getNext() (k any, val any, ok bool)
}

type dataIterator struct {
	start  float64
	end    float64
	offset float64
}

func (d *dataIterator) getNext() (k any, val any, ok bool) {
	hasnext := d.offset < d.end
	k = d.offset
	val = d.offset + d.start
	d.offset++
	return k, val, hasnext
}

type Return struct {
	Var []any
}

func (r *Return) Error() string {
	return fmt.Sprintf("return: %v", r.Var)
}

func ValueOfReturn(e error) []any {
	r, ok := e.(*Return)
	if ok {
		return r.Var
	}
	return nil
}

const ()

type Result struct {
	Err  any `json:"err"`
	Data any `json:"data"`
}
