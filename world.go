package genCell

const (
	Loop SimType = iota
	//Have sides connect to themselves instead of loop?
	Bounce SimType = iota
	//Make patterns symmetric to be more efficient?
	Mirror SimType = iota
	//Make pattern count how many of the different neighbors there are
	Count SimType = iota
)

type SimType int

var SimulationType SimType

type Node struct {
	NodeType  int
	Neighbors []*Node
}

type World struct {
	getNewType      []int
	getNewTypeCount []interface{}
	Nodes, n2       *[][]*Node
	colors, w, h    int
	oddW, oddH      bool
}

func (world *World) Advance() {
	world.AdvanceAndPoll(nil)
}

func GetWorld(w, h, colors int, getNewType []int,
	getNewTypeCount []interface{}) *World {
	oddW, oddH := false, false
	if SimulationType == Mirror || SimulationType == Count {
		oddW = w%2 == 1
		oddH = h%2 == 1
		w = w/2 + w%2
		h = h/2 + h%2
	}
	nodes := make([][]*Node, w)
	n2 := make([][]*Node, w)
	for i := 0; i < w; i++ {
		nodes[i] = make([]*Node, h)
		n2[i] = make([]*Node, h)
		for j := 0; j < h; j++ {
			nodes[i][j] = &Node{0, nil}
			n2[i][j] = &Node{0, nil}
		}
	}

	world := &World{getNewType, getNewTypeCount,
		&nodes, &n2, colors, w, h, oddW, oddH}

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			nodes[i][j].Neighbors, n2[i][j].Neighbors =
				world.getSurrounding(i, j)
		}
	}
	return world
}

func (world *World) getSurrounding(X, Y int) ([]*Node, []*Node) {
	i := 0
	nodes := make([]*Node, 3*3)
	n2 := make([]*Node, 3*3)
	for x := X - 1; x <= X+1; x++ {
		for y := Y - 1; y <= Y+1; y++ {
			wx, wy := slide(x, world.w, world.oddW),
				slide(y, world.h, world.oddH)
			nodes[i] = (*world.Nodes)[wx][wy]
			n2[i] = (*world.n2)[wx][wy]
			i++
		}
	}
	return nodes, n2
}

func slide(v, max int, odd bool) int {
	if SimulationType == Mirror || SimulationType == Count {
		if v >= max {
			v = 2*max - v - 1
		} else if v < 0 {
			v = -v
			if odd {
				v++
			}
		}
	} else if SimulationType == Bounce {
		if v >= max {
			v = 2*max - v - 1
		} else if v < 0 {
			v = -v
		}
	} else {
		if v >= max {
			v -= max
		} else if v < 0 {
			v += max
		}
	}
	return v
}
