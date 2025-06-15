package utils

import (
	"strings"
	"unicode"
)

// TrimIfNeeded checks if a string needs trimming by checking if the first or last character is whitespace.
// It only calls strings.TrimSpace() if necessary, which is more efficient than always trimming.
func TrimIfNeeded(s string) string {
	if len(s) == 0 {
		return s
	}
	
	if unicode.IsSpace(rune(s[0])) || unicode.IsSpace(rune(s[len(s)-1])) {
		return strings.TrimSpace(s)
	}
	
	return s
}