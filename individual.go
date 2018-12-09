package main

import (
	"math/rand"
)

type individual struct {
	tour    []int
	fitness float64
}

func newIndividual(tour []int) *individual {
	indiv := &individual{
		tour: tour,
	}
	indiv.fitness = fitness(indiv.tour)
	return indiv
}

func randomIndividual() *individual {
	perm := rand.Perm(len(nodes) - 1)
	for i, u := range perm {
		perm[i] = u + 1
	}
	return newIndividual(append([]int{0}, perm...))
}

func (indiv *individual) copy() *individual {
	tour := make([]int, len(indiv.tour))
	copy(tour, indiv.tour)
	return &individual{
		tour:    tour,
		fitness: indiv.fitness,
	}
}

func (indiv *individual) shuffle(num int) {
	n := len(indiv.tour)
	for k := 0; k < num; k++ {
		i := rand.Intn(n)
		j := rand.Intn(n)
		if j == i {
			if j > 0 {
				j--
			} else {
				j++
			}
		}
		indiv.tour[i], indiv.tour[j] = indiv.tour[j], indiv.tour[i]
	}
	indiv.fitness = fitness(indiv.tour)
}

// FITNESS

func fitness(tour []int) (fit float64) {
	i, j := tour[0], tour[len(tour)-1]
	fit += dist(i, j)
	for k := 0; k < len(tour)-1; k++ {
		i, j = tour[k], tour[k+1]
		fit += dist(i, j)
	}
	return
}

// SIMILARITY

func similarity(indiv1, indiv2 *individual) float64 {
	generatePerm := func(tour []int) []int {
		n := len(tour)
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			j := (i + 1) % n
			perm[tour[i]] = tour[j]
		}
		return perm
	}

	next1 := generatePerm(indiv1.tour)
	next2 := generatePerm(indiv2.tour)

	commonEdgeCount := 0
	for v, nextV := range next1 {
		if v == next2[nextV] || nextV == next2[v] {
			commonEdgeCount++
		}
	}
	return float64(commonEdgeCount) / float64(len(next1))
}
