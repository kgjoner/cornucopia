package htypes

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
)

type Country string

func ParseCountry(str string) (Country, error) {
	if str == "" {
		return "", nil
	}

	_, exists := countries[str]
	if !exists {
		for key, name := range countries {
			if strings.EqualFold(name, str) {
				return Country(key), nil
			}
		}

		return "", apperr.NewValidationError("country does not exist")
	}

	return Country(str), nil
}

func (c Country) IsValid() error {
	if c.IsZero() {
		return nil
	}

	_, exists := countries[string(c)]
	if !exists {
		return apperr.NewValidationError("country does not exist; certify that the country code is correct or that country name was parsed before validation check")
	}

	return nil
}

func (c Country) IsZero() bool {
	return c == ""
}

func (c Country) Name() string {
	name := countries[string(c)]
	if name == "" {
		for _, value := range countries {
			if strings.EqualFold(value, string(c)) {
				name = value
				break
			}
		}
	}

	return name
}

func (c Country) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Name())
}

func (c *Country) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*c, err = ParseCountry(s)
	return err
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
