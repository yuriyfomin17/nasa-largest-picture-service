package config

import "os"

type Config struct {
	HTTPAddr   string
	NasaAPIKey string
}

func Read() (Config, error) {
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
	}
	return config, nil
}
