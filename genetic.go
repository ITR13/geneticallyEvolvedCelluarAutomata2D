package genCell

import (
	"math"
	"math/rand"

	"fmt"

	"github.com/MaxHalford/gago"
)

var MutRate float64

func GetGago(w, h, colors int, eval func(*World) float64) *gago.GA {
	MutRate = 0.15
	var redirects []int
	total := 1
	maxUsed := 1
	var geneCount int
	for i := 0; i < 9; i++ {
		total *= colors
	}
	if SimulationType == Mirror {
		fmt.Println("Setting up redirects")
		redirects = make([]int, total)
		layout := [][]int{[]int{0, 0, 0}, []int{0, 0, 0}, []int{0, 0, 1}}
		for i := 1; i < total; i++ {
			if i%1048576 == 1 {
				fmt.Printf("Redirects %f%% done...\n", float64(i)/float64(total))
			}
			if redirects[i] == 0 {
				redirects[i] = maxUsed
				c := 0
				for x := 0; x < 3; x++ {
					for y := 2; y >= 0; y-- {
						c *= colors
						c += layout[x][y]
					}
				}
				redirects[c] = maxUsed
				c = 0
				for x := 2; x >= 0; x-- {
					for y := 0; y < 3; y++ {
						c *= colors
						c += layout[x][y]
					}
				}
				c = 0
				redirects[c] = maxUsed
				for x := 2; x >= 0; x-- {
					for y := 2; y >= 0; y-- {
						c *= colors
						c += layout[x][y]
					}
				}
				redirects[c] = maxUsed
				maxUsed++
			}
			func() {
				for x := 0; x < 3; x++ {
					for y := 0; y < 3; y++ {
						layout[x][y]++
						if layout[x][y] == colors {
							layout[x][y] = 0
						} else {
							return
						}
					}
				}
			}()
		}
		geneCount = maxUsed
	} else if SimulationType == Count {
		for n := 8 + colors; n > 9; n-- {
			maxUsed *= n
		}
		for r := colors - 1; r > 1; r-- {
			maxUsed /= r
		}
		geneCount = 10
	} else {
		maxUsed = total
		geneCount = maxUsed
	}
	fmt.Printf("Total of %d possible combinations, whichof %d are used\n",
		total, maxUsed)
	return &gago.GA{
		FF{w, h, colors, total, &redirects, eval},
		Initializer{colors},
		gago.Topology{
			2,
			8,
			48,
			geneCount,
		},
		gago.ModGenerational{
			gago.SelElitism{},
			Crossover{maxUsed, colors},
			Mutator{colors, &redirects},
			2,
		},
		256,
		gago.MigShuffle{},
		gago.Individual{}, 0, 0, nil,
	}
}

type Initializer struct {
	colors int
}

func (initializer Initializer) Apply(indi *gago.Individual, rng *rand.Rand) {
	if SimulationType == Count {
		var SetGenesRecursively func(*[]interface{}, int, int)
		SetGenesRecursively = func(current *[]interface{}, maxI, maxD int) {
			if maxD <= 0 {
				for i := 0; i <= maxI; i++ {
					(*current)[i] = rng.Int() % initializer.colors
				}
			} else {
				for i := 0; i <= maxI; i++ {
					cur := make([]interface{}, maxI-i+1)
					SetGenesRecursively(&cur, maxI-i, maxD-1)
					(*current)[i] = cur
				}
			}
		}
		genome := []interface{}(indi.Genome)
		SetGenesRecursively(&genome, 9, initializer.colors)
		indi.Genome = genome
	} else {
		for i := range indi.Genome {
			indi.Genome[i] = rng.Int() % initializer.colors
		}
	}
}

type FF struct {
	w, h, colors, total int
	redirects           *[]int
	eval                func(*World) float64
}

func (ff FF) GetGetNewType(genome gago.Genome) []int {
	var casted []int
	if SimulationType == Mirror {
		casted = make([]int, ff.total)
		for i := range casted {
			casted[i] = genome[(*ff.redirects)[i]].(int)
		}

	} else if SimulationType == Count {
		panic("Not used when counting")
	} else {
		casted = make([]int, ff.total)
		for i := range genome {
			casted[i] = genome[i].(int)
		}
	}
	casted[0] = 0
	return casted
}

func (ff FF) Apply(genome gago.Genome) float64 {
	var getNewType []int
	if SimulationType != Count {
		getNewType = ff.GetGetNewType(genome)
	}

	w := GetWorld(ff.w, ff.h, ff.colors, getNewType, genome)
	return ff.eval(w)
}

type Crossover struct {
	geneCount, colors int
}

func (crossover Crossover) Apply(ind1, ind2 gago.Individual,
	r *rand.Rand) (gago.Individual, gago.Individual) {
	g1, g2 := make([]interface{}, crossover.geneCount),
		make([]interface{}, crossover.geneCount)
	if SimulationType == Count {
		var SetGenesRecursively func(*[]interface{}, []interface{},
			[]interface{}, int, int)
		SetGenesRecursively = func(current *[]interface{}, par1,
			par2 []interface{}, maxI, maxD int) {
			if maxD <= 0 {
				for i := 0; i <= maxI; i++ {
					if r.Int()%2 == 0 {
						//fmt.Printf("%d\t%d\t%d\n", i, len(*current), len(par1))
						(*current)[i] = par1[i]
					} else {
						//fmt.Printf("%d\t%d\t%d\n", i, len(*current), len(par2))
						(*current)[i] = par2[i]
					}
				}
			} else {
				for i := 0; i <= maxI; i++ {
					cur := make([]interface{}, maxI-i+1)
					SetGenesRecursively(&cur, par1[i].([]interface{}),
						par2[i].([]interface{}), maxI-i, maxD-1)
					(*current)[i] = cur
				}
			}
		}
		SetGenesRecursively(&g1, ind1.Genome, ind2.Genome, 9, crossover.colors)
		SetGenesRecursively(&g2, ind1.Genome, ind2.Genome, 9, crossover.colors)
	} else {
		for i := range ind1.Genome {
			if r.Int()%2 == 0 {
				g1[i], g2[i] = ind2.Genome[i], ind1.Genome[i]
			} else {
				g1[i], g2[i] = ind1.Genome[i], ind2.Genome[i]
			}
		}
	}
	return gago.Individual{g1, math.Inf(1), false, "-"},
		gago.Individual{g2, math.Inf(1), false, "-"}
}

type Mutator struct {
	colors    int
	redirects *[]int
}

func (mutator Mutator) Apply(ind *gago.Individual, r *rand.Rand) {
	//Using this instead of built in mut-rate to increase with fail
	if r.Float64() < MutRate {
		if SimulationType == Count {
			var SetGenesRecursively func(*[]interface{}, int, int)
			SetGenesRecursively = func(current *[]interface{}, maxI, maxD int) {
				if maxD <= 0 {
					for i := 0; i <= maxI; i++ {
						if r.Float64() < 0.001 || (MutRate > 20 && r.Float64() < 0.05) {
							(*current)[i] = r.Int() % mutator.colors
						}
					}
				} else {
					for i := 0; i <= maxI; i++ {
						cur := (*current)[i].([]interface{})
						SetGenesRecursively(&cur, maxI-i, maxD-1)
						(*current)[i] = cur
					}
				}
			}
			genome := []interface{}(ind.Genome)
			SetGenesRecursively(&genome, 9, mutator.colors)
			ind.Genome = genome
		} else {
			if r.Float64() < 0.0005 || (MutRate > 20 && r.Float64() < 0.005*math.Log2(MutRate/20)) {
				field := make([][]int, 3)
				for x := 0; x < 3; x++ {
					field[x] = make([]int, 3)
					for y := 0; y < 3; y++ {
						field[x][y] = r.Int() % mutator.colors
					}
				}
				getC := func() int {
					c := 0
					for x := 0; x < 3; x++ {
						for y := 0; y < 3; y++ {
							c *= mutator.colors
							c += field[x][y]
						}
					}
					return c
				}
				myC := getC()
				if SimulationType == Mirror {
					myC = (*mutator.redirects)[myC]
				}

				gen := ind.Genome[myC]

				from := r.Int() % 9
				toAmount := r.Int()%200 + 10
				for i := 0; i < toAmount; i++ {
					change := r.Int() % 9
					for change == from {
						change = r.Int() % 9
					}
					field[change/3][change%3] += r.Int() % (mutator.colors - 1)
					field[change/3][change%3] = field[change/3][change%3] % (mutator.colors)
					myC = getC()
					if SimulationType == Mirror {
						myC = (*mutator.redirects)[myC]
					}
					ind.Genome[myC] = gen
				}
			} else {
				chance := r.Float64()*2/5 - 0.05
				if chance < 0 {
					chance += 0.2
				}
				for i := range ind.Genome {
					if r.Float64() <= chance {
						ind.Genome[i] = r.Int() % mutator.colors
					}
				}
			}
		}
	}
}
