package validations

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

func Number(value reflect.Value, validations map[string][]string) error {
	for name, args := range validations {
		validationFn := numberValidations[name]
		if validationFn == nil {
			return nil
			// return normalizederr.NewValidationError("Unknown desired validation for type number")
		}

		err := validationFn(int(value.Int()), args)
		if err != nil {
			return err
		}
	}

	return nil
}

var numberValidations = map[string]func(int, []string) error{
	"required": requiredNum,
	"min":      min,
}

func requiredNum(num int, _ []string) error {
	if num == 0 {
		return normalizederr.NewValidationError("Required.")
	}
	return nil
}

func min(num int, args []string) error {
	limit, _ := strconv.Atoi(args[0])

	if num < limit {
		msg := fmt.Sprintf("Must be higher than %v", limit)
		return normalizederr.NewValidationError(msg)
	}
	return nil
}
