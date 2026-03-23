package prim

import (
	"strings"

	"github.com/kgjoner/cornucopia/v3/apperr"
	"github.com/kgjoner/cornucopia/v3/sanitizer"
)

type PhoneNumber string

func ParsePhoneNumber(str string) (PhoneNumber, error) {
	if str == "" {
		return "", nil
	}

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
			return nil, apperr.NewValidationError("invalid length for Brazil phone")
		}

		return &PhoneNumberParts{
			CountryCode: s[1:3],
			AreaCode:    s[3:5],
			Number:      s[5:],
		}, nil

	} else if strings.HasPrefix(s, "+1") {
		if len(s) != 12 {
			return nil, apperr.NewValidationError("invalid length for US phone")
		}

		return &PhoneNumberParts{
			CountryCode: s[1:2],
			AreaCode:    s[2:5],
			Number:      s[5:],
		}, nil
	}

	return nil, apperr.NewValidationError("phone out of region")
}

func (p PhoneNumber) IsZero() bool {
	return p == ""
}

func (p PhoneNumber) String() string {
	return string(p)
}

func (p *PhoneNumber) UnmarshalText(text []byte) error {
	parsed, err := ParsePhoneNumber(string(text))
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}
