package game

import (
	"fmt"
	"math"
)

var endian = Endian()

// type Msgpacket struct {
// 	Move float64 `json:'mv,omitempty'`
// 	Say  string  `json:'say,omitempty'`
// 	Do   []byte  `json:'do,omitempty'`
// }

const (
	x uint16 = 1 << 15
)

const (
	MaxChatSize = 8
)

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

type MessageType byte

// Network Messages
const (
	MessageSize                = 1
	_              MessageType = 2
	MsgPrefixLogin MessageType = iota
	MsgPrefixDPad
	MsgPrefixRepeat
	MsgPrefixUpdateGps
	MsgPrefixUpdatePlayers
	MsgPrefixPlayerMessage
	MsgPrefixErrorMessage
	MsgPrefixBufferSize
	MsgPrefixPlayerLogoff
	MsgPrefixPing
	MsgPrefixPong
	MsgPrefixWorld
)

var Sizemap = msgSize
var msgSize = [256]int{
	MsgPrefixLogin:         255,
	MsgPrefixDPad:          1,
	MsgPrefixRepeat:        24,
	MsgPrefixUpdateGps:     16,
	MsgPrefixUpdatePlayers: 24,
	MsgPrefixPlayerMessage: 1024,
	MsgPrefixErrorMessage:  1024,
	MsgPrefixBufferSize:    8,
	MsgPrefixPing:          255,
	MsgPrefixPong:          255,
}

func MsgSize(b byte) int {
	return msgSize[b]
}

var msgmap = [256]string{
	MsgPrefixDPad:          "[DPAD]",
	MsgPrefixRepeat:        "[REPEAT]",
	MsgPrefixUpdateGps:     "[UPDATESELF]",
	MsgPrefixUpdatePlayers: "[UPDATEPLAYERS]",
}

func MSG(b byte) string {
	return msgmap[b]
}

func GPS(gps [2]int64) string {
	return fmt.Sprintf("(%d,%d)", gps[0], gps[1])
}
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
	Other     uint16
}

func MessageEncode(m Message) []byte {
	var out uint64
	var buf = make([]byte, 8)

	out |= uint64(m.Dpad)
	fmt.Printf("[dpad] %020d: %s\n", out, Sprint64(out))
	out |= uint64(m.Key) << 16
	fmt.Printf("[ key] %020d: %s\n", out, Sprint64(out))
	out |= uint64(m.Keymod) << 32
	fmt.Printf("[ mod] %020d: %s\n", out, Sprint64(out))
	out |= uint64(m.Other) << 48

	endian.PutUint64(buf, out)
	return buf
}

func MessageDecode(b []byte) Message {
	out := endian.Uint64(b)
	return Message{
		Dpad:      byte(out & math.MaxUint8),
		UnusedPad: byte(out >> 8 & math.MaxUint8),
		Key:       int16(out >> 16 & math.MaxInt16),
		Keymod:    int16(out >> 32 & math.MaxInt16),
		Other:     uint16(out >> 48 & math.MaxUint16),
	}
}

const (
	ENT_NONE byte = iota
	ENT_HUMAN
	ENT_LAST // only for ENT_LAST-1
)
