package expr

import "testing"

func TestSwitchCases(t *testing.T) {
	c := DefaultEnv.NewContext(nil)

	parseAndExec(DefaultEnv, `
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
for k,v in ['1','2','3'] do
        str += string(k) + v 
end;
m2 = {};
for k,v in {name:'1','age':'2'} do
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

`, c)
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
	assertDeepEqual(t, c, "m2", map[string]any{"name": "1", "age": "2"})

}

func BenchmarkIFF(b *testing.B) {
	c := DefaultEnv.NewContext(nil)
	v, err := DefaultEnv.ParseValue(`if name == 'xiaoli' then age =1 else age =2 end `)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {

		c.ExecValue(v)
	}
}
