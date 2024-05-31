// src/cmd/main/main.go
package main

import (
	"encoding/csv"
	"fmt"
	"forecast/internal/logic"
	"log"
	"os"
	"strconv"
)

func main() {

	// Get transition matrix from "../data/dago-matrix.csv"

	var matrix [][]float64
	matrix = getMatrix("../data/dago-matrix.csv")

	naiveCalculator := logic.NaiveCalculator{}
	result := naiveCalculator.GetProbability(matrix, 1, 3)
	fmt.Println(result)
}

func getMatrix(path string) [][]float64 {
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	// Convert the records into a 2D slice of integers
	matrix := make([][]float64, len(records))
	for i, record := range records {
		matrix[i] = make([]float64, len(record))
		for j, value := range record {
			matrix[i][j], err = strconv.ParseFloat(value, 64)
			if err != nil {
				log.Fatalf("failed to convert string to int: %s", err)
			}
		}
	}

	return matrix
}
