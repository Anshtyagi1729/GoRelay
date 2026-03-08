package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port                  string   `json:"port"`
	Backends              []string `json:"backends"`
	Rate_limit            int      `json:"rate_limit"`
	Window_seconds        int      `json:"window_second"`
	Timeout_second        int      `json:"timeout_second"`
	Health_check_interval int      `json:"health_check_interval"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
