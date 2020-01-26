package world

import (
	"log"
	"math/rand"
	"sort"
	"testing"
)

func Test(t *testing.T) {
	w := New()

	// spwn new npc with random uid
	w.Update(NewStaticEntity(0, 0, 300*rand.Float64(), 300*rand.Float64()))
	w.Update(NewStaticEntity(0, 0, 300*rand.Float64(), 300*rand.Float64()))
	w.Update(NewStaticEntity(0, 0, 300*rand.Float64(), 300*rand.Float64()))
	w.Update(NewStaticEntity(0, 0, 300*rand.Float64(), 300*rand.Float64()))
	w.Update(NewStaticEntity(0, 0, 300*rand.Float64(), 300*rand.Float64()))

	beings := (w.SnapshotBeings())
	sort.Sort(beings)
	log.Println(beings.String())
}
