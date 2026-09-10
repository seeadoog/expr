package interpreter

// OpCode 字节码操作码
type OpCode byte

const (
	// 栈操作
	OpPush      OpCode = iota // 压入常量到栈：A=常量索引
	OpPop                     // 弹出栈顶
	OpLoadVar                 // 从 Context 加载全局变量：A=名称索引
	OpStoreVar                // 存储到 Context 全局变量：A=名称索引（保留栈顶值）
	OpLoadFast                // 加载局部变量：A=局部索引
	OpStoreFast               // 存储局部变量：A=局部索引（保留栈顶值）

	// 算术运算
	OpAdd // +
	OpSub // -
	OpMul // *
	OpDiv // /
	OpMod // %
	OpPow // **
	OpNeg // 一元负号

	// 比较运算
	OpEq   // ==
	OpNeq  // !=
	OpLt   // <
	OpLte  // <=
	OpGt   // >
	OpGte  // >=
	OpEqT  // === 类型转换相等
	OpNeqT // !== 类型转换不等

	// 逻辑运算
	OpNot // !

	// 位运算
	OpBitAnd // &
	OpBitOr  // |
	OpBitXor // ^
	OpBitNot // ~
	OpShl    // <<
	OpShr    // >>

	// 对象/索引访问
	OpGetAttr  // obj.field 读：A=名称索引
	OpSetAttr  // obj.field = v：A=名称索引
	OpGetIndex // arr[i] 读
	OpSetIndex // arr[i] = v
	OpSlice    // arr[st:ed]

	// 函数调用
	OpCall       // 全局函数调用：A=调用信息索引
	OpCallMethod // 方法调用 obj.method()：A=调用信息索引

	// 控制流
	OpJump        // 无条件跳转：A=目标
	OpJumpIfTrue  // 栈顶为真跳转（消费栈顶）：A=目标
	OpJumpIfFalse // 栈顶为假跳转（消费栈顶）：A=目标

	// 复合字面量
	OpMakeArray // 创建数组：A=元素个数
	OpMakeMap   // 创建 map：A=键值对个数

	// 特殊
	OpNotNil // !! 非空断言
	OpHalt   // 停止执行，返回栈顶
)

var opNames = [...]string{
	OpPush:      "PUSH",
	OpPop:       "POP",
	OpLoadVar:   "LOAD_VAR",
	OpStoreVar:  "STORE_VAR",
	OpLoadFast:  "LOAD_FAST",
	OpStoreFast: "STORE_FAST",

	OpAdd: "ADD",
	OpSub: "SUB",
	OpMul: "MUL",
	OpDiv: "DIV",
	OpMod: "MOD",
	OpPow: "POW",
	OpNeg: "NEG",

	OpEq:   "EQ",
	OpNeq:  "NEQ",
	OpLt:   "LT",
	OpLte:  "LTE",
	OpGt:   "GT",
	OpGte:  "GTE",
	OpEqT:  "EQT",
	OpNeqT: "NEQT",

	OpNot: "NOT",

	OpBitAnd: "BIT_AND",
	OpBitOr:  "BIT_OR",
	OpBitXor: "BIT_XOR",
	OpBitNot: "BIT_NOT",
	OpShl:    "SHL",
	OpShr:    "SHR",

	OpGetAttr:  "GET_ATTR",
	OpSetAttr:  "SET_ATTR",
	OpGetIndex: "GET_INDEX",
	OpSetIndex: "SET_INDEX",
	OpSlice:    "SLICE",

	OpCall:       "CALL",
	OpCallMethod: "CALL_METHOD",

	OpJump:        "JUMP",
	OpJumpIfTrue:  "JUMP_IF_TRUE",
	OpJumpIfFalse: "JUMP_IF_FALSE",

	OpMakeArray: "MAKE_ARRAY",
	OpMakeMap:   "MAKE_MAP",

	OpNotNil: "NOT_NIL",
	OpHalt:   "HALT",
}

func (op OpCode) String() string {
	if int(op) < len(opNames) && opNames[op] != "" {
		return opNames[op]
	}
	return "UNKNOWN"
}
