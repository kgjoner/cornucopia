package prim

import "strings"

type Locale string

const (
	Portuguese Locale = "pt-br"
	English    Locale = "en-us"
	Spanish    Locale = "es-es"
	French     Locale = "fr-fr"
	German     Locale = "de-de"
	Japanese   Locale = "ja-jp"
	Korean     Locale = "ko-kr"
)

func (l Locale) Enumerate() any {
	return []Locale{
		Portuguese,
		English,
		Spanish,
		French,
		German,
		Japanese,
		Korean,
	}
}

// Base returns the base language code: "pt-br" → "pt"
func (l Locale) Base() string {
    tag := string(l)
    if i := strings.IndexByte(tag, '-'); i != -1 {
        return tag[:i]
    }
    return tag
}
