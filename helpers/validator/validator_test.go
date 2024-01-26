package validator

import (
	"strings"
	"testing"
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
