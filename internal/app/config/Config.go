package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr   string
	NasaAPIKey string
	DSN        string
}

func Read() (*Config, error) {
	var config Config
	httpAddr, exists := os.LookupEnv("HTTP_ADDR")
	if exists {
		config.HTTPAddr = httpAddr
	} else {
		config.HTTPAddr = ":8080"
	}
	nasaAPIKey, exists := os.LookupEnv("NASA_API_KEY")
	if exists {
		config.NasaAPIKey = nasaAPIKey
	} else {
		return nil, fmt.Errorf("NASA_API_KEY environment variable not set")
	}

	dsn, exists := os.LookupEnv("DSN")
	if exists {
		config.DSN = dsn
	} else {
		config.DSN = "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	}
	return &config, nil
}
