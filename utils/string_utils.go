package utils

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TrimIfNeeded checks if a string needs trimming by checking if the first or last character is whitespace.
// It only calls strings.TrimSpace() if necessary, which is more efficient than always trimming.
func TrimIfNeeded(s string) string {
	if len(s) == 0 {
		return s
	}

	firstRune, _ := utf8.DecodeRuneInString(s)
	lastRune, _ := utf8.DecodeLastRuneInString(s)

	if unicode.IsSpace(firstRune) || unicode.IsSpace(lastRune) {
		return strings.TrimSpace(s)
	}

	return s
}

func TrimTrailingDecimal(s *string) {
	b := []byte(*s)
	if len(b) > 2 && b[len(b)-2] == '.' && b[len(b)-1] == '0' {
		b = b[:len(b)-2] // "truncate" the last two bytes
	}
	// Only convert to string now, if needed
	*s = string(b)
}

func FormattedStringToInt(s *string) (int, error) {
	TrimTrailingDecimal(s)
	return strconv.Atoi(*s)
}
