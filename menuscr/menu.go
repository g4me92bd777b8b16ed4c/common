package menuscr

import (
	"github.com/faiface/pixel"
)

type Menu struct {
	target pixel.Target
}

func New(target pixel.Target) *Menu {
	return &Menu{}
}
