package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"forecast/internal/logic"
	"math"
	"math/big"
	"os"
	"path/filepath"
)

type HistoricalDatum struct {
	Datetime           string
	WeatherDescription string
}

func main() {
	fmt.Println("Stress testing here we go!")

	// parse the csv
	cwd, _ := os.Getwd()

	data_path := filepath.Join(cwd, "../test/categorized_dago.csv")
	var data []HistoricalDatum = parseCSV(data_path)

	// train_data_length
	days_for_training := 14
	train_data_length := days_for_training * 4

	set := make(map[int]bool)
	var correct []int = make([]int, 4)
	var total []int = make([]int, 4)

	for i := 1; i <= 100000; i++ {
		random_int, _ := rand.Int(rand.Reader, big.NewInt(int64(len(data)-train_data_length)))
		for set[int(random_int.Int64())] {
			random_int, _ = rand.Int(rand.Reader, big.NewInt(int64(len(data)-train_data_length)))
		}

		set[int(random_int.Int64())+train_data_length] = true

		// test_data := data[int(random_int.Int64()+40)]
		train_data := data[int(random_int.Int64()) : int(random_int.Int64())+train_data_length]

		// get the transition matrix from the training data
		matrix := getMatrix(train_data)
		// fmt.Println("Transition Matrix: ", matrix)

		current_weather := getIndex(data[int(random_int.Int64())+train_data_length-1].WeatherDescription)
		// get the probability vector from the test data
		// using naive
		naive := logic.NaiveCalculator{}
		naive_probability := naive.GetProbability(matrix, 4, current_weather)

		// using optimized
		optimized := logic.OptimizedCalculator{}
		opt_probability := optimized.GetProbability(matrix, 4, current_weather)

		// compare the results
		// fmt.Println("Naive: ", naive.GetProbability(matrix, 10, 0))
		// fmt.Println("Optimized: ", optimized.GetProbability(matrix, 10, 0))

		prediction_weather := getIndex(data[int(random_int.Int64())+train_data_length].WeatherDescription)

		if math.Abs(naive_probability[0]-opt_probability[0]) > 1e-9 {
			fmt.Println("Different results!")
			fmt.Println("Naive: ", naive_probability)
			fmt.Println("Optimized: ", opt_probability)
			return
		} else {
			fmt.Println("Same results!")
		}

		// get the highest andthe second highest probability index
		highest_index := 0
		second_highest_index := 0
		for i := 1; i < 4; i++ {
			if naive_probability[i] > naive_probability[highest_index] {
				second_highest_index = highest_index
				highest_index = i
			} else if naive_probability[i] > naive_probability[second_highest_index] {
				second_highest_index = i
			}
		}

		total[prediction_weather]++
		if prediction_weather == highest_index || prediction_weather == second_highest_index {
			// fmt.Println("Correct!")
			correct[prediction_weather]++
		} else {
			// fmt.Println("Incorrect!")
		}

	}

	for i := 0; i < 4; i++ {
		// accuracy
		fmt.Println("Accuracy for ", getWeather(i), " : ", float64(correct[i])/float64(total[i]))
	}

}

func parseCSV(data_path string) []HistoricalDatum {
	var historicalData []HistoricalDatum

	file, err := os.Open(data_path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return nil
	}

	// Assuming the CSV has a header row
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}

		historicalDatum := HistoricalDatum{
			Datetime:           record[0],
			WeatherDescription: record[1],
		}

		historicalData = append(historicalData, historicalDatum)
	}

	return historicalData
}

func getMatrix(data []HistoricalDatum) [][]float64 {
	// make the probability transition matrix

	matrix := make([][]float64, 4)
	for i := range matrix {
		matrix[i] = make([]float64, 4)
	}

	// 0 is Terang
	// 1 is Berawan
	// 2 is Mendung
	// 3 is Hujan
	for i := 0; i < len(data)-1; i++ {
		switch data[i].WeatherDescription {
		case "Terang":
			switch data[i+1].WeatherDescription {
			case "Terang":
				matrix[0][0]++
			case "Berawan":
				matrix[0][1]++
			case "Mendung":
				matrix[0][2]++
			case "Hujan":
				matrix[0][3]++
			}
		case "Berawan":
			switch data[i+1].WeatherDescription {
			case "Terang":
				matrix[1][0]++
			case "Berawan":
				matrix[1][1]++
			case "Mendung":
				matrix[1][2]++
			case "Hujan":
				matrix[1][3]++
			}
		case "Mendung":
			switch data[i+1].WeatherDescription {
			case "Terang":
				matrix[2][0]++
			case "Berawan":
				matrix[2][1]++
			case "Mendung":
				matrix[2][2]++
			case "Hujan":
				matrix[2][3]++
			}
		case "Hujan":
			switch data[i+1].WeatherDescription {
			case "Terang":
				matrix[3][0]++
			case "Berawan":
				matrix[3][1]++
			case "Mendung":
				matrix[3][2]++
			case "Hujan":
				matrix[3][3]++
			}
		}
	}

	for i := 0; i < 4; i++ {
		total := 0.0
		for j := 0; j < 4; j++ {
			total += matrix[i][j]
		}

		for j := 0; j < 4; j++ {
			if total != 0 {
				matrix[i][j] /= total
			} else {
				matrix[i][j] = 0
			}
		}
	}

	return matrix
}

func getIndex(weather string) int {
	switch weather {
	case "Terang":
		return 0
	case "Berawan":
		return 1
	case "Mendung":
		return 2
	case "Hujan":
		return 3
	}
	return -1
}

func getWeather(index int) string {
	switch index {
	case 0:
		return "Terang"
	case 1:
		return "Berawan"
	case 2:
		return "Mendung"
	case 3:
		return "Hujan"
	}
	return ""
}
