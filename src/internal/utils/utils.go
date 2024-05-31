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
	fmt.Println(url)
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
