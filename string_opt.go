package expr

import (
	"strconv"
	"strings"
)

type StringOpt struct {
	res map[string]string
}

// parse opts like  a=b,c=d,e=f
func NewStringOpt(str string) *StringOpt {
	kv := make(map[string]string)
	for _, s := range strings.Split(str, ",") {
		k, v, ok := strings.Cut(s, "=")
		if ok {
			kv[k] = v
		}
	}
	return &StringOpt{kv}
}

func (o *StringOpt) GetIntDef(key string, def int) int {
	res := o.res[key]
	if res == "" {
		return def
	}
	i, err := strconv.Atoi(res)
	if err != nil {
		return def
	}
	return i
}
func (o *StringOpt) GetStrDef(key string, def string) string {
	res, ok := o.res[key]
	if !ok {
		return def
	}
	return res
}
func (o *StringOpt) GetBoolDef(key string, def bool) bool {
	res := o.res[key]
	if res == "" {
		return def
	}
	b, err := strconv.ParseBool(res)
	if err != nil {
		return def
	}
	return b
}

func (o *StringOpt) GetFloatDef(key string, def float64) float64 {
	res := o.res[key]
	if res == "" {
		return def
	}
	f, err := strconv.ParseFloat(res, 64)
	if err != nil {
		return def
	}
	return f
}
