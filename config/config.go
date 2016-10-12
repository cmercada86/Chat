package config

import (
	"encoding/json"
	"log"
	"os"
)

var config *Config

type Config struct {
	GogClientID     string `json:"google_client_id"`
	GogClientSecret string `json:"google_client_secret"`
	RedisURL        string `json:"redis_url"`
}

func GetConfig() *Config {
	return config
}

func ReadConfFile(filename string) {

	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Could not open config file:", filename, err)
	}
	defer configFile.Close()

	jp := json.NewDecoder(configFile)
	if err = jp.Decode(config); err != nil {
		log.Fatal("Error parsing config file", err)
	}

}
