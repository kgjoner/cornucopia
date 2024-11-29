package htypes

import (
	"strings"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	"github.com/kgjoner/cornucopia/utils/sanitizer"
)

type PhoneNumber string

func ParsePhoneNumber(str string) (*PhoneNumber, error) {
	res := PhoneNumber(str)
	return &res, res.IsValid()
}

func (p *PhoneNumber) IsValid() error {
	if p.IsZero() {
		return nil
	}

	_, err := p.Format()
	return err
}

type FormattedPhoneNumber struct {
	CountryCode string
	AreaCode    string
	Number      string
}

func (p *PhoneNumber) Format() (*FormattedPhoneNumber, error) {
	s := "+" + sanitizer.Digit(p.String())
	*p = PhoneNumber(s)

	if strings.HasPrefix(s, "+55") {
		if len(s) < 13 && len(s) > 14 {
			return nil, normalizederr.NewValidationError("invalid length for Brazil phone")
		}

		return &FormattedPhoneNumber{
			CountryCode: s[1:3],
			AreaCode:    s[3:5],
			Number:      s[5:],
		}, nil

	} else if strings.HasPrefix(s, "+1") {
		if len(s) != 12 {
			return nil, normalizederr.NewValidationError("invalid length for US phone")
		}

		return &FormattedPhoneNumber{
			CountryCode: s[1:2],
			AreaCode:    s[2:5],
			Number:      s[5:],
		}, nil
	}

	return nil, normalizederr.NewValidationError("phone out of region")
}

func (a PhoneNumber) IsZero() bool {
	return a == ""
}

func (a PhoneNumber) String() string {
	return string(a)
}
