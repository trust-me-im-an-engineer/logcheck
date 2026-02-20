package rules

import (
	"strings"
)

// StartsWithLowercase validates rule 1.
// msg must have at least 1 byte!
func StartsWithLowercase(msg string) bool {
	return msg[0] >= 'a' && msg[0] <= 'z'
}

// IndexIllegalCharacter validates rules 2 and 3.
// It returns position of first illegal character in msg, or -1 if all characters are legal.
// Only English letters, numbers and spaces considered legal [a-zA-Z0-9 ].
func IndexIllegalCharacter(msg string) int {
	for i := range msg {
		b := msg[i]
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == ' ') {
			return i
		}
	}
	return -1
}

// FindSensitiveName validates rule 4.
// It returns position of first sensitive keyword in name, or (-1, "") if no sensitive keywords found.
func FindSensitiveName(name string, sensitiveKeywords map[string]struct{}) (int, string) {
	lowerName := strings.ToLower(name)
	for kw := range sensitiveKeywords {
		if i := strings.Index(lowerName, kw); i != -1 {
			return i, kw
		}
	}
	return -1, ""
}
