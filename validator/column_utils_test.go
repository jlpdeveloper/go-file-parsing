package validator

import (
	"reflect"
	"testing"
)

func TestPreprocessColumns(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no whitespace",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "leading whitespace",
			input:    []string{" a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "trailing whitespace",
			input:    []string{"a ", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "both leading and trailing whitespace",
			input:    []string{" a ", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "multiple whitespace characters",
			input:    []string{"  a  ", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty string",
			input:    []string{"", "b", "c"},
			expected: []string{"", "b", "c"},
		},
		{
			name:     "only whitespace",
			input:    []string{" ", "b", "c"},
			expected: []string{"", "b", "c"},
		},
		{
			name:     "mixed",
			input:    []string{"a", " b ", "  c  "},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PreprocessColumns(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
