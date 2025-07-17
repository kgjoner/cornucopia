package htypes

import (
	"encoding/json"
	"net/mail"
	"strings"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type Email string

func ParseEmail(str string) (Email, error) {
	if str == "" {
		return "", nil
	}

	email := Email(strings.ToLower(str))
	return email, email.IsValid()
}

func (e Email) IsValid() error {
	if e.IsZero() {
		return nil
	}

	str := string(e)

	// Check if email is in lowercase format
	if str != strings.ToLower(str) {
		return normalizederr.NewValidationError("email must be in lowercase format; use ParseEmail to normalize")
	}

	// Fast basic checks
	if !strings.Contains(str, "@") ||
		strings.Count(str, "@") != 1 ||
		strings.HasPrefix(str, "@") ||
		strings.HasSuffix(str, "@") ||
		!strings.Contains(str, ".") {
		return normalizederr.NewValidationError("must be a valid email")
	}

	// Use standard library for comprehensive validation
	_, err := mail.ParseAddress(str)
	if err != nil {
		return normalizederr.NewValidationError("must be a valid email")
	}

	return nil
}

func (e Email) IsZero() bool {
	return e == ""
}

func (e Email) String() string {
	return string(e)
}

func (e *Email) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*e, err = ParseEmail(s)
	return err
}

// Deprecated: Use ParseEmail instead
//
// Turn all letters to lowercase
func (e Email) Normalize() Email {
	return Email(strings.ToLower(string(e)))
}
