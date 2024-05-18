package pwdgen

import (
	"math/rand"
	"strings"

	"github.com/kgjoner/cornucopia/utils/sliceman"
)

// Generate a string with desired length including random runes based on selected sets. The set options are: lower, upper, number and special. It will include at least one rune of each set. If no set is provided, it will include all of them. The minimum length is equal to the variety of selected sets, i.e., 4 if all of them are selected.
func Generate(length int, sets ...string) string {
	var password strings.Builder

	if len(sets) == 0 {
		sets = []string{"lower", "upper", "number", "special"}
	}

	fullSet := "";
	if (sliceman.IndexOf(sets, "lower") != -1) {
		lowerCharSet := "abcdedfghijklmnopqrst"
		fullSet += lowerCharSet

		index := rand.Intn(len(lowerCharSet))
		password.WriteString(string(lowerCharSet[index]))
	}

	if (sliceman.IndexOf(sets, "upper") != -1) {
		upperCharSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		fullSet += upperCharSet

		index := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[index]))
	}

	if (sliceman.IndexOf(sets, "number") != -1) {
		numberSet := "0123456789"
		fullSet += numberSet

		index := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[index]))
	}

	if (sliceman.IndexOf(sets, "special") != -1) {
		specialCharSet := "!@#$%&*"
		fullSet += specialCharSet

		index := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[index]))
	}

	for i := 0; i < (length - len(sets)); i++ {
		index := rand.Intn(len(fullSet))
		password.WriteString(string(fullSet[index]))
	}

	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
