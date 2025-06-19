package loan_info

import (
	"go-file-parsing/utils"
	"go-file-parsing/validator"
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

func validateFormattedInt(s *string, parseError error, rangeCheck func(int) error) error {
	var err error
	i, err := utils.FormattedStringToInt(s)
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

	// Get a map from the pool
	result := vCtx.GetMap()
	var err error
	dtiStr := utils.TrimIfNeeded(cols[colDTI])
	err = validateFormattedInt(&dtiStr, ErrDTINotNumber, func(i int) error {
		if i > 20 {
			return ErrDTITooHigh
		}
		result["dti"] = dtiStr
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}

	// Check home ownership
	homeOwnership := utils.TrimIfNeeded(cols[12])
	if homeOwnership != "MORTGAGE" && homeOwnership != "OWN" {
		validator.PutMap(result)
		return nil, ErrHomeOwnershipInvalid
	}
	result["homeOwnership"] = homeOwnership

	return result, nil
}

// Rule 7: Established Credit History
// earliest_cr_line not null and is > 10 years ago.
func hasEstablishedCreditHistory(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	earliestCRLine := utils.TrimIfNeeded(cols[colEarliestCrLine])

	if earliestCRLine == "" {
		return nil, ErrEarliestCrLineEmpty
	}
	result := vCtx.GetMap()
	var err error
	// Parse the date in format YYYY-MM
	workTime, err := time.Parse("2006-01", earliestCRLine)
	if err != nil {
		validator.PutMap(result)
		return nil, ErrEarliestCrLineFormat
	}

	// Check if the date is more than 10 years ago
	if workTime.After(tenYearsAgo) {
		validator.PutMap(result)
		return nil, ErrEarliestCrLineTooRecent
	}
	result["earliestCrLine"] = earliestCRLine

	return result, nil
}

// Rule 8: Healthy FICO Score
// fico_range_low >= 660 and fico_range_high <= 850.
func hasHealthyFICOScore(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	var err error
	// Get a map from the pool
	result := vCtx.GetMap()
	ficoStr := utils.TrimIfNeeded(cols[colFICORangeLow])
	err = validateFormattedInt(&ficoStr, ErrFICORangeLowNotNumber, func(i int) error {
		if i < 660 {
			return ErrFICORangeLowTooLow
		}
		result["ficoRangeLow"] = ficoStr
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}

	ficoHighStr := utils.TrimIfNeeded(cols[colFICORangeHigh])
	err = validateFormattedInt(&ficoHighStr, ErrFICORangeHighNotNumber, func(i int) error {
		if i > 850 {
			return ErrFICORangeHighTooHigh
		}
		result["ficoRangeHigh"] = ficoHighStr
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	return result, nil
}

// Rule 9: Has Sufficient Accounts
// total_acc >= 5 and open_acc >= 2.
func hasSufficientAccounts(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	result := vCtx.GetMap()
	var err error
	totalAcc := utils.TrimIfNeeded(cols[colTotalAcc])
	err = validateFormattedInt(&totalAcc, ErrTotalAccNotNumber, func(i int) error {
		if i < 5 {
			return ErrTotalAccTooFew
		}
		result["totalAcc"] = totalAcc
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}

	openAcc := utils.TrimIfNeeded(cols[colOpenAcc])
	err = validateFormattedInt(&openAcc, ErrOpenAccNotNumber, func(i int) error {
		if i < 2 {
			return ErrOpenAccTooFew
		}
		result["openAcc"] = openAcc
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	return result, nil
}

// Rule 11: No Public Record or Bankruptcies
// pub_rec == 0 and pub_rec_bankruptcies == 0 and tax_liens == 0.
func hasNoPublicRecordOrBankruptcies(vCtx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	var err error
	// Get a map from the pool
	result := vCtx.GetMap()
	pubRec := utils.TrimIfNeeded(cols[colPubRec])
	err = validateFormattedInt(&pubRec, ErrPubRecNotNumber, func(i int) error {
		if i != 0 {
			return ErrPubRecNotZero
		}
		result["pubRec"] = pubRec
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}
	pubRecBankruptcies := utils.TrimIfNeeded(cols[colPubRecBankruptcies])
	err = validateFormattedInt(&pubRecBankruptcies, ErrPubRecBankruptciesNotNumber, func(i int) error {
		if i != 0 {
			return ErrPubRecBankruptciesNotZero
		}
		result["pubRecBankruptcies"] = pubRecBankruptcies
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}

	taxLiens := utils.TrimIfNeeded(cols[colTaxLiens])
	err = validateFormattedInt(&taxLiens, ErrTaxLiensNotNumber, func(i int) error {
		if i != 0 {
			return ErrTaxLiensNotZero
		}
		result["taxLiens"] = taxLiens
		return nil
	})
	if err != nil {
		validator.PutMap(result)
		return nil, err
	}

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
	err := validateFormattedInt(&annualIncStr, ErrAnnualIncNotNumber, func(i int) error {
		if i <= 30000 {
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
