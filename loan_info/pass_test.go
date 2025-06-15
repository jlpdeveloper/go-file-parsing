package loan_info

import (
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"testing"
)

func TestPassExtraData(t *testing.T) {
	testCases := []struct {
		name     string
		cols     []string
		expected map[string]string
	}{
		{
			name: "all fields present with Joint App",
			cols: createColumnsWithValues(map[int]string{
				69: "Joint App",
				70: "100000",
				73: "2",
				74: "5000",
				92: "3000",
			}),
			expected: map[string]string{
				"application_type": "Joint App",
				"annual_inc_joint": "100000",
				"acc_now_delinq":   "2",
				"tot_coll_amt":     "5000",
				"avg_cur_bal":      "3000",
			},
		},
		{
			name: "all fields present with Individual application",
			cols: createColumnsWithValues(map[int]string{
				69: "Individual",
				70: "100000", // This should not be included in the result
				73: "0",
				74: "1000",
				92: "2500",
			}),
			expected: map[string]string{
				"application_type": "Individual",
				"acc_now_delinq":   "0",
				"tot_coll_amt":     "1000",
				"avg_cur_bal":      "2500",
			},
		},
		{
			name: "some fields missing",
			cols: createColumnsWithValues(map[int]string{
				69: "Individual",
				92: "2500",
			}),
			expected: map[string]string{
				"application_type": "Individual",
				"avg_cur_bal":      "2500",
			},
		},
		{
			name: "empty fields",
			cols: createColumnsWithValues(map[int]string{
				69: "",
				70: "",
				73: "",
				74: "",
				92: "",
			}),
			expected: map[string]string{},
		},
		{
			name:     "not enough columns",
			cols:     make([]string, 50), // Less than the required columns
			expected: map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
			}

			result, err := passExtraData(ctx, tc.cols)

			// We don't expect errors from this function
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check if the result has the expected number of entries
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d entries in result, got %d", len(tc.expected), len(result))
			}

			// Check if all expected entries are in the result with the correct values
			for key, expectedValue := range tc.expected {
				if result[key] != expectedValue {
					t.Errorf("expected %s to be '%s', got '%s'", key, expectedValue, result[key])
				}
			}
		})
	}
}

// Helper function to create a slice of columns with specific values at specific indices
func createColumnsWithValues(values map[int]string) []string {
	// Find the maximum index
	maxIndex := 0
	for idx := range values {
		if idx > maxIndex {
			maxIndex = idx
		}
	}

	// Create a slice with enough capacity
	cols := make([]string, maxIndex+1)

	// Set the values
	for idx, val := range values {
		cols[idx] = val
	}

	return cols
}
