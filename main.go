package main

import (
	"bufio"
	"context"
	"fmt"
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"go-file-parsing/validator"
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
	file, _ := os.Open(filename)
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
	errChan := make(chan validator.RowError, 100)
	var rowCount int64 = 0
	scanner := bufio.NewScanner(file)
	wg := &sync.WaitGroup{}
	for scanner.Scan() {
		if rowCount == 0 && conf.HasHeader {
			rowCount++
			continue
		}
		rowVal := <-pool
		wg.Add(1)
		go func(row string) {
			defer wg.Done()
			rowErr := rowVal.Validate(row)
			if rowErr != nil {
				println(rowErr.Error())
				errChan <- validator.RowError{
					Row:   rowCount,
					Error: rowErr,
				}
			}
			pool <- rowVal
		}(scanner.Text())
		rowCount++
	}

	wg.Wait()
	println("CSV parsing complete.")
	close(errChan)
	for err := range errChan {
		println(fmt.Sprintf("error on line: %d, error: %s", err.Row, err.Error.Error()))
	}

}
