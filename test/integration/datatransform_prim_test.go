package integration_test

import (
	"testing"
	"time"

	"github.com/kgjoner/cornucopia/v2/datatransform"
	"github.com/kgjoner/cornucopia/v2/prim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestToNullTimeWithZeroTime verifies that a zero time.Time produces an
// invalid (null) sql.NullTime.
func TestToNullTimeWithZeroTime(t *testing.T) {
	result := datatransform.ToNullTime(time.Time{})
	assert.False(t, result.Valid, "zero time should produce an invalid NullTime")

	// Matches prim.NullTime zero behaviour — a zero NullTime is a zero time.Time.
	var nt prim.NullTime
	assert.True(t, nt.IsZero())
}

// TestToNullTimeWithValidTime verifies that a non-zero time produces a valid
// sql.NullTime carrying the exact same time value.
func TestToNullTimeWithValidTime(t *testing.T) {
	ts := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	result := datatransform.ToNullTime(ts)
	assert.True(t, result.Valid)
	assert.Equal(t, ts, result.Time)
}

// TestToNullStringWithPrimEmail verifies conversions between prim.Email and
// sql.NullString through datatransform.ToNullString.
func TestToNullStringWithPrimEmail(t *testing.T) {
	email := prim.Email("user@example.com")
	result := datatransform.ToNullString(email.String())
	assert.True(t, result.Valid)
	assert.Equal(t, "user@example.com", result.String)

	empty := prim.Email("")
	result = datatransform.ToNullString(empty.String())
	assert.False(t, result.Valid, "empty email should produce an invalid NullString")
}

// TestToNullStringWithPrimPhoneNumber verifies that a non-zero PhoneNumber
// becomes a valid NullString and a zero one becomes invalid.
func TestToNullStringWithPrimPhoneNumber(t *testing.T) {
	phone := prim.PhoneNumber("+5511999999999")
	result := datatransform.ToNullString(phone.String())
	assert.True(t, result.Valid)
	assert.Equal(t, "+5511999999999", result.String)

	zero := prim.PhoneNumber("")
	result = datatransform.ToNullString(zero.String())
	assert.False(t, result.Valid)
}

// TestToNullRawMessageWithPrimPrice verifies that a populated Price produces a
// valid NullRawMessage and a nil Price produces an invalid one.
func TestToNullRawMessageWithPrimPrice(t *testing.T) {
	price := prim.NewPrice()
	require.NoError(t, price.UpsertCurrency(prim.BRL, 5000))

	result := datatransform.ToNullRawMessage(price)
	assert.True(t, result.Valid)
	assert.NotEmpty(t, result.RawMessage)

	// nil map (zero value for Price) should be treated as null.
	var nilPrice prim.Price
	result = datatransform.ToNullRawMessage(nilPrice)
	assert.False(t, result.Valid)
}

// TestToRawMessageWithPrimDate verifies that prim.Date marshals to a quoted
// YYYY-MM-DD string and a zero Date marshals to the JSON null literal.
func TestToRawMessageWithPrimDate(t *testing.T) {
	d := prim.NewDate(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC))
	raw := datatransform.ToRawMessage(d)
	assert.Equal(t, `"2025-01-15"`, string(raw))

	zeroDate := prim.Date{}
	raw = datatransform.ToRawMessage(zeroDate)
	assert.Equal(t, "null", string(raw))
}

// TestToStringArrayWithPrimLocales verifies that datatransform.ToStringArray
// correctly extracts the underlying string values from prim.Locale slice elements.
func TestToStringArrayWithPrimLocales(t *testing.T) {
	locales := []prim.Locale{prim.Portuguese, prim.English, prim.Spanish}
	result := datatransform.ToStringArray(locales)
	assert.Equal(t, []string{"pt-br", "en-us", "es-es"}, result)
}

// TestToInt64ArrayWithPaginationOffsets verifies that ToInt64Array works with
// integer values produced by prim.Pagination.Offset calculations.
func TestToInt64ArrayWithPaginationOffsets(t *testing.T) {
	pages := []prim.Pagination{
		{Page: 0, Limit: 10},
		{Page: 1, Limit: 10},
		{Page: 2, Limit: 10},
	}

	offsets := make([]int, len(pages))
	for i, p := range pages {
		offsets[i] = p.Offset()
	}

	result := datatransform.ToInt64Array(offsets)
	assert.Equal(t, []int64{0, 10, 20}, result)
}
