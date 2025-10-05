package apperr

import (
	"fmt"
	"strings"
)

type MapError struct {
	details map[string]string
}

func (e *MapError) Error() string {
	if len(e.details) == 0 {
		return ""
	}
	var sb strings.Builder

	// Write each field error on a new line.
	for field, msg := range e.details {
		sb.WriteString(fmt.Sprintf("\n- %s: %s", field, msg))
	}

	return sb.String()
}

func NewMapError(details map[string]string) error {
	return &MapError{details: details}
}
