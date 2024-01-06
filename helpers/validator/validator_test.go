package validator

import (
	"strings"
	"testing"
)

type StructSample struct {
	name string `validate:"required"`
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
	err = Validate(structSample)
	if err == nil || !strings.Contains(err.Error(), "name: ") {
		t.Errorf("Expected error, got nil")
		return
	} else if !strings.Contains(err.Error(), "name: Required") {
		t.Errorf("Expected Required error, got %s", err)
	}

	return
}
