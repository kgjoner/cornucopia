package hash

import (
	"crypto/sha256"
	"fmt"
)

// Hash args with Sha256 algorythm. It will concatenate their string values and return a base16 hash of it.
func From(args ...any) string {
	var argsStr string
	for _, arg := range args {
		argsStr += fmt.Sprintf("%v", arg)
	}

	h := sha256.New()
	h.Write([]byte(argsStr))
	return fmt.Sprintf("%x", h.Sum(nil))
}