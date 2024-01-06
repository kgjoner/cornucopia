package htypes

import (
	"time"
)

type NullTime struct {
	time.Time
}

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`null`), nil
	}

	return t.Time.MarshalJSON()
}
