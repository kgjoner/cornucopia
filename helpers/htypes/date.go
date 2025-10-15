package htypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Date represents a date-only value that marshals/unmarshals to/from YYYY-MM-DD format
type Date struct {
	time.Time
}

const DateFormat = "2006-01-02"

func NewDate(t time.Time) Date {
	return Date{Time: t}
}

// ParseDate creates a new Date from a string in YYYY-MM-DD format
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return Date{}, err
	}
	return Date{Time: t}, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format(DateFormat))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	// Handle null values
	if string(data) == "null" {
		*d = Date{}
		return nil
	}

	// Remove quotes from JSON string
	str := strings.Trim(string(data), `"`)

	// Handle empty string
	if str == "" {
		*d = Date{}
		return nil
	}

	// Parse the date
	t, err := time.Parse(DateFormat, str)
	if err != nil {
		return fmt.Errorf("cannot parse date %q: %w", str, err)
	}

	*d = Date{Time: t}
	return nil
}

func (d Date) String() string {
	if d.Time.IsZero() {
		return ""
	}
	return d.Time.Format(DateFormat)
}

func (d Date) IsZero() bool {
	return d.Time.IsZero()
}

func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format(DateFormat), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = Date{Time: v}
		return nil
	case string:
		if v == "" {
			*d = Date{}
			return nil
		}
		t, err := time.Parse(DateFormat, v)
		if err != nil {
			return fmt.Errorf("cannot parse date %q: %w", v, err)
		}
		*d = Date{Time: t}
		return nil
	case []byte:
		return d.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %T into Date", value)
	}
}
