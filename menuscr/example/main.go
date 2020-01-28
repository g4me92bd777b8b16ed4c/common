package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

func main() {
	sys := &System{}
	pixelgl.Run(sys.run)
}

type System struct{}

func (s *System) run() {
	config := pixelgl.WindowConfig{
		VSync:     true,
		Resizable: true,
		Bounds:    pixel.R(0, 0, 500, 500),
	}

	win, err := pixelgl.NewWindow(config)
	if err != nil {
		log.Fatalln(err)
	}

	var mousepos pixel.Vec
	var txt *text.Text
	txt = mkfont(24342, 22.0, "font/ka1.ttf")

	var (
		dt      = 0.0
		last    = time.Now()
		fps     = 0
		lastfps = 0
		second  = time.Tick(time.Second)
		tick    = time.Tick(time.Second / 2)
	)

	dicepic, err := loadPictureFromFile("assets/assets/spritesheet/dice-icons.png")
	if err != nil {
		panic(err)
	}
	heartpic, err := loadPictureFromFile("assets/assets/spritesheet/heart2.png") // 27x23
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(dicepic, dicepic.Bounds())
	heart := pixel.NewSprite(heartpic, heartpic.Bounds())
	var frames []pixel.Rect = make([]pixel.Rect, 48)
	maxx := dicepic.Bounds().W()
	maxy := dicepic.Bounds().H()
	wid, height := 36.0, 36.0
	i := 0
	for iy := 0.5; iy+height <= maxy; iy = iy + height + 0.5 {
		for ix := 0.5; ix+wid <= maxx; ix = ix + wid + 0.5 {
			if iy > maxy {
				iy = maxy
			}
			if ix > maxx {
				ix = maxx
			}
			frames[i] = pixel.R(ix, iy, ix+wid, iy+height)
			log.Println(i, frames[i])
			i++
		}
	}
	dicebox := []pixel.Rect{
		frames[2], frames[3], frames[4], frames[5], frames[41], frames[34],
	}

	// 2 3 4 5 41 34
	i = 0
	x := 0
	rand.Seed(time.Now().UnixNano())
	sprite.Set(dicepic, dicebox[0])
	for !win.Closed() {
		win.Clear(colors[22])
		txt.Clear()
		fps++
		dt = time.Since(last).Seconds()
		last = time.Now()

		// // back up
		// if win.JustPressed(pixelgl.KeyEnter) {
		// 	//i = (i + len(dicebox) - 1) % len(dicebox)

		// }

		// up one
		select {
		default:
		case <-second:
			lastfps = fps
			fps = 0
		case <-tick:
			i = rand.Intn(len(dicebox))
			sprite.Set(dicepic, dicebox[i])
		}

		mousepos = win.MousePosition().Sub(win.Bounds().Center())
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			//
			v := v2abs(mousepos)
			if v.X < 20 && v.Y < 20 {
				i = (i + 1) % len(dicebox)
				sprite.Set(dicepic, dicebox[i])
			}
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			x++
		}
		fmt.Fprintf(txt, "Mouse: %003.0f, %003.0f\n", mousepos.X, mousepos.Y)
		fmt.Fprintf(txt, "Win:   %003.0f, %003.0f\n", win.Bounds().Max.X, win.Bounds().Max.Y)
		fmt.Fprintf(txt, "FPS: 	 %03d DT=%03.0f ms\n", lastfps, dt*1000)
		fmt.Fprintf(txt, "Frame: %d (%s) Color: %d", i, dicebox[i], x)
		txt.Color = colors[i]
		numHearts := 20.0
		for i := 0.0; i < numHearts; i++ {
			heart.Draw(win, pixel.IM.Moved(pixel.V(100+(i*18.0), 150)))
		}
		txt.Draw(win, pixel.IM.Moved(pixel.V(5, 100)))
		sprite.DrawColorMask(win, pixel.IM.Moved(win.Bounds().Center()), colors[22])
		win.Update()
	}
}

//29 green

func v2abs(v pixel.Vec) pixel.Vec {
	return pixel.V(math.Floor(math.Abs(v.X)), math.Floor(math.Abs(v.Y)))
}
