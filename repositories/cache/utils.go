package cache

import (
	"github.com/kgjoner/cornucopia/v2/utils/hash"
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
	return k.prefix + ":" + hash.From(args...)
}
