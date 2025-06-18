package loan_info

import (
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"testing"
	"time"
)

func mockGetMap() map[string]string {
	return make(map[string]string)
}

func TestHasEmploymentInfo(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid employment info",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "Engineer", "5 years"},
			wantErr: false,
		},
		{
			name:    "empty employment title",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "5 years"},
			wantErr: true,
			errMsg:  "employment title is empty",
		},
		{
			name:    "empty employment length",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "Engineer", ""},
			wantErr: true,
			errMsg:  "employment length is empty",
		},
		{
			name:    "both empty",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", ""},
			wantErr: true,
			errMsg:  "employment title is empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
				GetMap: mockGetMap,
			}

			result, err := hasEmploymentInfo(ctx, tc.cols)

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
			if result["empTitle"] != tc.cols[10] {
				t.Errorf("expected empTitle '%s', got '%s'", tc.cols[10], result["empTitle"])
			}
			if result["empLength"] != tc.cols[11] {
				t.Errorf("expected empLength '%s', got '%s'", tc.cols[11], result["empLength"])
			}
		})
	}
}

func TestHasLowDTIAndHomeOwnership(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid DTI and home ownership",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "MORTGAGE", "50000", "", "", "", "", "", "", "", "", "", "", "15"},
			wantErr: false,
		},
		{
			name:    "DTI not a number",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "MORTGAGE", "50000", "", "", "", "", "", "", "", "", "", "", "abc"},
			wantErr: true,
			errMsg:  "DTI is not a number",
		},
		{
			name:    "DTI too high",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "MORTGAGE", "50000", "", "", "", "", "", "", "", "", "", "", "25"},
			wantErr: true,
			errMsg:  "DTI is not less than 20",
		},
		{
			name:    "invalid home ownership",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "RENT", "50000", "", "", "", "", "", "", "", "", "", "", "15"},
			wantErr: true,
			errMsg:  "home ownership is not MORTGAGE or OWN",
		},
		{
			name:    "annual income not a number",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "MORTGAGE", "abc", "", "", "", "", "", "", "", "", "", "", "15"},
			wantErr: true,
			errMsg:  "annual income is not a number",
		},
		{
			name:    "annual income too low",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "MORTGAGE", "30000", "", "", "", "", "", "", "", "", "", "", "15"},
			wantErr: true,
			errMsg:  "annual income is not greater than 40,000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
				GetMap: mockGetMap,
			}

			result, err := hasLowDTI(ctx, tc.cols)

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
			if result["dti"] != tc.cols[24] {
				t.Errorf("expected dti '%s', got '%s'", tc.cols[24], result["dti"])
			}
			if result["homeOwnership"] != tc.cols[12] {
				t.Errorf("expected homeOwnership '%s', got '%s'", tc.cols[12], result["homeOwnership"])
			}
			if result["annualInc"] != tc.cols[13] {
				t.Errorf("expected annualInc '%s', got '%s'", tc.cols[13], result["annualInc"])
			}
		})
	}
}

func TestHasEstablishedCreditHistory(t *testing.T) {
	// Get a date more than 10 years ago for valid test case
	moreThan10YearsAgo := time.Now().AddDate(-11, 0, 0).Format("2006-01")
	// Get a date less than 10 years ago for invalid test case
	lessThan10YearsAgo := time.Now().AddDate(-5, 0, 0).Format("2006-01")

	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid credit history",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", moreThan10YearsAgo},
			wantErr: false,
		},
		{
			name:    "empty credit history",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
			wantErr: true,
			errMsg:  "earliest credit line is empty",
		},
		{
			name:    "invalid date format",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "invalid-date"},
			wantErr: true,
			errMsg:  "earliest credit line is not in valid format (YYYY-MM)",
		},
		{
			name:    "credit history too recent",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", lessThan10YearsAgo},
			wantErr: true,
			errMsg:  "earliest credit line is not more than 10 years ago",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
				GetMap: mockGetMap,
			}

			result, err := hasEstablishedCreditHistory(ctx, tc.cols)

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
			if result["earliestCrLine"] != tc.cols[26] {
				t.Errorf("expected earliestCrLine '%s', got '%s'", tc.cols[26], result["earliestCrLine"])
			}
		})
	}
}

func TestHasHealthyFICOScore(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid FICO score",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "700", "800"},
			wantErr: false,
		},
		{
			name:    "FICO range low not a number",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "abc", "800"},
			wantErr: true,
			errMsg:  "FICO range low is not a number",
		},
		{
			name:    "FICO range high not a number",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "700", "abc"},
			wantErr: true,
			errMsg:  "FICO range high is not a number",
		},
		{
			name:    "FICO range low too low",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "650", "800"},
			wantErr: true,
			errMsg:  "FICO range low is less than 660",
		},
		{
			name:    "FICO range high too high",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "700", "860"},
			wantErr: true,
			errMsg:  "FICO range high is greater than 850",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
				GetMap: mockGetMap,
			}

			result, err := hasHealthyFICOScore(ctx, tc.cols)

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
			if result["ficoRangeLow"] != tc.cols[27] {
				t.Errorf("expected ficoRangeLow '%s', got '%s'", tc.cols[27], result["ficoRangeLow"])
			}
			if result["ficoRangeHigh"] != tc.cols[28] {
				t.Errorf("expected ficoRangeHigh '%s', got '%s'", tc.cols[28], result["ficoRangeHigh"])
			}
		})
	}
}

func TestHasSufficientAccounts(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid account numbers",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "5", "", "", "", "10"},
			wantErr: false,
		},
		{
			name:    "total accounts not a number",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "5", "", "", "", "abc"},
			wantErr: true,
			errMsg:  "total accounts is not a number",
		},
		{
			name:    "open accounts not a number",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "abc", "", "", "", "10"},
			wantErr: true,
			errMsg:  "open accounts is not a number",
		},
		{
			name:    "total accounts too low",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "5", "", "", "", "4"},
			wantErr: true,
			errMsg:  "total accounts is less than 5",
		},
		{
			name:    "open accounts too low",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "1", "", "", "", "10"},
			wantErr: true,
			errMsg:  "open accounts is less than 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
				GetMap: mockGetMap,
			}

			result, err := hasSufficientAccounts(ctx, tc.cols)

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
			if result["totalAcc"] != tc.cols[36] {
				t.Errorf("expected totalAcc '%s', got '%s'", tc.cols[36], result["totalAcc"])
			}
			if result["openAcc"] != tc.cols[32] {
				t.Errorf("expected openAcc '%s', got '%s'", tc.cols[32], result["openAcc"])
			}
		})
	}
}

func TestHasNoPublicRecordOrBankruptcies(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid - no public records or bankruptcies",
			cols:    createTestCols(33, "0", 109, "0", 110, "0"),
			wantErr: false,
		},
		{
			name:    "invalid - public records not a number",
			cols:    createTestCols(33, "abc", 109, "0", 110, "0"),
			wantErr: true,
			errMsg:  "public records is not a number",
		},
		{
			name:    "invalid - public record bankruptcies not a number",
			cols:    createTestCols(33, "0", 109, "abc", 110, "0"),
			wantErr: true,
			errMsg:  "public record bankruptcies is not a number",
		},
		{
			name:    "invalid - tax liens not a number",
			cols:    createTestCols(33, "0", 109, "0", 110, "abc"),
			wantErr: true,
			errMsg:  "tax liens is not a number",
		},
		{
			name:    "invalid - public records not zero",
			cols:    createTestCols(33, "1", 109, "0", 110, "0"),
			wantErr: true,
			errMsg:  "public records is not zero",
		},
		{
			name:    "invalid - public record bankruptcies not zero",
			cols:    createTestCols(33, "0", 109, "1", 110, "0"),
			wantErr: true,
			errMsg:  "public record bankruptcies is not zero",
		},
		{
			name:    "invalid - tax liens not zero",
			cols:    createTestCols(33, "0", 109, "0", 110, "1"),
			wantErr: true,
			errMsg:  "tax liens is not zero",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
				GetMap: mockGetMap,
			}

			result, err := hasNoPublicRecordOrBankruptcies(ctx, tc.cols)

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
			if result["pubRec"] != tc.cols[33] {
				t.Errorf("expected pubRec '%s', got '%s'", tc.cols[33], result["pubRec"])
			}
			if result["pubRecBankruptcies"] != tc.cols[109] {
				t.Errorf("expected pubRecBankruptcies '%s', got '%s'", tc.cols[109], result["pubRecBankruptcies"])
			}
			if result["taxLiens"] != tc.cols[110] {
				t.Errorf("expected taxLiens '%s', got '%s'", tc.cols[110], result["taxLiens"])
			}
		})
	}
}

func TestIsVerifiedWithIncome(t *testing.T) {
	testCases := []struct {
		name    string
		cols    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid - Source Verified with sufficient income",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "40000", "Source Verified"},
			wantErr: false,
		},
		{
			name:    "valid - Verified with sufficient income",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "40000", "Verified"},
			wantErr: false,
		},
		{
			name:    "invalid - Not Verified",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "40000", "Not Verified"},
			wantErr: true,
			errMsg:  "verification status is not Source Verified or Verified",
		},
		{
			name:    "invalid - empty verification status",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "40000", ""},
			wantErr: true,
			errMsg:  "verification status is not Source Verified or Verified",
		},
		{
			name:    "invalid - annual income not a number",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "abc", "Verified"},
			wantErr: true,
			errMsg:  "annual income is not a number",
		},
		{
			name:    "invalid - annual income too low",
			cols:    []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "25000", "Verified"},
			wantErr: true,
			errMsg:  "annual income is not greater than 30,000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &validator.RowValidatorContext{
				Config: &config.ParserConfig{},
				GetMap: mockGetMap,
			}

			result, err := isVerifiedWithIncome(ctx, tc.cols)

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
			if result["verificationStatus"] != tc.cols[14] {
				t.Errorf("expected verificationStatus '%s', got '%s'", tc.cols[14], result["verificationStatus"])
			}
			if result["annualInc"] != tc.cols[13] {
				t.Errorf("expected annualInc '%s', got '%s'", tc.cols[13], result["annualInc"])
			}
		})
	}
}

// Helper function to create test columns with specific values at specific indices
func createTestCols(indices ...interface{}) []string {
	// Find the maximum index to determine the slice size
	maxIndex := 0
	for i := 0; i < len(indices); i += 2 {
		index := indices[i].(int)
		if index > maxIndex {
			maxIndex = index
		}
	}

	// Create a slice with empty strings
	cols := make([]string, maxIndex+1)
	for i := range cols {
		cols[i] = ""
	}

	// Set the specified values at the specified indices
	for i := 0; i < len(indices); i += 2 {
		index := indices[i].(int)
		value := indices[i+1].(string)
		cols[index] = value
	}

	return cols
}
