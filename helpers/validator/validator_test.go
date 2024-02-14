package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StructSample struct {
	Name string `validate:"required"`
	Age  int
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

	return
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
