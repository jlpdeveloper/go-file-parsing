package validator

import (
	"context"
	"fmt"
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"strings"
	"sync"
)

type RowError struct {
	Row int64
	Err string
}

func (r RowError) Error() string {
	return fmt.Sprintf("Row:%d Error:%s", r.Row, r.Err)
}

type ColValidator func(cols *[]string, errChan chan<- RowError)

type RowValidator interface {
	Validate(row string) error
	//Validate(row string, errChan <-chan RowError) error
}

type CsvRowValidator struct {
	config      *config.ParserConfig
	cacheClient cache.DistributedCache
}

func (c CsvRowValidator) Validate(row string) error {
	cols := strings.Split(row, c.config.Delimiter)
	colValidators := []ColValidator{
		c.isValidSize,
	}
	errorChan := make(chan RowError)
	baseCtx := context.Background()
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()
	wg := &sync.WaitGroup{}
	go func() {
		select {
		case <-errorChan:
			cancel()
		}
	}()
	for _, colValidator := range colValidators {
		wg.Add(1)
		go func() {
			ch := make(chan bool)
			defer wg.Done()
			go func() {
				colValidator(&cols, errorChan)
				ch <- true
			}()
			select {
			case <-ctx.Done():
				wg.Done()
				return
			case <-ch:
				return
			}
		}()
	}
	wg.Wait()
	close(errorChan)
	for err := range errorChan {
		return err
	}
	return nil
}

func (c CsvRowValidator) isValidSize(cols *[]string, errChan chan<- RowError) {
	if len(*cols) != c.config.ExpectedColumns {
		errChan <- RowError{
			Row: 1,
			Err: "Invalid size",
		}
	}
}

func NewCsvRowValidator(conf *config.ParserConfig) CsvRowValidator {
	return CsvRowValidator{
		config: conf,
	}
}

func NewCsvRowValidatorPool(conf *config.ParserConfig, cache cache.DistributedCache, poolSize int) chan CsvRowValidator {
	pool := make(chan CsvRowValidator, poolSize)
	for i := 0; i < poolSize; i++ {
		pool <- CsvRowValidator{
			config:      conf,
			cacheClient: cache,
		}
	}
	return pool
}
