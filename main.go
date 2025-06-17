package main

import (
	"bufio"
	"fmt"
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"go-file-parsing/loan_info"
	"go-file-parsing/validator"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	cacheClient, err := cache.New()
	if err != nil {
		panic(err)
	}
	defer cacheClient.Close()
	start := time.Now()
	//readWriteValkey(cacheClient)
	//parseFile("sample.csv", cacheClient)
	parseFile("data/accepted_2007_to_2018Q4.csv", cacheClient)
	end := time.Now()
	log.Printf("Time elapsed: %s", end.Sub(start))
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
	// Create a pool of validators
	pool := loan_info.NewRowValidatorPool(&conf, cacheClient, 500)
	// Ensure validators are closed when function exits
	defer loan_info.CloseValidatorPool(pool)

	// Create a channel to receive errors
	errChan := make(chan validator.RowError, 100)
	var rowCount int64 = 0
	scanner := bufio.NewScanner(file)
	var errors []validator.RowError
	errorWg := &sync.WaitGroup{}
	errorWg.Add(1)
	go func() {
		defer errorWg.Done()
		for err := range errChan {
			errors = append(errors, err)
		}
	}()
	const maxScannerBufferSize = 1024 * 1024 // 1MB buffer
	buf := make([]byte, maxScannerBufferSize)
	scanner.Buffer(buf, maxScannerBufferSize)
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
			id, rowErr := rowVal.Validate(row)
			if rowErr != nil {
				//log.Println(rowErr.Error())
				errChan <- validator.RowError{
					Row:   rowNum,
					Id:    id,
					Error: rowErr,
				}
			}
			pool <- rowVal
		}(scanner.Text(), currentRow)
		if currentRow%10000 == 0 {
			log.Printf("Processed %d rows\n", currentRow)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
			fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
			fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
			fmt.Printf("\tNumGC = %v\n", m.NumGC)

		}
		rowCount++
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error scanning file: %v", err)
		}
	}

	wg.Wait()
	log.Println("CSV parsing complete.")
	close(errChan)
	errorWg.Wait()
	log.Println(fmt.Sprintf("Error Size:%d", len(errors)))
	//for _, err := range errors {
	//	//Cleanup all the bad data
	//	_ = cacheClient.Delete(context.Background(), err.Id)
	//	//log.Println(fmt.Sprintf("error on line: %d, error: %s", err.Row, err.Error.Error()))
	//}

}
