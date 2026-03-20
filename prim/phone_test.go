package prim_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/v3/prim"
	"github.com/stretchr/testify/assert"
)

func TestPhoneParse(t *testing.T) {
	input := "+55 (11) 99999-9999"
	res, err := prim.ParsePhoneNumber(input)

	assert.Nil(t, err)
	assert.Equal(t, prim.PhoneNumber("+5511999999999"), res)
}
