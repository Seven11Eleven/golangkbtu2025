package homeworks

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	stringsCount := make(map[string]int)

	for _, val := range words {
		stringsCount[val]++
	}
	return stringsCount
}

func main() {
	wc.Test(WordCount)
}
