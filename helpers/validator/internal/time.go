package validations

import (
	"fmt"
	"reflect"
	"time"

	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
)

func Time(value reflect.Value, validations map[string][]string) error {
	for name, args := range validations {
		validationFn := timeValidations[name]
		if validationFn == nil {
			return apperr.NewInternalError(fmt.Sprintf("unknown \"%s\" validation for type time", name))
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
		return apperr.NewValidationError("required")
	}
	return nil
}
