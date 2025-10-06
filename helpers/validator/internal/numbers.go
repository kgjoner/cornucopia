package validations

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
)

func Number(value reflect.Value, validations map[string][]string) error {
	for name, args := range validations {
		validationFn := numberValidations[name]
		if validationFn == nil {
			return apperr.NewInternalError(fmt.Sprintf("unknown \"%s\" validation for type number", name))
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
	"max":      max,
}

func requiredNum(num int, _ []string) error {
	if num == 0 {
		return apperr.NewValidationError("required")
	}
	return nil
}

func min(num int, args []string) error {
	limit, _ := strconv.Atoi(args[0])

	if num < limit {
		msg := fmt.Sprintf("must be higher than %v", limit)
		return apperr.NewValidationError(msg)
	}
	return nil
}

func max(num int, args []string) error {
	limit, _ := strconv.Atoi(args[0])

	if num > limit {
		msg := fmt.Sprintf("must be equal or less than %v", limit)
		return apperr.NewValidationError(msg)
	}
	return nil
}
