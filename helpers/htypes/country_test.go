package htypes_test

import (
	"testing"

	"github.com/kgjoner/cornucopia/helpers/htypes"
)

func TestCountry(t *testing.T) {
	tests := []struct {
		name    string
		country htypes.Country
		wantErr bool
	}{
		{"Valid Country", "BR", false},
		{"Valid Country", "brazil", false},
		{"Invalid Country", "XX", true},
		{"Invalid Country", "Bra", true},
		{"Zero Country", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.country.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("Country.IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.country.IsZero() && tt.country.Name() == "" && !tt.wantErr {
				t.Errorf("Country.Name() returned empty for non-zero country")
			}
		})
	}
}