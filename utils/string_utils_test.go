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