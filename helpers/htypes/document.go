package htypes

import (
	"strings"

	"github.com/kgjoner/cornucopia/helpers/validator"
)

type Document string

func (d Document) IsValid() error {
	parsed := d.Parse()
	return validator.Validate(parsed)
} 

func (d Document) IsZero() bool {
	return d == ""
}

func (d Document) String() string {
	return string(d)
}

type ParsedDocument struct {
	Number string `validate:"required"`
	Kind   string `validate:"required,oneof=cpf cnpj passport"`
}

func (d Document) Parse() ParsedDocument {
	parts := strings.Split(string(d), "_")
	if len(parts) != 2 {
		return ParsedDocument{}
	}

	return ParsedDocument {
		Number: parts[1],
		Kind: parts[0],
	}
}
