package common

import (
	"fmt"

	"github.com/g4me92bd777b8b16ed4c/common/types"
)

var endian = Endian()

// type Msgpacket struct {
// 	Move float64 `json:'mv,omitempty'`
// 	Say  string  `json:'say,omitempty'`
// 	Do   []byte  `json:'do,omitempty'`
// }

const (
	UP            byte = 1 << 0
	DOWN          byte = 1 << 1
	LEFT          byte = 1 << 2
	RIGHT         byte = 1 << 3
	UPLEFT        byte = UP | LEFT
	UPRIGHT       byte = UP | RIGHT
	DOWNLEFT      byte = DOWN | LEFT
	DOWNRIGHT     byte = DOWN | RIGHT
	ALLDIR        byte = UP | DOWN | LEFT | RIGHT
	UPDOWN        byte = UP | DOWN
	LEFTRIGHT     byte = LEFT | RIGHT
	UPDOWNLEFT    byte = UP | DOWN | LEFT
	UPDOWNRIGHT   byte = UP | DOWN | RIGHT
	LEFTDOWNRIGHT byte = LEFT | DOWN | RIGHT
	LEFTUPRIGHT   byte = LEFT | UP | RIGHT
)

func MSG(b byte) string {
	return types.Type(b).String()
}

// func GPS(gps [2]int64) string {
// 	return fmt.Sprintf("(%d,%d)", gps[0], gps[1])
// }

func DPAD(b byte) string {
	switch b {
	case UP:
		return "up"
	case DOWN:
		return "down"
	case LEFT:
		return "left"
	case RIGHT:
		return "right"
	case UPDOWN:
		return "up+down"
	case LEFTRIGHT:
		return "left+right"
	case ALLDIR:
		return "alldirections"
	case UPLEFT:
		return "upleft"
	case UPRIGHT:
		return "upright"
	case DOWNLEFT:
		return "downleft"
	case DOWNRIGHT:
		return "downright"
	case UPDOWNLEFT:
		return "updownleft"
	case LEFTDOWNRIGHT:
		return "leftdownright"
	case UPDOWNRIGHT:
		return "updownright"
	case LEFTUPRIGHT:
		return "leftupright"
	case 0:
		return "none"
	default:
		return fmt.Sprintf("invalid dir: %d", b)
	}
}

// Sprint64 returns a binary representation of 64 bit unsigned integer, but in 8 chunks of 8 bits
func Sprint64(i uint64) string {
	buf := make([]byte, 8)
	endian.PutUint64(buf, i)
	return fmt.Sprintf("%08b\n", buf)
}

// SprintByte returns a binary representation of 8 bits
func SprintByte(c byte) string {
	return fmt.Sprintf("%08b\n", c)
}

type Message struct {
	Dpad      byte
	UnusedPad byte
	Key       int16
	Keymod    int16
	Action    uint16
	Count     uint16
}
