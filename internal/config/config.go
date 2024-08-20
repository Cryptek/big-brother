package config

import (
	"big-brother/internal/models"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func LoadConfig(configFilePath string) (*models.Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg models.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
