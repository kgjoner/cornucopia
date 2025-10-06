package htypes_test

import (
	"encoding/json"
	"testing"

	"github.com/kgjoner/cornucopia/v2/helpers/htypes"
)

func TestCountry(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		wantValid        bool
		wantNormalized   string
		wantName         string
		wantUnmarshalErr bool
	}{
		{
			name:             "Valid country code",
			input:            "BR",
			wantNormalized:   "BR",
			wantValid:        true,
			wantName:         "Brazil",
			wantUnmarshalErr: false,
		},
		{
			name:             "Valid country name (exact case)",
			input:            "Brazil",
			wantNormalized:   "BR",
			wantValid:        true,
			wantName:         "Brazil",
			wantUnmarshalErr: false,
		},
		{
			name:             "Valid country name (lowercase)",
			input:            "brazil",
			wantNormalized:   "BR",
			wantValid:        true,
			wantName:         "Brazil",
			wantUnmarshalErr: false,
		},
		{
			name:             "Invalid country code",
			input:            "XX",
			wantValid:        false,
			wantNormalized:   "",
			wantName:         "",
			wantUnmarshalErr: true,
		},
		{
			name:             "Invalid country name",
			input:            "Bra",
			wantValid:        false,
			wantNormalized:   "",
			wantName:         "",
			wantUnmarshalErr: true,
		},
		{
			name:             "Empty country",
			input:            "",
			wantValid:        true,
			wantNormalized:   "",
			wantName:         "",
			wantUnmarshalErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, err := htypes.ParseCountry(tt.input)
			if (err == nil) != tt.wantValid {
				t.Errorf("Country.IsValid() error = %v, wantValid %v", err, tt.wantValid)
			}

			if string(country) != tt.wantNormalized {
				t.Errorf("Country after validation = %q, want %q", string(country), tt.wantNormalized)
			}

			// Test Name()
			if name := country.Name(); name != tt.wantName {
				t.Errorf("Country.Name() = %q, want %q", name, tt.wantName)
			}

			// Test IsZero()
			if zero := country.IsZero(); zero != (tt.wantNormalized == "") {
				t.Errorf("Country.IsZero() = %v, want %v", zero, tt.wantNormalized == "")
			}

			// Test MarshalJSON()
			jsonData, err := json.Marshal(country)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}
			expectedJSON := `"` + tt.wantName + `"`
			if string(jsonData) != expectedJSON {
				t.Errorf("json.Marshal() = %s, want %s", jsonData, expectedJSON)
			}

			// Test UnmarshalJSON()
			var unmarshalledCountry htypes.Country
			err = json.Unmarshal([]byte(`"`+tt.input+`"`), &unmarshalledCountry)
			if (err != nil) != tt.wantUnmarshalErr {
				t.Errorf("json.Unmarshal() error = %v, wantUnmarshalErr %v", err, tt.wantUnmarshalErr)
			}

			if !tt.wantUnmarshalErr && string(unmarshalledCountry) != tt.wantNormalized {
				t.Errorf("json.Unmarshal() result = %q, want %q", string(unmarshalledCountry), tt.wantNormalized)
			}
		})
	}
}

func TestCountryEdgeCases(t *testing.T) {
	// Test invalid JSON format
	t.Run("Invalid JSON", func(t *testing.T) {
		var country htypes.Country
		err := json.Unmarshal([]byte(`"unclosed`), &country)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})

	// Test JSON round-trip
	t.Run("JSON round-trip", func(t *testing.T) {
		original := htypes.Country("BR")

		// Marshal
		jsonData, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		// Unmarshal
		var roundTrip htypes.Country
		err = json.Unmarshal(jsonData, &roundTrip)
		if err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if roundTrip != original {
			t.Errorf("Round-trip failed: got %q, want %q", roundTrip, original)
		}
	})
}
