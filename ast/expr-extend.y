// expr-extend.y
// 表达式语法解析器
// 支持算术运算、逻辑运算、变量访问、函数调用、流程控制等
%{
package ast

import (
    //"reflect"
)

%}

// Token 定义
%token IDENT NUMBER STRING BOOL NIL EQ AND OR NOTEQ GT GTE LT LTE ORR ACC IF ELSE FOR IN ACC2 CONST LAMB ADDEQ SUBEQ MULEQ DIVEQ VARIADIC AS

// 运算符优先级和结合性定义（从低到高）
%left IDENT EXIT
%left IF ELSE END DO ELSEIF THEN FOR SWITCH CASE DEFAULT FUNCTION ANNO BREAK
%left ';'
%left LAMB
%left AS
%right '='

%right '?'
%left ':'
%right ORR
%right OR
%right AND
%left EQ   NOTEQ  GT GTE LT LTE IN  EQT NOTEQT
%left ADDEQ SUBEQ MULEQ DIVEQ ADDADD SUBSUB
%left '+' '-'
%left '*' '/'
%left '%'
%left '&' '|'
%right '!'
%right '^'
%left '@' NONIL
%left ACC '[' ']'
%left VARIADIC
%right UMINUS
%right ACC2
%right CONST


%%

// 起始规则：设置语法树的根节点
Input:
	   Expr { yylex.(Setter).SetRoot($1.node); $$.node = $1.node }
	;

Expr:
	  // 算术运算符
	  Expr '+' Expr       { $$.node = &Binary{Op:"+", L: $1.node, R: $3.node} }
	| Expr '-' Expr       { $$.node = &Binary{Op:"-", L: $1.node, R: $3.node} }
	| Expr '*' Expr       { $$.node = &Binary{Op:"*", L: $1.node, R: $3.node} }
	| Expr '/' Expr       { $$.node = &Binary{Op:"/", L: $1.node, R: $3.node} }
	| Expr '^' Expr       { $$.node = &Binary{Op:"^", L: $1.node, R: $3.node} }
	| Expr '%' Expr       { $$.node = &Binary{Op:"%", L: $1.node, R: $3.node} }

	// 位运算符
	| Expr '&' Expr       { $$.node = &Binary{Op:"&", L: $1.node, R: $3.node} }
	| Expr '|' Expr       { $$.node = &Binary{Op:"|", L: $1.node, R: $3.node} }

	// 比较运算符
	| Expr  EQ Expr       { $$.node = &Binary{Op:"==", L: $1.node, R: $3.node} }
	| Expr  EQT Expr      { $$.node = &Binary{Op:"===", L: $1.node, R: $3.node} }
	| Expr  NOTEQT Expr   { $$.node = &Binary{Op:"!==", L: $1.node, R: $3.node} }
	| Expr  NOTEQ Expr    { $$.node = &Binary{Op:"!=", L: $1.node, R: $3.node} }
	| Expr  GT Expr       { $$.node = &Binary{Op:">", L: $1.node, R: $3.node} }
	| Expr  GTE Expr      { $$.node = &Binary{Op:">=", L: $1.node, R: $3.node} }
	| Expr  LT Expr       { $$.node = &Binary{Op:"<", L: $1.node, R: $3.node} }
	| Expr  LTE Expr      { $$.node = &Binary{Op:"<=", L: $1.node, R: $3.node} }

	// 逻辑运算符
	| Expr  AND Expr      { $$.node = &Binary{Op:"&&", L: $1.node, R: $3.node} }
	| Expr  OR Expr       { $$.node = &Binary{Op:"||", L: $1.node, R: $3.node} }

	// 顺序执行（分号分隔）
	| Expr  ';' Expr      { $$.node = &Binary{Op:";", L: $1.node, R: $3.node} }
	| Expr ';'            { $$.node  = $1.node }

	// 赋值运算符
	| Ident  '=' Expr     { $$.node = &Set{L: $1.node, R: $3.node} }
	| Var  '=' Expr       { $$.node = &Set{L: $1.node, R: $3.node} }
	| ArrIndex  '=' Expr  { $$.node = &Set{L: $1.node, R: $3.node} }

	// 复合赋值运算符（+=）
	| Ident  ADDEQ Expr     { $$.node = &Binary{Op:"+=", L: $1.node, R: $3.node} }
	| Var  ADDEQ Expr       { $$.node = &Binary{Op:"+=", L: $1.node, R: $3.node} }
	| ArrIndex  ADDEQ Expr  { $$.node = &Binary{Op:"+=", L: $1.node, R: $3.node} }

	// 复合赋值运算符（-=）
	| Ident  SUBEQ Expr     { $$.node = &Binary{Op:"-=", L: $1.node, R: $3.node} }
	| Var  SUBEQ Expr       { $$.node = &Binary{Op:"-=", L: $1.node, R: $3.node} }
	| ArrIndex  SUBEQ Expr  { $$.node = &Binary{Op:"-=", L: $1.node, R: $3.node} }

	// 复合赋值运算符（*=）
	| Ident  MULEQ Expr     { $$.node = &Binary{Op:"*=", L: $1.node, R: $3.node} }
	| Var  MULEQ Expr       { $$.node = &Binary{Op:"*=", L: $1.node, R: $3.node} }
	| ArrIndex  MULEQ Expr  { $$.node = &Binary{Op:"*=", L: $1.node, R: $3.node} }

	// 复合赋值运算符（/=）
	| Ident  DIVEQ Expr     { $$.node = &Binary{Op:"/=", L: $1.node, R: $3.node} }
	| Var  DIVEQ Expr       { $$.node = &Binary{Op:"/=", L: $1.node, R: $3.node} }
	| ArrIndex  DIVEQ Expr  { $$.node = &Binary{Op:"/=", L: $1.node, R: $3.node} }

	// 其他二元运算符
	| Expr  ORR Expr      { $$.node = &Binary{Op: "orr", L: $1.node, R: $3.node } }
	| Expr  AS Ident      { $$.node = &Binary{Op: "as", L: $1.node, R: $3.node } }
	| Expr IN Expr        { $$.node = &Binary{Op: "in", L: $1.node, R: $3.node } }

	// 一元运算符
	| '!' Expr            { $$.node = &Unary{Op:"!", X: $2.node}  }
	| Expr ADDADD         { $$.node = &Unary{Op:"++", X: $1.node}  }
	| Expr SUBSUB         { $$.node = &Unary{Op:"--", X: $1.node}  }
	| Expr VARIADIC       { $$.node = &Unary{Op:"...", X: $1.node}  }
	| '-' Expr  %prec UMINUS { $$.node = &Unary{Op:"-", X: $2.node} }

	// 三元运算符（条件表达式）
	| Expr '?' Expr ':' Expr { $$.node = &Ternary{C:$1.node, L:$3.node, R:$5.node} }

	// 非空断言
	| Expr NONIL          { $$.node = &NotNil{$1.node} }
	| Expr '@'            { $$.node = &NotNil{$1.node} }

	// Lambda 表达式：{x, y} => expr
	| '{' Ids '}' LAMB  Expr  { $$.node = &Lambda{L: $2.strs, R:$5.node } }

	// Lambda 表达式：expr => { body }
	| Expr LAMB '{' Expr '}'  { $$.node = &Lambda2{L: $1.node, R:$4.node } }

	// for 循环：for k, v in expr do ... end（带键值）
	| FOR IDENT ',' IDENT IN Expr  DO Expr END { $$.node = &ForRange{KName: $2.str, VName:$4.str, Var:$6.node, Do:$8.node } }

	// for 循环：for v in expr do ... end（仅值）
	| FOR IDENT IN Expr  DO Expr END { $$.node = &ForRange{VName:$2.str, Var:$4.node, Do:$6.node } }

	// if-then-end（无 else）
	| IF Expr THEN EExpr END   { $$.node = &IfElse{Cond: $2.node, Then: $4.node } }

	// if-then-else-end
	| IF Expr THEN EExpr ELSE EExpr END   { $$.node = &IfElse{Cond: $2.node, Then: $4.node, Else: $6.node } }

	// if-then-elseif...-end
	| IF Expr THEN EExpr Elseifs END { $$.node = &IfElse{Cond: $2.node, Then: $4.node, Elseifs: $5.nodes } }

	// if-then-elseif...-else-end
	| IF Expr THEN EExpr Elseifs ELSE EExpr END { $$.node = &IfElse{Cond: $2.node, Then: $4.node, Elseifs: $5.nodes, Else: $7.node} }

	// switch-case 语句（无 default）
	| SWITCH Expr Cases END { $$.node = &Switch{Var: $2.node, Cases: $3.nodes} }

	// switch-case 语句（带 default）
	| SWITCH Expr Cases DEFAULT ':' EExpr END { $$.node = &Switch{Var: $2.node, Cases: $3.nodes, Default: $6.node} }

	// break 语句
	| BREAK { $$.node = &Break{} }

	// Lambda 简写：x => expr（单参数）
	| IDENT LAMB Expr { $$.node = &Lambda{L:[]string{$1.str}, R:$3.node } }

	// 常量表达式
	| CONST Expr { $$.node = &Const{L: $2.node} }

	// 标识符和基本表达式
	| Ident               { $$.node = $1.node }
	| Primary             { $$.node = $1.node }

	// 函数定义：function name(params) body end
	| FUNCTION Ident '(' Idss ')' EExpr END  { $$.node = &Set{L:$2.node, R:&Lambda{L: $4.strs, R:$6.node }} }

	// 匿名函数：function(params) body end
	| FUNCTION '(' Idss ')' EExpr END { $$.node = &Lambda{L: $3.strs, R:$5.node } }

	// 范围表达式：1~10
	| NUMBER '~' NUMBER { $$.node = &RangeValue{L: $1.num, R: $3.num} }

	;

// elseif 子句列表
Elseifs:
    elseif  { $$.nodes = []Node{$1.node} }
    | Elseifs elseif { $$.nodes = append($1.nodes, $2.node) }

// 单个 elseif 子句
elseif:
    ELSEIF Expr THEN EExpr  { $$.node = &ElseIf{Cond: $2.node, Then: $4.node} }

// case 子句列表
Cases:
    Case { $$.nodes = []Node{$1.node} }
    | Cases Case { $$.nodes = append($1.nodes, $2.node) }

// 单个 case 子句
Case:
    CASE ArgList ':' EExpr { $$.node = &Case{Vars: $2.nodes, Do: $4.node} }

// 可选表达式（可能为空）
EExpr:
    { $$.node = nil }
    | Expr { $$.node = $1.node }


// 标识符列表（至少一个）
Ids:
    Ident { $$.strs = []string{$1.str} }
    | Ids ',' Ident { $$.strs = append($1.strs, $3.str) }

// 标识符列表（可为空）
Idss:
       { $$.strs = []string{} }
    | Ident { $$.strs = []string{$1.str} }
    | Idss ',' Ident { $$.strs = append($1.strs, $3.str) }

// 变量访问：expr.ident 或 expr.[expr]
Var:
    Expr ACC Ident  %prec ACC  { $$.node = &Access{L: $1.node, R:$3.node} }
    | Expr ACC ArrIndex  %prec ACC  { $$.node = &Access{L: $1.node, R:$3.node} }

// 属性访问
Acc:
    Expr ACC Expr %prec ACC { $$.node = &Access{L: $1.node, R:$3.node} }

// 标识符
Ident:
    IDENT { $$.node = &Variable{Name: $1.str}; $$.str = $1.str }

// 基本表达式
Primary:
	  // 字面量
	  NUMBER         { $$.node = &Number{Val: $1.num} }
	| BOOL           { $$.node = &Bool{Val: $1.boolean} }
	| STRING         { $$.node = &String{Val: $1.str}; $$.str = $1.str }
	| NIL            { $$.node = &Nil{} }

	// 函数调用
	| Ident '(' ArgListOpt ')' { $$.node = &Call{Name: $1.str, Args: $3.nodes} }

	// 括号表达式
	| '(' Expr ')'   { $$.node = $2.node }

	// Map 定义：{key: value, ...}
	| '{' KvsOpt '}' { $$.node = &MapSet{Kvs: $2.kvs}  }

	// 数组定义：[elem1, elem2, ...]
	| '[' ArgListOpt ']' { $$.node = &ArrDef{V:$2.nodes}  }

	// 数组索引访问
	| ArrIndex { $$.node = $1.node }

	// 属性访问
	| Acc { $$.node = $1.node }

	// 切片操作：expr[start:end]
	| Expr '[' SliceN ':' SliceN ']' { $$.node = &SliceCut{V: $1.node, St: $3.node, Ed:$5.node} }

	;




// 切片边界（可为空）
SliceN:
    /*empty*/ { $$.node = nil }
    | Expr { $$.node = $1.node }

// 数组索引：expr[index]
ArrIndex:
    Expr '[' Expr ']' { $$.node = &ArrAccess{L:$1.node, R:$3.node} }

// 可选键值对列表
KvsOpt:
    /*empty*/ { $$.kvs = nil  }
    | Kvs { $$.kvs = $1.kvs  }
    ;

// 键值对列表
Kvs:
    Kv  { $$.kvs = []KV{$1.kv} }
    | Kvs ',' Kv { $$.kvs = append($1.kvs, $3.kv) }
    | Kvs ','  { $$.kvs = $1.kvs }
    ;

// 键值对
Kv:
    STRING ':' Expr { $$.kv = KV{ K:$1.str, V: $3.node} }
    | IDENT ':' Expr { $$.kv = KV{ K:$1.str, V: $3.node} }
    ;

// 可选参数列表
ArgListOpt:
	  /* empty */         { $$.nodes = nil }
	| ArgList             { $$.nodes = $1.nodes }
	;

// 参数列表（至少一个）
ArgList:
	  Expr                { $$.nodes = []Node{ $1.node } }
	| ArgList ',' Expr    { $$.nodes = append($1.nodes, $3.node)  }
	| ArgList ','         { $$.nodes = $1.nodes  }
	;

%%
