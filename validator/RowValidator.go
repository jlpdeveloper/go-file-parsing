package validator

import (
	"fmt"
	"go-file-parsing/config"
	"strings"
)

type RowError struct {
	Row int64
	Err string
}

type RowValidator interface {
	Validate(row string) error
	//Validate(row string, errChan <-chan RowError) error
}

type CsvRowValidator struct {
	config  *config.Config
	columns []string
}

func (c CsvRowValidator) Validate(row string) error {
	c.columns = strings.Split(row, c.config.Delimiter)
	c.isValidSize()
	return nil
}

func (c CsvRowValidator) isValidSize() {
	if len(c.columns) != c.config.ExpectedColumns {
		println(fmt.Sprintf("Id:%s Invalid size", c.columns[0]))
	}
}

func NewCsvRowValidator(conf *config.Config) CsvRowValidator {
	return CsvRowValidator{
		config: conf,
	}
}

func NewCsvRowValidatorPool(conf *config.Config, poolSize int) chan CsvRowValidator {
	pool := make(chan CsvRowValidator, poolSize)
	for i := 0; i < poolSize; i++ {
		pool <- CsvRowValidator{
			config: conf,
		}
	}
	return pool
}
