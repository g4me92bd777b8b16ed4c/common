package statuseffect

import (
	"errors"
	"log"
	"math/rand"

	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/g4me92bd777b8b16ed4c/common"
	"github.com/g4me92bd777b8b16ed4c/common/types"
	"github.com/g4me92bd777b8b16ed4c/common/world"
)

type StatusManager struct {
	m  map[uint64]StatusEffect
	mu sync.Mutex
}

// StatusEffect ... using Type we can know what to give it..?
type StatusEffect interface {
	Type() types.Type
	StatusEffect(dt float64, v interface{}) (float64, error)
}

func NewManager() *StatusManager {
	return &StatusManager{m: make(map[uint64]StatusEffect)}
}

func genuuid() uint64 {
	return rand.Uint64()
}
func (s *StatusManager) Push(a common.PlayerAction) {
	switch types.Type(a.Action) {
	case types.ActionManastorm:
		//log.Println("StatusMan got Manastorm")
		u := genuuid()

		go func(u uint64, s *StatusManager) {
			d := time.Second * 3
			<-time.After(d)
			log.Println("deleting manastorm:", u)
			s.mu.Lock()
			delete(s.m, u)
			s.mu.Unlock()
		}(u, s)
		s.mu.Lock()
		s.m[u] = NewManastorm(a.ID, pixel.V(float64(a.At[0]), float64(a.At[1])))
		s.mu.Unlock()
	}
}
func (s *StatusManager) Update(dt float64, players world.Beings) []uint64 {
	deadthings := []uint64{}
	for _, player := range players {
		for _, m := range s.m {
			oldhp := player.Health()
			hp, err := m.StatusEffect(dt, player)
			if err != nil {
				log.Println("Error status effect:", err, player, dt)
			}
			if oldhp != hp {
				log.Printf("Player %d took %02.0f damage", player.ID(), hp-oldhp)
			}
			if hp == 0 {
				deadthings = append(deadthings, player.ID())
			}
		}
	}
	return deadthings
}

type manastorm struct {
	at pixel.Vec

	// from player (immune from bad effects or target of good effects)
	from uint64

	level int
	// immune players from bad effects or target good
	immune []uint64
}

func (m manastorm) Type() types.Type {
	return types.ActionManastorm
}

func (m *manastorm) StatusEffect(dt float64, v interface{}) (float64, error) {
	area := pixel.R(0, 0, 100, 100)
	c := area.Center()
	damage := float64(m.level+1) * 10
	switch v.(type) {
	case world.Being:
		if v.(world.Being).ID() == m.from {
			return v.(world.Being).Health(), nil
		}
		if area.Moved(m.at.Sub(c)).Contains(pixel.V(v.(world.Being).X(), v.(world.Being).Y())) {
			nowhp := v.(world.Being).DealDamage(m.from, damage*dt)
			//log.Println("StatusEffect: MANA, hit player:", v.(world.Being).ID(), damage*dt, "DMG , now at HP", nowhp)
			return nowhp, nil
		}
		return v.(world.Being).Health(), nil
		// case types.Typer:
		// 	// if it were to effect another effect or something?
		// 	return nil
	}
	return 0, errors.New("v is wrong interface type")

}
func NewManastorm(from uint64, at pixel.Vec) StatusEffect {
	return &manastorm{
		at:   at,
		from: from,
	}
}
