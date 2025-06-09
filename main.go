package main

import (
	"bufio"
	"context"
	"fmt"
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"log"
	"os"
	"sync"
)

func main() {
	cacheClient, err := cache.New()
	if err != nil {
		panic(err)
	}
	defer cacheClient.Close()

	readWriteValkey(cacheClient)
	parseFile("sample.csv", cacheClient)

}

func readWriteValkey(cacheClient cache.DistributedCache) {
	ctx := context.Background()
	err := cacheClient.Set(ctx, "Hello", "World")
	if err != nil {
		panic(err)
	}
	println("Write Done")
	val, _ := cacheClient.Get(ctx, "Hello")

	println(val)
	//cacheClient.Do(ctx, cacheClient.B().Hsetnx().Key("test").Field("total").Value("1").Build())
	//cacheClient.Do(ctx, cacheClient.B().Hincrby().Key("test").Field("total").Increment(10).Build())
}

func parseFile(filename string, cacheClient cache.DistributedCache) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		fcErr := file.Close()
		if fcErr != nil {
			panic(fcErr)
		}
	}()
	conf, err := config.LoadParserConfig("config.json")
	if err != nil {
		panic(err)
	}
	pool := validator.NewCsvRowValidatorPool(&conf, cacheClient, 100)
	errChan := make(chan validator.RowError)
	var rowCount int64 = 0
	scanner := bufio.NewScanner(file)
	wg := &sync.WaitGroup{}
	for scanner.Scan() {
		currentRow := rowCount
		if currentRow == 0 && conf.HasHeader {
			rowCount++
			continue
		}
		rowVal := <-pool
		wg.Add(1)
		go func(row string, rowNum int64) {
			defer wg.Done()
			rowErr := rowVal.Validate(row)
			if rowErr != nil {
				log.Println(rowErr.Error())
				errChan <- validator.RowError{
					Row:   rowNum,
					Error: rowErr,
				}
			}
			pool <- rowVal
		}(scanner.Text(), currentRow)
		rowCount++
	}
	var errors []validator.RowError
	errorWg := &sync.WaitGroup{}
	errorWg.Add(1)
	go func() {
		defer errorWg.Done()
		for err := range errChan {
			errors = append(errors, err)
		}
	}()
	wg.Wait()
	log.Println("CSV parsing complete.")
	close(errChan)
	for _, err := range errors {
		log.Println(fmt.Sprintf("error on line: %d, error: %s", err.Row, err.Error.Error()))
	}

}
