package main

import (
	"container/heap"
	"fmt"
	"math"
	"time"
)

type Node struct {
	level       int
	path        []int
	bound       int
	currentCost int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].bound < pq[j].bound
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func reduceMatrix(graph [][]int) ([][]int, int) {
	n := len(graph)
	reduced := make([][]int, n)
	for i := range graph {
		reduced[i] = make([]int, n)
		copy(reduced[i], graph[i])
	}

	lowerBound := 0

	// Reduce rows
	for i := 0; i < n; i++ {
		minValue := math.MaxInt
		for j := 0; j < n; j++ {
			if reduced[i][j] < minValue {
				minValue = reduced[i][j]
			}
		}
		if minValue != math.MaxInt {
			lowerBound += minValue
			for j := 0; j < n; j++ {
				if reduced[i][j] != math.MaxInt {
					reduced[i][j] -= minValue
				}
			}
		}
	}

	// Reduce columns
	for j := 0; j < n; j++ {
		minValue := math.MaxInt
		for i := 0; i < n; i++ {
			if reduced[i][j] < minValue {
				minValue = reduced[i][j]
			}
		}
		if minValue != math.MaxInt {
			lowerBound += minValue
			for i := 0; i < n; i++ {
				if reduced[i][j] != math.MaxInt {
					reduced[i][j] -= minValue
				}
			}
		}
	}

	return reduced, lowerBound
}

func calculateBound(graph [][]int, path []int) int {
	n := len(graph)
	reduced, lowerBound := reduceMatrix(graph)

	for i := 0; i < len(path)-1; i++ {
		reduced[path[i]][path[i+1]] = math.MaxInt
		for j := 0; j < n; j++ {
			reduced[path[i]][j] = math.MaxInt
			reduced[j][path[i+1]] = math.MaxInt
		}
	}

	_, additionalCost := reduceMatrix(reduced)
	return lowerBound + additionalCost
}

func branchAndBound(graph [][]int, maxIterations int) (int, []int) {
	n := len(graph)
	_, initialLowerBound := reduceMatrix(graph)

	pq := &PriorityQueue{}
	heap.Init(pq)

	initialNode := &Node{
		level:       0,
		path:        []int{0},
		bound:       initialLowerBound,
		currentCost: 0,
	}
	heap.Push(pq, initialNode)

	bestCost := math.MaxInt
	var bestPath []int
	iterationCount := 0

	for pq.Len() > 0 {
		if iterationCount >= maxIterations {
			fmt.Println("Maximum iterations reached. Terminating early.")
			break
		}

		currentNode := heap.Pop(pq).(*Node)
		iterationCount++

		if currentNode.bound >= bestCost {
			continue
		}

		if currentNode.level == n-1 {
			finalCost := currentNode.currentCost
			if finalCost < bestCost {
				bestCost = finalCost
				bestPath = append([]int{}, currentNode.path...)
			}
			continue
		}

		currentCity := currentNode.path[len(currentNode.path)-1]
		for nextCity := 0; nextCity < n; nextCity++ {
			if contains(currentNode.path, nextCity) {
				continue
			}

			nextPath := append([]int{}, currentNode.path...)
			nextPath = append(nextPath, nextCity)
			nextCost := currentNode.currentCost + graph[currentCity][nextCity]
			nextBound := calculateBound(graph, nextPath)

			if nextBound < bestCost {
				nextNode := &Node{
					level:       currentNode.level + 1,
					path:        nextPath,
					bound:       nextBound,
					currentCost: nextCost,
				}
				heap.Push(pq, nextNode)
			}
		}
	}
	fmt.Printf("Iteration count: %d\n", iterationCount)
	return bestCost, bestPath
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func calculateDistance(graph [][]int, path []int) int {
	n := len(path)
	distance := 0
	for i := 0; i < n-1; i++ {
		distance += graph[path[i]][path[i+1]]
	}
	return distance
}

var Inf int = math.MaxInt

func main() {
	// 15
	graph := [][]int{
		{Inf, 75987, 76499, 76503, 76657, 76726, 76554, 76917, 77002, 77221, 77562, 77770, 77774, 78402, 78439},
		{75627, Inf, 909, 913, 1067, 1136, 964, 1328, 1412, 1631, 1972, 2181, 2184, 2813, 2849, 2978, 3051},
		{76194, 892, Inf, 1141, 1296, 1364, 1192, 1556, 1641, 1859, 2201, 2409, 2412, 3041, 3077, 3206},
		{76204, 902, 1080, Inf, 299, 368, 550, 913, 998, 1217, 1558, 1767, 1770, 2399, 2435, 2563},
		{76359, 1057, 1234, 299, Inf, 173, 704, 1068, 1153, 1371, 1713, 1921, 1925, 2553, 2589},
		{76427, 1125, 1303, 368, 154, Inf, 773, 1137, 1221, 1440, 1781, 1990, 1993, 2622, 2658},
		{76209, 907, 1084, 623, 777, 846, Inf, 517, 602, 821, 1162, 1371, 1374, 2003, 2039},
		{76572, 1270, 1448, 987, 1141, 1210, 517, Inf, 271, 490, 666, 874, 877, 1506, 1542},
		{76657, 1355, 1533, 1072, 1226, 1295, 602, 271, Inf, 218, 560, 768, 771, 1400, 1436},
		{76876, 1574, 1751, 1290, 1444, 1513, 821, 490, 218, Inf, 341, 549, 553, 1181, 1218},
		{77217, 1915, 2093, 1632, 1786, 1855, 1162, 666, 560, 341, Inf, 208, 211, 840, 876},
		{77425, 2123, 2301, 1840, 1994, 2513, 1371, 874, 768, 549, 208, Inf, 3, 707, 743},
		{77429, 2127, 2305, 1843, 1998, 2516, 1374, 877, 771, 553, 211, 3, Inf, 711, 747},
		{78007, 2704, 2882, 2421, 2575, 3094, 1952, 1455, 1349, 1130, 789, 656, 660, Inf, 117},
		{78094, 2792, 2969, 2508, 2662, 3181, 2039, 1542, 1436, 1218, 876, 743, 747, 117, Inf},
	}

	maxIterations := 10000000 // Укажите максимальное количество итераций

	t := time.Now()
	bestCost, bestPath := branchAndBound(graph, maxIterations)
	ts := time.Since(t)
	fmt.Printf("Time: %v\n", ts)
	fmt.Printf("Minimum cost: %d\n", bestCost)
	fmt.Printf("Best path: %v\n", bestPath)
	fmt.Printf("Calc path: %v\n", calculateDistance(graph, bestPath))
	fmt.Println()

}
