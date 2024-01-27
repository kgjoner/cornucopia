package htypes

import (
	"strings"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type PhoneNumber string

func ParsePhoneNumber(str string) (PhoneNumber, error) {
	if str == "" {
		return PhoneNumber(""), nil
	}

	j := 0
	parsedBytes := []byte(str)
	for _, b := range parsedBytes {
		if '0' <= b && b <= '9' {
			parsedBytes[j] = b
			j++
		}
	}

	res := PhoneNumber("+" + string(parsedBytes[:j]))
	return res, res.IsValid()
}

func (p PhoneNumber) IsValid() error {
	if p == "" {
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
			return nil, normalizederr.NewValidationError("Invalid length for Brazil phone.")
		}

		return &PhoneNumberParts{
			CountryCode: s[1:3],
			AreaCode:    s[3:5],
			Number:      s[5:],
		}, nil

	} else if strings.HasPrefix(s, "+1") {
		if len(s) != 12 {
			return nil, normalizederr.NewValidationError("Invalid length for US phone.")
		}

		return &PhoneNumberParts{
			CountryCode: s[1:2],
			AreaCode:    s[2:5],
			Number:      s[5:],
		}, nil
	}

	return nil, normalizederr.NewValidationError("Phone out of region.")
}

func (a PhoneNumber) IsZero() bool {
	return a == ""
}

func (a PhoneNumber) String() string {
	return string(a)
}
