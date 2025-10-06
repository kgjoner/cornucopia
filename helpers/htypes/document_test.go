package htypes_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/v2/helpers/htypes"
	"github.com/stretchr/testify/assert"
)

func TestDocumentParsing(t *testing.T) {
	// Simple CPF
	unformatted := "024.969.460-31"
	doc, err := htypes.ParseDocument(unformatted)

	assert.Nil(t, err)
	assert.Equal(t, "cpf:02496946031", doc.String())

	// Simple CNPJ
	unformatted = "86.978.987/0001-39"
	doc, err = htypes.ParseDocument(unformatted)

	assert.Nil(t, err)
	assert.Equal(t, "cnpj:86978987000139", doc.String())

	// Empty
	unformatted = ""
	doc, err = htypes.ParseDocument(unformatted)

	assert.Nil(t, err)
	assert.Equal(t, "", doc.String())
}
