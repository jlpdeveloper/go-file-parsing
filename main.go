package main

import (
	"bufio"
	"context"
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"os"
)

func main() {
	cacheClient, err := cache.New()
	if err != nil {
		panic(err)
	}
	defer cacheClient.Close()

	readWriteValkey(cacheClient)
	parseFile("sample.csv")

}

func readWriteValkey(cacheClient cache.DistributedCache) {
	ctx := context.Background()
	cacheClient.Set(ctx, "Hello", "World")
	println("Write Done")
	val, _ := cacheClient.Get(ctx, "Hello")

	println(val)
	//cacheClient.Do(ctx, cacheClient.B().Hsetnx().Key("test").Field("total").Value("1").Build())
	//cacheClient.Do(ctx, cacheClient.B().Hincrby().Key("test").Field("total").Increment(10).Build())
}

func parseFile(filename string) {
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

	var rowCount int64 = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if rowCount == 0 && conf.HasHeader {
			rowCount++
			continue
		}
		rowVal := validator.NewCsvRowValidator(&conf)
		_ = rowVal.Validate(scanner.Text())
		//println(scanner.Text())
		rowCount++
	}
}
