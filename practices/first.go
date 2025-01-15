package practices

import (
	"errors"
	"strings"
	"unicode"
)

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	var result strings.Builder
	var prev rune
	var escaping bool

	runes := []rune(input)
	for _, r := range runes {
		if escaping {
			if r != '\\' && !unicode.IsDigit(r) {
				return "", errors.New("некорректная строка")
			}
			result.WriteRune(r)
			escaping = false
			prev = r
			continue
		}

		if r == '\\' {
			escaping = true
			continue
		}

		if unicode.IsDigit(r) {
			if prev == 0 || unicode.IsDigit(prev) {
				return "", errors.New("некорректная строка")
			}

			digit := int(r - '0')
			if digit > 1 {
				result.WriteString(strings.Repeat(string(prev), digit-1))
			}
			prev = 0
		} else {
			result.WriteRune(r)
			prev = r
		}
	}

	if escaping {
		return "", errors.New("некорректная строка")
	}

	return result.String(), nil
}
