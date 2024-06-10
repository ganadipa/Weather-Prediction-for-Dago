package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"forecast/internal/logic"
	"forecast/internal/types"
	"forecast/internal/utils"
	"math"
	"math/big"
	"os"
	"path/filepath"
)

func main() {
	// fmt.Println("Stress testing here we go!")

	var number int
	fmt.Scan(&number)

	// Get the current working directory
	cwd, _ := os.Getwd()

	// Get the data
	data_path := filepath.Join(cwd, "../test/categorized_dago.csv")
	var data []types.HistoricalDatum = parseCSV(data_path)

	// train data length
	days_for_training := 14
	train_data_length := days_for_training * 4

	set := make(map[int]bool)
	var correct []int = make([]int, 4)
	var total []int = make([]int, 4)

	for i := 1; i <= 10000; i++ {
		random_int, _ := rand.Int(rand.Reader, big.NewInt(int64(len(data)-train_data_length-number)))
		for set[int(random_int.Int64())] {
			random_int, _ = rand.Int(rand.Reader, big.NewInt(int64(len(data)-train_data_length-number)))
		}

		set[int(random_int.Int64())] = true
		rd := int(random_int.Int64()) + number

		// get the train data then reverse it because the data in the csv is reversed.
		train_data := data[rd : rd+train_data_length]
		for i, j := 0, len(train_data)-1; i < j; i, j = i+1, j-1 {
			train_data[i], train_data[j] = train_data[j], train_data[i]
		}

		// Make the transition matrix
		matrix := utils.GetMatrix(train_data)

		// get the current weather
		current_weather := utils.GetIndex(data[rd].WeatherDescription)

		// using naive
		naive := logic.NaiveCalculator{}
		naive_probability := naive.GetProbability(matrix, 4, current_weather, 1)[0]

		// using optimized
		optimized := logic.OptimizedCalculator{}
		opt_probability := optimized.GetProbability(matrix, 4, current_weather, 1)[0]

		// compare the results
		prediction_weather := utils.GetIndex(data[rd-number].WeatherDescription)

		if math.Abs(naive_probability[0]-opt_probability[0]) > 1e-9 {
			fmt.Println("Different results!")
			fmt.Println("Naive: ", naive_probability)
			fmt.Println("Optimized: ", opt_probability)
			return
		} else {
			// fmt.Println("Same results!")
		}

		// get the highest andthe second highest probability index
		first, second := utils.GetFirstAndSecondHighestIndex(naive_probability)

		total[first]++
		total[second]++
		if prediction_weather == first || prediction_weather == second {
			// fmt.Println("Correct!")
			correct[first]++
			correct[second]++
		} else {
			// fmt.Println("Incorrect!")
		}

	}

	for i := 0; i < 4; i++ {
		// accuracy
		fmt.Println("Accuracy for ", utils.GetWeather(i), " : ", float64(correct[i])/float64(total[i]))
	}

}

func parseCSV(data_path string) []types.HistoricalDatum {
	var historicalData []types.HistoricalDatum

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

		historicalDatum := types.HistoricalDatum{
			Datetime:           record[0],
			WeatherDescription: record[1],
		}

		historicalData = append(historicalData, historicalDatum)
	}

	return historicalData
}
