package world

import (
	"math"

	"github.com/beefsack/go-astar"
	"github.com/faiface/pixel"
	"github.com/g4me92bd777b8b16ed4c/common/pathfinder"
	"github.com/g4me92bd777b8b16ed4c/common/types"
)

func (t *Tile) PathNeighborCost(to pathfinder.Pather) float64 {
	return to.(*Tile).MovementCost
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	return t.ManhattanDistance(to)
}

func (t *Tile) PathTo(t2 *Tile) (pathfinder.Pathers, float64, bool) {
	return pathfinder.Path(t, t2)
}

type Tile struct {
	X, Y         float64
	MovementCost float64
	Type         types.Type
	tileset      Tileset
}

func (t *Tile) SetTileset(newtileset Tileset) {
	t.tileset = newtileset
}
func (t *Tile) Copy(newtileset Tileset) *Tile {
	return &Tile{
		X:            t.X,
		Y:            t.Y,
		tileset:      newtileset,
		MovementCost: t.MovementCost,
	}
}

type Tileset map[pixel.Vec]*Tile

func (t *Tile) PathNeighbors() pathfinder.Pathers {
	return pathfinder.Pathers{
		t.Up(),
		t.Right(),
		t.Down(),
		t.Left(),
	}
}

func (t *Tile) Up() *Tile {
	return t.tileset[t.Pos().Add(pixel.V(0, 1))]
}
func (t *Tile) Down() *Tile {
	return t.tileset[t.Pos().Add(pixel.V(0, -1))]
}
func (t *Tile) Left() *Tile {
	return t.tileset[t.Pos().Add(pixel.V(-1, 0))]
}
func (t *Tile) Right() *Tile {
	return t.tileset[t.Pos().Add(pixel.V(1, 0))]
}

func (t *Tile) Pos() pixel.Vec {
	return pixel.V(t.X, t.Y)
}

func (t *Tile) ManhattanDistance(to astar.Pather) float64 {
	return Distance(t.Pos(), to.(*Tile).Pos())
}

// Distance between two vectors
func Distance(v1, v2 pixel.Vec) float64 {
	r := pixel.Rect{v1, v2}.Norm()
	v1 = r.Min
	v2 = r.Max
	h := (v1.X - v2.X) * (v1.X - v2.X)
	v := (v1.Y - v2.Y) * (v1.Y - v2.Y)
	return Sqrt(h + v)
}

func ManDist(v1, v2 pixel.Vec) float64 {
	r := pixel.Rect{v1, v2}.Norm()
	v1 = r.Min
	v2 = r.Max
	return math.Abs(v2.X-v1.X) + math.Abs(v2.Y-v1.Y)
}
func Sqrt(x float64) float64 {
	z := float64(2.)
	s := float64(0)
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
		if math.Abs(z-s) < 1e-10 {
			break
		}
		s = z
	}
	return z
}

func (w *World) NewTile(tiletype types.Type, x, y float64) *Tile {
	t := &Tile{
		X:            x,
		Y:            y,
		MovementCost: 0,
		tileset:      w.tiles,
		Type:         tiletype,
	}
	w.tiles[pixel.V(x, y)] = t
	return t
}
