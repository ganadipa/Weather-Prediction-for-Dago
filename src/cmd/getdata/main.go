package main

import (
	"fmt"
	"forecast/internal/types"
	"forecast/internal/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the API key from environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in .env file")
	}

	// Get the last day int the data/dago.csv
	// If the file doesn't exist, set the last day to "2024-05-25"
	toDate := utils.GetLastDayFromDataset()
	toDate = utils.GetDaysAfterNum(toDate, 1)
	fromDate := utils.GetDaysAfterNum(toDate, -30)

	fmt.Println("Fetching data from", fromDate, "to", toDate)

	res := utils.FetchData(apiKey, fromDate, toDate)

	// Unmarshal the JSON response
	var weatherResponse types.WeatherResponse
	utils.Parse(res, &weatherResponse)

	// Write at data/dago.csv
	utils.AppendWeatherResponseToDataset(weatherResponse)
}
