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

// func branchAndBound(graph [][]int) (int, []int) {
// 	n := len(graph)
// 	_, initialLowerBound := reduceMatrix(graph)

// 	pq := &PriorityQueue{}
// 	heap.Init(pq)

// 	initialNode := &Node{
// 		level:       0,
// 		path:        []int{0},
// 		bound:       initialLowerBound,
// 		currentCost: 0,
// 	}
// 	heap.Push(pq, initialNode)

// 	bestCost := math.MaxInt
// 	var bestPath []int

// 	for pq.Len() > 0 {
// 		currentNode := heap.Pop(pq).(*Node)

// 		if currentNode.bound >= bestCost {
// 			continue
// 		}

// 		if currentNode.level == n-1 {
// 			finalCost := currentNode.currentCost
// 			if finalCost < bestCost {
// 				bestCost = finalCost
// 				bestPath = append([]int{}, currentNode.path...)
// 			}
// 			continue
// 		}

// 		currentCity := currentNode.path[len(currentNode.path)-1]
// 		for nextCity := 0; nextCity < n; nextCity++ {
// 			if contains(currentNode.path, nextCity) {
// 				continue
// 			}

// 			nextPath := append([]int{}, currentNode.path...)
// 			nextPath = append(nextPath, nextCity)
// 			nextCost := currentNode.currentCost + graph[currentCity][nextCity]
// 			nextBound := calculateBound(graph, nextPath)

// 			if nextBound < bestCost {
// 				nextNode := &Node{
// 					level:       currentNode.level + 1,
// 					path:        nextPath,
// 					bound:       nextBound,
// 					currentCost: nextCost,
// 				}
// 				heap.Push(pq, nextNode)
// 			}
// 		}
// 	}

// 	return bestCost, bestPath
// }

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

var Inf int = math.MaxInt

func main() {
	// graph := [][]int{
	// 	{Inf, 47, 22, 46, 29},
	// 	{34, Inf, 25, 34, 19},
	// 	{18, 18, Inf, 33, 7},
	// 	{38, 27, 24, Inf, 38},
	// 	{21, 14, 6, 27, Inf},
	// }
	// graph := [][]int{
	// 	{Inf, 1, 2, 3, 4},
	// 	{14, Inf, 15, 16, 5},
	// 	{13, 20, Inf, 17, 6},
	// 	{12, 19, 18, Inf, 7},
	// 	{11, 10, 9, 8, Inf},
	// }
	// graph := [][]int{
	// 	{Inf, 1, 2, 3, 4},
	// 	{0, Inf, 15, 16, 5},
	// 	{0, 20, Inf, 17, 6},
	// 	{0, 19, 18, Inf, 7},
	// 	{0, 10, 9, 8, Inf},
	// }
	// graph := [][]int{
	// 	{Inf, 0, 22, 46, 29},
	// 	{34, Inf, 25, 34, 19},
	// 	{18, 0, Inf, 33, 7},
	// 	{38, 0, 24, Inf, 38},
	// 	{21, 0, 6, 27, Inf},
	// }

	// 12
	// graph := [][]int{
	// 	{Inf, 755, 1701, 1866, 1675, 2424, 2652, 2707, 2522, 2732, 3497, 3763},
	// 	{712, Inf, 1578, 1581, 1390, 2043, 2855, 2910, 2237, 2906, 3212, 3142},
	// 	{1985, 1912, Inf, 542, 726, 1460, 2191, 2246, 1573, 1783, 2548, 3168},
	// 	{1987, 2028, 542, Inf, 729, 1644, 2194, 2249, 1576, 1785, 2550, 3171},
	// 	{1797, 1838, 727, 729, Inf, 736, 1548, 1603, 1051, 1260, 2025, 2646},
	// 	{2316, 2357, 1246, 1248, 603, Inf, 1146, 1201, 1500, 1710, 2309, 2930},
	// 	{3071, 2867, 2125, 2783, 2137, 1166, Inf, 134, 2007, 2431, 2997, 3617},
	// 	{3920, 3716, 2749, 2751, 2106, 2015, 134, Inf, 1975, 2400, 2965, 3586},
	// 	{2644, 2685, 1573, 1576, 1051, 1634, 2563, 2618, Inf, 531, 1296, 1916},
	// 	{3283, 3323, 1783, 1785, 1260, 1843, 3037, 3092, 531, Inf, 859, 1479},
	// 	{2580, 2621, 2854, 2857, 2331, 4186, 3932, 3987, 1602, 1165, Inf, 777},
	// 	{3416, 3457, 3291, 3294, 2768, 4623, 4369, 4424, 2039, 1602, 2506, Inf},
	// }

	// 13
	// graph := [][]int{
	// 	{Inf, 75987, 76499, 76503, 76657, 76726, 76554, 76917, 77002, 77221, 77562, 77770, 77774},
	// 	{75627, Inf, 909, 913, 1067, 1136, 964, 1328, 1412, 1631, 1972, 2181, 2184, 2813, 2849},
	// 	{76194, 892, Inf, 1141, 1296, 1364, 1192, 1556, 1641, 1859, 2201, 2409, 2412, 3041},
	// 	{76204, 902, 1080, Inf, 299, 368, 550, 913, 998, 1217, 1558, 1767, 1770, 2399},
	// 	{76359, 1057, 1234, 299, Inf, 173, 704, 1068, 1153, 1371, 1713, 1921, 1925},
	// 	{76427, 1125, 1303, 368, 154, Inf, 773, 1137, 1221, 1440, 1781, 1990, 1993},
	// 	{76209, 907, 1084, 623, 777, 846, Inf, 517, 602, 821, 1162, 1371, 1374},
	// 	{76572, 1270, 1448, 987, 1141, 1210, 517, Inf, 271, 490, 666, 874, 877},
	// 	{76657, 1355, 1533, 1072, 1226, 1295, 602, 271, Inf, 218, 560, 768, 771},
	// 	{76876, 1574, 1751, 1290, 1444, 1513, 821, 490, 218, Inf, 341, 549, 553},
	// 	{77217, 1915, 2093, 1632, 1786, 1855, 1162, 666, 560, 341, Inf, 208, 211},
	// 	{77425, 2123, 2301, 1840, 1994, 2513, 1371, 874, 768, 549, 208, Inf, 3},
	// 	{77429, 2127, 2305, 1843, 1998, 2516, 1374, 877, 771, 553, 211, 3, Inf},
	// }

	// 14
	// graph := [][]int{
	// 	{Inf, 75987, 76499, 76503, 76657, 76726, 76554, 76917, 77002, 77221, 77562, 77770, 77774, 78402},
	// 	{75627, Inf, 909, 913, 1067, 1136, 964, 1328, 1412, 1631, 1972, 2181, 2184, 2813, 2849, 2978},
	// 	{76194, 892, Inf, 1141, 1296, 1364, 1192, 1556, 1641, 1859, 2201, 2409, 2412, 3041, 3077},
	// 	{76204, 902, 1080, Inf, 299, 368, 550, 913, 998, 1217, 1558, 1767, 1770, 2399, 2435},
	// 	{76359, 1057, 1234, 299, Inf, 173, 704, 1068, 1153, 1371, 1713, 1921, 1925, 2553},
	// 	{76427, 1125, 1303, 368, 154, Inf, 773, 1137, 1221, 1440, 1781, 1990, 1993, 2622},
	// 	{76209, 907, 1084, 623, 777, 846, Inf, 517, 602, 821, 1162, 1371, 1374, 2003},
	// 	{76572, 1270, 1448, 987, 1141, 1210, 517, Inf, 271, 490, 666, 874, 877, 1506},
	// 	{76657, 1355, 1533, 1072, 1226, 1295, 602, 271, Inf, 218, 560, 768, 771, 1400},
	// 	{76876, 1574, 1751, 1290, 1444, 1513, 821, 490, 218, Inf, 341, 549, 553, 1181},
	// 	{77217, 1915, 2093, 1632, 1786, 1855, 1162, 666, 560, 341, Inf, 208, 211, 840},
	// 	{77425, 2123, 2301, 1840, 1994, 2513, 1371, 874, 768, 549, 208, Inf, 3, 707},
	// 	{77429, 2127, 2305, 1843, 1998, 2516, 1374, 877, 771, 553, 211, 3, Inf, 711},
	// 	{78007, 2704, 2882, 2421, 2575, 3094, 1952, 1455, 1349, 1130, 789, 656, 660, Inf},
	// }

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

	// graph := [][]int{
	// 	{Inf, 75987, 76499, 76503, 76657, 76726, 76554, 76917, 77002, 77221, 77562, 77770, 77774, 78402, 78439, 78567, 78641, 78925, 77208, 77613, 77883, 78212, 88377, 88466, 88466},
	// 	{75627, Inf, 909, 913, 1067, 1136, 964, 1328, 1412, 1631, 1972, 2181, 2184, 2813, 2849, 2978, 3051, 3335, 1618, 2023, 2294, 2623, 12556, 12645, 12645},
	// 	{76194, 892, Inf, 1141, 1296, 1364, 1192, 1556, 1641, 1859, 2201, 2409, 2412, 3041, 3077, 3206, 3279, 3563, 1846, 2251, 2522, 2851, 12784, 12873, 12873},
	// 	{76204, 902, 1080, Inf, 299, 368, 550, 913, 998, 1217, 1558, 1767, 1770, 2399, 2435, 2563, 2637, 2921, 1204, 1609, 1879, 1854, 12142, 12231, 12231},
	// 	{76359, 1057, 1234, 299, Inf, 173, 704, 1068, 1153, 1371, 1713, 1921, 1925, 2553, 2589, 2718, 2791, 3076, 1358, 1763, 1355, 1659, 12296, 12385, 12385},
	// 	{76427, 1125, 1303, 368, 154, Inf, 773, 1137, 1221, 1440, 1781, 1990, 1993, 2622, 2658, 2787, 2860, 3144, 1427, 1612, 1181, 1486, 12365, 12454, 12454},
	// 	{76209, 907, 1084, 623, 777, 846, Inf, 517, 602, 821, 1162, 1371, 1374, 2003, 2039, 2167, 2241, 2525, 808, 1059, 1329, 2459, 11746, 11835, 11835},
	// 	{76572, 1270, 1448, 987, 1141, 1210, 517, Inf, 271, 490, 666, 874, 877, 1506, 1542, 1671, 1744, 2028, 477, 783, 1053, 2184, 11249, 11338, 11338},
	// 	{76657, 1355, 1533, 1072, 1226, 1295, 602, 271, Inf, 218, 560, 768, 771, 1400, 1436, 1565, 1638, 1922, 271, 884, 1155, 2285, 11143, 11232, 11232},
	// 	{76876, 1574, 1751, 1290, 1444, 1513, 821, 490, 218, Inf, 341, 549, 553, 1181, 1218, 1346, 1420, 1704, 490, 1102, 1373, 2503, 10925, 11013, 11013},
	// 	{77217, 1915, 2093, 1632, 1786, 1855, 1162, 666, 560, 341, Inf, 208, 211, 840, 876, 1005, 1078, 1362, 831, 1278, 1428, 2713, 10583, 10672, 10672},
	// 	{77425, 2123, 2301, 1840, 1994, 2513, 1371, 874, 768, 549, 208, Inf, 3, 707, 743, 872, 946, 1230, 1039, 1489, 1344, 2629, 10450, 10539, 10539},
	// 	{77429, 2127, 2305, 1843, 1998, 2516, 1374, 877, 771, 553, 211, 3, Inf, 711, 747, 875, 949, 1233, 1043, 1493, 1348, 2633, 10454, 10543, 10543},
	// 	{78007, 2704, 2882, 2421, 2575, 3094, 1952, 1455, 1349, 1130, 789, 656, 660, Inf, 117, 336, 364, 603, 1620, 2070, 1926, 3211, 9824, 9913, 9913},
	// 	{78094, 2792, 2969, 2508, 2662, 3181, 2039, 1542, 1436, 1218, 876, 743, 747, 117, Inf, 424, 481, 559, 1708, 2158, 2013, 3298, 9780, 9869, 9869},
	// 	{78222, 2920, 3098, 2637, 2791, 3310, 2167, 1671, 1565, 1346, 1005, 872, 875, 387, 424, Inf, 626, 910, 1836, 2286, 2141, 3426, 10131, 10220, 10220},
	// 	{78296, 2994, 3172, 2710, 2865, 3052, 2241, 1744, 1638, 1420, 1078, 946, 949, 364, 481, 626, Inf, 967, 1910, 2029, 1884, 3169, 10188, 10277, 10277},
	// 	{78580, 3278, 3456, 2994, 3149, 3667, 2525, 2028, 1922, 1704, 1362, 1230, 1233, 603, 559, 910, 967, Inf, 2194, 2644, 2499, 3784, 9434, 9523, 9523},
	// 	{76863, 1561, 1738, 1277, 1431, 1500, 808, 477, 271, 490, 831, 1039, 1043, 1671, 1708, 1836, 1910, 2194, Inf, 1090, 1360, 2490, 11415, 11504, 11504},
	// 	{77268, 1966, 2143, 1937, 1723, 1569, 1059, 783, 884, 1102, 1278, 1487, 1490, 2120, 2157, 2285, 2254, 2643, 1090, Inf, 555, 1686, 11864, 11952, 11952},
	// 	{77566, 2264, 2442, 1507, 1293, 1139, 1329, 1053, 1155, 1373, 1329, 1537, 1540, 2120, 2157, 2285, 2254, 2643, 1360, 555, Inf, 1255, 11864, 11953, 11953},
	// 	{77914, 2612, 2789, 1854, 1640, 1486, 2259, 2184, 2285, 2503, 2679, 2887, 2891, 3521, 3557, 3686, 3655, 4044, 2490, 1686, 1255, Inf, 13264, 13353, 13353},
	// 	{88028, 12499, 12676, 12215, 12369, 12888, 11746, 11249, 11143, 10925, 10583, 10450, 10454, 9824, 9780, 10131, 10188, 9434, 11415, 11865, 11720, 13005, Inf, 126, 126},
	// 	{88117, 12588, 12765, 12304, 12458, 12977, 11835, 11338, 11232, 11013, 10672, 10539, 10543, 9913, 9869, 10220, 10277, 9523, 11504, 11953, 11809, 13094, 126, Inf, 0},
	// 	{88117, 12588, 12765, 12304, 12458, 12977, 11835, 11338, 11232, 11013, 10672, 10539, 10543, 9913, 9869, 10220, 10277, 9523, 11504, 11953, 11809, 13094, 126, 0, Inf},
	// }

	maxIterations := 10000000 // Укажите максимальное количество итераций

	t := time.Now()
	bestCost, bestPath := branchAndBound(graph, maxIterations)
	ts := time.Since(t)
	fmt.Printf("Time: %v\n", ts)
	fmt.Printf("Minimum cost: %d\n", bestCost)
	fmt.Printf("Best path: %v\n", bestPath)
	fmt.Printf("Calc path: %v\n", calculateDistance(graph, bestPath))
	fmt.Println()

	//_ = graph12
	// t := time.Now()
	// bestCost, bestPath := branchAndBound(graph12)
	// ts := time.Since(t)
	// fmt.Printf("Time: %v\n", ts)
	// fmt.Printf("Minimum cost: %d\n", bestCost)
	// fmt.Printf("Best path: %v\n", bestPath)
	// fmt.Printf("Calc path: %v\n", calculateDistance(graph12, bestPath))
	// fmt.Println()

	// _ = graph13
	// t = time.Now()
	// bestCost, bestPath = branchAndBound(graph13)
	// ts = time.Since(t)
	// fmt.Printf("Time: %v\n", ts)
	// fmt.Printf("Minimum cost: %d\n", bestCost)
	// fmt.Printf("Best path: %v\n", bestPath)
	// fmt.Printf("Calc path: %v\n", calculateDistance(graph13, bestPath))
	// fmt.Println()

	//	_ = graph14
	// t = time.Now()
	// bestCost, bestPath = branchAndBound(graph14)
	// ts = time.Since(t)
	// fmt.Printf("Time: %v\n", ts)
	// fmt.Printf("Minimum cost: %d\n", bestCost)
	// fmt.Printf("Best path: %v\n", bestPath)
	// fmt.Printf("Calc path: %v\n", calculateDistance(graph14, bestPath))
	// fmt.Println()

	// _ = graph15
	// t = time.Now()
	// bestCost, bestPath = branchAndBound(graph15)
	// ts = time.Since(t)
	// fmt.Printf("Time: %v\n", ts)
	// fmt.Printf("Minimum cost: %d\n", bestCost)
	// fmt.Printf("Best path: %v\n", bestPath)
	// fmt.Printf("Calc path: %v\n", calculateDistance(graph15, bestPath))
	// fmt.Println()

}

func calculateDistance(graph [][]int, path []int) int {
	n := len(path)
	distance := 0
	for i := 0; i < n-1; i++ {
		distance += graph[path[i]][path[i+1]]
	}
	return distance
}
