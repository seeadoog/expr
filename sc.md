阅读expr 目录的代码，这是一个规则引擎，给我整理成一个文档，文档分为两个部分
- 开发者文档： 详细列出给开发者暴露的 api ，调用方式


- 内置函数介绍 介绍规则引擎支持的内置函数
介绍内置函数，分为两个部分：函数介绍要有入参数，出参数介绍。和使用说明，
1 全局函数：
例如 time_now()   : 获取当前时间 返回 time.Time 类型

2 对象函数，介绍规则引擎支持的对象类型包含的子函数。例如： 

map[string]interface{}::len() float64  : 返回map 类型的长度

time.Time::unix()float64  : 返回time类型的时间戳，精确到秒

string::split(num float64)[]any  返回string 切割后的数组。
        

你可以补充下你觉得需要写到文档中的部分，将文档写入expr_help.md

/*
req.a == 'xxx' && req.b = 'xxx'? req.bcd = 'x5c90':_ ;
req.b == 'xxx56'? req.def = '${req.a}/${req.b or 1}';

server.name = 'sxf' ;
server.age = 'sxf'  ;
log_dir = 'xxx'     ;
log.1.file = log_dir + '/xxxx'  ;
log.2.file = log_dir + '/xxxx'  ;


*/


优化下 Env.ParseValueFromNode 函数，改成调用 Env.parseValueFromNode，将里面的不同类型的解析函数拆分，用 RegisterParseFunc 注册对应的handler,
拆分的handler 写入 parser_handler.go. 每一个handler 写到一个单独函数 。最后在 init 函数中RegisterParseFunc 注册这些函数，这些函数不用导出，小写开头即可
