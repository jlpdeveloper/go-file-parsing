package main

import (
	"bufio"
	"context"
	"fmt"
	"go-file-parsing/cache"
	"go-file-parsing/config"
	"go-file-parsing/loan_info"
	"go-file-parsing/validator"
	"log"
	"os"
	"runtime"
	"strconv"
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

	// Default file to parse
	fileToProcess := "data/accepted_2007_to_2018Q4.csv"

	// Check if a file was specified as a command-line argument
	if len(os.Args) > 1 {
		fileToProcess = os.Args[1]
		log.Printf("Using file specified by command-line argument: %s", fileToProcess)
	} else {
		log.Printf("No file specified, using default: %s", fileToProcess)
	}

	parseFile(fileToProcess, cacheClient)
	end := time.Now()
	log.Printf("Time elapsed: %s", end.Sub(start))
}

func NewErrChan(cache cache.DistributedCache, size int, wg *sync.WaitGroup) chan validator.RowError {
	errChan := make(chan validator.RowError, size)
	errWorkerPool := make(chan func(validator.RowError), size)
	for i := 0; i < size; i++ {
		errWorkerPool <- func(err validator.RowError) {
			cacheErr := cache.Set(context.Background(), fmt.Sprintf("err:row%s:id%s", strconv.FormatInt(err.Row, 10), err.Id), err.Error.Error())
			if cacheErr != nil {
				log.Printf("Error writing to cache: %v", cacheErr)
			}
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range errChan {
			worker := <-errWorkerPool
			go func(e validator.RowError) {
				worker(e)
				errWorkerPool <- worker
			}(err)
		}
	}()
	return errChan
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
	chanWg := &sync.WaitGroup{}
	cacheChan := validator.NewCacheChannel(cacheClient, chanWg, 10000)
	// Create a pool of validators
	pool := loan_info.NewRowValidatorPool(&conf, cacheChan, 1000)
	// Ensure validators are closed when function exits
	defer loan_info.CloseValidatorPool(pool)

	// Create a channel to receive errors
	errChan := NewErrChan(cacheClient, 10000, chanWg)
	var rowCount int64 = 0
	scanner := bufio.NewScanner(file)
	const maxScannerBufferSize = 1024 * 1024 // 1MB buffer
	buf := make([]byte, maxScannerBufferSize)
	scanner.Buffer(buf, maxScannerBufferSize)
	wg := &sync.WaitGroup{}

	times := make([]int, 10)
	prevTime := time.Now()
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

			now := time.Now()
			diffMs := now.Sub(prevTime).Milliseconds()
			times = append(times, int(diffMs))
			prevTime = now

		}
		rowCount++

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file: %v", err)
	}

	wg.Wait()
	log.Println("CSV parsing complete.")
	close(errChan)
	close(cacheChan)
	chanWg.Wait()
	log.Println("Finished writing to cache.")
	avgTime := 0
	for _, t := range times {
		avgTime += t
	}
	avgTime /= len(times)
	log.Printf("Average time per 10,000 rows: %dms", avgTime)
	log.Printf("Total rows: %d", rowCount)
}
