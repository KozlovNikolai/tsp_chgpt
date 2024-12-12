package pkg

func CalculateCost(mx [][]int) int {
	rows := len(mx)
	cols := len(mx[0])
	var sum int
	// создаем результирующую матрицу на одну колонку больше
	resultMx := make([][]int, rows)
	for i := 0; i < rows; i++ {
		resultMx[i] = make([]int, cols+1)
	}
	// заполняем заголовки
	for i := 0; i < rows; i++ {
		resultMx[i][0] = mx[i][0]
	}
	for j := 0; j < cols; j++ {
		resultMx[0][j] = mx[0][j]
	}
	// идем по строкам исключая строку с заголовками и ищем минимумы в каждой строке
	for i := 1; i < rows; i++ {
		min := INF
		// идем по ячейкам строки исключая заголовок строки
		for j := 1; j < cols; j++ {
			// находим минимум в строке
			if mx[i][j] < INF {
				if mx[i][j] < min {
					min = mx[i][j]
				}
			}
		}
		// записываем результат в конец строки
		resultMx[i][cols] = min
		sum += min
	}
	return sum
}
