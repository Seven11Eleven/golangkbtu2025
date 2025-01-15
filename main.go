package main

import (
	"fmt"
	"github.com/Seven11Eleven/golangkbtu2025/homeworks"
)

func main() {
	fmt.Println(homeworks.AtoiBase("125", "0123456789"))
	fmt.Println(homeworks.AtoiBase("1111101", "01"))
	fmt.Println(homeworks.AtoiBase("7D", "0123456789ABCDEF"))
	fmt.Println(homeworks.AtoiBase("uoi", "choumi"))
	fmt.Println(homeworks.AtoiBase("bbbbbab", "-ab"))
}

/*
$ go run .
125
125
125
125
0
$
*/
