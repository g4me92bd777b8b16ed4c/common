package types

import "testing"
import "fmt"

func TestOne(t *testing.T) {
	for _, v := range AllTypes {
		fmt.Printf("Type #%4d: %q\n", v, v.String())
	}
}
