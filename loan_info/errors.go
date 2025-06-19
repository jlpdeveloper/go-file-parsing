package loan_info

import "errors"

// Loan amount validation errors
var (
	ErrLoanAmountNotNumber      = errors.New("loan amount is not a number")
	ErrLoanAmountNotPositive    = errors.New("loan amount is not a positive number")
	ErrFundingAmountNotNumber   = errors.New("funding amount is not a number")
	ErrFundingAmountNotPositive = errors.New("funding amount is not a positive number")
	ErrFundingInvAmtNotNumber   = errors.New("funding inv amt is not a number")
	ErrFundingInvAmtNotPositive = errors.New("funding inv amt is not a positive number")
	ErrFundingInvAmtNotEqual    = errors.New("funding inv amt is not equal to funding amount")
)

// Interest rate validation errors
var (
	ErrInterestRateNotNumber  = errors.New("interest rate is not a number")
	ErrInterestRateOutOfRange = errors.New("interest rate is not between 5% and 35%")
)

// Term validation errors
var (
	ErrTermNotNumber  = errors.New("term is not a number")
	ErrTermOutOfRange = errors.New("term is not between 12 and 72 months")
)

// Grade and subgrade validation errors
var (
	ErrGradeInvalid            = errors.New("grade must be a single letter from A to G")
	ErrSubgradeInvalid         = errors.New("subgrade must be the grade letter followed by a number from 1 to 5")
	ErrGradeForSubgradeInvalid = errors.New("invalid grade for subgrade validation")
)

// Employment info validation errors
var (
	ErrEmpTitleEmpty  = errors.New("employment title is empty")
	ErrEmpLengthEmpty = errors.New("employment length is empty")
)

// DTI and home ownership validation errors
var (
	ErrDTINotNumber         = errors.New("DTI is not a number")
	ErrDTITooHigh           = errors.New("DTI is not less than 20")
	ErrHomeOwnershipInvalid = errors.New("home ownership is not MORTGAGE or OWN")
)

// Income validation errors
var (
	ErrAnnualIncNotNumber = errors.New("annual income is not a number")
	ErrAnnualIncTooLow30K = errors.New("annual income is not greater than 30,000")
)

// Credit history validation errors
var (
	ErrEarliestCrLineEmpty     = errors.New("earliest credit line is empty")
	ErrEarliestCrLineFormat    = errors.New("earliest credit line is not in valid format (YYYY-MM)")
	ErrEarliestCrLineTooRecent = errors.New("earliest credit line is not more than 10 years ago")
)

// FICO score validation errors
var (
	ErrFICORangeLowNotNumber  = errors.New("FICO range low is not a number")
	ErrFICORangeHighNotNumber = errors.New("FICO range high is not a number")
	ErrFICORangeLowTooLow     = errors.New("FICO range low is less than 660")
	ErrFICORangeHighTooHigh   = errors.New("FICO range high is greater than 850")
)

// Account validation errors
var (
	ErrTotalAccNotNumber = errors.New("total accounts is not a number")
	ErrOpenAccNotNumber  = errors.New("open accounts is not a number")
	ErrTotalAccTooFew    = errors.New("total accounts is less than 5")
	ErrOpenAccTooFew     = errors.New("open accounts is less than 2")
)

// Public record validation errors
var (
	ErrPubRecNotNumber             = errors.New("public records is not a number")
	ErrPubRecBankruptciesNotNumber = errors.New("public record bankruptcies is not a number")
	ErrTaxLiensNotNumber           = errors.New("tax liens is not a number")
	ErrPubRecNotZero               = errors.New("public records is not zero")
	ErrPubRecBankruptciesNotZero   = errors.New("public record bankruptcies is not zero")
	ErrTaxLiensNotZero             = errors.New("tax liens is not zero")
)

// Verification status validation errors
var (
	ErrVerificationStatusInvalid = errors.New("verification status is not Source Verified or Verified")
)
