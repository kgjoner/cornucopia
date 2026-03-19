package i18n

type Language string

const (
	Portuguese Language = "pt-br"
	English    Language = "en-us"
	Spanish    Language = "es-es"
	French     Language = "fr-fr"
	German     Language = "de-de"
	Japanese   Language = "ja-jp"
	Korean     Language = "ko-kr"
)

func (s Language) Enumerate() any {
	return []Language{
		Portuguese,
		English,
		Spanish,
		French,
		German,
		Japanese,
		Korean,
	}
}
