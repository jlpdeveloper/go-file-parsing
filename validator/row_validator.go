package validator

import (
	"fmt"
	"go-file-parsing/config"
	"golang.org/x/sync/errgroup"
	"strings"
)

type CsvRowValidator struct {
	config        *config.ParserConfig
	colValidators []ColValidator
	cacheChan     chan CacheData
	closed        bool
}

func (c *CsvRowValidator) Validate(row string) (string, error) {
	if c.closed {
		return "", fmt.Errorf("validator is closed")
	}

	//Split the columns, then set the first value to the raw data string (for debug purposes)
	cols := PreprocessColumns(strings.Split(row, c.config.Delimiter))
	id := cols[0]

	vCtx := RowValidatorContext{
		Config: c.config,
		GetMap: getMap,
	}

	m := vCtx.GetMap()
	m["id"] = id
	m["raw"] = row
	c.cacheChan <- CacheData{
		Id:   id,
		Data: m,
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
				c.cacheChan <- CacheData{
					Id:   id,
					Data: data,
				}
			}
			return nil
		})
	}

	return id, g.Wait() // returns the first error (if any), cancels other goroutines
}

// Close closes the validator and releases resources.
// It should be called when the validator is no longer needed.
func (c *CsvRowValidator) Close() {
	if !c.closed {
		c.closed = true
	}
}
