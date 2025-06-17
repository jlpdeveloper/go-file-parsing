package loan_info

import (
	"go-file-parsing/utils"
	"go-file-parsing/validator"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Column index constants for CSV fields
const (
	colEmpTitle           = 10
	colEmpLength          = 11
	colAnnualInc          = 13
	colVerificationStatus = 14
	colDTI                = 24
	colEarliestCrLine     = 26
	colFICORangeLow       = 27
	colFICORangeHigh      = 28
	colOpenAcc            = 44
	colPubRec             = 45
	colTotalAcc           = 48
	colPubRecBankruptcies = 121
	colTaxLiens           = 122
)

var tenYearsAgo = time.Now().AddDate(-10, 0, 0)

var intPool = &sync.Pool{
	New: func() interface{} {
		return new(int)
	},
}

var strPool = &sync.Pool{
	New: func() interface{} {
		return new(string)
	},
}

var timePool = &sync.Pool{
	New: func() interface{} {
		return new(time.Time)
	},
}

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
	defer func() {
		if recover() != nil {
			validator.PutMap(result)
		}
	}()
	result["empTitle"] = empTitle
	result["empLength"] = empLength

	return result, nil
}

// Rule 6: Low DTI and Home Ownership
// dti < 20, home_ownership in [MORTGAGE, OWN], and annual_inc > 40,000.
func hasLowDTIAndHomeOwnership(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	workStr := strPool.Get().(*string)
	*workStr = utils.TrimIfNeeded(cols[colDTI])
	workInt := intPool.Get().(*int)
	var result map[string]string
	defer func() {
		intPool.Put(workInt)
		strPool.Put(workStr)
		// If we're returning an error, return the map to the pool
		if result != nil && recover() != nil {
			validator.PutMap(result)
			result = nil
		}
	}()
	var err error
	*workInt, err = strconv.Atoi(strings.Split(*workStr, ".")[0])
	if err != nil {
		return nil, ErrDTINotNumber
	}

	if *workInt >= 30 {
		return nil, ErrDTITooHigh
	}
	// Get a map from the pool
	result = vCtx.GetMap()
	result["dti"] = *workStr
	*workStr = utils.TrimIfNeeded(cols[colAnnualInc])
	utils.TrimTrailingDecimal(workStr)
	*workInt, err = strconv.Atoi(*workStr)
	if err != nil {
		validator.PutMap(result)
		return nil, ErrAnnualIncNotNumber
	}

	if *workInt <= 30000 {
		validator.PutMap(result)
		return nil, ErrAnnualIncTooLow40K
	}

	result["annualInc"] = *workStr

	return result, nil
}

// Rule 7: Established Credit History
// earliest_cr_line not null and is > 10 years ago.
func hasEstablishedCreditHistory(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	workStr := strPool.Get().(*string)
	*workStr = utils.TrimIfNeeded(cols[colEarliestCrLine])

	if *workStr == "" {
		strPool.Put(workStr)
		return nil, ErrEarliestCrLineEmpty
	}
	workTime := timePool.Get().(*time.Time)
	var result map[string]string
	var err error
	defer func() {
		strPool.Put(workStr)
		timePool.Put(workTime)

	}()

	// Parse the date in format YYYY-MM
	*workTime, err = time.Parse("Jan-2006", *workStr)
	if err != nil {
		validator.PutMap(result)
		return nil, ErrEarliestCrLineFormat
	}

	// Check if the date is more than 10 years ago
	if workTime.After(tenYearsAgo) {
		validator.PutMap(result)
		return nil, ErrEarliestCrLineTooRecent
	}

	// Create a copy of the map to return
	result = vCtx.GetMap()
	result["earliestCrLine"] = *workStr

	return result, nil
}

// Rule 8: Healthy FICO Score
// fico_range_low >= 660 and fico_range_high <= 850.
func hasHealthyFICOScore(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	workStr := strPool.Get().(*string)
	*workStr = utils.TrimIfNeeded(cols[colFICORangeLow])
	utils.TrimTrailingDecimal(workStr)
	workInt := intPool.Get().(*int)
	var err error
	// Get a map from the pool
	result := vCtx.GetMap()
	defer func() {
		strPool.Put(workStr)
		intPool.Put(workInt)

	}()

	*workInt, err = strconv.Atoi(*workStr)
	if err != nil {
		validator.PutMap(result)
		return nil, ErrFICORangeLowNotNumber
	}

	if *workInt < 660 {
		validator.PutMap(result)
		return nil, ErrFICORangeLowTooLow
	}
	result["ficoRangeLow"] = *workStr

	*workStr = utils.TrimIfNeeded(cols[colFICORangeHigh])
	utils.TrimTrailingDecimal(workStr)
	*workInt, err = strconv.Atoi(*workStr)
	if err != nil {
		validator.PutMap(result)
		return nil, ErrFICORangeHighNotNumber
	}
	if *workInt > 850 {
		validator.PutMap(result)
		return nil, ErrFICORangeHighTooHigh
	}
	result["ficoRangeHigh"] = *workStr

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
	defer func() {
		if recover() != nil {
			validator.PutMap(result)
		}
	}()
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
	defer func() {
		if recover() != nil {
			validator.PutMap(result)
		}
	}()
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
	defer func() {
		if recover() != nil {
			validator.PutMap(result)
		}
	}()
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
	defer func() {
		if recover() != nil {
			validator.PutMap(result)
		}
	}()
	result["verificationStatus"] = verificationStatus
	result["annualInc"] = annualIncStr

	return result, nil
}
