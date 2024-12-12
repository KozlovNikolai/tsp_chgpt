package pkg

import (
	"fmt"
)

// func CalculateCost(graph [][]int) int {
// 	lowerBound := 0
// 	n := len(graph)
// 	// Reduce rows
// 	for i := 0; i < n; i++ {
// 		minValue := math.MaxInt
// 		for j := 0; j < n; j++ {
// 			if graph[i][j] < minValue {
// 				minValue = graph[i][j]
// 			}
// 		}
// 		if minValue != math.MaxInt {
// 			lowerBound += minValue
// 		}
// 	}
// 	// PrintMatrix(graph)
// 	// fmt.Printf("initial lower bound: %d\n", lowerBound)
// 	return lowerBound
// }

func PrintMatrix(mx [][]int) {
	for i := 0; i < len(mx); i++ {
		for j := 0; j < len(mx[0]); j++ {
			if mx[i][j] == INF {
				fmt.Printf("\tInf\t")
			} else {
				fmt.Printf("\t%d\t", mx[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func CloneMx(mx [][]int) [][]int {
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

func RemoveCellFromMatrixByIdx(matrix [][]int, row int, col int) [][]int {
	if row < 0 || row >= len(matrix) || col < 0 || col >= len(matrix[0]) {
		panic("Invalid row or column index")
	}

	// Инициализируем новый массив для результата
	result := make([][]int, 0, len(matrix)-1)

	for i := 0; i < len(matrix); i++ {
		if i == row {
			continue
		}
		newRow := make([]int, 0, len(matrix[0])-1)
		for j := 0; j < len(matrix[i]); j++ {
			if j == col {
				continue
			}
			newRow = append(newRow, matrix[i][j])
		}
		result = append(result, newRow)
	}
	return result
}

func PrintNode(node *Node) {
	fmt.Printf("id:%d, parentID:%d, level:%d, cost:%d\n", node.ID, node.ParentID, node.Level, node.Cost)
	fmt.Printf("path: %+v\n", node.Path)
	fmt.Println("Matrix:")
	PrintMatrix(node.Matrix)
	fmt.Println()
}

func SetNaming(mx [][]int) [][]int {
	lenRows := len(mx)
	lenCols := len(mx[0])

	mmx := make([][]int, lenRows+1)
	for i := range mmx {
		mmx[i] = make([]int, lenCols+1)
		// заполняем заголовки столбцов:
		if i == 0 {
			for j := range mmx[i] {
				mmx[i][j] = j
			}
		} else {
			mmx[i][0] = i
			for j := 1; j < len(mmx[i]); j++ {
				mmx[i][j] = mx[i-1][j-1]
			}
		}

	}
	return mmx
}
