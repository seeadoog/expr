package expr

import "testing"

func BenchmarkExecNew(b *testing.B) {

	v, err := DefaultEnv.ParseValue(`
 a= 1;
if a == 1 then b = 1 end;
a= 1;
if a == 1 then b = 1 end;
a= 1;
if a == 1 then b = 1 end;
a= 1;
if a == 1 then b = 1 end ;
str = 'hello';

str2 = str + 'world';
str2 = str + 'world';
str2 = str + 'world';
`)
	if err != nil {
		b.Fatal(err)
	}
	c := DefaultEnv.NewContext(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ExecValue(v)
	}
}
