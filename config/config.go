package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type PBTechConfig struct {
	Filter string `json:"filter"`
}
type ComputerLoungeConfig struct {
	TitleFilter string `json:"title_filter"`
}

type Config struct {
	SearchTerm           string               `json:"search_term"`
	PBTechConfig         PBTechConfig         `json:"pb_tech"`
	ComputerLoungeConfig ComputerLoungeConfig `json:"computer_lounge"`
}

func ParseConfig(configFp string) (*Config, error) {
	config := &Config{}
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
