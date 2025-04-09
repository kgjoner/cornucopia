package datatransform_test

import (
	"encoding/json"
	"testing"

	"github.com/kgjoner/cornucopia/utils/datatransform"
	"github.com/stretchr/testify/assert"
)

type Kind string

type OriginalStruct struct {
	Name   string `json:"jname"`
	Number int    `json:"jnumber"`
}

func TestStringArray(t *testing.T) {
	orig := []Kind{"kind-one"}
	res := datatransform.ToStringArray(orig)

	expected := []string{"kind-one"}

	assert.Equal(t, expected, res)
}

func TestNullRawMessage(t *testing.T) {
	// Test with a valid struct
	orig := OriginalStruct{
		Name:   "OrigName",
		Number: 10,
	}
	res := datatransform.ToNullRawMessage(orig)

	expected := json.RawMessage(`{"jname":"OrigName","jnumber":10}`)

	assert.Equal(t, true, res.Valid)
	assert.Equal(t, expected, res.RawMessage)

	// Test with a zero struct
	zero := OriginalStruct{}
	res = datatransform.ToNullRawMessage(zero)

	assert.Equal(t, false, res.Valid)
}
