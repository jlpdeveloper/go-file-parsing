package loan_info

import (
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"strconv"
	"strings"
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

func TestHasValidTerm(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid term - middle of range",
			cols:    []string{"id", "name", "1000", "500", "500", "12 months"},
			wantErr: false,
		},
		{
			name:    "valid term - minimum",
			cols:    []string{"id", "name", "1000", "500", "500", "12 months"},
			wantErr: false,
		},
		{
			name:    "valid term - maximum",
			cols:    []string{"id", "name", "1000", "500", "500", "72 months"},
			wantErr: false,
		},
		{
			name:    "valid term - with extra spaces",
			cols:    []string{"id", "name", "1000", "500", "500", "  24 months  "},
			wantErr: false,
		},
		{
			name:    "non-numeric term",
			cols:    []string{"id", "name", "1000", "500", "500", "abc months"},
			wantErr: true,
			errMsg:  "term is not a number",
		},
		{
			name:    "term below minimum",
			cols:    []string{"id", "name", "1000", "500", "500", "11 months"},
			wantErr: true,
			errMsg:  "term is not between 12 and 72 months",
		},
		{
			name:    "term above maximum",
			cols:    []string{"id", "name", "1000", "500", "500", "73 months"},
			wantErr: true,
			errMsg:  "term is not between 12 and 72 months",
		},
		{
			name:    "negative term",
			cols:    []string{"id", "name", "1000", "500", "500", "-5 months"},
			wantErr: true,
			errMsg:  "term is not between 12 and 72 months",
		},
		{
			name:    "missing months suffix",
			cols:    []string{"id", "name", "1000", "500", "500", "12"},
			wantErr: false,
		},
		{
			name:    "empty term",
			cols:    []string{"id", "name", "1000", "500", "500", ""},
			wantErr: true,
			errMsg:  "term is not a number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
			}

			result, err := hasValidTerm(ctx, tc.cols)

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

			// Verify the returned map contains the expected term value
			// Process the term value using the same logic as the function
			termStr := strings.TrimSpace(tc.cols[5])

			// Remove the "months" suffix using the same logic as the function
			if strings.HasSuffix(strings.ToLower(termStr), "months") {
				monthsIndex := strings.LastIndex(strings.ToLower(termStr), "months")
				if monthsIndex > 0 {
					termStr = termStr[:monthsIndex]
				}
			}

			termStr = strings.TrimSpace(termStr)

			term, err := strconv.Atoi(termStr)
			if err != nil {
				// If there's an error, we'll use 0 as the term value for comparison
				term = 0
			}

			if result["term"] != strconv.Itoa(term) {
				t.Errorf("expected term '%s', got '%s'", strconv.Itoa(term), result["term"])
			}
		})
	}
}
