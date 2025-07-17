package htypes_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/helpers/htypes"
	"github.com/stretchr/testify/assert"
)

func TestPhoneParse(t *testing.T) {
	input := "+55 (11) 99999-9999"
	res, err := htypes.ParsePhoneNumber(input)

	assert.Nil(t, err)
	assert.Equal(t, htypes.PhoneNumber("+5511999999999"), res)
}
