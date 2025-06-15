package loan_info

import (
	"errors"
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

func hasValidLoanAmount(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	loanAmount, err := strconv.Atoi(cols[2])
	if err != nil {
		return nil, errors.New("loan amount is not a number")
	}
	if loanAmount <= 0 {
		return nil, errors.New("loan amount is not a positive number")
	}
	fundingAmount, err := strconv.Atoi(cols[3])
	if err != nil {
		return nil, errors.New("funding amount is not a number")
	}
	if fundingAmount <= 0 {
		return nil, errors.New("funding amount is not a positive number")
	}
	fundingInvAmt, err := strconv.Atoi(cols[4])
	if err != nil {
		return nil, errors.New("funding inv amt is not a number")
	}
	if fundingInvAmt <= 0 {
		return nil, errors.New("funding inv amt is not a positive number")
	}
	if fundingInvAmt != fundingAmount {
		return nil, errors.New("funding inv amt is not equal to funding amount")
	}
	return map[string]string{
		"loanAmount":    cols[2],
		"fundingAmount": cols[3],
		"fundingInvAmt": cols[4],
	}, nil
}

func hasValidInterestRate(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	rate, err := strconv.ParseFloat(cols[6], 64)
	if err != nil {
		return nil, errors.New("interest rate is not a number")
	}
	if rate < 5 || rate > 35 {
		return nil, errors.New("interest rate is not between 5% and 35%")
	}
	return map[string]string{
		"interestRate": cols[6],
	}, nil
}

func hasValidTerm(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
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
		return nil, errors.New("term is not a number")
	}
	if term < 12 || term > 72 {
		return nil, errors.New("term is not between 12 and 72 months")
	}
	return map[string]string{
		"term": strconv.Itoa(term),
	}, nil
}

func hasValidGradeSubgrade(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	// Grade is in column 8 (index 8, 0-based)
	grade := utils.TrimIfNeeded(cols[8])

	// Subgrade is in column 9 (index 9, 0-based)
	subgrade := utils.TrimIfNeeded(cols[9])

	// Check if grade is a single letter from A to G using precompiled regex
	if !gradeRegex.MatchString(grade) {
		return nil, errors.New("grade must be a single letter from A to G")
	}

	// Check if subgrade matches the pattern of grade letter followed by a number from 1 to 5
	// using the precompiled regex for the specific grade
	if regex, exists := subgradeRegexes[grade]; exists {
		if !regex.MatchString(subgrade) {
			return nil, errors.New("subgrade must be the grade letter followed by a number from 1 to 5")
		}
	} else {
		return nil, errors.New("invalid grade for subgrade validation")
	}

	return map[string]string{
		"grade":    grade,
		"subgrade": subgrade,
	}, nil
}
