package hash

import "testing"

func TestFromDeterministic(t *testing.T) {
	h1 := From("abc", 123, true)
	h2 := From("abc", 123, true)
	if h1 != h2 {
		t.Fatalf("expected deterministic hash, got %q and %q", h1, h2)
	}
	if len(h1) != 64 {
		t.Fatalf("expected SHA256 hex length 64, got %d", len(h1))
	}
}

func TestFromDifferentInputs(t *testing.T) {
	h1 := From("abc")
	h2 := From("abcd")
	if h1 == h2 {
		t.Fatalf("expected different inputs to produce different hashes")
	}
}
