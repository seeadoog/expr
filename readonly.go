package expr

type ReadOnlyMap map[string]interface{}

type ReadOnlyArray []interface{}

// NewReadOnlyVal Convert val to Readonly. change properties of readonly val is disabled.
func NewReadOnlyVal(val any) any {
	return convertToReadOnlyMap(val)
}

func convertExprReadOnlyVal(val Val) Val {
	cs, ok := val.(*constraint)
	if ok {
		cs.value = convertToReadOnlyMap(cs.value)
	}
	return val
}

func convertToReadOnlyMap(v any) any {
	switch x := v.(type) {
	case map[string]interface{}:
		res := make(ReadOnlyMap, len(x))
		for k, v := range x {
			res[k] = convertToReadOnlyMap(v)
		}
		return res
	case []any:
		res := make(ReadOnlyArray, len(x))
		for i, v := range x {
			res[i] = convertToReadOnlyMap(v)
		}
		return res
	default:
		return v
	}
}
