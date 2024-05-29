package logic

func GetProbability(transition_square_matrix [][]float64, at_period int, initial_state int) []float64 {
	// Initialize the probability vector
	var probability []float64
	probability = make([]float64, len(transition_square_matrix))

	// Using naive dp
	var dp [][]float64
	dp = make([][]float64, len(transition_square_matrix))
	for i := range dp {
		dp[i] = make([]float64, at_period+1)
	}

	for i := 0; i < len(transition_square_matrix); i++ {
		dp[i][0] = transition_square_matrix[initial_state][i]
	}

	for i := 1; i <= at_period; i++ {
		for j := 0; j < len(transition_square_matrix); j++ {
			for k := 0; k < len(transition_square_matrix); k++ {
				dp[j][i] += dp[k][i-1] * transition_square_matrix[k][j]
			}
		}
	}

	for i := 0; i < len(transition_square_matrix); i++ {
		probability[i] = dp[i][at_period]
	}

	return probability
}
