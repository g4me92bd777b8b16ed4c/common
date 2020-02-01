package common

import "fmt"
import "testing"

func TestDir(t *testing.T) {
	dir := UPLEFT
	fmt.Println(DIR(dir).Vec())
	cc := DIR(dir).Vec()
	if cc.Y != 1 || cc.X != -1 {
		t.Log("byte:", SprintByte(dir))
		t.FailNow()
	}
}
