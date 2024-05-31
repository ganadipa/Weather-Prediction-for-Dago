package logic

type OptimizedCalculator struct{}

func (o OptimizedCalculator) GetProbability(transition_square_matrix [][]float64, at_period int, initial_state int) []float64 {
	// Initialize the probability vector
	var probability []float64 = make([]float64, len(transition_square_matrix))

	// Using optimized dp w/ divide and conquer approach (matrix exponentiation)

	return probability
}
