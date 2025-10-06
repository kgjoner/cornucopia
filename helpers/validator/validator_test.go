package validator

import (
	"errors"
	"strings"
	"testing"

	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
	"github.com/stretchr/testify/assert"
)

type StructSample struct {
	Name string `validate:"required"`
	Age  int
}

func (s StructSample) IsValid() error {
	return nil
}

func TestRequiredField(t *testing.T) {
	err := Validate("", "required")
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	} else if !strings.Contains(err.Error(), "Required") {
		t.Errorf("Expected Required error, got %s", err)
	}

	structSample := StructSample{}
	err = Validate(structSample, "required")
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	} else if !strings.Contains(err.Error(), "Required") {
		t.Errorf("Expected Required error, got %s", err)
	}

	structSample.Age = 21
	err = Validate(structSample, "required")
	if err == nil || !strings.Contains(err.Error(), "Name: ") {
		t.Errorf("Expected error, got %v", err)
		return
	} else if !strings.Contains(err.Error(), "Name: Required") {
		t.Errorf("Expected Required error, got %s", err)
	}
}

func TestPasswordValidation(t *testing.T) {
	validations := []string{"required", "min=8", "atLeastOne=letter number specialChar"}

	err := Validate("", validations...)
	assert.Contains(t, err.Error(), "Required")

	err = Validate("1234", validations...)
	assert.Contains(t, err.Error(), "at least 8 char")

	err = Validate("12345678", validations...)
	assert.Contains(t, err.Error(), "at least 1 letter 1 number 1 special char")

	err = Validate("1234ABCD", validations...)
	assert.Contains(t, err.Error(), "at least 1 letter 1 number 1 special char")

	err = Validate("Abcdefg!", validations...)
	assert.Contains(t, err.Error(), "at least 1 letter 1 number 1 special char")

	err = Validate("Abc1234!", validations...)
	assert.Nil(t, err)
}

func TestPointerSubstruct(t *testing.T) {
	type ParentStruct struct {
		Name   string `validate:"required"`
		Sample *StructSample
	}

	err := Validate(ParentStruct{Name: "name"})
	assert.Nil(t, err)
}

// Write a test for the Enum validation

type SliceEnum string

const (
	EnumValue1 SliceEnum = "value1"
	EnumValue2 SliceEnum = "value2"
)

func (e SliceEnum) Enumerate() any {
	return []Enum{EnumValue1, EnumValue2}
}

type StructEnum string

func (e StructEnum) Enumerate() any {
	return struct {
		Value1 StructEnum
		Value2 StructEnum
	}{
		Value1: "value1",
		Value2: "value2",
	}
}

func TestEnumValidation(t *testing.T) {
	// Test Enum with slice enumeration
	type WithSliceEnum struct {
		Value SliceEnum
	}

	valid1 := WithSliceEnum{Value: EnumValue1}
	err := Validate(valid1)
	assert.Nil(t, err)

	var appErr *apperr.AppError
	invalid1 := WithSliceEnum{Value: "invalid"}
	err = Validate(invalid1)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, appErr.Code, apperr.InvalidData)

	// Test Enum with slice enumeration
	type WithStructEnum struct {
		Value StructEnum
	}

	valid2 := WithStructEnum{Value: StructEnum("value1")}
	err = Validate(valid2)
	assert.Nil(t, err)

	invalid2 := WithStructEnum{Value: StructEnum("invalid")}
	err = Validate(invalid2)
	assert.NotNil(t, err)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, appErr.Code, apperr.InvalidData)
}
