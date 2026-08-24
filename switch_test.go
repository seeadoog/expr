package expr

import (
	"testing"
)

func TestSwitchCases(t *testing.T) {
	c := DefaultEnv.NewContext(nil)
	str := `
if a== 5 then c=1 else c=2 end;
switch a 
case 1: name = 1
case 2: name = 2
default: name = 3
end;
d = 5;
e = 6;
if d == 5 then
        if e == 6 then
            f = 1
        elseif e == 7 then
                f = 2
    end 
end ;
e = 7;
if d == 5 then
        if e == 6 then
            g = 1
        elseif e ==7 then
                g = 2
    end 
end ;

if ab  == nil then 
        switch true 
        case aa == nil:
                aac = 1
        case 2:
                aa = 1
        default:
                aa = 3
        end 
end;

e = 8;
if d == 5 then
        if e == 6 then
            h = 1
        elseif e ==7 then
                h = 2
    elseif e ==8 then
                h = 3
        elseif e == 9 then
    end ;
end ;
str = '';
for k,v in const ['1','2','3'] do
        str += string(k) + v 
end;
m2 = {};
for k,v in const {name:'1','age':'2'} do
        m2[k] = v 
end;

if aa ==3 then
        
else
        aacc = 13
end ;
sss = 1;
ccd = switch sss case 0: 'start' case 1: 'continue' case 2: 'end' default: 'unknown' end;
ccd2 = switch 3 case 0: 'start' case 1: 'continue' case 2: 'end' default: 'unknown' end;

ds1 = if true then 1 else 2 end;

func $add ($1,$2)
	$1+$2
end;

addv = $add(3,5);

$add2 = func($1,$2) $1 + $2 end;
addv2 = $add2(3,4);


swv = 30;
switch swv 
case 0~30:
case 30~60:
	sw1 = 1;
case 60~100:
	
end;

swv = 100;
switch swv 
case 0~30:
case 30~60:
	sw2 = 1;
case 60~100:
default:
     sw2 = 2
end;

app_id_rules = const {
	name: func()
       funcexec = 1
    end,
    "50043455": func()
	    funcexec = 2
	end,
};

call(app_id_rules["name"]);
call(func() callVal  = 3 end);
maps = {a:1,b:1};
for _,v in maps do 
	mapsv = v;
end;

for v in const [1,3,4,5] do
	brv = v ; 
	if v == 3 then 
		break
	end ;
end;
mulsw = 3;
switch mulsw
case 1,2,3:
	mulv = 1;
default:
	mulv = 2;
end;

switch mulsw
case 0~1, 2~4:
	mulv2 = 1;
case 5~9:
    mulv2 = 2;
default:
    mulv2 = 3;
end;

`
	parseAndExec(DefaultEnv, str, c)
	assertEqual(t, c, "c", 2.0)
	assertEqual(t, c, "name", 3.0)
	assertEqual(t, c, "f", 1.0)
	assertEqual(t, c, "g", 2.0)
	assertEqual(t, c, "aac", 1.0)
	assertEqual(t, c, "h", 3.0)
	assertEqual(t, c, "str", "011223")
	assertEqual(t, c, "aacc", 13.0)
	assertEqual(t, c, "ccd", "continue")
	assertEqual(t, c, "ccd2", "unknown")
	assertEqual(t, c, "ds1", 1.0)
	assertEqual(t, c, "addv", 8.0)
	assertEqual(t, c, "addv2", 7.0)
	assertEqual(t, c, "sw1", 1.0)
	assertEqual(t, c, "sw2", 2.0)
	assertEqual(t, c, "funcexec", 1.0)
	assertEqual(t, c, "callVal", 3.0)
	assertEqual(t, c, "mapsv", 1.0)
	assertEqual(t, c, "brv", 3.0)
	assertEqual(t, c, "mulv", 1.0)
	assertEqual(t, c, "mulv2", 1.0)
	assertDeepEqual(t, c, "m2", map[string]any{"name": "1", "age": "2"})

	//res := testing.Benchmark(func(b *testing.B) {
	//	v, err := DefaultEnv.ParseValue(str)
	//	if err != nil {
	//		b.Fatal(err)
	//	}
	//	b.ReportAllocs()
	//	for i := 0; i < b.N; i++ {
	//		c.ExecValue(v)
	//	}
	//})
	//fmt.Println(res)
	//
	//fmt.Println("allocs/op:", res.AllocsPerOp())
	//
	//fmt.Println("bytes/op:", res.AllocedBytesPerOp())
}

func BenchmarkIFF(b *testing.B) {
	c := DefaultEnv.NewContext(nil)
	b.ReportAllocs()
	v, err := DefaultEnv.ParseValue(`

app_id_rules = const {
name: func()
   a= 1
end,

"50043455": func()
	a= 2
end,

"0899s3234": func()
    a= 3
end

};

call(app_id_rules["name"])



`)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {

		c.ExecValue(v)

	}

}

func TestConstLambda(t *testing.T) {

	e, err := DefaultEnv.parseValueV(`
app_id_rules = const {
	n1: func()
       a= 1
    end,
    "n2": func()
	    a= 2
	end,
};

$hd = app_id_rules[name];
$hd();

app_id_rules.age = 1;


rule_table2 = const{
	s578902: {
		name: func()
				b = 1
			  end
	}
};

`)
	if err != nil {
		t.Fatal(err)
	}
	c := DefaultEnv.NewContext(nil)
	c.SetByString("name", "n1")
	c.ExecValue(e)
	assertEqual(t, c, "a", 1.0)
	c.SetByString("name", "n2")
	c.ExecValue(e)
	assertEqual(t, c, "a", 2.0)
	assertEqual(t, c, "app_id_rules.age", nil)
}
