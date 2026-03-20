package dbhandler

import (
	"testing"
	"time"
)

type scanTarget struct {
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

func TestStructScan_NormalizesKeyAndTimestamp(t *testing.T) {
	var dst scanTarget
	scanner := Struct(&dst)

	err := scanner.Scan([]byte(`{"created_at":"2024-01-02T03:04:05","name":"john"}`))
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if dst.Name != "john" {
		t.Fatalf("expected name john, got %q", dst.Name)
	}
	if dst.CreatedAt.IsZero() {
		t.Fatal("expected createdAt to be parsed")
	}
}

func TestStructArrayScan_NormalizesEachElement(t *testing.T) {
	var dst []scanTarget
	scanner := StructArray(&dst)

	err := scanner.Scan([]byte(`[{"created_at":"2024-01-02T03:04:05","name":"a"},{"created_at":"2024-01-03T03:04:05","name":"b"}]`))
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if len(dst) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(dst))
	}
	if dst[0].Name != "a" || dst[1].Name != "b" {
		t.Fatalf("unexpected names: %+v", dst)
	}
}

func TestParseJSONArray(t *testing.T) {
	elems, err := parseJSONArray([]byte(`[{"a":1},{"b":2}]`))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(elems) != 2 {
		t.Fatalf("expected 2 elems, got %d", len(elems))
	}
}
