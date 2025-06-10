package validator

import (
	"context"
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
	//First we create a channel to write to the cache
	ctx := context.Background()
	//Create a channel for writing to the cache. Defer closing it
	cacheChan := make(chan map[string]string, 100)
	defer close(cacheChan)
	//Spin off a new goroutine that will write to the cache as it processes from the channel
	go func() {
		for data := range cacheChan {
			for key, value := range data {
				_ = c.cacheClient.SetField(ctx, key, value, row)
			}
		}
	}()
	//Split the columns, then set the first value to the raw data string (for debug purposes)
	cols := strings.Split(row, c.config.Delimiter)
	_ = c.cacheClient.SetField(ctx, cols[0], "raw", row)
	_ = c.cacheClient.SetField(ctx, cols[0], "id", cols[0])

	vCtx := RowValidatorContext{
		Config: c.config,
	}
	var g errgroup.Group

	for _, validator := range c.colValidators {
		f := validator // capture
		g.Go(func() error {
			data, err := f(&vCtx, cols)
			if err != nil {
				return err
			}
			if data != nil {
				cacheChan <- data
			}
			return nil
		})
	}

	return g.Wait() // returns first error (if any), cancels other goroutines
}
