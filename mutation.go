package main

func (indiv *individual) swap2OptIfOptimal(i, j int) bool {
	n := len(indiv.tour)
	k := (i + 1) % n
	l := (j + 1) % n

	oldPartialDist := dist(indiv.tour[i], indiv.tour[k]) + dist(indiv.tour[j], indiv.tour[l])
	newPartialDist := dist(indiv.tour[i], indiv.tour[j]) + dist(indiv.tour[k], indiv.tour[l])

	if newPartialDist < oldPartialDist {
		// swap tour section
		mid := (j - k + 1) / 2
		for m := 0; m < mid; m++ {
			indiv.tour[k+m], indiv.tour[j-m] = indiv.tour[j-m], indiv.tour[k+m]
		}
		// update fitness
		indiv.fitness += newPartialDist - oldPartialDist

		return true
	}

	return false
}

func (indiv *individual) optimize2OptOnce() bool {
	n := len(indiv.tour)

	for i := 0; i < n-2; i++ {
		maxJ := n
		if i == 0 {
			maxJ = n - 1
		}
		for j := i + 2; j < maxJ; j++ {
			if indiv.swap2OptIfOptimal(i, j) {
				return true
			}
		}
	}

	return false
}

const num2OptForMutation = 1

func (indiv *individual) mutation() {
	for i := 0; i < num2OptForMutation; i++ {
		if !indiv.optimize2OptOnce() {
			return
		}
	}
}
