package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//const filePath string = "./data/users.txt"

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, _ := os.Open(filePath)
	defer file.Close()

	user := User{}
	browsers := make(map[string]bool, 1000)
	var email string
	fmt.Fprintln(out, "found users:")
	var isAndroid, isMSIE bool
	var i int
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		user.UnmarshalJSON(sc.Bytes())
		isAndroid = false
		isMSIE = false
		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
			} else {
				continue
			}

			browsers[browser] = true
		}
		if isAndroid && isMSIE {
			email = strings.Replace(user.Email, "@", " [at] ", -1)
			fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", i, user.Name, email))
		}
		i++
	}
	fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
}
