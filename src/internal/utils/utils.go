package utils

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"forecast/internal/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func GetLastDayFromDataset() string {
	defaultDate := "2024-05-25"
	filePath := "../data/dago-complete.csv"

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return defaultDate
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return defaultDate
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(bufio.NewReader(file))

	var lastDate string
	for {
		// Read the next record
		record, err := reader.Read()

		if err != nil {
			// Check for end of file
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading CSV record:", err)
			return defaultDate
		}

		// Update the lastDate with the datetime column value
		if len(record) > 0 {
			lastDate = record[0]
		}
	}

	if lastDate == "" {
		return defaultDate
	}

	lastDate = strings.Split(lastDate, ":")[0]

	// Return the last datetime value
	return lastDate
}

func GetDaysAfterNum(date string, num int) string {
	// Parse the given date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}

	// Subtract 10 days from the parsed date
	newDate := parsedDate.AddDate(0, 0, num)

	// Format the new date back to string
	return newDate.Format("2006-01-02")
}

func GetDaysWithHourAfterNumDays(datetime string, num int) string {
	// Parse the given date
	parsedDate, err := time.Parse("2006-01-02:15", datetime)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}

	// Subtract 10 days from the parsed date
	newDate := parsedDate.AddDate(0, 0, num)

	// Format the new date back to string
	return newDate.Format("2006-01-02:15")
}

func AppendWeatherResponseToDataset(weatherResponse types.WeatherResponse) {
	dirPath := "../data"
	filePath := filepath.Join(dirPath, "dago-complete.csv")

	// Check if the directory exists, if not, create it
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	// Check if the file exists, if not, create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		file.Close()
	}

	// Open the file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Append each Data record from the WeatherResponse in reverse order
	for i := len(weatherResponse.Data) - 1; i >= 0; i-- {
		data := weatherResponse.Data[i]
		if timeToInt(data.Datetime)%6 != 0 {
			continue
		}

		record := []string{data.Datetime, data.Weather.Description,
			strconv.FormatFloat(data.Humidity, 'f', -1, 64),
			strconv.FormatFloat(data.Temperature, 'f', -1, 64),
			strconv.FormatFloat(data.Precipitation, 'f', -1, 64),
			strconv.FormatFloat(data.WindSpeed, 'f', -1, 64),
		}
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record to file:", err)
		}
	}
}

func FetchData(apiKey, fromDate, toDate string) []byte {
	// Dago, Bandung latitute and longitude
	lat := "-6.8852"
	lon := "107.6136"

	url := fmt.Sprintf("https://api.weatherbit.io/v2.0/history/hourly?lat=%s&lon=%s&start_date=%s&end_date=%s&tz=local&key=%s", lat, lon, fromDate, toDate, apiKey)
	// fmt.Println(url)
	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	return body
}

func Parse(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
}

func timeToInt(time string) int {
	hour := strings.Split(time, ":")[1]
	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return 0
	}
	return hourInt
}

func GetMatrix(data []types.HistoricalDatum) [][]float64 {
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

func GetIndex(weather string) int {
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

func GetWeather(index int) string {
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

func GetFirstAndSecondHighestIndex(probability []float64) (int, int) {
	highest_index := 0
	second_highest_index := 0
	for i := 1; i < 4; i++ {
		if probability[i] > probability[highest_index] {
			second_highest_index = highest_index
			highest_index = i
		} else if probability[i] > probability[second_highest_index] {
			second_highest_index = i
		}
	}
	return highest_index, second_highest_index
}

func GetCurrentWeather() types.HistoricalDatum {
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

	lat := "-6.8852"
	lon := "107.6136"

	// Get the weather data from the API
	url := fmt.Sprintf("https://api.weatherbit.io/v2.0/current?lat=%s&lon=%s&key=%s", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Unmarshal the JSON response
	var weatherResponse types.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	var ret types.HistoricalDatum
	ret.Datetime = weatherResponse.Data[0].Datetime
	ret.WeatherDescription = Categorize(weatherResponse.Data[0].Weather.Description)

	return ret
}

func Categorize(weather string) string {
	categorizer := map[string]string{
		"Clear Sky":                    "Terang",
		"Few clouds":                   "Terang",
		"Scattered clouds":             "Berawan",
		"Broken clouds":                "Berawan",
		"Overcast clouds":              "Mendung",
		"Fog":                          "Mendung",
		"Haze":                         "Mendung",
		"Light rain":                   "Hujan",
		"Moderate rain":                "Hujan",
		"Heavy rain":                   "Hujan",
		"Thunderstorm with heavy rain": "Hujan",
	}

	return categorizer[weather]
}
