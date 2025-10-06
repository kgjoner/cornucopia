package pwdgen

import (
	cryptorand "crypto/rand"
	"math/big"
	"strings"

	"github.com/kgjoner/cornucopia/utils/sliceman"
)

// GeneratePassword creates a string with desired length including random runes based
// on selected sets. The set options are: lower, upper, number and special. It will
// include at least one rune of each set. If no set is provided, it will include all
// of them. The minimum length is equal to the variety of selected sets, i.e., 4 if
// all of them are selected.
func GeneratePassword(length int, sets ...string) string {
	var password strings.Builder

	if len(sets) == 0 {
		sets = []string{"lower", "upper", "number", "special"}
	}

	fullSet := ""
	if sliceman.IndexOf(sets, "lower") != -1 {
		lowerCharSet := "abcdefghijklmnopqrstuvwxyz"
		fullSet += lowerCharSet

		index := secureRandomInt(len(lowerCharSet))
		password.WriteString(string(lowerCharSet[index]))
	}

	if sliceman.IndexOf(sets, "upper") != -1 {
		upperCharSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		fullSet += upperCharSet

		index := secureRandomInt(len(upperCharSet))
		password.WriteString(string(upperCharSet[index]))
	}

	if sliceman.IndexOf(sets, "number") != -1 {
		numberSet := "0123456789"
		fullSet += numberSet

		index := secureRandomInt(len(numberSet))
		password.WriteString(string(numberSet[index]))
	}

	if sliceman.IndexOf(sets, "special") != -1 {
		specialCharSet := "!@#$%&*"
		fullSet += specialCharSet

		index := secureRandomInt(len(specialCharSet))
		password.WriteString(string(specialCharSet[index]))
	}

	for i := 0; i < (length - len(sets)); i++ {
		index := secureRandomInt(len(fullSet))
		password.WriteString(string(fullSet[index]))
	}

	inRune := []rune(password.String())
	for i := len(inRune) - 1; i > 0; i-- {
		j := secureRandomInt(i + 1)
		inRune[i], inRune[j] = inRune[j], inRune[i]
	}

	return string(inRune)
}

// GenerateNumericCode creates a secure numeric code of the specified length.
func GenerateNumericCode(length int) string {
	const numberSet = "0123456789"

	code := make([]byte, length)
	for i := 0; i < length; i++ {
		index := secureRandomInt(len(numberSet))
		code[i] = numberSet[index]
	}

	return string(code)
}

// secureRandomInt generates a cryptographically secure random integer
func secureRandomInt(max int) int {
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic("Failed to generate secure random number")
	}
	return int(n.Int64())
}
