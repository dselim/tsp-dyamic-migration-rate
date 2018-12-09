package main

import (
	"math/rand"
)

func removeOneFrom(slice []int, val int) []int {
	for i, v := range slice {
		if v == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func contains(slice []int, elem int) bool {
	for _, val := range slice {
		if val == elem {
			return true
		}
	}
	return false
}

func chooseRandom(slice []int) int {
	return slice[rand.Intn(len(slice))]
}

func randomSlice(k, n, exclude int) []int {
	r := make([]int, 0, k)
	for len(r) < k {
		i := rand.Intn(n)
		for i == exclude || contains(r, i) {
			i = rand.Intn(n)
		}
		r = append(r, i)
	}
	return r
}
