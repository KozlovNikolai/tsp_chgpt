package pkg

import "math"

var INF int = math.MaxInt / 2

var Graph = [][]int{
	{INF, 4, 9, 5},
	{6, INF, 4, 8},
	{9, 4, INF, 9},
	{5, 8, 9, INF},
}
