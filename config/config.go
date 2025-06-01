package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Delimiter       string
	ExpectedColumns int
	HasHeader       bool
}

func LoadConfig(filename string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
