package validations

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"

	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
)

func String(value reflect.Value, validations map[string][]string) error {
	for name, args := range validations {
		validationFn := stringValidations[name]
		if validationFn == nil {
			return apperr.NewInternalError(fmt.Sprintf("unknown \"%s\" validation for type string", name))
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
	"wordID":     wordID,
}

func requiredStr(str string, _ []string) error {
	if str == "" {
		return apperr.NewValidationError("required")
	}
	return nil
}

func atLeastOne(str string, options []string) error {
	if str == "" {
		return nil
	}

	bodyMsg := ""
	doesMatch := true
	for _, opt := range options {
		var rgxStr string
		switch opt {
		case "letter":
			rgxStr = "[a-zA-Z]"
			bodyMsg += " 1 letter"
		case "lowercase":
			rgxStr = "[a-z]"
			bodyMsg += " 1 lowercase"
		case "uppercase":
			rgxStr = "[A-Z]"
			bodyMsg += " 1 uppercase"
		case "number":
			rgxStr = "[0-9]"
			bodyMsg += " 1 number"
		case "specialChar":
			rgxStr = `[@#$%^&*\-_+=!?(){}[\]]`
			bodyMsg += " 1 special character"
		}

		if doesMatch {
			rgx := regexp.MustCompile(rgxStr)
			doesMatch = rgx.MatchString(str)
		}
	}

	if doesMatch {
		return nil
	}

	msg := fmt.Sprintf("must have at least%v.", bodyMsg)
	return apperr.NewValidationError(msg)
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

	msg := fmt.Sprintf("must have %v characters", target)
	return apperr.NewValidationError(msg)
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

	msg := fmt.Sprintf("must have at maximum %v characters", target)
	return apperr.NewValidationError(msg)
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

	msg := fmt.Sprintf("must have at least %v characters", target)
	return apperr.NewValidationError(msg)
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

	msg := fmt.Sprintf("must be one of %v", options)
	return apperr.NewValidationError(msg)
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

	msg := "must have only alphanumeric characters, hyphen and underscore"
	return apperr.NewValidationError(msg)
}

func uri(str string, _ []string) error {
	if str == "" {
		return nil
	}

	u, err := url.Parse(str)
	if err != nil {
		return err
	} else if u.Scheme != "" && u.Host != "" {
		return nil
	}

	msg := "must have a valid uri format"
	return apperr.NewValidationError(msg)
}

func wordID(str string, _ []string) error {
	if str == "" {
		return nil
	}

	doesMatch, err := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, str)
	if err != nil {
		return err
	} else if !doesMatch {
		msg := "must have only letters, numbers, period, underscore, and hyphen"
		return apperr.NewValidationError(msg)
	}
	
	doesMatch, err = regexp.MatchString(`[._-]{2}`, str)
	if err != nil {
		return err
	} else if !doesMatch {
		return nil
	}

	msg := "must not have consecutive period, underscore, and hyphen"
	return apperr.NewValidationError(msg)
}
