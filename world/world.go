package world

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/faiface/pixel"
	"github.com/g4me92bd777b8b16ed4c/common/types"

	// "io"
	"log"
	"math/rand"
	//"sync"
)

// Beings Sort entities export
type Beings []Being

func (entities Beings) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "beings: %04d", entities.Len())
	for i := 0; i < entities.Len(); i++ {
		fmt.Fprintf(&buf, "\n\t[%d] PID: %d (%02.2f, %02.2f)",
			entities[i].Type(), entities[i].ID(), entities[i].X(), entities[i].Y())
	}

	return buf.String()
}

// Being / Entity: a thing with distinct and independent existence. also see Object
type Being interface {
	ID() uint64
	X() float64
	Y() float64
	Type() types.Type
	MoveTo([2]float64)

	Health() float64
	SetHealth(float64)
	DealDamage(from uint64, damage float64) (healthAfter float64)
}

// World for entities to live in
type World struct {
	entities map[uint64]Being
	mu       sync.Mutex

	tiles   map[pixel.Vec]*Tile
	objects map[uint64]Object
}

type worldOut struct {
	Entities map[uint64]Being
}

type Object interface {
	Pos() pixel.Vec
	Name() string
	Type() types.Type
}

// New Empty World
func New() *World {
	gob.Register(worldOut{})
	gob.Register(World{})

	return &World{entities: make(map[uint64]Being), objects: make(map[uint64]Object), tiles: make(map[pixel.Vec]*Tile)}
}

// Update Being
func (w *World) DealDamage(fromPlayer uint64, toPlayer uint64, amount float64) (healthAfter float64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.entities[toPlayer] == nil {
		return 0
	}
	healthAfter = w.entities[toPlayer].DealDamage(fromPlayer, amount)
	return healthAfter
}

// Update Being
func (w *World) Update(b Being) (isNew bool) {
	//log.Println("UPDATING BEING:", b.ID(), b.Type().String(), "HP", b.Health())
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.entities[b.ID()] == nil {
		log.Printf("Making new entity in world map %02x: id %d type %d (%s) HP %2.0f", &b, b.ID(), b.Type(), b.Type().String(), b.Health())
		w.entities[b.ID()] = b
		isNew = true
		return
	}
	//log.Println("Updating entity:", b.Type().String(), b.ID(), b.Health(), "HP", w.entities[b.ID()].X(), w.entities[b.ID()].Y())
	w.entities[b.ID()].SetHealth(b.Health())
	w.entities[b.ID()].MoveTo([2]float64{b.X(), b.Y()})
	//log.Println("Updated entity:", b.Type().String(), b.ID(), b.Health(), "HP", b.X(), b.Y())
	return
}

// Update Being
func (w *World) Remove(playerid uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entities, playerid)
}
func (w *World) BatchFn(fn func(b Being) error) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range w.entities {
		if err := fn(w.entities[i]); err != nil {
			return err
		}
	}
	return nil
}

// Batch update Beings
func (w *World) Batch(b Beings) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range b {
		if w.entities[b[i].ID()] == nil {
			log.Println("Making new entity in world map (same ptr):", b[i])
			w.entities[b[i].ID()] = b[i]
		}
		w.entities[b[i].ID()].MoveTo([2]float64{b[i].X(), b[i].Y()})
	}
}

func (w *World) Get(uid uint64) Being {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.entities[uid]
}
func (w *World) GetTile(at pixel.Vec) *Tile {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.tiles[at]
}
func (w *World) SetTile(t *Tile) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.tiles[t.Pos()] = t
}

func (w *World) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.entities)
}

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
	// not sorting here, but sortable somehow
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

	// abs
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
	typ  types.Type
}

func (s *staticEntity) Health() float64                             { return 100 }
func (s *staticEntity) SetHealth(float64)                           {}
func (s *staticEntity) DealDamage(from uint64, dam float64) float64 { return 100 }
func (s *staticEntity) ID() uint64                                  { return s.id }
func (s *staticEntity) X() float64                                  { return s.x }
func (s *staticEntity) Y() float64                                  { return s.y }
func (s *staticEntity) Type() types.Type                            { return s.typ }
func (s *staticEntity) MoveTo(v [2]float64) {
	s.x = v[0]
	s.y = v[1]
}

// NewStaticEntity returns a ... if static entity uid or beingtype are zero they will randomly be generated
func NewStaticEntity(uid uint64, beingtype types.Type, atx, aty float64) Being {
	if uid == 0 {
		uid = rand.Uint64()
		log.Println("New UID:", uid)
	}
	if beingtype == 0 {
		panic("zero static entity type")
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
	PT  types.Type
}

func (p P) Encode(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(p); err != nil {
		return 0, err
	}
	n = copy(b[:], buf.Bytes())
	return n, nil
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
func (p P) Type() types.Type {
	return p.PT
}
func (p P) MoveTo(d [2]float64) {
	p.PX = d[0]
	p.PY = d[1]
}

// type Plist []P

// func (p *Plist) Beings() Beings {
// 	b := make(Beings, len(*p))
// 	for i := range *p {
// 		b[i] = (*p)[i]
// 	}
// 	return b
// }

// func (entities Beings) Encodable() []P {
// 	x := make([]P, len(entities))
// 	for i := range entities {
// 		x[i].PID = (entities)[i].ID()
// 		x[i].PX = (entities)[i].X()
// 		x[i].PY = (entities)[i].Y()
// 		x[i].PT = (entities)[i].Type()
// 	}
// 	//	fmt.Println("ENCODING BEINGS->", x)
// 	return x
// }
// func (entities Beings) Encode(to []byte) (n int, err error) {

// 	buf := new(bytes.Buffer)
// 	if err = gob.NewEncoder(buf).Encode(entities.Encodable()); err != nil {
// 		return 0, err
// 	}
// 	n = copy(to[0:], buf.Bytes())
// 	if n != buf.Len() {
// 		return n, fmt.Errorf("beings encode: short write want %d got %d", buf.Len(), n)
// 	}
// 	//log.Printf("encoding=%02x", to[:n])
// 	return n, nil
// }

// func (entities *Beings) Decode(from []byte) (err error) {
// 	//	log.Printf("decoding=%02x", from)
// 	var enc = Plist{}
// 	if err := gob.NewDecoder(bytes.NewReader(from)).Decode(&enc); err != nil && err != io.EOF {
// 		return err
// 	}
// 	*entities = (enc).Beings()
// 	//log.Printf("decoded: %d entities", len(*entities))
// 	return nil
// }

// func (entities Beings) Type() byte {
// 	return
// }
