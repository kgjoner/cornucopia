package htypes_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/v2/helpers/htypes"
	"github.com/kgjoner/cornucopia/v2/helpers/validator"
	"github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
	empty := htypes.Email("")
	err := empty.IsValid()
	assert.Nil(t, err)

	invalid := htypes.Email("test.com")
	err = invalid.IsValid()
	assert.NotNil(t, err)

	valid := htypes.Email("test@test.com")
	err = valid.IsValid()
	assert.Nil(t, err)

	insideStr := struct {
		Name  string
		Email htypes.Email `validate:"required"`
	}{
		"Dummy",
		empty,
	}
	err = validator.Validate(insideStr)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Email: Required")
}
