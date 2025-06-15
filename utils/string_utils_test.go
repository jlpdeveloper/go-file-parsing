package utils

import (
	"testing"
)

func TestTrimIfNeeded(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no whitespace",
			input:    "abc",
			expected: "abc",
		},
		{
			name:     "leading whitespace",
			input:    " abc",
			expected: "abc",
		},
		{
			name:     "trailing whitespace",
			input:    "abc ",
			expected: "abc",
		},
		{
			name:     "both leading and trailing whitespace",
			input:    " abc ",
			expected: "abc",
		},
		{
			name:     "multiple whitespace characters",
			input:    "  abc  ",
			expected: "abc",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    " ",
			expected: "",
		},
		{
			name:     "whitespace in the middle",
			input:    "a b c",
			expected: "a b c",
		},
		{
			name:     "leading multi-byte whitespace (non-breaking space U+00A0)",
			input:    "\u00A0abc",
			expected: "abc",
		},
		{
			name:     "trailing multi-byte whitespace (non-breaking space U+00A0)",
			input:    "abc\u00A0",
			expected: "abc",
		},
		{
			name:     "leading and trailing multi-byte whitespace (non-breaking space U+00A0)",
			input:    "\u00A0abc\u00A0",
			expected: "abc",
		},
		{
			name:     "leading multi-byte whitespace (em space U+2003)",
			input:    "\u2003abc",
			expected: "abc",
		},
		{
			name:     "trailing multi-byte whitespace (em space U+2003)",
			input:    "abc\u2003",
			expected: "abc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := TrimIfNeeded(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
