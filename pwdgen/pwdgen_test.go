package pwdgen

import "testing"

func containsAny(s, chars string) bool {
	for _, r := range s {
		for _, c := range chars {
			if r == c {
				return true
			}
		}
	}
	return false
}

func TestGeneratePasswordDefaultSets(t *testing.T) {
	pwd := GeneratePassword(16)
	if len(pwd) != 16 {
		t.Fatalf("expected password length 16, got %d", len(pwd))
	}
	if !containsAny(pwd, "abcdefghijklmnopqrstuvwxyz") {
		t.Fatal("expected at least one lower char")
	}
	if !containsAny(pwd, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Fatal("expected at least one upper char")
	}
	if !containsAny(pwd, "0123456789") {
		t.Fatal("expected at least one number")
	}
	if !containsAny(pwd, "!@#$%&*") {
		t.Fatal("expected at least one special char")
	}
}

func TestGeneratePasswordSpecificSets(t *testing.T) {
	pwd := GeneratePassword(8, Lower, Number)
	if len(pwd) != 8 {
		t.Fatalf("expected password length 8, got %d", len(pwd))
	}
	if !containsAny(pwd, "abcdefghijklmnopqrstuvwxyz") {
		t.Fatal("expected at least one lower char")
	}
	if !containsAny(pwd, "0123456789") {
		t.Fatal("expected at least one number")
	}
}

func TestGenerateNumericCode(t *testing.T) {
	code := GenerateNumericCode(6)
	if len(code) != 6 {
		t.Fatalf("expected length 6, got %d", len(code))
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			t.Fatalf("expected only numeric chars, got %q", code)
		}
	}
}
