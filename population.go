package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"sort"
)

type population []*individual

func (popul population) newGeneration() population {
	newPopul := make(population, 0, len(popul))

	elite := popul[0].copy()
	elite.mutation()
	newPopul = newPopul.insert(elite)

	for len(newPopul) < len(popul) {
		i := popul.selection()
		j := popul.selectionExcept(i)
		offspring := crossover(popul[i], popul[j])

		offspring.mutation()

		newPopul = newPopul.insert(offspring)
	}
	return newPopul
}

// tournament selection

func (popul population) selection() int {
	sum := 0.0
	for _, indiv := range popul {
		sum += indiv.fitness
	}
	r := rand.Float64() * sum
	i := 0
	for r > 0 {
		r -= popul[i].fitness
		i++
	}
	return i - 1
}

func (popul population) selectionExcept(index int) int {
	i := popul.selection()
	if i == index {
		if i > 0 {
			return i - 1
		}
		return i + 1
	}
	return i
}

func randomPopulation(size int) population {
	popul := make(population, size)
	for i := 0; i < size; i++ {
		popul[i] = randomIndividual()
	}
	popul.sort()
	return popul
}

func randomPopulationOptimized(size int) population {
	popul := make(population, size)
	for i := 0; i < size; i++ {
		popul[i] = randomIndividual()
		for m := 0; m < 10; m++ {
			popul[i].mutation()
		}
	}
	popul.sort()
	return popul
}

// implement sort.Interface
func (popul population) Len() int           { return len(popul) }
func (popul population) Swap(i, j int)      { popul[i], popul[j] = popul[j], popul[i] }
func (popul population) Less(i, j int) bool { return popul[i].fitness < popul[j].fitness }

func (popul population) sort() {
	sort.Sort(popul)
}

func (popul population) insert(indiv *individual) population {
	for i, indiv2 := range popul {
		if indiv.fitness < indiv2.fitness {
			// insert
			newPopul := append(popul, nil)
			copy(newPopul[i+1:], popul[i:])
			newPopul[i] = indiv
			return newPopul
		}
	}
	return append(popul, indiv)
}

func (popul population) removeRandom() (*individual, population) {
	i := rand.Intn(len(popul))

	indiv := popul[i]

	copy(popul[i:], popul[i+1:])
	popul[len(popul)-1] = nil
	return indiv, popul[:len(popul)-1]
}

func (popul population) copy() population {
	newPopul := make(population, len(popul))
	for i, indiv := range popul {
		newPopul[i] = indiv.copy()
	}
	return newPopul
}

func (popul population) diversity() float64 {
	similSum := 0.0

	samples := make([]int, 4)
	samples[0] = 0
	samples[3] = len(popul) - 1
	samples[1] = samples[3] / 3
	samples[2] = samples[1] * 2
	for k, i := range samples[:3] {
		for _, j := range samples[k+1:] {
			similSum += similarity(popul[i], popul[j])
		}
	}

	return 1.0 - similSum/6.0
}

// ISLANDS

func (popul population) generateIslands(numIslands int) []population {
	populSizePerIsland := len(popul) / numIslands
	populRemainder := len(popul) % numIslands

	indivs := popul.copy()

	islands := make([]population, numIslands)
	for i := 0; i < numIslands; i++ {
		islandSize := populSizePerIsland
		if populRemainder > 0 {
			islandSize++
			populRemainder--
		}

		island := make(population, islandSize)
		for j := 0; j < islandSize; j++ {
			island[j], indivs = indivs.removeRandom()
		}
		island.sort()
		islands[i] = island
	}

	return islands
}

func randomIslandsOptimized(populSize, numIslands int) []population {
	popul := randomPopulationOptimized(populSize)
	return popul.generateIslands(numIslands)
}

func islandsDiversity(islands []population) float64 {
	numIslands := len(islands)
	similSum := 0.0
	counter := 0.0
	for k, island1 := range islands[:numIslands-1] {
		for _, island2 := range islands[k+1:] {
			similSum += similarity(island1[0], island2[0])
			counter++
		}
	}
	return 1.0 - similSum/counter
}

func islandsBestFitness(islands []population) float64 {
	bestFit := islands[0][0].fitness
	for _, island := range islands[1:] {
		if island[0].fitness < bestFit {
			bestFit = island[0].fitness
		}
	}
	return bestFit
}

// IO UTILS

func islandsToFile(islands []population, filepath string) {
	tours := make([][][]int, len(islands))
	for i, popul := range islands {
		tours[i] = make([][]int, len(popul))
		for j, indiv := range popul {
			tours[i][j] = indiv.tour
		}
	}
	bytes, err := json.Marshal(tours)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func islandsFromFile(filepath string) []population {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var tours [][][]int
	err = json.Unmarshal(bytes, &tours)
	if err != nil {
		panic(err)
	}
	islands := make([]population, len(tours))
	for i, t := range tours {
		islands[i] = make(population, 0, len(t))
		for _, tour := range t {
			islands[i] = islands[i].insert(newIndividual(tour))
		}
	}
	return islands
}
