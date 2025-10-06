package htypes_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kgjoner/cornucopia/helpers/htypes"
)

type TestStruct struct {
	Address htypes.Address `json:"address,omitempty"`
}

func (t TestStruct) MarshalJSON() ([]byte, error) {
	var pointer *htypes.Address
	if !t.Address.IsZero() {
		pointer = &t.Address
	}

	return json.Marshal(&struct {
		Address *htypes.Address `json:"address,omitempty"`
	}{
		Address: pointer,
	})
}

func TestOptionalAddressInStruct(t *testing.T) {
	empty := TestStruct{}
	data, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("Failed to marshal address: %v", err)
	}

	// Verify that the Address field is omitted when empty
	if string(data) != "{}" {
		t.Errorf("Expected empty JSON object, got %s. Address is zero: %v", data, reflect.ValueOf(empty.Address).IsZero())
	}
}
