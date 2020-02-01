package pathfinder

import (
	"github.com/beefsack/go-astar"
)

func Path(t1, t2 Pather) ([]astar.Pather, float64, bool) {
	// t1 and t2 are *Tile objects from inside the world.
	return astar.Path(t1, t2)
	// path is a slice of Pather objects which you can cast back to *Tile.
	// distance
	// and if false, dont even try
}

type Pathers = []astar.Pather
type Pather = astar.Pather
