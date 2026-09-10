package interpreter

import (
	"fmt"
	"strings"
)

// Instruction 字节码指令
// 采用双操作数设计：A 通常是主索引（常量/名称/跳转目标），B 是辅助数据（如参数个数）
type Instruction struct {
	Op OpCode // 操作码
	A  uint32 // 主操作数
	B  uint32 // 辅助操作数
}

// CallInfo 函数调用信息（存于独立表，避免指令内位打包）
type CallInfo struct {
	Name string // 函数名
	Hash uint64 // 函数名哈希
	Argc int    // 参数个数
}

// NameRef 名称引用（变量名/属性名等，携带预算好的 hash）
type NameRef struct {
	Name string
	Hash uint64
}

// ByteCode 编译后的字节码
type ByteCode struct {
	Instructions []Instruction // 指令序列
	Constants    []any         // 常量池
	Names        []NameRef     // 名称表（变量/属性名，含 hash）
	Calls        []CallInfo    // 函数调用信息表
	MaxLocals    int           // 最大局部变量数量

	// 编译期去重缓存（运行期不使用）
	constIndex map[any]uint32
	nameIndex  map[string]uint32
}

// NewByteCode 创建新的字节码对象
func NewByteCode() *ByteCode {
	return &ByteCode{
		Instructions: make([]Instruction, 0, 64),
		Constants:    make([]any, 0, 32),
		Names:        make([]NameRef, 0, 16),
		Calls:        make([]CallInfo, 0, 8),
		constIndex:   make(map[any]uint32, 32),
		nameIndex:    make(map[string]uint32, 16),
	}
}

// Emit 发射一条指令，返回指令位置
func (bc *ByteCode) Emit(op OpCode, a uint32) int {
	bc.Instructions = append(bc.Instructions, Instruction{Op: op, A: a})
	return len(bc.Instructions) - 1
}

// EmitAB 发射一条带两个操作数的指令
func (bc *ByteCode) EmitAB(op OpCode, a, b uint32) int {
	bc.Instructions = append(bc.Instructions, Instruction{Op: op, A: a, B: b})
	return len(bc.Instructions) - 1
}

// AddConstant 添加常量到常量池，返回索引（可去重的标量类型会去重）
func (bc *ByteCode) AddConstant(val any) uint32 {
	// 仅对可比较且常见的标量做去重，避免不可比较类型（slice/map）导致 panic
	switch val.(type) {
	case float64, int, string, bool, nil, uint64:
		if idx, ok := bc.constIndex[val]; ok {
			return idx
		}
		idx := uint32(len(bc.Constants))
		bc.Constants = append(bc.Constants, val)
		bc.constIndex[val] = idx
		return idx
	default:
		idx := uint32(len(bc.Constants))
		bc.Constants = append(bc.Constants, val)
		return idx
	}
}

// AddName 添加名称到名称表，返回索引
func (bc *ByteCode) AddName(name string, hash uint64) uint32 {
	if idx, ok := bc.nameIndex[name]; ok {
		return idx
	}
	idx := uint32(len(bc.Names))
	bc.Names = append(bc.Names, NameRef{Name: name, Hash: hash})
	bc.nameIndex[name] = idx
	return idx
}

// AddCall 添加函数调用信息，返回索引
func (bc *ByteCode) AddCall(name string, hash uint64, argc int) uint32 {
	idx := uint32(len(bc.Calls))
	bc.Calls = append(bc.Calls, CallInfo{Name: name, Hash: hash, Argc: argc})
	return idx
}

// PatchJump 修补跳转指令的目标地址
func (bc *ByteCode) PatchJump(offset int, target int) {
	bc.Instructions[offset].A = uint32(target)
}

// CurrentOffset 返回当前指令偏移（即下一条指令的位置）
func (bc *ByteCode) CurrentOffset() int {
	return len(bc.Instructions)
}

// Disassemble 反汇编字节码（用于调试）
func (bc *ByteCode) Disassemble() string {
	var sb strings.Builder
	sb.WriteString("=== Bytecode Disassembly ===\n")

	sb.WriteString("Constants:\n")
	for i, c := range bc.Constants {
		fmt.Fprintf(&sb, "  [%d] %#v\n", i, c)
	}

	sb.WriteString("Names:\n")
	for i, n := range bc.Names {
		fmt.Fprintf(&sb, "  [%d] %s (hash=%d)\n", i, n.Name, n.Hash)
	}

	if len(bc.Calls) > 0 {
		sb.WriteString("Calls:\n")
		for i, call := range bc.Calls {
			fmt.Fprintf(&sb, "  [%d] %s argc=%d\n", i, call.Name, call.Argc)
		}
	}

	fmt.Fprintf(&sb, "MaxLocals: %d\n", bc.MaxLocals)

	sb.WriteString("Instructions:\n")
	for i, inst := range bc.Instructions {
		fmt.Fprintf(&sb, "  %04d  %-14s  A=%-4d B=%-4d", i, inst.Op.String(), inst.A, inst.B)
		// 补充可读注释
		switch inst.Op {
		case OpPush:
			if int(inst.A) < len(bc.Constants) {
				fmt.Fprintf(&sb, "  ; %#v", bc.Constants[inst.A])
			}
		case OpLoadVar, OpStoreVar, OpGetAttr, OpSetAttr:
			if int(inst.A) < len(bc.Names) {
				fmt.Fprintf(&sb, "  ; %s", bc.Names[inst.A].Name)
			}
		case OpCall, OpCallMethod:
			if int(inst.A) < len(bc.Calls) {
				fmt.Fprintf(&sb, "  ; %s()", bc.Calls[inst.A].Name)
			}
		case OpJump, OpJumpIfTrue, OpJumpIfFalse:
			fmt.Fprintf(&sb, "  ; -> %d", inst.A)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
