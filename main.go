package main

import (
	"fmt"
	"github.com/Seven11Eleven/golangkbtu2025/practices"
)

func runTests() {
	tests := []struct {
		input    string
		expected string
		error    bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"3abc", "", true},
		{"45", "", true},
		{"aaa10b", "", true},
		{"aaa0b", "aab", false},
		{"", "", false},
		{"d\n5abc", "d\n\n\n\n\nabc", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"qwe\\\\\\5", "qwe\\\\\\", false},
		{"qw\ne", "", true},
	}

	for _, test := range tests {
		result, err := practices.Unpack(test.input)
		if test.error {
			if err == nil {
				fmt.Printf("FAIL: Input: %q -> Expected error but got result: %q\n", test.input, result)
			} else {
				fmt.Printf("PASS: Input: %q -> Expected error and got error: %v\n", test.input, err)
			}
		} else {
			if err != nil || result != test.expected {
				fmt.Printf("FAIL: Input: %q -> Expected: %q, Got: %q, Error: %v\n", test.input, test.expected, result, err)
			} else {
				fmt.Printf("PASS: Input: %q -> Output: %q\n", test.input, result)
			}
		}
	}
}

func main() {
	runTests()
}
