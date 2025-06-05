package config

import (
	"encoding/json"
	"os"
)

type ParserConfig struct {
	Delimiter       string
	ExpectedColumns int
	HasHeader       bool
}

func LoadParserConfig(filename string) (ParserConfig, error) {
	var cfg ParserConfig
	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
