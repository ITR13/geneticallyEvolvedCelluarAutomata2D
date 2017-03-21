/*
    This file is part of InvertoTanks.

    Foobar is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    InvertoTanks is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with InvertoTanks.  If not, see <http://www.gnu.org/licenses/>.
*/
	
package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ITR13/geneticallyEvolvedCelluarAutomata2D"
	"github.com/MaxHalford/gago"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	simSizeX int = 95
	simSizeY int = 95
)

var currentSetup *genCell.WindowSetup
var colors []*sdl.Color

var primes []int

func main() {
	primes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541}

	fmt.Println("Starting")
	genCell.SimulationType = genCell.Count
	colors = make([]*sdl.Color, 5)
	colors[0] = &sdl.Color{255, 255, 255, 255}
	colors[1] = &sdl.Color{0, 255, 0, 0}
	colors[2] = &sdl.Color{0, 127, 0, 0}
	colors[3] = &sdl.Color{0, 51, 102, 255}
	colors[4] = &sdl.Color{255, 165, 0, 255}
	fmt.Println("Set starting values")

	ga := genCell.GetGago(simSizeX, simSizeY, 5, EvaluateFlower)
	fmt.Println("Got GA")
	ga.Initialize()
	fmt.Println("Initialized")
	score := math.Inf(1)
	increases := 0
	failedToIncrease := 0
	go StartDrawFlower()

	backup, err := genCell.Read("./backup.bytes")
	if err == nil {
		for true {
			DrawFlower(backup, ga)
		}
	}

	for true {
		ga.Enhance()
		if score > ga.Best.Fitness {
			increases++
			failedToIncrease = 0
			genCell.Save("./output/Flower"+strconv.Itoa(increases)+".bytes", ga.Best)
			genCell.Save("./temp.bytes", ga.Best)
			score = ga.Best.Fitness
			genCell.MutRate = 0.15
		} else {
			failedToIncrease++
			if failedToIncrease%4 == 0 {
				genCell.MutRate *= 1.25
			}
		}
		fmt.Printf("Best of generation %d: %f"+
			"\t(Increases: %d, Generations since last increase: %d)\n",
			ga.Generations, ga.Best.Fitness, increases, failedToIncrease)
		if makeNewWindow {
			DrawFlower(ga.Best, ga)
			makeNewWindow = false
		}
	}
}

func EvaluateFlower(w *genCell.World) float64 {
	score := float64(0)
	(*w.Nodes)[0][0].NodeType = 1
	cont := true
	for days := 0; days < 8; days++ {
		totHours := primes[(days*12)%100] * (days*12/100 + 1)
		for hours := 0; hours < totHours; hours++ {
			tmp, end, _, _ :=
				GetPoints(w, hours >= totHours/2)
			score += tmp
			if end || !cont {
				cont = false
			} else {
				w.Advance()
			}

			score += 0.25
		}
	}
	return score
}

func GetPoints(w *genCell.World, night bool) (float64, bool, int, int) {
	score := float64(0)
	end := true
	green := 0
	other := 0
	for x := range *w.Nodes {
		for y := range (*w.Nodes)[x] {
			n := (*w.Nodes)[x][y]
			if n.NodeType != 0 {
				end = false
				score += 0.05
				if n.NodeType == 1 || n.NodeType == 2 {
					green++
				} else if n.NodeType == 3 || n.NodeType == 4 {
					other++
				}
				if night {
					if n.NodeType == 1 {
						score += 1.9
					} else if n.NodeType == 2 {
						score += 0.8
					} else if n.NodeType == 4 {
						n.NodeType = 3
					}
				} else {
					if n.NodeType == 1 {
						score -= 2
					} else if n.NodeType == 2 {
						score -= 1
					} else if n.NodeType == 3 {
						n.NodeType = 4
					}
				}
			} else if night {
				score += 0.01
			} else {
				score += 1
			}
		}
	}
	return score, end, green, other
}

var makeNewWindow bool

func StartDrawFlower() {
	fmt.Println("Ready to signal window-making!")
	for true {
		fmt.Scanln()
		if makeNewWindow {
			fmt.Println("Already signaling!")
		} else {
			fmt.Println("Signaling to make a new window!")
		}
		makeNewWindow = true
	}
}

func DrawFlower(ind gago.Individual, ga *gago.GA) {
	ff := ga.Ff.(genCell.FF)
	var genes []int
	if genCell.SimulationType != genCell.Count {
		genes = ff.GetGetNewType(ind.Genome)
	}
	world := genCell.GetWorld(simSizeX, simSizeY, 5, genes, ind.Genome)
	if genCell.SimulationType == genCell.Mirror ||
		genCell.SimulationType == genCell.Count {
		(*world.Nodes)[0][0].NodeType = 1
	} else {
		(*world.Nodes)[simSizeX/2][simSizeY/2].NodeType = 1
	}
	setup, err := world.GetSetup(8, 8)
	if err != nil {
		fmt.Printf("Failed to set up window:\n%v\n", err)
		return
	}
	fmt.Println("Set up new window!")
	defer setup.Destroy()

	drawWorld := setup.GetDrawWorld(colors)

	quit = false
	i := 0
	drawWorld.Draw(setup.Renderer)
	for i := 0; i < 10 && !quit; i++ {
		fmt.Println(9 - i)
		delay(1000)
	}
	prev := time.Now()
	totalScore := float64(0)
	for !quit {
		totHours := primes[(i*12)%100] * (i*12/100 + 1)
		for hours := 0; hours < totHours && !quit; hours++ {
			score, end, green, other :=
				GetPoints(world, hours >= totHours/2)
			totalScore += score
			if end {
				quit = true
			}
			fmt.Printf("Day %d, Hour %d\nnight: %t\nScore: "+
				"%.02f (today %.02f)\ngreen: %d, other: %d\n",
				i, hours, hours >= totHours/2, totalScore, score, green, other)
			pollEvents()
			drawWorld.Draw(setup.Renderer)
			world.AdvanceAndPoll(pollEvents)

			now := time.Now()
			d := uint32(now.Sub(prev).Seconds() * 400)
			if d <= 0 {
				delay(400)
			} else if d >= 400 {
				pollEvents()
			} else {
				delay(400 - d)
			}
			prev = time.Now()
		}
		i++
	}
}

var quit bool

func delay(ms uint32) {
	for ms > 20 {
		sdl.Delay(20)
		pollEvents()
		if quit {
			return
		}
		ms -= 20
	}
	if ms > 0 {
		sdl.Delay(ms)
		pollEvents()
	}
}

func pollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) { //Add window resizing
		case *sdl.QuitEvent:
			fmt.Print("Exited: ")
			fmt.Println(t.Timestamp)
			quit = true
		}
	}
}
