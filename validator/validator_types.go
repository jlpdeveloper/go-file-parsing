package validator

import (
	"go-file-parsing/cache"
	"go-file-parsing/config"
)

type RowError struct {
	Row   int64
	Id    string
	Error error
}
type ColValidator func(*RowValidatorContext, []string) (map[string]string, error)

type RowValidatorContext struct {
	Config *config.ParserConfig
	GetMap func() map[string]string
}

type RowValidator interface {
	Validate(row string) (string, error)
}

func New(conf *config.ParserConfig, cache cache.DistributedCache, colValidators []ColValidator) CsvRowValidator {
	return CsvRowValidator{
		config:        conf,
		cacheClient:   cache,
		colValidators: colValidators,
	}
}
