package htypes

import "github.com/kgjoner/cornucopia/helpers/validator"

type Email string

func (e Email) IsValid() error {
	return validator.Validate(string(e), "email")
}
