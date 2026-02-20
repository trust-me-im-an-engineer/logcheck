package analyser

import (
	"strings"
)

// startsWithLowercase validates rule 1.
func startsWithLowercase(msg string) bool {
	return msg[0] >= 'a' && msg[0] <= 'z'
}

// indexIllegalCharacter validates rules 2 and 3.
// It returns position of first illegal character in msg, or -1 if all characters are legal.
// Only English letters, numbers and spaces considered legal [a-zA-Z0-9 ].
func indexIllegalCharacter(msg string) int {
	for i := range msg {
		b := msg[i]
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == ' ') {
			return i
		}
	}
	return -1
}

// indexSensitiveName validates rule 4.
// It returns position of first sensitive keyword in name, or (-1, "") if no sensitive keywords found.
func indexSensitiveName(name string, sensitiveKeywords map[string]struct{}) (int, string) {
	for kw := range sensitiveKeywords {
		if i := strings.Index(strings.ToLower(name), kw); i != -1 {
			return i, kw
		}
	}
	return -1, ""
}
