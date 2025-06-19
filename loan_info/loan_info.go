package loan_info

import (
	"go-file-parsing/config"
	"go-file-parsing/validator"
)

var validators = []validator.ColValidator{
	isValidSize,
	hasValidLoanAmount,
	hasValidInterestRate,
	hasValidTerm,
	hasEmploymentInfo,
	hasEstablishedCreditHistory,
	hasHealthyFICOScore,
	hasSufficientAccounts,
	isVerifiedWithIncome,
	hasValidGradeSubgrade,
	hasLowDTI,
	passExtraData,
}

func NewRowValidatorPool(conf *config.ParserConfig, cacheChan chan validator.CacheData, poolSize int) chan validator.CsvRowValidator {
	pool := make(chan validator.CsvRowValidator, poolSize)
	for i := 0; i < poolSize; i++ {
		pool <- validator.New(conf, cacheChan, validators)
	}
	return pool
}

// CloseValidatorPool closes all validators in the pool to prevent resource leaks.
// It should be called when the application exits.
func CloseValidatorPool(pool chan validator.CsvRowValidator) {
	close(pool)
	for v := range pool {
		v.Close()
	}
}
