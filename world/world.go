package world

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"sync"
)

// Beings Sort entities export
type Beings []Being

var Type byte = 13 // set if changed?

func (entities Beings) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "beings: %04d", entities.Len())
	for i := 0; i < entities.Len(); i++ {
		fmt.Fprintf(&buf, "\n\t[%d] PID: %x (%02.2f, %02.2f)",
			entities[i].Type(), entities[i].ID(), entities[i].X(), entities[i].Y())
	}

	return buf.String()
}

// Being / Entity: a thing with distinct and independent existence.
type Being interface {
	ID() uint64
	X() float64
	Y() float64
	Type() byte // can increase once we have 255 types i guess
}

// World for entities to live in
type World struct {
	entities map[uint64]Being
	mu       sync.Mutex
}

// New Empty World
func New() *World {
	return &World{entities: make(map[uint64]Being)}
}

// Update Being
func (w *World) Update(b Being) {
	w.entities[b.ID()] = b
}

// Batch update Beings
func (w *World) Batch(b Beings) {
	for i := range b {
		w.entities[b[i].ID()] = b[i]
	}
}

// SnapshotBeings exports an unsorted list of beings
func (w *World) SnapshotBeings() Beings {
	b := make(Beings, len(w.entities))
	n := 0
	for _, v := range w.entities {
		b[n] = v
		n++
	}
	//sort.Sort(b)
	return b
}

func (entities Beings) Len() int {
	return len(entities)
}

func (entities Beings) Swap(x, y int) {
	entities[x], entities[y] = entities[y], entities[x]
}

func (entities Beings) Less(x, y int) bool {

	x1 := entities[x].X()
	x2 := entities[y].X()
	y1 := entities[x].Y()
	y2 := entities[y].Y()
	if x1 < 0 {
		x1 *= -1
	}
	if x2 < 0 {
		x2 *= -1
	}
	if y1 < 0 {
		y1 *= -1
	}
	if y2 < 0 {
		y2 *= -1
	}
	return x1+y1 < x2+y2
}

// here is an entity
type staticEntity struct {
	id   uint64
	x, y float64
	typ  byte
}

func (s *staticEntity) ID() uint64 { return s.id }
func (s *staticEntity) X() float64 { return s.x }
func (s *staticEntity) Y() float64 { return s.y }
func (s *staticEntity) Type() byte { return s.typ }

// NewStaticEntity returns a ... if static entity uid or beingtype are zero they will randomly be generated
func NewStaticEntity(uid uint64, beingtype byte, atx, aty float64) Being {
	if uid == 0 {
		uid = rand.Uint64()
		log.Println("New UID:", uid)
	}
	if beingtype == 0 {
		beingtype = byte(rand.Intn(10))
		log.Printf("UID %d has new Beingtype: %d", uid, beingtype)
	}
	return &staticEntity{
		id:  uid,
		x:   atx,
		y:   aty,
		typ: beingtype,
	}
}

func (entities *Beings) Encode(to []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(entities)
	n = copy(to, buf.Bytes())
	if n != buf.Len() {
		return n, fmt.Errorf("beings encode: short write want %d got %d", buf.Len(), n)
	}
	return n, nil
}
func (entities *Beings) Decode(from []byte) (err error) {
	return gob.NewDecoder(bytes.NewReader(from)).Decode(entities)
}

func (entities Beings) Type() byte {
	return Type
}
