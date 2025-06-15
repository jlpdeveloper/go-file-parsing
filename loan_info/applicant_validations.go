package loan_info

import (
	"go-file-parsing/utils"
	"go-file-parsing/validator"
	"strconv"
	"strings"
	"time"
)

// Column index constants for CSV fields
const (
	colEmpTitle           = 10
	colEmpLength          = 11
	colHomeOwnership      = 12
	colAnnualInc          = 13
	colVerificationStatus = 14
	colDTI                = 36
	colEarliestCrLine     = 38
	colFICORangeLow       = 39
	colFICORangeHigh      = 40
	colOpenAcc            = 44
	colPubRec             = 45
	colTotalAcc           = 48
	colPubRecBankruptcies = 121
	colTaxLiens           = 122
)

// Rule 5: Has Employment Info
// Non-empty emp_title and emp_length is not null.
func hasEmploymentInfo(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	empTitle := utils.TrimIfNeeded(cols[colEmpTitle])
	empLength := utils.TrimIfNeeded(cols[colEmpLength])

	if empTitle == "" {
		return nil, ErrEmpTitleEmpty
	}

	if empLength == "" {
		return nil, ErrEmpLengthEmpty
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["empTitle"] = empTitle
	result["empLength"] = empLength

	return result, nil
}

// Rule 6: Low DTI and Home Ownership
// dti < 20, home_ownership in [MORTGAGE, OWN], and annual_inc > 40,000.
func hasLowDTIAndHomeOwnership(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	dtiStr := utils.TrimIfNeeded(cols[colDTI])
	homeOwnership := strings.ToUpper(utils.TrimIfNeeded(cols[colHomeOwnership]))
	annualIncStr := utils.TrimIfNeeded(cols[colAnnualInc])

	dti, err := strconv.ParseFloat(dtiStr, 64)
	if err != nil {
		return nil, ErrDTINotNumber
	}

	if dti >= 20 {
		return nil, ErrDTITooHigh
	}

	if homeOwnership != "MORTGAGE" && homeOwnership != "OWN" {
		return nil, ErrHomeOwnershipInvalid
	}

	annualInc, err := strconv.ParseFloat(annualIncStr, 64)
	if err != nil {
		return nil, ErrAnnualIncNotNumber
	}

	if annualInc <= 40000 {
		return nil, ErrAnnualIncTooLow40K
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["dti"] = dtiStr
	result["homeOwnership"] = homeOwnership
	result["annualInc"] = annualIncStr

	return result, nil
}

// Rule 7: Established Credit History
// earliest_cr_line not null and is > 10 years ago.
func hasEstablishedCreditHistory(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	earliestCrLine := utils.TrimIfNeeded(cols[colEarliestCrLine])

	if earliestCrLine == "" {
		return nil, ErrEarliestCrLineEmpty
	}

	// Parse the date in format YYYY-MM
	crDate, err := time.Parse("2006-01", earliestCrLine)
	if err != nil {
		return nil, ErrEarliestCrLineFormat
	}

	// Check if the date is more than 10 years ago
	tenYearsAgo := time.Now().AddDate(-10, 0, 0)
	if crDate.After(tenYearsAgo) {
		return nil, ErrEarliestCrLineTooRecent
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["earliestCrLine"] = earliestCrLine

	return result, nil
}

// Rule 8: Healthy FICO Score
// fico_range_low >= 660 and fico_range_high <= 850.
func hasHealthyFICOScore(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	ficoRangeLowStr := utils.TrimIfNeeded(cols[colFICORangeLow])
	ficoRangeHighStr := utils.TrimIfNeeded(cols[colFICORangeHigh])

	ficoRangeLow, err := strconv.Atoi(ficoRangeLowStr)
	if err != nil {
		return nil, ErrFICORangeLowNotNumber
	}

	ficoRangeHigh, err := strconv.Atoi(ficoRangeHighStr)
	if err != nil {
		return nil, ErrFICORangeHighNotNumber
	}

	if ficoRangeLow < 660 {
		return nil, ErrFICORangeLowTooLow
	}

	if ficoRangeHigh > 850 {
		return nil, ErrFICORangeHighTooHigh
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["ficoRangeLow"] = ficoRangeLowStr
	result["ficoRangeHigh"] = ficoRangeHighStr

	return result, nil
}

// Rule 9: Has Sufficient Accounts
// total_acc >= 5 and open_acc >= 2.
func hasSufficientAccounts(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	totalAccStr := utils.TrimIfNeeded(cols[colTotalAcc])
	openAccStr := utils.TrimIfNeeded(cols[colOpenAcc])

	totalAcc, err := strconv.Atoi(totalAccStr)
	if err != nil {
		return nil, ErrTotalAccNotNumber
	}

	openAcc, err := strconv.Atoi(openAccStr)
	if err != nil {
		return nil, ErrOpenAccNotNumber
	}

	if totalAcc < 5 {
		return nil, ErrTotalAccTooFew
	}

	if openAcc < 2 {
		return nil, ErrOpenAccTooFew
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["totalAcc"] = totalAccStr
	result["openAcc"] = openAccStr

	return result, nil
}

// Rule 10: Stable Employment
// emp_length in [5 years, 6 years, 7 years, 8 years, 9 years, 10+ years].
func hasStableEmployment(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	empLength := utils.TrimIfNeeded(cols[colEmpLength])
	validEmpLengths := map[string]bool{
		"5 years":   true,
		"6 years":   true,
		"7 years":   true,
		"8 years":   true,
		"9 years":   true,
		"10+ years": true,
	}

	if !validEmpLengths[empLength] {
		return nil, ErrEmpLengthNotStable
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["empLength"] = empLength

	return result, nil
}

// Rule 11: No Public Record or Bankruptcies
// pub_rec == 0 and pub_rec_bankruptcies == 0 and tax_liens == 0.
func hasNoPublicRecordOrBankruptcies(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	pubRecStr := utils.TrimIfNeeded(cols[colPubRec])
	pubRecBankruptciesStr := utils.TrimIfNeeded(cols[colPubRecBankruptcies])
	taxLiensStr := utils.TrimIfNeeded(cols[colTaxLiens])

	pubRec, err := strconv.Atoi(pubRecStr)
	if err != nil {
		return nil, ErrPubRecNotNumber
	}

	pubRecBankruptcies, err := strconv.Atoi(pubRecBankruptciesStr)
	if err != nil {
		return nil, ErrPubRecBankruptciesNotNumber
	}

	taxLiens, err := strconv.Atoi(taxLiensStr)
	if err != nil {
		return nil, ErrTaxLiensNotNumber
	}

	if pubRec != 0 {
		return nil, ErrPubRecNotZero
	}

	if pubRecBankruptcies != 0 {
		return nil, ErrPubRecBankruptciesNotZero
	}

	if taxLiens != 0 {
		return nil, ErrTaxLiensNotZero
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["pubRec"] = pubRecStr
	result["pubRecBankruptcies"] = pubRecBankruptciesStr
	result["taxLiens"] = taxLiensStr

	return result, nil
}

// Rule 12: Verified with Income
// verification_status in [Source Verified, Verified] and annual_inc > 30,000.
func isVerifiedWithIncome(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	verificationStatus := utils.TrimIfNeeded(cols[colVerificationStatus])
	annualIncStr := utils.TrimIfNeeded(cols[colAnnualInc])

	validVerificationStatuses := map[string]bool{
		"Source Verified": true,
		"Verified":        true,
	}

	if !validVerificationStatuses[verificationStatus] {
		return nil, ErrVerificationStatusInvalid
	}

	annualInc, err := strconv.ParseFloat(annualIncStr, 64)
	if err != nil {
		return nil, ErrAnnualIncNotNumber
	}

	if annualInc <= 30000 {
		return nil, ErrAnnualIncTooLow30K
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	result["verificationStatus"] = verificationStatus
	result["annualInc"] = annualIncStr

	return result, nil
}
