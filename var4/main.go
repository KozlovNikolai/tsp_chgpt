package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
)

type Node struct {
	level       int
	path        []int
	bound       int
	currentCost int
	parent      *Node
	matrix      [][]int
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

func reduceMatrix(mx [][]int, path []int) [][]int {
	// printMatrix(graph)

	mxr := RemoveCellFromMatrixByIdx(mx, 0)
	// printMatrix(reduced)
	return mxr
}
func printMatrix(mx [][]int) {
	for i := 0; i < len(mx); i++ {
		for j := 0; j < len(mx[0]); j++ {
			if mx[i][j] == math.MaxInt {
				fmt.Printf("\tInf\t")
			} else {
				fmt.Printf("\t%d\t", mx[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// func calculateBound(mx [][]int, path []int) int {

// 	// // Mark edges to prevent cycles
// 	// for i := 0; i < len(path)-1; i++ {
// 	// 	from, to := path[i], path[i+1]
// 	// 	reduced[from][to] = math.MaxInt
// 	// 	for j := 0; j < n; j++ {
// 	// 		reduced[from][j] = math.MaxInt
// 	// 		reduced[j][to] = math.MaxInt
// 	// 	}
// 	// }

// 	// Mark return edge to prevent premature cycle closure
// 	if len(path) == n {
// 		fmt.Printf("Mark Infinity Cell: (%d,%d)\n", path[len(path)-1], path[0])
// 		mx[path[len(path)-1]][path[0]] = math.MaxInt
// 	}

// 	// Compute the lower bound for the reduced matrix
// 	_, lowerBound := reduceMatrix(reduced)
// 	fmt.Printf("lower bound: %d\n", lowerBound)
// 	// Add the costs of the selected edges
// 	for i := 0; i < len(path)-1; i++ {
// 		fmt.Printf("up LB (%d,%d):%d\n", path[i], path[i+1], mx[path[i]][path[i+1]])
// 		lowerBound += mx[path[i]][path[i+1]]
// 	}

// 	return lowerBound, reduced
// }

func calculateBound(graph [][]int) int {
	lowerBound := 0
	n := len(graph)
	// Reduce rows
	for i := 0; i < n; i++ {
		minValue := math.MaxInt
		for j := 0; j < n; j++ {
			if graph[i][j] < minValue {
				minValue = graph[i][j]
			}
		}
		if minValue != math.MaxInt {
			lowerBound += minValue
		}
	}
	printMatrix(graph)
	fmt.Printf("initial lower bound: %d\n", lowerBound)
	return lowerBound
}

func cloneMx(mx [][]int) [][]int {
	lenRows := len(mx)
	lenCols := len(mx[0])
	mxClone := make([][]int, lenRows)
	for i := range mxClone {
		mxClone[i] = make([]int, lenCols)
	}
	for i := 0; i < lenRows; i++ {
		for j := 0; j < lenCols; j++ {
			mxClone[i][j] = mx[i][j]
		}
	}
	return mxClone
}

func branchAndBound(graph [][]int) (int, []int) {

	initialLowerBound := calculateBound(graph)
	pq := &PriorityQueue{}
	heap.Init(pq)

	treeNodes := []*Node{}
	leafNodes := []*Node{}

	initialNode := &Node{
		level:       0,
		path:        []int{0},
		bound:       initialLowerBound,
		currentCost: 0,
		parent:      nil,
		matrix:      cloneMx(graph),
	}
	heap.Push(pq, initialNode)
	treeNodes = append(treeNodes, initialNode)
	leafNodes = append(leafNodes, initialNode)

	bestCost := math.MaxInt
	var bestPath []int
	step := 0

	for len(leafNodes) > 0 {
		step++
		n := len(graph)
		// Output current state of leaves
		fmt.Printf("Step %d:\n", step)
		fmt.Printf("Current leaves:\n")
		for _, leaf := range leafNodes {
			fmt.Printf("(%v):%d\n", leaf.path, leaf.bound)
		}

		// Find leaf with minimum bound
		minLeafIndex := 0
		for i, leaf := range leafNodes {
			if leaf.bound < leafNodes[minLeafIndex].bound {
				minLeafIndex = i
			}
		}
		currentNode := leafNodes[minLeafIndex]
		// извлекаем из массива листов лист, из которого будет дальнейший рост...
		leafNodes = append(leafNodes[:minLeafIndex], leafNodes[minLeafIndex+1:]...)
		// работаем с извлеченным листом как с точкой роста
		fmt.Printf("Selected leaf: (%v):%d\n", currentNode.path, currentNode.bound)

		if currentNode.level == n-1 {
			finalCost := currentNode.currentCost + graph[currentNode.path[len(currentNode.path)-1]][0]
			if finalCost < bestCost {
				bestCost = finalCost
				bestPath = append([]int{}, currentNode.path...)
				bestPath = append(bestPath, 0) // Completing the cycle
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
			nextMatrix := reduceMatrix(cloneMx(currentNode.matrix), nextPath)
			nextBound := calculateBound(currentNode.matrix)

			if nextBound < bestCost {
				nextNode := &Node{
					level:       currentNode.level + 1,
					path:        nextPath,
					bound:       nextBound,
					currentCost: nextCost,
					parent:      currentNode,
					matrix:      nextMatrix,
				}
				leafNodes = append(leafNodes, nextNode)
				treeNodes = append(treeNodes, nextNode)
			}
		}
	}

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

func main() {
	graph := [][]int{
		{math.MaxInt, 4, 9, 5},
		{6, math.MaxInt, 4, 8},
		{9, 4, math.MaxInt, 9},
		{5, 8, 9, math.MaxInt},
	}

	bestCost, bestPath := branchAndBound(graph)

	fmt.Printf("Minimum cost: %d\n", bestCost)
	fmt.Printf("Best path: %v\n", bestPath)
}

func RemoveCellFromMatrixByIdx(mx [][]int, idxRow int, idxCol int) [][]int {
	mt := RemoveRowFromMatrixByIndex(mx, idxRow)
	resultMx := RemoveColFromMatrixByIndex(mt, idxCol)
	return resultMx
}

func RemoveRowFromMatrixByIndex(mx [][]int, nameIndex int) [][]int {
	lenRows := len(mx)
	var resultMx [][]int
	for i := 0; i < lenRows-1; i++ {
		if i < nameIndex {
			resultMx = append(resultMx, mx[i])
		} else {
			resultMx = append(resultMx, mx[i+1])
		}
	}
	return resultMx
}

func IdxByName(m [][]int, rowName, colName int) (rowIdx, colIdx int, ok bool) {
	for i, v := range m {
		if v[0] == rowName {
			rowIdx = i
			break
		}
	}
	if rowIdx == 0 {
		return 0, 0, false
	}
	for j, v := range m[0] {
		if v == colName {
			colIdx = j
			break
		}
	}
	if colIdx == 0 {
		return 0, 0, false
	}
	return rowIdx, colIdx, true
}

func RemoveColFromMatrixByIndex(mx [][]int, nameIndex int) [][]int {
	lenRows := len(mx)
	var resultMx [][]int
	for i := 0; i < lenRows; i++ {
		resultMx = append(resultMx, mx[i][:nameIndex])
		resultMx[i] = append(resultMx[i], mx[i][nameIndex+1:]...)
	}
	return resultMx
}

func RemoveCellFromMatrixByName(mx [][]int, nameRow int, nameCol int) [][]int {
	idxRow, idxCol, ok := IdxByName(mx, nameRow, nameCol)
	if !ok {
		log.Println("Первый: не могу получить индексы из имени !!!")
	}
	mt := RemoveRowFromMatrixByIndex(mx, idxRow)
	resultMx := RemoveColFromMatrixByIndex(mt, idxCol)
	return resultMx
}
