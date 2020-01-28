package plugint

import (
	"fmt"
	"plugin"

	"github.com/faiface/pixel"
	"github.com/g4me92bd777b8b16ed4c/common/codec"
	"github.com/g4me92bd777b8b16ed4c/common/world"
)

type PluginUpdateFunc = func(float64, *world.World, *uint8, *codec.Codec, pixel.Target)
type PluginDrawFunc = func(pixel.Target)

func RegisterPlugin(w *world.World, c *codec.Codec, target pixel.Target, pluginPath string) (updateFn PluginUpdateFunc, drawFn PluginDrawFunc, err error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, nil, err
	}
	name, err := p.Lookup("Name")
	if err != nil {
		return nil, nil, err
	}
	pluginName := name.(func() string)()

	// if plugin needs Init
	init, err := p.Lookup("Init")
	if err == nil {
		println(pluginName, "Initializing")
		init.(func(*world.World, *codec.Codec))(w, c)
	}
	var ok bool
	p_updateFn, err := p.Lookup("Update")
	if err == nil {
		updateFn, ok = p_updateFn.(PluginUpdateFunc)
		if !ok {
			return nil, nil, fmt.Errorf("wrong type: %T", p_updateFn)
		}
	}
	p_drawFn, err := p.Lookup("Draw")
	if err == nil {
		drawFn, ok = p_drawFn.(PluginDrawFunc)
		if !ok {
			return nil, nil, fmt.Errorf("wrong type: %T", p_drawFn)
		}
	}
	return updateFn, drawFn, nil
}
