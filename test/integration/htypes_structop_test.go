package integration_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/v2/helpers/htypes"
	"github.com/kgjoner/cornucopia/v2/utils/structop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Email htypes.Email
	Phone htypes.PhoneNumber
}

func TestHtypesStructopIntegration(t *testing.T) {
	target := TestStruct{
		Email: "original@example.com",
		Phone: "+5511999999999",
	}

	// Test with map[string]any (simulating JSON unmarshaling)
	updateMap := map[string]any{
		"Email": "updated@example.com",
		"Phone": "+5511888888888",
	}

	err := structop.New(&target).UpdateViaMap(updateMap)
	require.NoError(t, err, "UpdateViaMap should not return an error")

	// Verify the values were properly updated
	assert.Equal(t, htypes.Email("updated@example.com"), target.Email)
	assert.Equal(t, htypes.PhoneNumber("+5511888888888"), target.Phone)

	// Verify the values were properly parsed and are valid
	assert.NoError(t, target.Email.IsValid(), "Email should be valid")
	assert.NoError(t, target.Phone.IsValid(), "Phone should be valid")
}

func TestHtypesStructopIntegrationWithInvalidData(t *testing.T) {
	target := TestStruct{
		Email: "original@example.com",
		Phone: "+5511999999999",
	}

	// Test with invalid data - this should fail during conversion
	updateMap := map[string]any{
		"Email": "invalid-email",
		"Phone": "+123", // short phone number that will fail validation
	}

	err := structop.New(&target).UpdateViaMap(updateMap)
	assert.Error(t, err, "UpdateViaMap should return an error when validation fails during conversion")

	// TECHNICAL DEBT: Current behavior violates atomicity principle
	// When validation fails, some fields may be updated while others remain unchanged,
	// leaving the object in an inconsistent state. This creates unpredictable behavior
	// depending on field processing order.
	//
	// TODO: Implement proper transaction-like behavior:
	// - Option 1: Validate all fields first, then update (atomic operation)
	// - Option 2: Use backup/restore pattern on validation failure
	// - Option 3: Return detailed results about which fields succeeded/failed
	//
	// For now, documenting the current behavior for integration testing purposes.

	// Current behavior: Email gets updated because it's processed first, even though invalid
	assert.Equal(t, htypes.Email("invalid-email"), target.Email)
	// Phone might not be updated if Email validation fails first (depends on map iteration order)
	// In Go 1.11+, map iteration order is randomized, making this behavior unpredictable

	// Verify the updated email is properly detected as invalid
	assert.Error(t, target.Email.IsValid(), "Email should be invalid")
}
