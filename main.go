package main

import (
	"fmt"
	"math"
)

type TSPSolver struct {
	n           int
	graph       [][]int
	minCost     int
	bestPath    []int
	currentPath []int
	visited     []bool
}

func NewTSPSolver(graph [][]int) *TSPSolver {
	n := len(graph)
	return &TSPSolver{
		n:           n,
		graph:       graph,
		minCost:     math.MaxInt,
		bestPath:    make([]int, n+1),
		currentPath: make([]int, n+1),
		visited:     make([]bool, n),
	}
}

func (solver *TSPSolver) solveRecursively(level int, cost int) {
	if level == solver.n {
		finalCost := cost + solver.graph[solver.currentPath[level-1]][solver.currentPath[0]]
		if finalCost < solver.minCost {
			solver.minCost = finalCost
			copy(solver.bestPath, solver.currentPath)
			solver.bestPath[solver.n] = solver.currentPath[0]
		}
		return
	}

	for i := 0; i < solver.n; i++ {
		if !solver.visited[i] {
			solver.visited[i] = true
			solver.currentPath[level] = i

			if level == 0 || cost+solver.graph[solver.currentPath[level-1]][i] < solver.minCost {
				nextCost := cost
				if level > 0 {
					nextCost += solver.graph[solver.currentPath[level-1]][i]
				}
				solver.solveRecursively(level+1, nextCost)
			}

			solver.visited[i] = false
		}
	}
}

func (solver *TSPSolver) Solve() (int, []int) {
	solver.solveRecursively(0, 0)
	return solver.minCost, solver.bestPath
}

func main() {
	// graph := [][]int{
	// 	{math.MaxInt, 755, 1701, 1866, 1675, 2424, 2652, 2707, 2522, 2732, 3497, 3763},
	// 	{712, math.MaxInt, 1578, 1581, 1390, 2043, 2855, 2910, 2237, 2906, 3212, 3142},
	// 	{1985, 1912, math.MaxInt, 542, 726, 1460, 2191, 2246, 1573, 1783, 2548, 3168},
	// 	{1987, 2028, 542, math.MaxInt, 729, 1644, 2194, 2249, 1576, 1785, 2550, 3171},
	// 	{1797, 1838, 727, 729, math.MaxInt, 736, 1548, 1603, 1051, 1260, 2025, 2646},
	// 	{2316, 2357, 1246, 1248, 603, math.MaxInt, 1146, 1201, 1500, 1710, 2309, 2930},
	// 	{3071, 2867, 2125, 2783, 2137, 1166, math.MaxInt, 134, 2007, 2431, 2997, 3617},
	// 	{3920, 3716, 2749, 2751, 2106, 2015, 134, math.MaxInt, 1975, 2400, 2965, 3586},
	// 	{2644, 2685, 1573, 1576, 1051, 1634, 2563, 2618, math.MaxInt, 531, 1296, 1916},
	// 	{3283, 3323, 1783, 1785, 1260, 1843, 3037, 3092, 531, math.MaxInt, 859, 1479},
	// 	{2580, 2621, 2854, 2857, 2331, 4186, 3932, 3987, 1602, 1165, math.MaxInt, 777},
	// 	{3416, 3457, 3291, 3294, 2768, 4623, 4369, 4424, 2039, 1602, 2506, math.MaxInt},
	// }
	graph := [][]int{
		{math.MaxInt, 755, 1701, 1866, 1675, 2424, 2652, 2707, 2522, 2732, 3497, 3763},
		{0, math.MaxInt, 1578, 1581, 1390, 2043, 2855, 2910, 2237, 2906, 3212, 3142},
		{0, 1912, math.MaxInt, 542, 726, 1460, 2191, 2246, 1573, 1783, 2548, 3168},
		{0, 2028, 542, math.MaxInt, 729, 1644, 2194, 2249, 1576, 1785, 2550, 3171},
		{0, 1838, 727, 729, math.MaxInt, 736, 1548, 1603, 1051, 1260, 2025, 2646},
		{0, 2357, 1246, 1248, 603, math.MaxInt, 1146, 1201, 1500, 1710, 2309, 2930},
		{0, 2867, 2125, 2783, 2137, 1166, math.MaxInt, 134, 2007, 2431, 2997, 3617},
		{0, 3716, 2749, 2751, 2106, 2015, 134, math.MaxInt, 1975, 2400, 2965, 3586},
		{0, 2685, 1573, 1576, 1051, 1634, 2563, 2618, math.MaxInt, 531, 1296, 1916},
		{0, 3323, 1783, 1785, 1260, 1843, 3037, 3092, 531, math.MaxInt, 859, 1479},
		{0, 2621, 2854, 2857, 2331, 4186, 3932, 3987, 1602, 1165, math.MaxInt, 777},
		{0, 3457, 3291, 3294, 2768, 4623, 4369, 4424, 2039, 1602, 2506, math.MaxInt},
	}

	solver := NewTSPSolver(graph)
	minCost, bestPath := solver.Solve()

	fmt.Printf("Minimum cost: %d\n", minCost)
	fmt.Printf("Best path: %v\n", bestPath)
}
