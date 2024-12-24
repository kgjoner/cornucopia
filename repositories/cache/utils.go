package cache

import (
	"crypto/sha256"
	"fmt"
)

type KeyGen struct {
	prefix string
}

func NewKeyGen(prefix string) *KeyGen {
	return &KeyGen{
		prefix,
	}
}

// Generate a key in the form prefix:hash, where prefix is the one chosen when KeyGen was created and hash is a sha256 of all args.
func (k KeyGen) Key(args ...any) string {
	return k.prefix + ":" + hashArgs(args)
}

func hashArgs(args ...any) string {
	var argsStr string
	for _, arg := range args {
		argsStr += fmt.Sprintf("%v", arg)
	}

	h := sha256.New()
	h.Write([]byte(argsStr))
	return fmt.Sprintf("%x", h.Sum(nil))
}
