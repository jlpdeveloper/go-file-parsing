package loan_info

import (
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"strings"
	"testing"
)

func TestValidate_ColumnCountValidation(t *testing.T) {
	testCases := []struct {
		name            string
		row             string
		expectedColumns int
		wantErr         bool
	}{
		{
			name:            "correct column count",
			row:             "a,b,c",
			expectedColumns: 3,
			wantErr:         false,
		},
		{
			name:            "too few columns",
			row:             "a,b",
			expectedColumns: 3,
			wantErr:         true,
		},
		{
			name:            "too many columns",
			row:             "a,b,c,d",
			expectedColumns: 3,
			wantErr:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New(
				&config.ParserConfig{
					Delimiter:       ",",
					ExpectedColumns: tc.expectedColumns,
				},
				&validator.MockCache{},
				[]validator.ColValidator{isValidSize})

			id, err := v.Validate(tc.row)

			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}

			// Check that the ID (first column) is correctly returned
			expectedID := ""
			if len(tc.row) > 0 {
				parts := strings.Split(tc.row, ",")
				if len(parts) > 0 {
					expectedID = parts[0]
				}
			}

			if id != expectedID {
				t.Errorf("expected ID to be '%s', got '%s'", expectedID, id)
			}
		})
	}
}
