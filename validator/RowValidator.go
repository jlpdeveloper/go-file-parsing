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
	config *config.Config
}

func (c CsvRowValidator) Validate(row string) error {
	cols := strings.Split(row, c.config.Delimiter)
	c.isValidSize(&cols)
	return nil
}

func (c CsvRowValidator) isValidSize(cols *[]string) {
	if len(*cols) != c.config.ExpectedColumns {
		println(fmt.Sprintf("Id:%s Invalid size", (*cols)[0]))
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
