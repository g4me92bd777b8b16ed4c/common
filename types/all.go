package types

// logg "log"
// "os"
// "strings"

// one big enum?
// or...

//go:generate stringer -type Type
//go:generate sed -i s/Type(/UnknownType(/g type_string.go
// var AllTypes []Type

const MaxType uint16 = 64000 // for now!

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
	File
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
	PlayerAction
	PlayerMessage
	ErrorMessage
	BufferSize
	PlayerLogoff
	Ping
	Pong
	World
	RemoveEntity
	PlayerDeath
)

// Tiles
const (
	NoTile Type = iota + 256
	TileGrass
	TileWater
	TileRock
)

// Entities

const (
	Villager Type = iota + 1024
	Robot
	Alien
	Zombie
	GhostBomb
	Skeleton
	NPC1
	NPC2
)

type Typer interface {
	Type() Type
}

// Actions
const (
	ActionSlash Type = iota + 2048
	ActionThrow
	ActionUse
	ActionSpell
	ActionManastorm
	ActionMagicbullet
)

// Status Effects
const (

	// Positive ones
	StatusHaste Type = iota + 4096
	StatusBloodlust

	// Negative ones
	StatusPoisoned = iota + 5120
	StatusSleeping
	StatusTired
	StatusConfused
)
const BarXP = BarEXP
const (
	BarHealth Type = iota + 8192
	BarMana
	BarEXP
)
const (
	SomethingElse Type = iota + 16384
)
const (
	SomethingUndefined Type = iota + 32768
)

// // func Bytes2Type(b []byte) Type {
// // 	if len(b) > 2 {
// // 		return 0
// // 	}
// // 	var t Type
// // 	t |= Type(b[0]) << 8
// // 	t |= Type(b[1]) << 0
// // 	return t
// // }

// func Type2Bytes(t Type) [2]byte {
// 	return [2]byte{byte(t >> 8), byte(t << 8 >> 8)}
// }
