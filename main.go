package main

import (
	"fmt"
	"github.com/Seven11Eleven/golangkbtu2025/homeworks"
	"math"
)

func main() {
	for i := 1.0; i < 4; i++ {
		fmt.Printf("sqrt of %v using own implementation equals %v\n", i, homeworks.Sqrt(i))
	}

	for i := 1.0; i < 4; i++ {
		fmt.Printf("sqrt of %v using math package equals %v\n", i, math.Sqrt(i))
	}

	//mapa := make(map[string]string)

}
