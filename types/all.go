package types

import (
	logg "log"
	"os"
	"strings"
)

// one big enum?
// or...

//go:generate stringer -type Type
//go:generate sed -i s/Type(/UnknownType(/g type_string.go
var AllTypes []Type

const MaxType = 65000 // for now!

func init() {
	var log = logg.New(os.Stderr, "types: ", logg.Llongfile)

	//	log.SetPrefix("types: ")
	for i := 0; i < MaxType; i++ {
		s := Type(i).String()
		if strings.Contains(s, "UnknownType") {
			continue
		}
		log.Printf("Registering type: %d (%q)", i, s)
		AllTypes = append(AllTypes, Type(i))
	}
}

// Type
type Type uint16 // two byte, 16 bits, 65k types

func (t Type) Byte() byte {
	if t > 0xFF {
		return 0
	}
	return byte(t)
}
func (t Type) Uint16() uint16 {
	if t > 0xFFFF {
		return 0
	}
	return uint16(t)
}

const (
	_ Type = iota // 0 reserved for unset

	// Game things (nine reserved)
	Error
	Warn
	Info
	Update
	Asset
	Menu // options
	Player
)

// networking stuff
const (
	Login Type = iota + 20
	DPad
	Repeat
	UpdateGps
	UpdatePlayers
	PlayerMessage
	ErrorMessage
	BufferSize
	PlayerLogoff
	Ping
	Pong
	World
)

// under 255 are able to be byte()

// Entities

const (
	Human Type = iota + 1024
	Robot
	Alien
	Zombie
	Skeleton
	NPC1
	NPC2
)
