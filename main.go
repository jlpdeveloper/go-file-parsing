package main

import (
	"bufio"
	"go-file-parsing/config"
	"go-file-parsing/validator"
	"os"
)

func main() {
	println("Hello World")
	file, _ := os.Open("sample.csv")
	defer func() {
		fcErr := file.Close()
		if fcErr != nil {
			panic(fcErr)
		}
	}()
	conf, err := config.LoadConfig("config.json")
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
