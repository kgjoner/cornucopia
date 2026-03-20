package integration_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/v2/prim"
	"github.com/kgjoner/cornucopia/v2/sanitizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSanitizerDigitThenParsePhoneNumber verifies that sanitizer.Digit strips
// formatting characters from a raw phone string so that prim.ParsePhoneNumber
// can parse it into a canonical PhoneNumber value.
func TestSanitizerDigitThenParsePhoneNumber(t *testing.T) {
	cases := []struct {
		raw      string
		expected prim.PhoneNumber
	}{
		{"+55 (11) 99999-9999", "+5511999999999"},
		{"55 11 99999 9999", "+5511999999999"},
		{"(55)11999999999", "+5511999999999"},
	}

	for _, tc := range cases {
		digits := sanitizer.Digit(tc.raw)
		assert.True(t, sanitizer.IsDigitOnly(digits), "Digit() result should contain only digits for input %q", tc.raw)

		phone, err := prim.ParsePhoneNumber(tc.raw)
		require.NoError(t, err, "ParsePhoneNumber should succeed for %q", tc.raw)
		assert.Equal(t, tc.expected, phone)
		assert.Equal(t, prim.PhoneNumber("+"+digits), phone)
	}
}

// TestSanitizerDigitThenParseCPF verifies the sanitizer + CPF document parse pipeline.
func TestSanitizerDigitThenParseCPF(t *testing.T) {
	rawCPF := "024.969.460-31"

	digits := sanitizer.Digit(rawCPF)
	assert.Equal(t, "02496946031", digits)
	assert.True(t, sanitizer.IsDigitOnly(digits))

	doc, err := prim.ParseDocument(rawCPF)
	require.NoError(t, err)
	assert.Equal(t, "cpf:"+digits, doc.String())
}

// TestSanitizerDigitThenParseCNPJ verifies the sanitizer + CNPJ document parse pipeline.
func TestSanitizerDigitThenParseCNPJ(t *testing.T) {
	rawCNPJ := "86.978.987/0001-39"

	digits := sanitizer.Digit(rawCNPJ)
	assert.Equal(t, "86978987000139", digits)
	assert.True(t, sanitizer.IsDigitOnly(digits))

	doc, err := prim.ParseDocument(rawCNPJ)
	require.NoError(t, err)
	assert.Equal(t, "cnpj:"+digits, doc.String())
}

// TestSanitizerIsDigitOnlyWithDocumentParts verifies that document numbers
// stored inside a prim.Document consist solely of digit characters.
func TestSanitizerIsDigitOnlyWithDocumentParts(t *testing.T) {
	doc, err := prim.ParseDocument("024.969.460-31")
	require.NoError(t, err)

	parts, err := doc.Parts()
	require.NoError(t, err)

	assert.True(t, sanitizer.IsDigitOnly(parts.Number), "CPF number stored in Document must be all digits")
}

// TestSanitizerEmptyInputsAreHandledGracefully ensures that both functions
// behave correctly with empty strings.
func TestSanitizerEmptyInputsAreHandledGracefully(t *testing.T) {
	assert.Equal(t, "", sanitizer.Digit(""))
	assert.True(t, sanitizer.IsDigitOnly(""))

	phone, err := prim.ParsePhoneNumber("")
	require.NoError(t, err)
	assert.True(t, phone.IsZero())

	doc, err := prim.ParseDocument("")
	require.NoError(t, err)
	assert.True(t, doc.IsZero())
}
