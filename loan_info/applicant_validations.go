package loan_info

import (
	"errors"
	"go-file-parsing/utils"
	"go-file-parsing/validator"
	"strconv"
	"strings"
	"time"
)

// Rule 5: Has Employment Info
// Non-empty emp_title and emp_length is not null.
func hasEmploymentInfo(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	empTitle := utils.TrimIfNeeded(cols[10])
	empLength := utils.TrimIfNeeded(cols[11])

	if empTitle == "" {
		return nil, errors.New("employment title is empty")
	}

	if empLength == "" {
		return nil, errors.New("employment length is empty")
	}

	return map[string]string{
		"empTitle":  empTitle,
		"empLength": empLength,
	}, nil
}

// Rule 6: Low DTI and Home Ownership
// dti < 20, home_ownership in [MORTGAGE, OWN], and annual_inc > 40,000.
func hasLowDTIAndHomeOwnership(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	dtiStr := utils.TrimIfNeeded(cols[36])
	homeOwnership := strings.ToUpper(utils.TrimIfNeeded(cols[12]))
	annualIncStr := utils.TrimIfNeeded(cols[13])

	dti, err := strconv.ParseFloat(dtiStr, 64)
	if err != nil {
		return nil, errors.New("DTI is not a number")
	}

	if dti >= 20 {
		return nil, errors.New("DTI is not less than 20")
	}

	if homeOwnership != "MORTGAGE" && homeOwnership != "OWN" {
		return nil, errors.New("home ownership is not MORTGAGE or OWN")
	}

	annualInc, err := strconv.ParseFloat(annualIncStr, 64)
	if err != nil {
		return nil, errors.New("annual income is not a number")
	}

	if annualInc <= 40000 {
		return nil, errors.New("annual income is not greater than 40,000")
	}

	return map[string]string{
		"dti":           dtiStr,
		"homeOwnership": homeOwnership,
		"annualInc":     annualIncStr,
	}, nil
}

// Rule 7: Established Credit History
// earliest_cr_line not null and is > 10 years ago.
func hasEstablishedCreditHistory(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	earliestCrLine := utils.TrimIfNeeded(cols[38])

	if earliestCrLine == "" {
		return nil, errors.New("earliest credit line is empty")
	}

	// Parse the date in format YYYY-MM
	crDate, err := time.Parse("2006-01", earliestCrLine)
	if err != nil {
		return nil, errors.New("earliest credit line is not in valid format (YYYY-MM)")
	}

	// Check if the date is more than 10 years ago
	tenYearsAgo := time.Now().AddDate(-10, 0, 0)
	if crDate.After(tenYearsAgo) {
		return nil, errors.New("earliest credit line is not more than 10 years ago")
	}

	return map[string]string{
		"earliestCrLine": earliestCrLine,
	}, nil
}

// Rule 8: Healthy FICO Score
// fico_range_low >= 660 and fico_range_high <= 850.
func hasHealthyFICOScore(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	ficoRangeLowStr := utils.TrimIfNeeded(cols[39])
	ficoRangeHighStr := utils.TrimIfNeeded(cols[40])

	ficoRangeLow, err := strconv.Atoi(ficoRangeLowStr)
	if err != nil {
		return nil, errors.New("FICO range low is not a number")
	}

	ficoRangeHigh, err := strconv.Atoi(ficoRangeHighStr)
	if err != nil {
		return nil, errors.New("FICO range high is not a number")
	}

	if ficoRangeLow < 660 {
		return nil, errors.New("FICO range low is less than 660")
	}

	if ficoRangeHigh > 850 {
		return nil, errors.New("FICO range high is greater than 850")
	}

	return map[string]string{
		"ficoRangeLow":  ficoRangeLowStr,
		"ficoRangeHigh": ficoRangeHighStr,
	}, nil
}

// Rule 9: Has Sufficient Accounts
// total_acc >= 5 and open_acc >= 2.
func hasSufficientAccounts(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	totalAccStr := utils.TrimIfNeeded(cols[48])
	openAccStr := utils.TrimIfNeeded(cols[44])

	totalAcc, err := strconv.Atoi(totalAccStr)
	if err != nil {
		return nil, errors.New("total accounts is not a number")
	}

	openAcc, err := strconv.Atoi(openAccStr)
	if err != nil {
		return nil, errors.New("open accounts is not a number")
	}

	if totalAcc < 5 {
		return nil, errors.New("total accounts is less than 5")
	}

	if openAcc < 2 {
		return nil, errors.New("open accounts is less than 2")
	}

	return map[string]string{
		"totalAcc": totalAccStr,
		"openAcc":  openAccStr,
	}, nil
}

// Rule 10: Stable Employment
// emp_length in [5 years, 6 years, 7 years, 8 years, 9 years, 10+ years].
func hasStableEmployment(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	empLength := utils.TrimIfNeeded(cols[11])
	validEmpLengths := map[string]bool{
		"5 years":   true,
		"6 years":   true,
		"7 years":   true,
		"8 years":   true,
		"9 years":   true,
		"10+ years": true,
	}

	if !validEmpLengths[empLength] {
		return nil, errors.New("employment length is not stable (5-10+ years)")
	}

	return map[string]string{
		"empLength": empLength,
	}, nil
}

// Rule 11: No Public Record or Bankruptcies
// pub_rec == 0 and pub_rec_bankruptcies == 0 and tax_liens == 0.
func hasNoPublicRecordOrBankruptcies(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	pubRecStr := utils.TrimIfNeeded(cols[45])
	pubRecBankruptciesStr := utils.TrimIfNeeded(cols[121])
	taxLiensStr := utils.TrimIfNeeded(cols[122])

	pubRec, err := strconv.Atoi(pubRecStr)
	if err != nil {
		return nil, errors.New("public records is not a number")
	}

	pubRecBankruptcies, err := strconv.Atoi(pubRecBankruptciesStr)
	if err != nil {
		return nil, errors.New("public record bankruptcies is not a number")
	}

	taxLiens, err := strconv.Atoi(taxLiensStr)
	if err != nil {
		return nil, errors.New("tax liens is not a number")
	}

	if pubRec != 0 {
		return nil, errors.New("public records is not zero")
	}

	if pubRecBankruptcies != 0 {
		return nil, errors.New("public record bankruptcies is not zero")
	}

	if taxLiens != 0 {
		return nil, errors.New("tax liens is not zero")
	}

	return map[string]string{
		"pubRec":             pubRecStr,
		"pubRecBankruptcies": pubRecBankruptciesStr,
		"taxLiens":           taxLiensStr,
	}, nil
}

// Rule 12: Verified with Income
// verification_status in [Source Verified, Verified] and annual_inc > 30,000.
func isVerifiedWithIncome(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	verificationStatus := utils.TrimIfNeeded(cols[14])
	annualIncStr := utils.TrimIfNeeded(cols[13])

	validVerificationStatuses := map[string]bool{
		"Source Verified": true,
		"Verified":        true,
	}

	if !validVerificationStatuses[verificationStatus] {
		return nil, errors.New("verification status is not Source Verified or Verified")
	}

	annualInc, err := strconv.ParseFloat(annualIncStr, 64)
	if err != nil {
		return nil, errors.New("annual income is not a number")
	}

	if annualInc <= 30000 {
		return nil, errors.New("annual income is not greater than 30,000")
	}

	return map[string]string{
		"verificationStatus": verificationStatus,
		"annualInc":          annualIncStr,
	}, nil
}
