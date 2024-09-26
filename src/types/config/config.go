package config

import (
	"encoding/json"
	"l0/types/logger"
	"os"
)

type Config struct {
	DbPassword string
	DbUser     string
	DbHost     string
	DbPort     string
	KafkaHost  string
	KafkaPort  string
	KafkaTopic string
}

func GetConfig() (*Config, error) {
	data, err := os.ReadFile("/home/svovoniks/Desktop/L0_go/src/cfg.json")
	if err != nil {
		logger.Logger.Warn().
			Msg("Config file '.cfg' not found")
		return nil, err
	}

	var cfg Config

	if err = json.Unmarshal(data, &cfg); err != nil {
		logger.Logger.Warn().
			Msg("Couldn't parse config")
		return nil, err
	}

	return &cfg, nil

}
