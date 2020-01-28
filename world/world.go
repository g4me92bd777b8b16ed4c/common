package world

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
)

// Beings Sort entities export
type Beings []Being

var Type byte = 31 // set if changed?

func (entities Beings) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "beings: %04d", entities.Len())
	for i := 0; i < entities.Len(); i++ {
		fmt.Fprintf(&buf, "\n\t[%d] PID: %d (%02.2f, %02.2f)",
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

type worldOut struct {
	Entities map[uint64]Being
}

// New Empty World
func New() *World {
	gob.Register(worldOut{})
	gob.Register(World{})

	return &World{entities: make(map[uint64]Being)}
}

// Update Being
func (w *World) Update(b Being) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entities[b.ID()] = b
}

// Update Being
func (w *World) Remove(playerid uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entities, playerid)
}

// Batch update Beings
func (w *World) Batch(b Beings) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range b {
		w.entities[b[i].ID()] = b[i]
	}
}

func (w *World) Get(uid uint64) Being {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.entities[uid]
}

// SnapshotBeings exports an unsorted list of beings
func (w *World) SnapshotBeings() Beings {
	w.mu.Lock()
	defer w.mu.Unlock()
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

type P struct {
	PID uint64
	PX  float64
	PY  float64
	PT  byte
}

func (p P) ID() uint64 {
	return p.PID
}
func (p P) X() float64 {
	return p.PX
}
func (p P) Y() float64 {
	return p.PY
}
func (p P) Type() byte {
	return p.PT
}

type Plist []P

func (p *Plist) Beings() Beings {
	b := make(Beings, len(*p))
	for i := range *p {
		b[i] = (*p)[i]
	}
	return b
}

func (entities Beings) Encodable() []P {
	x := make([]P, len(entities))
	for i := range entities {
		x[i].PID = (entities)[i].ID()
		x[i].PX = (entities)[i].X()
		x[i].PY = (entities)[i].Y()
		x[i].PT = (entities)[i].Type()
	}
	//	fmt.Println("ENCODING BEINGS->", x)
	return x
}
func (entities Beings) Encode(to []byte) (n int, err error) {

	buf := new(bytes.Buffer)
	if err = gob.NewEncoder(buf).Encode(entities.Encodable()); err != nil {
		return 0, err
	}
	n = copy(to[0:], buf.Bytes())
	if n != buf.Len() {
		return n, fmt.Errorf("beings encode: short write want %d got %d", buf.Len(), n)
	}
	//log.Printf("encoding=%02x", to[:n])
	return n, nil
}

func (entities *Beings) Decode(from []byte) (err error) {
	//	log.Printf("decoding=%02x", from)
	var enc = Plist{}
	if err := gob.NewDecoder(bytes.NewReader(from)).Decode(&enc); err != nil && err != io.EOF {
		return err
	}
	*entities = (enc).Beings()
	//log.Printf("decoded: %d entities", len(*entities))
	return nil
}

func (entities Beings) Type() byte {
	return Type
}
