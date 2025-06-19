package loan_info

import (
	"go-file-parsing/utils"
	"go-file-parsing/validator"
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
	colOpenAcc            = 32
	colPubRec             = 33
	colTotalAcc           = 36
	colPubRecBankruptcies = 109
	colTaxLiens           = 110
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

func validateFormattedInt(s *string, parseError error, rangeCheck func(*int) error) error {
	var err error
	i := intPool.Get().(*int)
	defer intPool.Put(i)
	*i, err = utils.FormattedStringToInt(s)
	if err != nil {
		return parseError
	}
	return rangeCheck(i)

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
func hasLowDTI(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	workStr := strPool.Get().(*string)
	*workStr = utils.TrimIfNeeded(cols[colDTI])
	// Get a map from the pool
	result := vCtx.GetMap()
	defer func() {
		strPool.Put(workStr)
		// If we're returning an error, return the map to the pool
		if result != nil && recover() != nil {
			validator.PutMap(result)
			result = nil
		}
	}()
	var err error
	err = validateFormattedInt(workStr, ErrDTINotNumber, func(i *int) error {
		if *i >= 20 {
			return ErrDTITooHigh
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["dti"] = *workStr

	// Check home ownership
	homeOwnership := utils.TrimIfNeeded(cols[12])
	if homeOwnership != "MORTGAGE" && homeOwnership != "OWN" {
		validator.PutMap(result)
		return nil, ErrHomeOwnershipInvalid
	}
	result["homeOwnership"] = homeOwnership

	// Check annual income
	*workStr = utils.TrimIfNeeded(cols[colAnnualInc])
	err = validateFormattedInt(workStr, ErrAnnualIncNotNumber, func(i *int) error {
		if *i <= 40000 {
			return ErrAnnualIncTooLow40K
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
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
	result := vCtx.GetMap()
	var err error
	defer func() {
		strPool.Put(workStr)
		timePool.Put(workTime)

	}()

	// Parse the date in format YYYY-MM
	*workTime, err = time.Parse("2006-01", *workStr)
	if err != nil {
		validator.PutMap(result)
		return nil, ErrEarliestCrLineFormat
	}

	// Check if the date is more than 10 years ago
	if workTime.After(tenYearsAgo) {
		validator.PutMap(result)
		return nil, ErrEarliestCrLineTooRecent
	}
	result["earliestCrLine"] = *workStr

	return result, nil
}

// Rule 8: Healthy FICO Score
// fico_range_low >= 660 and fico_range_high <= 850.
func hasHealthyFICOScore(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	workStr := strPool.Get().(*string)
	*workStr = utils.TrimIfNeeded(cols[colFICORangeLow])
	utils.TrimTrailingDecimal(workStr)
	var err error
	// Get a map from the pool
	result := vCtx.GetMap()
	defer func() {
		strPool.Put(workStr)

	}()

	err = validateFormattedInt(workStr, ErrFICORangeLowNotNumber, func(i *int) error {
		if *i < 660 {
			return ErrFICORangeLowTooLow
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["ficoRangeLow"] = *workStr

	*workStr = utils.TrimIfNeeded(cols[colFICORangeHigh])
	utils.TrimTrailingDecimal(workStr)
	err = validateFormattedInt(workStr, ErrFICORangeHighNotNumber, func(i *int) error {
		if *i > 850 {
			return ErrFICORangeHighTooHigh
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["ficoRangeHigh"] = *workStr

	return result, nil
}

// Rule 9: Has Sufficient Accounts
// total_acc >= 5 and open_acc >= 2.
func hasSufficientAccounts(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	workStr := strPool.Get().(*string)
	result := vCtx.GetMap()
	var err error
	defer func() {
		strPool.Put(workStr)
	}()
	*workStr = utils.TrimIfNeeded(cols[colTotalAcc])
	err = validateFormattedInt(workStr, ErrTotalAccNotNumber, func(i *int) error {
		if *i < 5 {
			return ErrTotalAccTooFew
		}

		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["totalAcc"] = *workStr

	*workStr = utils.TrimIfNeeded(cols[colOpenAcc])
	err = validateFormattedInt(workStr, ErrOpenAccNotNumber, func(i *int) error {
		if *i < 2 {
			return ErrOpenAccTooFew
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["openAcc"] = *workStr

	return result, nil
}

// Rule 11: No Public Record or Bankruptcies
// pub_rec == 0 and pub_rec_bankruptcies == 0 and tax_liens == 0.
func hasNoPublicRecordOrBankruptcies(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	workStr := strPool.Get().(*string)
	var err error
	// Get a map from the pool
	result := vCtx.GetMap()
	defer func() {
		strPool.Put(workStr)
	}()
	*workStr = utils.TrimIfNeeded(cols[colPubRec])
	err = validateFormattedInt(workStr, ErrPubRecNotNumber, func(i *int) error {
		if *i != 0 {
			return ErrPubRecNotZero
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["pubRec"] = *workStr

	*workStr = utils.TrimIfNeeded(cols[colPubRecBankruptcies])
	err = validateFormattedInt(workStr, ErrPubRecBankruptciesNotNumber, func(i *int) error {
		if *i != 0 {
			return ErrPubRecBankruptciesNotZero
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["pubRecBankruptcies"] = *workStr

	*workStr = utils.TrimIfNeeded(cols[colTaxLiens])
	err = validateFormattedInt(workStr, ErrTaxLiensNotNumber, func(i *int) error {
		if *i != 0 {
			return ErrTaxLiensNotZero
		}
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["taxLiens"] = *workStr

	return result, nil
}

// Rule 12: Verified with Income
// verification_status in [Source Verified, Verified] and annual_inc > 30,000.
func isVerifiedWithIncome(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	verificationStatus := utils.TrimIfNeeded(cols[colVerificationStatus])
	annualIncStr := utils.TrimIfNeeded(cols[colAnnualInc])

	// Check verification status
	if verificationStatus != "Source Verified" && verificationStatus != "Verified" {
		return nil, ErrVerificationStatusInvalid
	}

	// Get a map from the pool
	result := vCtx.GetMap()
	err := validateFormattedInt(&annualIncStr, ErrAnnualIncNotNumber, func(i *int) error {
		if *i <= 30000 {
			return ErrAnnualIncTooLow30K
		}
		return nil
	})

	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	result["verificationStatus"] = verificationStatus
	result["annualInc"] = annualIncStr

	return result, nil
}
