package htypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// TIME
type NullTime struct {
	time.Time
}

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`null`), nil
	}

	return t.Time.MarshalJSON()
}

func (t *NullTime) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case time.Time:
		t.Time = src
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, []byte{})
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (t NullTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

// RAW MESSAGE
type NullRawMessage struct {
	RawMessage json.RawMessage
	Valid      bool
}

func (n *NullRawMessage) Scan(src interface{}) error {
	if src == nil {
		n.Valid = false
		return nil
	}
	switch src := src.(type) {
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		n.RawMessage, n.Valid = srcCopy, true
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, []byte{})
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (n NullRawMessage) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return []byte(n.RawMessage), nil
}
