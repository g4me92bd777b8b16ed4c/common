package common

import (
	//"bytes"
	//"encoding/gob"
	"log"

	"github.com/faiface/pixel"
	"github.com/g4me92bd777b8b16ed4c/common/types"
)

// Component type.
// Such as...
// Health
// Mana
// StatusEffects
// ArmorClass
// WeaponDamage
type Component interface {
	Type() types.Type
	Value() interface{}
	String() string // return short human readable form, Include "Component" or "Component " prefix
}

const (
	max64 uint64 = 1<<63 - 1
	max32 uint32 = 1<<31 - 1
	max16 uint16 = 1<<15 - 1
	max8  uint8  = 1<<7 - 1
)

type StatusEffect struct {
	StatusType types.Type
}

func (s StatusEffect) Type() types.Type {
	return s.StatusType
}

// type StatBarInterface interface {
// 	StatusEffects() []StatusEffect
// }

type StatBar struct {
	StatType types.Type
	Current  uint16
	Max      uint16
}

type StatBar64 struct {
	StatType types.Type
	Current  uint64
	Max      uint64
}

type HealthBar struct {
	StatBar
}
type ManaBar struct {
	StatBar
}
type ExpBar struct {
	StatBar64
}

func (s StatBar) Type() types.Type {
	return s.StatType
}
func (s StatBar) CurrentMax() float64 {
	return float64(s.Current) / float64(s.Max)
}
func (s StatBar64) Type() types.Type {
	return s.StatType
}
func (s StatBar64) CurrentMax() float64 {
	return float64(s.Current) / float64(s.Max)
}

func init() {

}

// Player is a Typer, Positioner, Pos(), and a world.Being (MoveTo, ID, X, Y)
//
// Weird entity names because methods should instead be used by users of this struct.
// Positions/Components arent guarded, and instead should be guarded by the World or Game
type Player struct {
	PID        uint64
	EntityType uint16
	PosX       float64
	PosY       float64
	// StatBars   []StatBar
	// StatBar64  []StatBar64
	HP float64
	MP float64
}

func (e Player) Health() float64 {
	return e.HP
}
func (e *Player) SetHealth(hp float64) {
	e.HP = hp
}
func (e *Player) DealDamage(from uint64, amount float64) float64 {
	if e.ID() == from && amount > 0 {
		return e.HP
	}
	if e.HP < amount {
		e.HP = 0
		return 0
	}
	e.HP -= amount
	return e.HP
}
func (p Player) ID() uint64 {
	return p.PID
}
func (p Player) Pos() pixel.Vec {
	return pixel.V(p.PosX, p.PosY)
}
func (p Player) X() float64 { return p.PosX }
func (p Player) Y() float64 { return p.PosY }

func (p *Player) MoveTo(xy [2]float64) {
	p.PosX = xy[0]
	p.PosY = xy[1]
}
func (p Player) Type() types.Type {
	return types.Type(p.EntityType)
}

// func (p Player) Encode(b []byte) (n int, err error) {
// 	buf := new(bytes.Buffer)
// 	err = gob.NewEncoder(buf).Encode(p)
// 	if err != nil {
// 		log.Println("Encoding game.Player:", err)
// 		return 0, err
// 	}
// 	n += copy(b[:], buf.Bytes())
// 	return n, nil
// }

type DIR byte

// 1000
// 0101
func (dd DIR) Vec() pixel.Vec {
	switch byte(dd) {
	case UP:
		return pixel.V(0, 1)
	case UPRIGHT:
		return pixel.V(1, 1)
	case UPLEFT:
		return pixel.V(-1, 1)
	case LEFT:
		return pixel.V(-1, 0)
	case RIGHT:
		return pixel.V(1, 0)
	case DOWN:
		return pixel.V(0, -1)
	case DOWNRIGHT:
		return pixel.V(1, -1)
	case DOWNLEFT:
		return pixel.V(-1, -1)
	default:
		return pixel.ZV
	}
}

func Vec2Dir(v pixel.Vec) byte {
	switch v {
	case pixel.V(0, 1):
		return UP
	case pixel.V(1, 1):
		return UPRIGHT
	case pixel.V(-1, 1):
		return LEFT
	case pixel.V(1, 0):
		return DOWN
	case pixel.V(1, 0):
		return DOWNRIGHT
	case pixel.V(-1, 0):
		return DOWNLEFT
	case pixel.V(1, 0):
		return RIGHT
	case pixel.V(-1, 0):
		return LEFT
	default:
		log.Println("UNKONWN ANGLE:", v)
		return 0
	}
}

type HealthReport struct {
	M map[uint64]HealthStatus // current/max health
}

type HealthStatus struct {
	ID   uint64
	From uint64
	Cur  float64
	Max  float64
	Dam  float64
}
