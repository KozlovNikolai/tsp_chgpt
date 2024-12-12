package main

import (
	"fmt"
	"tspchgpt/var5/pkg"
)

func branchAndBound(graph [][]int) (int, []int) {
	matrixNamed := pkg.SetNaming(graph)
	//	pkg.PrintMatrix(matrixNamed)
	initialLowerBound := pkg.CalculateCost(matrixNamed)
	r := pkg.NewRepo(&pkg.Node{
		ID:       0,
		Level:    0,
		Path:     []int{},
		Cost:     initialLowerBound,
		ParentID: 0,
		Matrix:   matrixNamed,
	})

	r.CreateLeaves()

	fmt.Printf("Current Node ID: %d, Next ID: %d\n", r.CurrentNodeID, r.NextID)
	for _, value := range r.Store {
		pkg.PrintNode(value)
	}

	//fmt.Printf("%+v\n", r.Store[0].Matrix)
	bestCost := 0
	bestPath := []int{0}
	return bestCost, bestPath
}

func main() {
	graph := pkg.Graph

	bestCost, bestPath := branchAndBound(graph)

	fmt.Printf("Minimum cost: %d\n", bestCost)
	fmt.Printf("Best path: %v\n", bestPath)
}
