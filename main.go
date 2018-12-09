package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	nodes []point
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) != 4 {
		log.Fatalf("usage: %s <tsp file> <island file> <output file>\n", os.Args[0])
	}

	nodes = parseFile(os.Args[1])

	islands := islandsFromFile(os.Args[2])

	csvHelper := newCSVHelper(os.Args[3])
	defer csvHelper.close()

	iterations := 20
	generations := 100
	// between 1 and numIslands
	migRates := []float64{1, 2, 3, 4}
	//migRates := []float64{1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5, 6}

	for m, migRate := range migRates {

		avgBestFit := 0.0
		for i := 0; i < iterations; i++ {
			log.Printf("running iteration %v/%v (migRate=%v %v/%v)\n", m*iterations+i+1, len(migRates)*iterations, migRate, i+1, iterations)
			t := time.Now()

			bestFitPerGen, _ := run(generations, islands, migRate, false)
			avgBestFit += bestFitPerGen[len(bestFitPerGen)-1]

			log.Printf("done after %v\n", time.Since(t))
		}
		avgBestFit /= float64(iterations)

		csvHelper.write([]string{
			strconv.FormatFloat(migRate, 'E', -1, 64),
			strconv.FormatFloat(avgBestFit, 'E', -1, 64),
		})
	}
}

func run(generations int, islands []population, migRate float64, dynamic bool) (bestFitPerGen []float64, divPerIslPerGen [][]float64) {
	numIslands := len(islands)

	islandsCopy := make([]population, numIslands)
	for i, popul := range islands {
		islandsCopy[i] = popul.copy()
	}
	islands = islandsCopy

	bestFitPerGen = make([]float64, generations+1)

	divPerIslPerGen = make([][]float64, numIslands)
	for i := 0; i < numIslands; i++ {
		divPerIslPerGen[i] = make([]float64, generations+1)
	}

	migrants := make([][]*individual, numIslands)
	for i := 0; i < numIslands; i++ {
		migrants[i] = make([]*individual, 0, len(islands[i]))
	}

	var wgMain sync.WaitGroup
	wgIslands := make([]sync.WaitGroup, numIslands)
	// launch islands
	for i := 0; i < numIslands; i++ {
		go func(i int) {
			clock := time.Now()
			if i == 0 {
				log.Printf("%5.1f%% done\n", 0.0)
			}

			for g := 0; g < generations; g++ {
				divPerIslPerGen[i][g] = islands[i].diversity()

				// new generation
				islands[i] = islands[i].newGeneration()

				// immigration
				for j, immigrant := range migrants[i] {
					_, islands[i] = islands[i].removeRandom()
					islands[i] = islands[i].insert(immigrant)

					migrants[i][j] = nil
				}
				migrants[i] = migrants[i][:0]

				if i == 0 && time.Since(clock) > 10*time.Second {
					clock = time.Now()
					log.Printf("%5.1f%% done\n", 100*float64(g+1)/float64(generations))
				}

				wgMain.Done()
				wgIslands[i].Add(1)
				wgIslands[i].Wait()
			}
			divPerIslPerGen[i][generations] = islands[i].diversity()
		}(i)
	}

	// synchronize islands
	for g := 0; g < generations; g++ {
		bestFitPerGen[g] = islandsBestFitness(islands)

		wgMain.Add(numIslands)
		wgMain.Wait()

		// migrations
		for i := 0; i < numIslands; i++ {
			immigrations := int(migRate)
			if dynamic {
				concentration := 1.0 - islands[i].diversity()
				immigrations = int(migRate * concentration)
			}

			sources := randomSlice(immigrations, numIslands, i)
			for _, j := range sources {
				k := islands[j].selection()
				migrants[i] = append(migrants[i], islands[j][k].copy())
			}
		}

		for i := 0; i < numIslands; i++ {
			wgIslands[i].Done()
		}
	}

	bestFitPerGen[generations] = islandsBestFitness(islands)
	return
}
