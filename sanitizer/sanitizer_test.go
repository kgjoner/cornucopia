package sanitizer

import "testing"

func TestDigit(t *testing.T) {
	if got := Digit("(+55) 11 99876-1234"); got != "5511998761234" {
		t.Fatalf("unexpected digits: %q", got)
	}
	if got := Digit("abc"); got != "" {
		t.Fatalf("expected empty for no digits, got %q", got)
	}
}

func TestIsDigitOnly(t *testing.T) {
	if !IsDigitOnly("012345") {
		t.Fatal("expected string with only digits to be true")
	}
	if IsDigitOnly("12a45") {
		t.Fatal("expected alpha char to make IsDigitOnly false")
	}
	if IsDigitOnly("12 45") {
		t.Fatal("expected whitespace to make IsDigitOnly false")
	}
}
