package htypes

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type Country string

func (c *Country) IsValid() error {
	if c.IsZero() {
		return nil
	}

	_, exists := countries[string(*c)]
	if !exists {
		for key, name := range countries {
			if strings.EqualFold(name, string(*c)) {
				*c = Country(key)
				return nil
			}
		}

		return normalizederr.NewValidationError("country does not exist")
	}

	return nil
}

func (c Country) IsZero() bool {
	return c == ""
}

func (c Country) Name() string {
	return countries[string(c)]
}

/* ================================================================================
	INIT
================================================================================ */

//go:embed assets/countries.json
var countriesJSON []byte

var countries map[string]string

func init() {
	countries = make(map[string]string)
	if err := json.Unmarshal(countriesJSON, &countries); err != nil {
		panic("failed to load countries: " + err.Error())
	}
}
