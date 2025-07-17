package sanitizer

// Keep only numeric characters
func Digit(str string) string {
	j := 0
	parsedBytes := []byte(str)
	for _, b := range parsedBytes {
		if '0' <= b && b <= '9' {
			parsedBytes[j] = b
			j++
		}
	}

	return string(parsedBytes[:j])
}

// Checks if a string contains only digit characters
func IsDigitOnly(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
