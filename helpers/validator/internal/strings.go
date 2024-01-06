package validations

import (
	"fmt"
	"reflect"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

func String(value reflect.Value, validations map[string][]string) error {
	for name, args := range validations {
		validationFn := stringValidations[name]
		if validationFn == nil {
			return nil
			// return normalizederr.NewValidationError("Unknown desired validation for type string")
		}

		err := validationFn(value.String(), args)
		if err != nil {
			return err
		}
	}

	return nil
}

var stringValidations = map[string]func(string, []string) error{
	"required": requiredStr,
	"oneof":    oneof,
}

func requiredStr(str string, _ []string) error {
	if str == "" {
		return normalizederr.NewValidationError("Required.")
	}
	return nil
}

func oneof(str string, options []string) error {
	if str == "" {
		return nil
	}

	for _, option := range options {
		if option == str {
			return nil
		}
	}

	msg := fmt.Sprintf("Must be one of %v", options)
	return normalizederr.NewValidationError(msg)
}
