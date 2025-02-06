package main

import (
	"fmt"
	"log"
	"strconv"
)

var ACTIONS = map[int]string{
	0: "wink",
	1: "double blink",
	2: "close your eyes",
	3: "jump",
}

func reverseSlice(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Handshake(code uint) []string {
	var result []string
	binaryCode := strconv.FormatInt(int64(code), 2)
	log.Print(binaryCode)

	for i, v := len(binaryCode)-1, 0; i != -1; i, v = i-1, v+1 {
		if binaryCode[i] == 48 {
			continue
		} else {
			if v == 4 {
				reverseSlice(result)
				break
			}
			result = append(result, ACTIONS[v])
		}
	}

	return result
}

func main() {
	var n uint = 9
	fmt.Println(Handshake(n))
}
