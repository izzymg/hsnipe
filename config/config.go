package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	SearchTerm string `json:"search_term"`
}

func NewConfig() *Config {
	return &Config{
		SearchTerm: "RTX 5070 ti",
	}
}

func ParseConfig(configFp string) (*Config, error) {
	config := NewConfig()
	file, err := os.Open(configFp)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	if config.SearchTerm == "" {
		return nil, fmt.Errorf("search term is empty")
	}

	return config, nil
}
