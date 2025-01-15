package homeworks

func AtoiBase(s string, base string) int {
	if !isValidBase(base) {
		return 0
	}

	baseMapa := make(map[rune]int)
	for i, char := range base {
		baseMapa[char] = i
	}

	res := 0
	for _, char := range s {
		res = res*len(base) + baseMapa[char]
	}
	return res
}

func isValidBase(base string) bool {
	if len(base) < 2 {
		return false
	}

	charSet := make(map[rune]bool)
	for _, char := range base {
		if char == '+' || char == '-' || charSet[char] {
			return false
		}
		charSet[char] = true
	}
	return true
}
