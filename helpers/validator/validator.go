package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	v "github.com/kgjoner/cornucopia/helpers/validator/internal"
	"github.com/kgjoner/cornucopia/utils/sliceman"
)

type Validator interface {
	IsValid() error
}

func Validate(value interface{}, validations ...string) error {
	if err := assertSelfValidation(value); err != nil {
		return err
	}

	var reflectValue reflect.Value
	if r, ok := value.(reflect.Value); ok {
		reflectValue = r
	} else {
		reflectValue = reflect.ValueOf(value)
	}

	if reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflect.Indirect(reflectValue)
	}

	if reflectValue.IsValid() && reflectValue.CanInterface() {
		switch t := reflectValue.Interface().(type) {
		case Enum:
			return validateEnum(t, validations)
		case time.Time:
			return validateTime(reflectValue, validations...)
		}
	}

	switch reflectValue.Kind() {
	case reflect.Struct:
		return validateStruct(reflectValue, validations)
	case reflect.Slice, reflect.Array:
		return validateArray(reflectValue, validations)
	case reflect.Map:
		return validateMap(reflectValue, validations)
	// TODO: handle interface kind validation
	case reflect.Interface:
		return nil
	default:
		return validatePrimitive(reflectValue, validations)
	}
}

func assertSelfValidation(primitive interface{}) error {
	p := primitive
	if v, ok := primitive.(reflect.Value); ok && v.CanInterface() &&
		(v.Kind() != reflect.Pointer || !v.IsNil()) {
		p = v.Interface()
	}

	if v, ok := p.(Validator); ok {
		return v.IsValid()
	}

	return nil
}

func validateArray(arr reflect.Value, validations []string) error {
	length := arr.Len()

	forwardedValidations := []string{}
	for i, validation := range validations {
		if strings.Contains(validation, "required") {
			if arr.IsZero() || length == 0 {
				return normalizederr.NewValidationError("Required.")
			}
			validations = sliceman.Remove(validations, i)
		} else {
			forwardedValidations = append(forwardedValidations, validation)
		}
	}

	for i := 0; i < length; i++ {
		v := arr.Index(i)
		err := Validate(v, forwardedValidations...)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateMap(mp reflect.Value, validations []string) error {
	length := mp.Len()

	forwardedValidations := []string{}
	for i, validation := range validations {
		if strings.Contains(validation, "required") {
			if mp.IsZero() || length == 0 {
				return normalizederr.NewValidationError("Required.")
			}
			validations = sliceman.Remove(validations, i)
		} else {
			forwardedValidations = append(forwardedValidations, validation)
		}
	}

	for _, key := range mp.MapKeys() {
		err := Validate(key)
		if err != nil {
			return err
		}

		v := mp.MapIndex(key)
		err = Validate(v, forwardedValidations...)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateEnum(enum Enum, validations []string) error {
	var availableOpt []string

	forwardedValidations := []string{}
	for _, validation := range validations {
		if strings.Contains(validation, "restrictenum=") {
			availableOpt = validationMap(validation)["restrictenum"]
		} else {
			forwardedValidations = append(forwardedValidations, validation)
		}
	}

	err := validatePrimitive(reflect.ValueOf(enum), forwardedValidations)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(enum.Enumerate())
	var length int

	switch v.Kind() {
	case reflect.Struct:
		length = v.NumField()
	case reflect.Slice, reflect.Array:
		length = v.Len()
	default:
		return normalizederr.NewValidationError("Invalid enum type. Must be a struct or slice.")
	}

	var enumerate []Enum
	isValid := false
	enumName := reflect.TypeOf(enum).Name()
	for i := 0; i < length; i++ {
		var field reflect.Value
		switch v.Kind() {
		case reflect.Struct:
			field = v.Field(i)
		case reflect.Slice, reflect.Array:
			field = v.Index(i)
		}

		validValue, ok := field.Interface().(Enum)
		if !ok {
			message := fmt.Sprintf("All possible values of a Enum must be a Enum as well. %s is not structured correctly.", enumName)
			return normalizederr.NewValidationError(message)
		}

		if len(availableOpt) != 0 && field.CanConvert(reflect.TypeOf("")) {
			if value, ok := field.Convert(reflect.TypeOf("")).Interface().(string); ok {
				for _, opt := range availableOpt {
					if opt == value {
						enumerate = append(enumerate, validValue)
						if enum == validValue {
							isValid = true
						}
						break
					}
				}
			}
			continue
		}

		enumerate = append(enumerate, validValue)
		if enum == validValue {
			isValid = true
		}
	}

	if !isValid {
		message := fmt.Sprintf("Invalid %s value. Must be one of: %v", enumName, enumerate)
		return normalizederr.NewValidationError(message)
	}

	return nil
}

func validatePrimitive(primitive reflect.Value, validations []string) error {
	if len(validations) == 0 {
		return nil
	}

	valMap := validationMap(validations...)

	switch primitive.Kind() {
	case reflect.String:
		return v.String(primitive, valMap)
	case reflect.Int:
		return v.Number(primitive, valMap)
	default:
		return normalizederr.NewValidationError("No accepted primitive type")
	}
}

func validateStruct(obj interface{}, validations []string) error {
	var objValue reflect.Value
	if r, ok := obj.(reflect.Value); ok {
		objValue = r
	} else {
		objValue = reflect.ValueOf(obj)
	}

	for _, validation := range validations {
		switch validation {
		case "required":
			if objValue.IsZero() {
				return normalizederr.NewValidationError("Required.")
			}
		case "ignore":
			return nil
		}
	}

	if objValue.IsZero() {
		return nil
	}

	validationsByField := extractValidationsByField(obj)
	errors := make(map[string]error)

	for field, validations := range validationsByField {
		fieldValue := objValue.FieldByName(field)
		err := Validate(fieldValue, validations...)
		if err != nil {
			errors[field] = err
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return normalizederr.NewValidationErrorFromMap(errors)
}

func validateTime(value reflect.Value, validations ...string) error {
	valMap := validationMap(validations...)
	return v.Time(value, valMap)
}

func extractValidationsByField(obj interface{}) map[string][]string {
	o := reflect.TypeOf(obj)
	if v, ok := obj.(reflect.Value); ok {
		o = v.Type()
	}

	validations := make(map[string][]string)

	for i := 0; i < o.NumField(); i++ {
		field := o.Field(i)

		fieldValidations := field.Tag.Get("validate")
		if fieldValidations == "" {
			validations[field.Name] = nil
		} else {
			validations[field.Name] = strings.Split(fieldValidations, ",")
		}
	}

	return validations
}

func validationMap(validations ...string) map[string][]string {
	validationMap := make(map[string][]string)
	for _, validationString := range validations {
		validationSlice := strings.Split(validationString, "=")
		if len(validationSlice) == 1 {
			validationMap[validationSlice[0]] = nil
		} else if len(validationSlice) == 2 {
			validationMap[validationSlice[0]] = strings.Split(validationSlice[1], " ")
		}
	}

	return validationMap
}
