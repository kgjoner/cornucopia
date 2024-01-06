package validations

import (
	"reflect"
	"time"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

func Time(value reflect.Value, validations map[string][]string) error {
	for name, args := range validations {
		validationFn := timeValidations[name]
		if validationFn == nil {
			return nil
			// return normalizederr.NewValidationError("Unknown desired validation for type time")
		}

		t, _ := value.Interface().(time.Time)
		err := validationFn(t, args)
		if err != nil {
			return err
		}
	}

	return nil
}

var timeValidations = map[string]func(time.Time, []string) error{
	"required": requiredTime,
}

func requiredTime(t time.Time, _ []string) error {
	if t.IsZero() {
		return normalizederr.NewValidationError("Required.")
	}
	return nil
}
