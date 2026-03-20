// Package integration_test contains cross-component integration tests that
// verify the interaction between multiple packages.
package integration_test

import (
	"errors"
	"testing"

	"github.com/kgjoner/cornucopia/v2/apperr"
	"github.com/kgjoner/cornucopia/v2/prim"
	"github.com/kgjoner/cornucopia/v2/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// userProfile represents a realistic struct combining multiple prim types, used
// to exercise multi-field validation together with the apperr MapError flow.
type userProfile struct {
	Name  string     `validate:"required"`
	Email prim.Email `validate:"required"`
	Phone prim.PhoneNumber
}

// TestValidatorWithPrimEmail verifies that prim.Email's IsValid integrates
// correctly with validator.Validate, producing apperr.Validation errors.
func TestValidatorWithPrimEmail(t *testing.T) {
	err := validator.Validate(prim.Email(""))
	assert.NoError(t, err, "empty email should be valid (not required)")

	err = validator.Validate(prim.Email("not-an-email"))
	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.Validation, appErr.Kind)
	assert.Equal(t, apperr.InvalidData, appErr.Code)

	err = validator.Validate(prim.Email("user@example.com"))
	assert.NoError(t, err, "valid email should pass")
}

// TestValidatorWithPrimPhoneNumber verifies that prim.PhoneNumber's IsValid
// integrates with the validator and produces the right apperr on failure.
func TestValidatorWithPrimPhoneNumber(t *testing.T) {
	err := validator.Validate(prim.PhoneNumber(""))
	assert.NoError(t, err, "empty phone should be valid")

	err = validator.Validate(prim.PhoneNumber("+5511999999999"))
	assert.NoError(t, err, "valid BR phone should pass")

	err = validator.Validate(prim.PhoneNumber("+12025550100"))
	assert.NoError(t, err, "valid US phone should pass")

	err = validator.Validate(prim.PhoneNumber("+44123456789"))
	require.Error(t, err, "unsupported country phone should fail")
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.Validation, appErr.Kind)
}

// TestValidatorWithPrimDocument verifies that prim.Document's IsValid (which
// calls the CPF/CNPJ checksum logic) integrates with the validator.
func TestValidatorWithPrimDocument(t *testing.T) {
	err := validator.Validate(prim.Document(""))
	assert.NoError(t, err, "empty document should be valid")

	validCPF, parseErr := prim.ParseDocument("024.969.460-31")
	require.NoError(t, parseErr)
	err = validator.Validate(validCPF)
	assert.NoError(t, err, "valid CPF should pass")

	// Malformed document (no colon prefix, not 11 or 14 digits)
	err = validator.Validate(prim.Document("unknown:123"))
	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.Validation, appErr.Kind)
}

// TestValidatorMultiFieldStructWithPrimTypes verifies that a struct with
// multiple prim-type fields accumulates all field errors into an apperr.MapError.
func TestValidatorMultiFieldStructWithPrimTypes(t *testing.T) {
	// A zero-value struct is intentionally skipped by the validator unless
	// "required" is explicitly passed as a top-level validation.
	empty := userProfile{}
	err := validator.Validate(empty)
	assert.NoError(t, err, "zero-value struct should be skipped without explicit required")

	// Passing "required" at the top level causes a zero struct to fail.
	err = validator.Validate(empty, "required")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")

	// Partial struct (Name set, Email still empty) → error only on Email.
	partial := userProfile{Name: "Alice"}
	err = validator.Validate(partial)
	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.Validation, appErr.Kind)
	assert.Contains(t, err.Error(), "Email")
	assert.NotContains(t, err.Error(), "Name: ")

	// All required fields filled → no error.
	valid := userProfile{Name: "Alice", Email: prim.Email("alice@example.com")}
	err = validator.Validate(valid)
	assert.NoError(t, err)
}

// TestValidatorWithLocaleEnum verifies that prim.Locale (which implements
// validator.Enum) validates correctly, rejecting unknown locale strings.
func TestValidatorWithLocaleEnum(t *testing.T) {
	err := validator.Validate(prim.Portuguese)
	assert.NoError(t, err)

	err = validator.Validate(prim.English)
	assert.NoError(t, err)

	err = validator.Validate(prim.Locale("klingon"))
	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.InvalidData, appErr.Code)
}

// TestValidatorWithMarketEnum verifies that prim.Market enum validation works
// end-to-end with the validator.
func TestValidatorWithMarketEnum(t *testing.T) {
	err := validator.Validate(prim.MarketBrazil)
	assert.NoError(t, err)

	err = validator.Validate(prim.MarketUSA)
	assert.NoError(t, err)

	err = validator.Validate(prim.Market("mars"))
	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.InvalidData, appErr.Code)
}

// TestValidatorWithCurrencyEnum verifies that prim.Currency enum validation
// integrates with the validator.
func TestValidatorWithCurrencyEnum(t *testing.T) {
	err := validator.Validate(prim.BRL)
	assert.NoError(t, err)

	err = validator.Validate(prim.USD)
	assert.NoError(t, err)

	err = validator.Validate(prim.Currency("DOGE"))
	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.InvalidData, appErr.Code)
}

// TestPriceValidationPipeline exercises the full price flow: creating a price,
// upserting a currency, applying a discount, and reading the values back —
// all using the validator internally via prim.Price methods.
func TestPriceValidationPipeline(t *testing.T) {
	price := prim.NewPrice()

	err := price.UpsertCurrency(prim.BRL, 1000)
	require.NoError(t, err)

	err = price.SetDiscount(200, prim.BRL)
	require.NoError(t, err)

	values, err := price.Values(prim.BRL)
	require.NoError(t, err)
	assert.Equal(t, 1000, values.FullPrice)
	assert.Equal(t, 800, values.SalePrice)

	// Validate the whole price map through the validator.
	err = validator.Validate(price)
	assert.NoError(t, err)

	// Discount that drives SalePrice below the 100-cent minimum should fail.
	err = price.SetDiscount(950, prim.BRL)
	require.Error(t, err)
}

// TestValidatorAddressStruct verifies that a fully-populated prim.Address
// passes struct validation, while a partial one reports missing required fields.
func TestValidatorAddressStruct(t *testing.T) {
	// Empty struct is zero value — validator skips validation for zero structs.
	err := validator.Validate(prim.Address{})
	assert.NoError(t, err, "zero-value address should be skipped by default")

	// Partial address with only some required fields filled.
	partial := prim.Address{
		Line1: "Rua das Flores",
		// Number, City, State, Country, ZipCode all missing
	}
	err = validator.Validate(partial, "required")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Number")
	assert.Contains(t, err.Error(), "City")
	assert.Contains(t, err.Error(), "State")

	// Fully populated valid address.
	full := prim.Address{
		Line1:   "Rua das Flores",
		Number:  "123",
		City:    "São Paulo",
		State:   "SP",
		Country: "BR",
		ZipCode: "01001000",
	}
	err = validator.Validate(full, "required")
	assert.NoError(t, err)
}
