// src/cmd/getdata/main.go
package main

import (
	"fmt"
	"forecast/internal"
)

func main() {
	fmt.Println("Start fetching data...")
	internal.FetchData()
}
