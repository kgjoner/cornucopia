package validations

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

func String(value reflect.Value, validations map[string][]string) error {
	for name, args := range validations {
		validationFn := stringValidations[name]
		if validationFn == nil {
			return normalizederr.NewValidationError("Unknown desired validation for type string")
		}

		err := validationFn(value.String(), args)
		if err != nil {
			return err
		}
	}

	return nil
}

var stringValidations = map[string]func(string, []string) error{
	"required":   requiredStr,
	"atLeastOne": atLeastOne,
	"length":     length,
	"max":        maxStr,
	"min":        minStr,
	"oneof":      oneof,
	"slug":       slug,
	"uri":        uri,
	"wordId":     wordId,
}

func requiredStr(str string, _ []string) error {
	if str == "" {
		return normalizederr.NewValidationError("Required.")
	}
	return nil
}

func atLeastOne(str string, options []string) error {
	if str == "" {
		return nil
	}

	rgxStr := "^"
	bodyMsg := ""
	for _, opt := range options {
		switch opt {
		case "letter":
			rgxStr += "(?=.*[a-zA-Z])"
			bodyMsg += " 1 letter"
		case "lowercase":
			rgxStr += "(?=.*[a-z])"
			bodyMsg += " 1 lowercase"
		case "uppercase":
			rgxStr += "(?=.*[A-Z])"
			bodyMsg += " 1 uppercase"
		case "number":
			rgxStr += "(?=.*[0-9])"
			bodyMsg += " 1 number"
		case "specialChar":
			rgxStr += `(?=.*[@#$%^&*\-_+=!?(){}[\]])`
			bodyMsg += " 1 special character"
		}
	}
	rgxStr += ".*$"

	doesMatch, err := regexp.MatchString(rgxStr, str)
	if err != nil {
		return err
	} else if doesMatch {
		return nil
	}

	msg := fmt.Sprintf("Must have at least%v.", bodyMsg)
	return normalizederr.NewValidationError(msg)
}

func length(str string, options []string) error {
	if str == "" {
		return nil
	}

	target, err := strconv.ParseInt(options[0], 10, 32)
	if err != nil {
		return err
	}

	if len(str) == int(target) {
		return nil
	}

	msg := fmt.Sprintf("Must have %v characters", target)
	return normalizederr.NewValidationError(msg)
}

func maxStr(str string, options []string) error {
	if str == "" {
		return nil
	}

	target, err := strconv.ParseInt(options[0], 10, 32)
	if err != nil {
		return err
	}

	if len(str) <= int(target) {
		return nil
	}

	msg := fmt.Sprintf("Must have at maximum %v characters", target)
	return normalizederr.NewValidationError(msg)
}

func minStr(str string, options []string) error {
	if str == "" {
		return nil
	}

	target, err := strconv.ParseInt(options[0], 10, 32)
	if err != nil {
		return err
	}

	if len(str) >= int(target) {
		return nil
	}

	msg := fmt.Sprintf("Must have at least %v characters", target)
	return normalizederr.NewValidationError(msg)
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

func slug(str string, _ []string) error {
	if str == "" {
		return nil
	}

	doesMatch, err := regexp.MatchString(`^[\w-]+$`, str)
	if err != nil {
		return err
	} else if doesMatch {
		return nil
	}

	msg := fmt.Sprintf("Must have only alphanumeric characters, hyphen and underscore")
	return normalizederr.NewValidationError(msg)
}

func uri(str string, _ []string) error {
	if str == "" {
		return nil
	}

	doesMatch, err := regexp.MatchString(
		`^[A-Za-z]+:\/\/[\w-]+(\.[\w-]+)+(:\d+)?(\/[\w-]+)*(\?\w+=[\w%]+(&\w+=[\w%]+)*)?$`,
		str,
	)
	if err != nil {
		return err
	} else if doesMatch {
		return nil
	}

	msg := fmt.Sprintf("Must have a valid uri format")
	return normalizederr.NewValidationError(msg)
}

func wordId(str string, _ []string) error {
	if str == "" {
		return nil
	}

	doesMatch, err := regexp.MatchString(`^[a-zA-Z0-9.]+$`, str)
	if err != nil {
		return err
	} else if doesMatch {
		return nil
	}

	msg := fmt.Sprintf("Must have only letters, numbers and period.")
	return normalizederr.NewValidationError(msg)
}
