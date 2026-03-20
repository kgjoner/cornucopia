package dbhandler

import "testing"

func TestNormalizeMapKeys(t *testing.T) {
	in := map[string]any{
		"created_at": map[string]any{
			"updated_at": "x",
		},
	}

	out, ok := normalizeMapKeys(in).(map[string]any)
	if !ok {
		t.Fatalf("expected map output")
	}
	if _, exists := out["createdat"]; !exists {
		t.Fatalf("expected key createdat, got %#v", out)
	}
}

func TestNormalizeTimeString(t *testing.T) {
	got := normalizeTimeString("2024-01-02T03:04:05")
	if got != "2024-01-02T03:04:05Z" {
		t.Fatalf("unexpected normalized time: %q", got)
	}

	unchanged := normalizeTimeString("2024-01-02")
	if unchanged != "2024-01-02" {
		t.Fatalf("unexpected changed date string: %q", unchanged)
	}
}
