package main

import (
	"math"
	"math/rand"
)

func mockTSP(size int) []point {
	points := make([]point, size)
	for i := 0; i < size; i++ {
		points[i] = point{rand.Float64(), rand.Float64()}
	}
	return points
}

type point struct {
	x float64
	y float64
}

func dist2(u, v int) float64 {
	p1 := nodes[u]
	p2 := nodes[v]
	return (p2.x-p1.x)*(p2.x-p1.x) + (p2.y-p1.y)*(p2.y-p1.y)
}

func dist(u, v int) float64 {
	return math.Sqrt(dist2(u, v))
}

type byTSPSize [][]point

func (tsps byTSPSize) Len() int           { return len(tsps) }
func (tsps byTSPSize) Swap(i, j int)      { tsps[i], tsps[j] = tsps[j], tsps[i] }
func (tsps byTSPSize) Less(i, j int) bool { return len(tsps[i]) < len(tsps[j]) }
