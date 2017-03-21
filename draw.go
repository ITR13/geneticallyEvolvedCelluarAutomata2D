/*
    This file is part of GeneticallyEvolvedGA.

    GeneticallyEvolvedGA is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    GeneticallyEvolvedGA is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with GeneticallyEvolvedGA.  If not, see <http://www.gnu.org/licenses/>.
*/
	
package genCell

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type DrawWorld struct {
	world                *World
	rect                 *sdl.Rect
	sizemultX, sizemultY int32
	midX, midY, w, h     int32
	colors               []*sdl.Color
}

type WindowSetup struct {
	Window               *sdl.Window
	Renderer             *sdl.Renderer
	World                *World
	sizemultX, sizemultY int32
}

func (world *World) GetSetup(sizeMultX, sizeMultY int) (*WindowSetup, error) {
	W, H := world.w*sizeMultX, world.w*sizeMultY
	if SimulationType == Mirror {
		if world.oddW {
			W--
		}
		W *= 2
		if world.oddH {
			H--
		}
		H *= 2
	}
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow("GenCell", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, W, H, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		renderer.Destroy()
		return nil, err
	}
	renderer.Clear()

	return &WindowSetup{window, renderer, world,
		int32(sizeMultX), int32(sizeMultY)}, nil
}

func (windowSetup *WindowSetup) Destroy() {
	windowSetup.Renderer.Destroy()
	windowSetup.Window.Destroy()
}

func (windowSetup *WindowSetup) GetDrawWorld(colors []*sdl.Color) *DrawWorld {
	return windowSetup.World.GetDrawWorld(
		windowSetup.sizemultX, windowSetup.sizemultY, colors)
}

func (world *World) GetDrawWorld(sizemultX, sizemultY int32, colors []*sdl.Color) *DrawWorld {
	w, h := int32(world.w), int32(world.h)
	if SimulationType == Mirror {
		w *= 2
		if world.oddW {
			w--
		}
		h *= 2
		if world.oddH {
			h--
		}
	}
	midX, midY := w*sizemultX/2, h*sizemultY/2
	fmt.Printf("%d, %d\t%d, %d\t%d, %d\n",
		sizemultX, sizemultY, midX, midY, int32(world.w), int32(world.h))
	return &DrawWorld{
		world, &sdl.Rect{0, 0, sizemultX, sizemultY}, sizemultX, sizemultY,
		midX, midY, int32(world.w), int32(world.h), colors}
}

func (d *DrawWorld) Draw(r *sdl.Renderer) {
	r.Clear()
	if SimulationType == Mirror || SimulationType == Count {
		for x := int32(0); x < d.w; x++ {
			for y := int32(0); y < d.h; y++ {
				c := d.colors[(*d.world.Nodes)[x][y].NodeType]
				r.SetDrawColor(c.R, c.G, c.B, c.A)
				for mx := int32(-1); mx < 2; mx += 2 {
					for my := int32(-1); my < 2; my += 2 {
						X, Y := d.midX+mx*x*d.sizemultX,
							d.midY+my*y*d.sizemultY
						if !d.world.oddW {
							X += ((mx - 1) / 2) * d.sizemultX
						}
						if !d.world.oddH {
							Y += ((my - 1) / 2) * d.sizemultX
						}
						d.rect.X = X
						d.rect.Y = Y
						r.FillRect(d.rect)
					}
				}
			}
		}
	} else {
		for x := int32(0); x < d.w; x++ {
			for y := int32(0); y < d.h; y++ {
				c := d.colors[(*d.world.Nodes)[x][y].NodeType]
				r.SetDrawColor(c.R, c.G, c.B, c.A)
				d.rect.X, d.rect.Y =
					x*d.sizemultX, y*d.sizemultY
				r.FillRect(d.rect)
			}
		}
	}
	r.Present()
}

func (world *World) AdvanceAndPoll(poll func()) {
	if poll != nil {
		poll()
	}
	for x := 0; x < world.w; x++ {
		for y := 0; y < world.h; y++ {
			neighbors := (*world.Nodes)[x][y].Neighbors
			if SimulationType == Count {
				current := world.getNewTypeCount
				c := make([]int, world.colors)
				for i := 0; i < len(neighbors); i++ {
					c[neighbors[i].NodeType]++
				}
				if c[0] == 9 {
					continue
				}
				for i := 0; i < len(c); i++ {
					current = current[c[i]].([]interface{})
				}
				(*world.n2)[x][y].NodeType =
					current[0].(int)
			} else {
				c := 0
				for i := 0; i < len(neighbors); i++ {
					c *= world.colors
					c += neighbors[i].NodeType
				}
				(*world.n2)[x][y].NodeType = world.getNewType[c]
			}
		}
		if poll != nil {
			poll()
		}
	}
	world.Nodes, world.n2 = world.n2, world.Nodes
}
