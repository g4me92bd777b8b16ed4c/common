package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	"github.com/faiface/pixel/pixelgl"
)

func TestOne(t *testing.T) {
	var (
		dpad   byte            // 8, but dpad is rightmost 4 bits only , what can we fit in other 4?
		key    int16           // 16, to fit in all the keys possible (one at a time though)
		keymod int16           // 16, a single key modifier
		other  uint16 = 0xffff // ouch
	)

	dpad = DOWN | LEFT
	key = int16(pixelgl.KeyF)
	keymod = int16(pixelgl.KeyLeftControl)
	fmt.Printf("[in] dir=%s key=%d=%s     mod=%d=%s other=%x\n", DPAD(dpad), key, pixelgl.KeyF.String(), keymod, pixelgl.KeyLeftControl.String(), other)

	var (
		out uint64
	)

	out |= uint64(dpad)
	fmt.Printf("[dpad] %020d: %s\n", out, Sprint64(out))
	out |= uint64(key) << 16
	fmt.Printf("[ key] %020d: %s\n", out, Sprint64(out))
	out |= uint64(keymod) << 32
	fmt.Printf("[ mod] %020d: %s\n", out, Sprint64(out))
	out |= uint64(other) << 48

	fmt.Printf("[othr] %020d: %s\n", out, Sprint64(out))
	fmt.Println("DECODE")

	var ndpad byte = byte(out & math.MaxUint8)
	var npadpad byte = byte(out >> 8 & math.MaxUint8)
	var nkey int16 = int16(out >> 16 & math.MaxInt16)
	var nmod int16 = int16(out >> 32 & math.MaxInt16)
	var nother uint16 = uint16(out >> 48 & math.MaxUint16)

	// mouse
	// var mouseX uint64 = math.Float64bits(1440.42)
	// var mouseY uint64 = math.Float64bits(4000.224)
	// // 128
	if ndpad != 6 || npadpad != 0 || nkey != 70 || nmod != 341 || nother != 65535 {
		fmt.Println("expected: dpad=6 (downleft), pad=0, nkey=70, nmod=341, nother=65535, got:")
		fmt.Printf("dpad=%d (%s), pad=%d, nkey=%d, nmod=%d, nother=%d\n", ndpad, DPAD(ndpad), npadpad, nkey, nmod, nother)
		t.FailNow()
	}

	// fmt.Printf("\tdpad+key=%08b\n\tX=%064b\n\tY=%08b\n\tall=%d\n", out, mouseX, mouseY, 64+64+8)

}

func TestTwo(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, tc := range [][2]float64{
		[2]float64{65757865876587658765786578652.0365436543654365430, -4.00},
		[2]float64{87657865876587658765768587646752.02764367547654367436540, -4.00},
		[2]float64{4657467547654756476542.0345254324352453254324530345235425432, -4.00},
		[2]float64{7654675435643654356432.0076587658764565745634523, -4.00},
		[2]float64{65435643564365848752.07654765476546754765476540, -4.00},
		[2]float64{646754765476542.76547654765400, -4.00},
		[2]float64{56365436453564345632.765467540765465740, -4.00},
	} {
		err := binary.Write(buf, endian, tc)
		if err != nil {
			fmt.Printf("err=%v \n%v==(%b) size=%d\n", err, tc, buf.Bytes(), buf.Len())
		}
	}

}

func TestCodecOne(t *testing.T) {

}
