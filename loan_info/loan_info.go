package loan_info

import (
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"go-file-parsing/validator"
)

var validators = []validator.ColValidator{
	isValidSize,
	hasValidLoanAmount,
	hasValidInterestRate,
}

func NewRowValidatorPool(conf *config.ParserConfig, cache cache.DistributedCache, poolSize int) chan validator.CsvRowValidator {
	pool := make(chan validator.CsvRowValidator, poolSize)
	for i := 0; i < poolSize; i++ {
		pool <- validator.New(conf, cache, validators)
	}
	return pool
}
