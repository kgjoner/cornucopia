package htypes

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type Country string

func (c Country) IsValid() error {
	data, _ := ioutil.ReadFile("./pkg/helpers/htypes/assets/countries.json")

	var countriesMap map[string]string
	json.Unmarshal(data, &countriesMap)

	_, exists := countriesMap[string(c)]
	if !exists {
		return normalizederr.NewValidationError("country does not exist")
	}

	return nil
}

func (c Country) IsZero() bool {
	return c == ""
}

func (c Country) Name() string {
	data, _ := ioutil.ReadFile("./assets/countries.json")

	var countriesMap map[string]string
	json.Unmarshal(data, &countriesMap)

	return countriesMap[string(c)]
}
