package logic

type NaiveCalculator struct{}

func (n NaiveCalculator) GetProbability(transition_square_matrix [][]float64, at_period int, initial_state int, num_moment int) [][]float64 {
	if num_moment <= 0 {
		return nil
	}

	// Initialize the probability vector
	var probability [][]float64 = make([][]float64, num_moment)
	for i := range probability {
		probability[i] = make([]float64, len(transition_square_matrix))
	}

	// Using naive dp
	var dp [][]float64 = make([][]float64, len(transition_square_matrix))
	for i := range dp {
		dp[i] = make([]float64, at_period+num_moment)
	}

	for i := 0; i < len(transition_square_matrix); i++ {
		dp[i][0] = 0
	}

	dp[initial_state][0] = 1

	for i := 1; i < at_period+num_moment; i++ {
		for j := 0; j < len(transition_square_matrix); j++ {
			for k := 0; k < len(transition_square_matrix); k++ {
				dp[j][i] += dp[k][i-1] * transition_square_matrix[k][j]
			}
		}
	}

	for i := 0; i < num_moment; i++ {
		for j := 0; j < len(transition_square_matrix); j++ {
			probability[i][j] = dp[j][at_period+i]
		}
	}

	return probability

}
