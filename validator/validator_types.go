package validator

import (
	"go-file-parsing/cache"
	"go-file-parsing/config"
)

type RowError struct {
	Row   int64
	Error error
}
type ColValidator func(*RowValidatorContext, *[]string) error

type RowValidatorContext struct {
	Config *config.ParserConfig
	Cache  cache.DistributedCache
}

type RowValidator interface {
	Validate(row string) error
}

func NewCsvRowValidatorPool(conf *config.ParserConfig, cache cache.DistributedCache, poolSize int) chan CsvRowValidator {
	pool := make(chan CsvRowValidator, poolSize)
	for i := 0; i < poolSize; i++ {
		pool <- CsvRowValidator{
			config:        conf,
			cacheClient:   cache,
			colValidators: validators,
		}
	}
	return pool
}

var validators = []ColValidator{
	isValidSize,
}
