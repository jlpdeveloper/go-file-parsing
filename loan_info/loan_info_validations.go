package loan_info

import (
	"go-file-parsing/utils"
	"go-file-parsing/validator"
	"regexp"
	"strconv"
	"strings"
)

// Precompiled regex patterns for grade and subgrade validation
var (
	gradeRegex      = regexp.MustCompile(`^[A-G]$`)
	subgradeRegexes = map[string]*regexp.Regexp{
		"A": regexp.MustCompile(`^A[1-5]$`),
		"B": regexp.MustCompile(`^B[1-5]$`),
		"C": regexp.MustCompile(`^C[1-5]$`),
		"D": regexp.MustCompile(`^D[1-5]$`),
		"E": regexp.MustCompile(`^E[1-5]$`),
		"F": regexp.MustCompile(`^F[1-5]$`),
		"G": regexp.MustCompile(`^G[1-5]$`),
	}
)

const (
	colLoanAmount    = 2
	colFundingAmount = 3
	colFundingInvAmt = 4
	colInterestRate  = 6
	colTerm          = 5
	colGrade         = 8
	colSubgrade      = 9
)

func convertAmtToInt(amt string) (int, error) {
	tmpStr := utils.TrimIfNeeded(amt)
	if strings.HasSuffix(tmpStr, ".0") {
		tmpStr = tmpStr[:len(tmpStr)-2]
	}
	return strconv.Atoi(tmpStr)
}

func hasValidLoanAmount(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {

	loanAmount, err := convertAmtToInt(cols[colLoanAmount])
	if err != nil {
		return nil, ErrLoanAmountNotNumber
	}
	if loanAmount <= 0 {
		return nil, ErrLoanAmountNotPositive
	}
	fundingAmount, err := convertAmtToInt(cols[colFundingAmount])
	if err != nil {
		return nil, ErrFundingAmountNotNumber
	}
	if fundingAmount <= 0 {
		return nil, ErrFundingAmountNotPositive
	}
	fundingInvAmt, err := convertAmtToInt(cols[colFundingInvAmt])
	if err != nil {
		return nil, ErrFundingInvAmtNotNumber
	}
	if fundingInvAmt <= 0 {
		return nil, ErrFundingInvAmtNotPositive
	}
	if fundingInvAmt != fundingAmount {
		return nil, ErrFundingInvAmtNotEqual
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["loanAmount"] = cols[colLoanAmount]
	result["fundingAmount"] = cols[colFundingAmount]
	result["fundingInvAmt"] = cols[colFundingInvAmt]

	return result, nil
}

func hasValidInterestRate(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	rate, err := strconv.ParseFloat(cols[6], 64)
	if err != nil {
		return nil, ErrInterestRateNotNumber
	}
	if rate < 5 || rate > 35 {
		return nil, ErrInterestRateOutOfRange
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["interestRate"] = cols[6]

	return result, nil
}

func hasValidTerm(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	// First trim spaces from the original string
	termStr := utils.TrimIfNeeded(cols[5])

	// Remove the " months" suffix, handling case where there might be spaces
	// Use strings.HasSuffix to check if the string ends with " months"
	if strings.HasSuffix(strings.ToLower(termStr), "months") {
		// Find the last occurrence of "months" and take everything before it
		monthsIndex := strings.LastIndex(strings.ToLower(termStr), "months")
		if monthsIndex > 0 {
			termStr = termStr[:monthsIndex]
		}
	}

	// Trim spaces again
	termStr = utils.TrimIfNeeded(termStr)
	term, err := strconv.Atoi(termStr)

	if err != nil {
		return nil, ErrTermNotNumber
	}
	if term < 12 || term > 72 {
		return nil, ErrTermOutOfRange
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["term"] = strconv.Itoa(term)

	return result, nil
}

func hasValidGradeSubgrade(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	// Grade is in column 8 (index 8, 0-based)
	originalGrade := utils.TrimIfNeeded(cols[8])
	grade := strings.ToUpper(originalGrade)

	// Subgrade is in column 9 (index 9, 0-based)
	originalSubgrade := utils.TrimIfNeeded(cols[9])
	subgrade := strings.ToUpper(originalSubgrade)

	// Check if grade is a single letter from A to G using precompiled regex
	if !gradeRegex.MatchString(grade) {
		return nil, ErrGradeInvalid
	}

	// Check if subgrade matches the pattern of grade letter followed by a number from 1 to 5
	// using the precompiled regex for the specific grade
	if regex, exists := subgradeRegexes[grade]; exists {
		if !regex.MatchString(subgrade) {
			return nil, ErrSubgradeInvalid
		}
	} else {
		return nil, ErrGradeForSubgradeInvalid
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["grade"] = originalGrade
	result["subgrade"] = originalSubgrade

	return result, nil
}
