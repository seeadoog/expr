package expr

import (
	"fmt"
	"testing"
)

func TestINva(t *testing.T) {
	v, err := DefaultEnv.ParseValue(`exit1 [true]`)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}
