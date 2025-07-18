package htypes

import (
	"encoding/json"
	"strings"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	"github.com/kgjoner/cornucopia/utils/sanitizer"
)

type PhoneNumber string

func ParsePhoneNumber(str string) (PhoneNumber, error) {
	s := "+" + sanitizer.Digit(str)
	p := PhoneNumber(s)
	return p, p.IsValid()
}

func (p PhoneNumber) IsValid() error {
	if p.IsZero() {
		return nil
	}

	_, err := p.Parts()
	return err
}

type PhoneNumberParts struct {
	CountryCode string
	AreaCode    string
	Number      string
}

func (p PhoneNumber) Parts() (*PhoneNumberParts, error) {
	s := string(p)
	if strings.HasPrefix(s, "+55") {
		if len(s) < 13 && len(s) > 14 {
			return nil, normalizederr.NewValidationError("invalid length for Brazil phone")
		}

		return &PhoneNumberParts{
			CountryCode: s[1:3],
			AreaCode:    s[3:5],
			Number:      s[5:],
		}, nil

	} else if strings.HasPrefix(s, "+1") {
		if len(s) != 12 {
			return nil, normalizederr.NewValidationError("invalid length for US phone")
		}

		return &PhoneNumberParts{
			CountryCode: s[1:2],
			AreaCode:    s[2:5],
			Number:      s[5:],
		}, nil
	}

	return nil, normalizederr.NewValidationError("phone out of region")
}

func (p PhoneNumber) IsZero() bool {
	return p == ""
}

func (p PhoneNumber) String() string {
	return string(p)
}

func (p *PhoneNumber) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*p, err = ParsePhoneNumber(s)
	return err
}
