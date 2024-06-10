// src/cmd/main/main.go
package main

import (
	"fmt"
	"forecast/internal/logic"
	"forecast/internal/types"
	"forecast/internal/utils"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Weather prediction for Dago, Bandung!\n\n")

	fmt.Print("Getting the last 14 days weather...")
	data := get14DaysWeather()
	fmt.Println(" Done!")

	fmt.Print("Transforming the data into transition matrix...")
	transition_matrix := utils.GetMatrix(data)
	fmt.Println(" Done!")

	fmt.Print("Getting the current weather...")
	current_weather := utils.GetCurrentWeather()
	current_weather_index := utils.GetIndex(current_weather.WeatherDescription)
	fmt.Println(" Done!\n\n")

	fmt.Println("Select your method! ")
	fmt.Print("Naive or Optimized (n/o)?")

	var choice string
	fmt.Scanln(&choice)
	for choice != "n" && choice != "o" {
		fmt.Println("Invalid choice, please input n or o")
		fmt.Scanln(&choice)
	}

	var probs [][]float64
	fmt.Print("OK! ")
	if choice == "n" {
		fmt.Println("Using Naive Method\n\n")
		naive := logic.NaiveCalculator{}
		probs = naive.GetProbability(transition_matrix, 4, current_weather_index, 4)
	} else {
		fmt.Println("Using Optimized Method\n\n")
		optimized := logic.OptimizedCalculator{}
		probs = optimized.GetProbability(transition_matrix, 4, current_weather_index, 4)
	}

	fmt.Println("From our calculation, our prediction is")
	for i := range probs {
		highest, second := utils.GetFirstAndSecondHighestIndex(probs[i])
		fmt.Print("In the next ", (i+1)*6, " hour, ")
		fmt.Println("It will be", utils.GetWeather(highest), "or", utils.GetWeather(second))
	}

}

func getApiKey() string {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the API key from environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in .env file")
	}

	return apiKey
}

func get14DaysWeather() []types.HistoricalDatum {
	apiKey := getApiKey()
	timeNow := time.Now().Format("2006-01-02:15")
	from := utils.GetDaysWithHourAfterNumDays(timeNow, -14)

	var bytes []byte = utils.FetchData(apiKey, from, timeNow)
	var data types.DescriptionOnlyData
	utils.Parse(bytes, &data)

	// categorize the weather

	var ret []types.HistoricalDatum
	for i := range data.Data {
		hour := strings.Split(data.Data[i].Datetime, ":")[1]
		intHour, err := strconv.Atoi(hour)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return nil
		}

		if (intHour % 6) == 0 {
			var current types.HistoricalDatum
			current.Datetime = data.Data[i].Datetime
			current.WeatherDescription = data.Data[i].Weather.Description
			ret = append(ret, current)
		}
	}

	for i := range ret {
		ret[i].WeatherDescription = utils.Categorize(ret[i].WeatherDescription)
	}

	return ret
}
