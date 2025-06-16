package validator

import (
	"go-file-parsing/config"
	"golang.org/x/sync/errgroup"
	"strings"
)

type CsvRowValidator struct {
	config        *config.ParserConfig
	colValidators []ColValidator
	cacheChan     chan CacheData
}

func (c *CsvRowValidator) Validate(row string) (string, error) {
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

	return cols[0], g.Wait() // returns the first error (if any), cancels other goroutines
}
