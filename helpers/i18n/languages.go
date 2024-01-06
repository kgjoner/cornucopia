package i18n

type Language string

type languageValues struct {
	PT_BR Language
	EN_US Language
}

func (s Language) Enumerate() any {
	return languageValues{
		"pt-br",
		"en-us",
	}
}

var LanguageValues = Language.Enumerate("").(languageValues)
