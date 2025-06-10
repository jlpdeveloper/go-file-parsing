package validator

import (
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"golang.org/x/sync/errgroup"
	"strings"
)

type CsvRowValidator struct {
	config        *config.ParserConfig
	cacheClient   cache.DistributedCache
	colValidators []ColValidator
}

func (c *CsvRowValidator) Validate(row string) error {
	cols := strings.Split(row, c.config.Delimiter)

	vCtx := RowValidatorContext{
		Config: c.config,
		Cache:  c.cacheClient,
	}

	var g errgroup.Group

	for _, validator := range c.colValidators {
		f := validator // capture
		g.Go(func() error {
			return f(&vCtx, cols)
		})
	}

	return g.Wait() // returns first error (if any), cancels other goroutines
}
