package interpreter

import (
	"fmt"

	"github.com/seeadoog/expr/ast"
)

// Compiler AST 到字节码的编译器
type Compiler struct {
	bc         *ByteCode
	localVars  map[string]int // 局部变量名 -> 局部槽索引
	localCount int            // 已分配的局部变量数量
	hashFunc   func(string) uint64
}

// NewCompiler 创建编译器。hashFunc 用于与宿主 Env 保持一致的哈希算法，
// 若为 nil 则使用内置默认哈希（此时无法与宿主变量表对接，仅供独立测试）。
func NewCompiler(hashFunc func(string) uint64) *Compiler {
	if hashFunc == nil {
		hashFunc = defaultHash
	}
	return &Compiler{
		bc:        NewByteCode(),
		localVars: make(map[string]int),
		hashFunc:  hashFunc,
	}
}

// defaultHash 仅用于无宿主的独立测试场景
func defaultHash(s string) uint64 {
	h := uint64(1469598103934665603) // FNV-1a offset basis
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

// Compile 编译 AST 节点，返回可执行字节码
func (c *Compiler) Compile(node ast.Node) (*ByteCode, error) {
	if err := c.compileNode(node); err != nil {
		return nil, err
	}
	c.bc.Emit(OpHalt, 0)
	c.bc.MaxLocals = c.localCount
	return c.bc, nil
}

func (c *Compiler) compileNode(node ast.Node) error {
	if node == nil {
		c.bc.Emit(OpPush, c.bc.AddConstant(nil))
		return nil
	}

	switch n := node.(type) {
	case *ast.Number:
		c.bc.Emit(OpPush, c.bc.AddConstant(n.Val))
	case *ast.String:
		c.bc.Emit(OpPush, c.bc.AddConstant(n.Val))
	case *ast.Bool:
		c.bc.Emit(OpPush, c.bc.AddConstant(n.Val))
	case *ast.Nil:
		c.bc.Emit(OpPush, c.bc.AddConstant(nil))
	case *ast.Variable:
		return c.compileVariable(n)
	case *ast.Binary:
		return c.compileBinary(n)
	case *ast.Unary:
		return c.compileUnary(n)
	case *ast.Call:
		return c.compileCall(n)
	case *ast.Access:
		return c.compileAccess(n)
	case *ast.ArrAccess:
		return c.compileArrAccess(n)
	case *ast.Set:
		return c.compileSet(n)
	case *ast.Ternary:
		return c.compileTernary(n)
	case *ast.ArrDef:
		return c.compileArrDef(n)
	case *ast.MapSet:
		return c.compileMapSet(n)
	case *ast.SliceCut:
		return c.compileSliceCut(n)
	case *ast.NotNil:
		return c.compileNotNil(n)
	case *ast.NodeList:
		return c.compileNodeList(n)
	default:
		return fmt.Errorf("interpreter: unsupported node type %T (%s)", node, node.String())
	}
	return nil
}

func (c *Compiler) compileVariable(n *ast.Variable) error {
	// 优先局部变量
	if localIdx, ok := c.localVars[n.Name]; ok {
		c.bc.Emit(OpLoadFast, uint32(localIdx))
		return nil
	}
	// 全局变量：名称表携带 hash，运行期直接用 hash 查 Context
	nameIdx := c.bc.AddName(n.Name, c.hashFunc(n.Name))
	c.bc.Emit(OpLoadVar, nameIdx)
	return nil
}

func (c *Compiler) compileBinary(n *ast.Binary) error {
	switch n.Op {
	case "&&":
		// 短路：左假则整体为左值，跳过右侧
		if err := c.compileNode(n.L); err != nil {
			return err
		}
		jumpIfFalse := c.bc.Emit(OpJumpIfFalse, 0)
		if err := c.compileNode(n.R); err != nil {
			return err
		}
		// 右侧计算完，跳到结尾
		jumpEnd := c.bc.Emit(OpJump, 0)
		// 左侧为假的落点：压入 false 作为结果
		c.bc.PatchJump(jumpIfFalse, c.bc.CurrentOffset())
		c.bc.Emit(OpPush, c.bc.AddConstant(false))
		c.bc.PatchJump(jumpEnd, c.bc.CurrentOffset())
		return nil
	case "||":
		// 短路：左真则整体为 true，跳过右侧
		if err := c.compileNode(n.L); err != nil {
			return err
		}
		jumpIfTrue := c.bc.Emit(OpJumpIfTrue, 0)
		if err := c.compileNode(n.R); err != nil {
			return err
		}
		jumpEnd := c.bc.Emit(OpJump, 0)
		c.bc.PatchJump(jumpIfTrue, c.bc.CurrentOffset())
		c.bc.Emit(OpPush, c.bc.AddConstant(true))
		c.bc.PatchJump(jumpEnd, c.bc.CurrentOffset())
		return nil
	}

	// 普通二元运算
	if err := c.compileNode(n.L); err != nil {
		return err
	}
	if err := c.compileNode(n.R); err != nil {
		return err
	}

	op, ok := binaryOps[n.Op]
	if !ok {
		return fmt.Errorf("interpreter: unsupported binary operator %q", n.Op)
	}
	c.bc.Emit(op, 0)
	return nil
}

var binaryOps = map[string]OpCode{
	"+":   OpAdd,
	"-":   OpSub,
	"*":   OpMul,
	"/":   OpDiv,
	"%":   OpMod,
	"**":  OpPow,
	"==":  OpEq,
	"!=":  OpNeq,
	"<":   OpLt,
	"<=":  OpLte,
	">":   OpGt,
	">=":  OpGte,
	"===": OpEqT,
	"!==": OpNeqT,
	"&":   OpBitAnd,
	"|":   OpBitOr,
	"^":   OpBitXor,
	"<<":  OpShl,
	">>":  OpShr,
}

func (c *Compiler) compileUnary(n *ast.Unary) error {
	if err := c.compileNode(n.X); err != nil {
		return err
	}
	switch n.Op {
	case "-":
		c.bc.Emit(OpNeg, 0)
	case "!":
		c.bc.Emit(OpNot, 0)
	case "~":
		c.bc.Emit(OpBitNot, 0)
	default:
		return fmt.Errorf("interpreter: unsupported unary operator %q", n.Op)
	}
	return nil
}

func (c *Compiler) compileCall(n *ast.Call) error {
	for _, arg := range n.Args {
		if err := c.compileNode(arg); err != nil {
			return err
		}
	}
	callIdx := c.bc.AddCall(n.Name, c.hashFunc(n.Name), len(n.Args))
	c.bc.Emit(OpCall, callIdx)
	return nil
}

func (c *Compiler) compileAccess(n *ast.Access) error {
	if err := c.compileNode(n.L); err != nil {
		return err
	}
	switch r := n.R.(type) {
	case *ast.Variable:
		nameIdx := c.bc.AddName(r.Name, c.hashFunc(r.Name))
		c.bc.Emit(OpGetAttr, nameIdx)
	case *ast.Call:
		for _, arg := range r.Args {
			if err := c.compileNode(arg); err != nil {
				return err
			}
		}
		callIdx := c.bc.AddCall(r.Name, c.hashFunc(r.Name), len(r.Args))
		c.bc.Emit(OpCallMethod, callIdx)
	default:
		return fmt.Errorf("interpreter: unsupported access right side %T", n.R)
	}
	return nil
}

func (c *Compiler) compileArrAccess(n *ast.ArrAccess) error {
	if err := c.compileNode(n.L); err != nil {
		return err
	}
	if err := c.compileNode(n.R); err != nil {
		return err
	}
	c.bc.Emit(OpGetIndex, 0)
	return nil
}

func (c *Compiler) compileSet(n *ast.Set) error {
	// 先求右值，压栈
	if err := c.compileNode(n.R); err != nil {
		return err
	}

	switch l := n.L.(type) {
	case *ast.Variable:
		if localIdx, ok := c.localVars[l.Name]; ok {
			c.bc.Emit(OpStoreFast, uint32(localIdx))
		} else if n.Const {
			idx := c.localCount
			c.localVars[l.Name] = idx
			c.localCount++
			c.bc.Emit(OpStoreFast, uint32(idx))
		} else {
			nameIdx := c.bc.AddName(l.Name, c.hashFunc(l.Name))
			c.bc.Emit(OpStoreVar, nameIdx)
		}
	case *ast.Access:
		// obj.field = val：栈序需为 [val, obj]，SET_ATTR 消费两者
		if v, ok := l.R.(*ast.Variable); ok {
			if err := c.compileNode(l.L); err != nil {
				return err
			}
			nameIdx := c.bc.AddName(v.Name, c.hashFunc(v.Name))
			c.bc.Emit(OpSetAttr, nameIdx)
		} else {
			return fmt.Errorf("interpreter: unsupported access assignment target %T", l.R)
		}
	case *ast.ArrAccess:
		// arr[idx] = val：栈序 [val, obj, idx]
		if err := c.compileNode(l.L); err != nil {
			return err
		}
		if err := c.compileNode(l.R); err != nil {
			return err
		}
		c.bc.Emit(OpSetIndex, 0)
	default:
		return fmt.Errorf("interpreter: unsupported assignment target %T", n.L)
	}
	return nil
}

func (c *Compiler) compileTernary(n *ast.Ternary) error {
	if err := c.compileNode(n.C); err != nil {
		return err
	}
	jumpIfFalse := c.bc.Emit(OpJumpIfFalse, 0)
	if err := c.compileNode(n.L); err != nil {
		return err
	}
	jumpEnd := c.bc.Emit(OpJump, 0)
	c.bc.PatchJump(jumpIfFalse, c.bc.CurrentOffset())
	if n.R != nil {
		if err := c.compileNode(n.R); err != nil {
			return err
		}
	} else {
		c.bc.Emit(OpPush, c.bc.AddConstant(nil))
	}
	c.bc.PatchJump(jumpEnd, c.bc.CurrentOffset())
	return nil
}

func (c *Compiler) compileArrDef(n *ast.ArrDef) error {
	for _, v := range n.V {
		if err := c.compileNode(v); err != nil {
			return err
		}
	}
	c.bc.Emit(OpMakeArray, uint32(len(n.V)))
	return nil
}

func (c *Compiler) compileMapSet(n *ast.MapSet) error {
	for _, kv := range n.Kvs {
		c.bc.Emit(OpPush, c.bc.AddConstant(kv.K))
		if err := c.compileNode(kv.V); err != nil {
			return err
		}
	}
	c.bc.Emit(OpMakeMap, uint32(len(n.Kvs)))
	return nil
}

func (c *Compiler) compileSliceCut(n *ast.SliceCut) error {
	if err := c.compileNode(n.V); err != nil {
		return err
	}
	if n.St != nil {
		if err := c.compileNode(n.St); err != nil {
			return err
		}
	} else {
		c.bc.Emit(OpPush, c.bc.AddConstant(nil))
	}
	if n.Ed != nil {
		if err := c.compileNode(n.Ed); err != nil {
			return err
		}
	} else {
		c.bc.Emit(OpPush, c.bc.AddConstant(nil))
	}
	c.bc.Emit(OpSlice, 0)
	return nil
}

func (c *Compiler) compileNotNil(n *ast.NotNil) error {
	if err := c.compileNode(n.N); err != nil {
		return err
	}
	c.bc.Emit(OpNotNil, 0)
	return nil
}

func (c *Compiler) compileNodeList(n *ast.NodeList) error {
	for i, node := range n.Ns {
		if err := c.compileNode(node); err != nil {
			return err
		}
		if i < len(n.Ns)-1 {
			c.bc.Emit(OpPop, 0) // 丢弃非末位表达式的结果
		}
	}
	return nil
}
