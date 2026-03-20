package prim_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/v3/prim"
	"github.com/kgjoner/cornucopia/v3/validator"
	"github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
	empty := prim.Email("")
	err := empty.IsValid()
	assert.Nil(t, err)

	invalid := prim.Email("test.com")
	err = invalid.IsValid()
	assert.NotNil(t, err)

	valid := prim.Email("test@test.com")
	err = valid.IsValid()
	assert.Nil(t, err)

	insideStr := struct {
		Name  string
		Email prim.Email `validate:"required"`
	}{
		"Dummy",
		empty,
	}
	err = validator.Validate(insideStr)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Email: required")
}
