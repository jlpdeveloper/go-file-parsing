package validator

import (
	"fmt"
	"go-file-parsing/config"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
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
	mu := sync.Mutex{}
	m := vCtx.GetMap()
	m["id"] = id
	m["raw"] = row

	var g errgroup.Group

	for _, validator := range c.colValidators {
		f := validator // capture
		g.Go(func() error {
			data, err := f(&vCtx, cols)
			if err != nil {
				return err
			}
			if data != nil {
				mu.Lock()
				for k, v := range data {
					m[k] = v
				}
				PutMap(data)
				mu.Unlock()
			}
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		PutMap(m)
		return id, err
	}

	c.cacheChan <- CacheData{
		Id:   id,
		Data: m,
	}

	return id, err // returns the first error (if any), cancels other goroutines
}

// Close closes the validator and releases resources.
// It should be called when the validator is no longer needed.
func (c *CsvRowValidator) Close() {
	if !c.closed {
		c.closed = true
	}
}
