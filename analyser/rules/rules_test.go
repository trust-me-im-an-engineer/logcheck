package rules

import (
	"testing"
)

func TestStartsWithLowercase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid lowercase", "hello", true},
		{"Valid single char", "z", true},
		{"Invalid uppercase", "Hello", false},
		{"Invalid number", "1hello", false},
		{"Invalid space", " hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StartsWithLowercase(tt.input); got != tt.expected {
				t.Errorf("StartsWithLowercase(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIndexIllegalCharacter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"All valid lowercase", "abc 123", -1},
		{"All valid mixed", "aB 9", -1},
		{"Illegal symbol at start", "!abc", 0},
		{"Illegal symbol in middle", "abc@123", 3},
		{"Illegal newline", "abc\n", 3},
		{"Empty string", "", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexIllegalCharacter(tt.input); got != tt.expected {
				t.Errorf("IndexIllegalCharacter(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFindSensitiveName(t *testing.T) {
	keywords := map[string]struct{}{
		"admin": {},
		"root":  {},
	}

	tests := []struct {
		name        string
		input       string
		wantIndex   int
		wantKeyword string
	}{
		{"No sensitive words", "user123", -1, ""},
		{"Exact match", "admin", 0, "admin"},
		{"Case insensitive match", "AdMiN", 0, "admin"},
		{"Match in middle", "the_root_user", 4, "root"},
		{"Empty string", "", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, gotKw := FindSensitiveName(tt.input, keywords)
			if gotIdx != tt.wantIndex || gotKw != tt.wantKeyword {
				t.Errorf("FindSensitiveName(%q) = (%v, %q); want (%v, %q)",
					tt.input, gotIdx, gotKw, tt.wantIndex, tt.wantKeyword)
			}
		})
	}
}
