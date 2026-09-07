package probe

import "testing"

func TestHasANSI(t *testing.T) {
	if HasANSI([]byte("plain text")) {
		t.Error("plain text should not have ANSI")
	}
	if !HasANSI([]byte("\x1b[31mred\x1b[0m")) {
		t.Error("colored text should have ANSI")
	}
	if !HasANSI([]byte("\x1b[1mbold\x1b[0m")) {
		t.Error("bold text should have ANSI")
	}
}

func TestStripANSI(t *testing.T) {
	input := []byte("\x1b[1;31mError:\x1b[0m file not found")
	got := string(StripANSI(input))
	want := "Error: file not found"
	if got != want {
		t.Errorf("StripANSI = %q, want %q", got, want)
	}
}

func TestHasColorCodes(t *testing.T) {
	if HasColorCodes([]byte("\x1b[1mbold only\x1b[0m")) {
		t.Error("bold-only should not count as color")
	}
	if !HasColorCodes([]byte("\x1b[31mred\x1b[0m")) {
		t.Error("SGR 31 is a color code")
	}
	if !HasColorCodes([]byte("\x1b[38;5;196mred\x1b[0m")) {
		t.Error("8-bit color should be detected")
	}
	if !HasColorCodes([]byte("\x1b[38;2;255;0;0mred\x1b[0m")) {
		t.Error("24-bit color should be detected")
	}
}

func TestClassifyColors(t *testing.T) {
	input := []byte("\x1b[31mred\x1b[0m \x1b[38;5;196m256red\x1b[0m \x1b[38;2;0;255;0mtrue\x1b[0m")
	c := ClassifyColors(input)
	if c.FourBit < 1 {
		t.Error("expected at least 1 four-bit color")
	}
	if c.EightBit < 1 {
		t.Error("expected at least 1 eight-bit color")
	}
	if c.TwentyFourBit < 1 {
		t.Error("expected at least 1 twenty-four-bit color")
	}
}
