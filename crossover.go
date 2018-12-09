package main

import (
	"math/rand"
)

func crossover(parent1, parent2 *individual) *individual {
	n := len(parent1.tour)

	neighbors := make(map[int][]int, n)

	for u := 0; u < n; u++ {
		neighbors[u] = make([]int, 0, 4)
	}

	addNeighbors := func(neighbors map[int][]int, tour []int) {
		for i, u := range tour {
			v := tour[(i+1)%n]
			w := tour[(i+2)%n]

			if !contains(neighbors[v], u) {
				neighbors[v] = append(neighbors[v], u)
			}
			if !contains(neighbors[v], w) {
				neighbors[v] = append(neighbors[v], w)
			}
		}
	}
	addNeighbors(neighbors, parent1.tour)
	addNeighbors(neighbors, parent2.tour)

	offspringTour := make([]int, 0, n)
	u := parent1.tour[0]

	remaining := make(map[int]bool, n)
	for i := 0; i < n; i++ {
		remaining[i] = true
	}

	for len(offspringTour) < n {
		// add u to offspring
		offspringTour = append(offspringTour, u)
		delete(remaining, u)

		// remove u from neighbors
		for _, v := range neighbors[u] {
			neighbors[v] = removeOneFrom(neighbors[v], u)
		}

		// if u has available neighbors
		if len(neighbors[u]) > 0 {
			// take the best neighbor
			v := neighbors[u][0]
			for _, w := range neighbors[u][1:] {
				if dist(u, w) < dist(u, v) {
					v = w
				}
			}
			u = v
		} else {
			for v := range remaining {
				u = v
				break
			}
		}
	}

	return newIndividual(offspringTour)
}

// LEGACY

const (
	lagacyCrossoverMinK = 0.2
	lagacyCrossoverMaxK = 0.3
)

func lagacyCrossoverMultiPoint(parent1, parent2 *individual) *individual {
	// copy content of parent tours so we don't alter the originals
	tour1 := append([]int(nil), parent1.tour...)
	tour2 := append([]int(nil), parent2.tour...)

	offspringTour := make([]int, 0, len(tour1))

	k := lagacyCrossoverMinK + rand.Float64()*(lagacyCrossoverMaxK-lagacyCrossoverMinK)
	offset := int(k * float64(len(tour1)))

	for len(tour1) > 0 {
		maxOffset := offset
		if len(tour1) < maxOffset {
			maxOffset = len(tour1)
		}

		section := tour1[:maxOffset]

		offspringTour = append(offspringTour, section...)
		tour1 = tour1[maxOffset:]
		for _, i := range section {
			tour2 = removeOneFrom(tour2, i)
		}

		tour1, tour2 = tour2, tour1
	}

	return newIndividual(offspringTour)
}

var (
	selectionScale = 4.0
)

func lagacyCrossoverSelection(n int) int {
	r := rand.Float64() * (float64(n) * (selectionScale + 1.0)) / 2
	for i := 0; i < n-1; i++ {
		r -= ((1-selectionScale)/float64(n-1))*float64(i) + selectionScale
		if r <= 0 {
			return i
		}
	}
	return n - 1
}
