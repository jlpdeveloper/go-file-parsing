package loan_info

import (
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"testing"
)

func TestHasValidLoanAmount(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid loan amounts",
			cols:    []string{"id", "name", "1000", "500", "500"},
			wantErr: false,
		},
		{
			name:    "non-numeric loan amount",
			cols:    []string{"id", "name", "abc", "500", "500"},
			wantErr: true,
			errMsg:  "loan amount is not a number",
		},
		{
			name:    "non-positive loan amount",
			cols:    []string{"id", "name", "0", "500", "500"},
			wantErr: true,
			errMsg:  "loan amount is not a positive number",
		},
		{
			name:    "negative loan amount",
			cols:    []string{"id", "name", "-100", "500", "500"},
			wantErr: true,
			errMsg:  "loan amount is not a positive number",
		},
		{
			name:    "non-numeric funding amount",
			cols:    []string{"id", "name", "1000", "abc", "500"},
			wantErr: true,
			errMsg:  "funding amount is not a number",
		},
		{
			name:    "non-positive funding amount",
			cols:    []string{"id", "name", "1000", "0", "500"},
			wantErr: true,
			errMsg:  "funding amount is not a positive number",
		},
		{
			name:    "negative funding amount",
			cols:    []string{"id", "name", "1000", "-100", "500"},
			wantErr: true,
			errMsg:  "funding amount is not a positive number",
		},
		{
			name:    "non-numeric funding inv amount",
			cols:    []string{"id", "name", "1000", "500", "abc"},
			wantErr: true,
			errMsg:  "funding inv amt is not a number",
		},
		{
			name:    "non-positive funding inv amount",
			cols:    []string{"id", "name", "1000", "500", "0"},
			wantErr: true,
			errMsg:  "funding inv amt is not a positive number",
		},
		{
			name:    "negative funding inv amount",
			cols:    []string{"id", "name", "1000", "500", "-100"},
			wantErr: true,
			errMsg:  "funding inv amt is not a positive number",
		},
		{
			name:    "funding inv amount not equal to funding amount",
			cols:    []string{"id", "name", "1000", "500", "600"},
			wantErr: true,
			errMsg:  "funding inv amt is not equal to funding amount",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
			}

			result, err := hasValidLoanAmount(ctx, tc.cols)

			// Check if error was expected
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if err.Error() != tc.errMsg {
					t.Errorf("expected error message '%s', got '%s'", tc.errMsg, err.Error())
				}
				return
			}

			// If no error was expected, check the result
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify the returned map contains the expected values
			if result["loanAmount"] != tc.cols[2] {
				t.Errorf("expected loanAmount '%s', got '%s'", tc.cols[2], result["loanAmount"])
			}
			if result["fundingAmount"] != tc.cols[3] {
				t.Errorf("expected fundingAmount '%s', got '%s'", tc.cols[3], result["fundingAmount"])
			}
			if result["fundingInvAmt"] != tc.cols[4] {
				t.Errorf("expected fundingInvAmt '%s', got '%s'", tc.cols[4], result["fundingInvAmt"])
			}
		})
	}
}

func TestHasValidInterestRate(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid interest rate",
			cols:    []string{"id", "name", "1000", "500", "500", "term", "10"},
			wantErr: false,
		},
		{
			name:    "valid interest rate at minimum",
			cols:    []string{"id", "name", "1000", "500", "500", "term", "05"},
			wantErr: false,
		},
		{
			name:    "valid interest rate at maximum",
			cols:    []string{"id", "name", "1000", "500", "500", "term", "35"},
			wantErr: false,
		},
		{
			name:    "non-numeric interest rate",
			cols:    []string{"id", "name", "1000", "500", "500", "term", "abc"},
			wantErr: true,
			errMsg:  "interest rate is not a number",
		},
		{
			name:    "interest rate below minimum",
			cols:    []string{"id", "name", "1000", "500", "500", "term", "04"},
			wantErr: true,
			errMsg:  "interest rate is not between 5% and 35%",
		},
		{
			name:    "interest rate above maximum",
			cols:    []string{"id", "name", "1000", "500", "500", "term", "36"},
			wantErr: true,
			errMsg:  "interest rate is not between 5% and 35%",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
			}

			result, err := hasValidInterestRate(ctx, tc.cols)

			// Check if error was expected
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if err.Error() != tc.errMsg {
					t.Errorf("expected error message '%s', got '%s'", tc.errMsg, err.Error())
				}
				return
			}

			// If no error was expected, check the result
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify the returned map contains the expected values
			if result["interestRate"] != tc.cols[6] {
				t.Errorf("expected interestRate '%s', got '%s'", tc.cols[6], result["interestRate"])
			}
		})
	}
}
