package logic

type OptimizedCalculator struct{}

func (o OptimizedCalculator) GetProbability(transition_square_matrix [][]float64, at_period int, initial_state int, num_moment int) [][]float64 {
	// Initialize the probability vector
	var probability [][]float64 = make([][]float64, 4)
	for i := range probability {
		probability[i] = make([]float64, len(transition_square_matrix))
	}

	// Using optimized dp w/ divide and conquer approach (matrix exponentiation)

	// Get the result of the matrix exponentiation
	var result [][]float64 = matrixExponentiation(transition_square_matrix, at_period)

	// Get the probability vector
	probability[0] = result[initial_state]
	for i := 1; i < num_moment; i++ {
		result := multiplySquareMatrix(result, transition_square_matrix)[initial_state]
		probability[i] = result
	}

	return probability
}

func matrixExponentiation(matrix [][]float64, power int) [][]float64 {
	if power == 1 {
		return matrix
	}

	sqrt := matrixExponentiation(matrix, power/2)
	if power%2 == 0 {
		return multiplySquareMatrix(sqrt, sqrt)
	} else {
		return multiplySquareMatrix(matrix, multiplySquareMatrix(sqrt, sqrt))
	}
}

func multiplySquareMatrix(matrix1 [][]float64, matrix2 [][]float64) [][]float64 {
	n := len(matrix1)
	result := make([][]float64, n)
	for i := 0; i < n; i++ {
		result[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				result[i][j] += matrix1[i][k] * matrix2[k][j]
			}
		}
	}
	return result
}
