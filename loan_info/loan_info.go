package loan_info

import (
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"go-file-parsing/validator"
)

var validators = []validator.ColValidator{
	//isValidSize,
	//hasValidLoanAmount,
	//hasValidInterestRate,
	//hasValidTerm,
	//hasEmploymentInfo,
	hasLowDTIAndHomeOwnership,
	//hasEstablishedCreditHistory,
	//hasHealthyFICOScore,
	//hasSufficientAccounts,
	//hasStableEmployment,
	//hasNoPublicRecordOrBankruptcies,
	//isVerifiedWithIncome,
	//hasValidGradeSubgrade,
	//passExtraData,
}

func NewRowValidatorPool(conf *config.ParserConfig, cache cache.DistributedCache, poolSize int) chan validator.CsvRowValidator {
	cacheChan := validator.NewCacheChannel(cache)
	pool := make(chan validator.CsvRowValidator, poolSize)
	for i := 0; i < poolSize; i++ {
		pool <- validator.New(conf, cacheChan, validators)
	}
	return pool
}
