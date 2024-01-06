package datatransform_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/utils/datatransform"
	"github.com/stretchr/testify/assert"
)

type Kind string

func TestStringArray(t *testing.T) {
	orig := []Kind{"kind-one"}
	res := datatransform.ToStringArray(orig)

	expected := []string{"kind-one"}

	assert.Equal(t, expected, res)
}
