package htypes

import (
	"fmt"
	"regexp"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type Email string

func (e Email) IsValid() error {
	str := string(e)
	if str == "" {
		return nil
	}

	doesMatch, err := regexp.MatchString(
		`^[a-z0-9!#$%&'*+/=?^_`+"`{|}~-]+(?:"+`\.[a-z0-9!#$%&'*+/=?^_`+"`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?"+`\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`,
		str,
	)
	if err != nil {
		return err
	} else if doesMatch {
		return nil
	}

	msg := fmt.Sprintf("Must be a valid email.")
	return normalizederr.NewValidationError(msg)
}

func (e Email) IsZero() bool {
	return e == ""
}
